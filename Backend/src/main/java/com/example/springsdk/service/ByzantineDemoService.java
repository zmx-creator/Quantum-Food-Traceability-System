package com.example.springsdk.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Control and status layer for the Byzantine demonstration.
 *
 * Each orderer embeds the QBA consensus plugin, whose HTTP listener answers on
 * ports 6070..6074 of the network host with three administration endpoints:
 *
 *   POST /fault?mode=none|tamper          arm or disarm one node
 *   POST /faultpeers?ids=1,2            which nodes are armed (the collusion map)
 *   GET  /status                          mode, injection counter, round state
 *   GET  /trace?since=&limit=             the QDS steps of the running round
 *
 * This service fans those calls out and folds the answers into one payload, so
 * the demonstration page shows the QDS round, the malicious node behaviour and
 * the ledger outcome without reading container logs. The only ledger writing it
 * performs is the single small transaction started by runRound().
 */
@Service
public class ByzantineDemoService {

    /** QBA consensus nodes, one per orderer. */
    public static final int NODE_COUNT = 5;

    /** Node i answers on ADMIN_PORT_BASE + i on the network host. */
    private static final int ADMIN_PORT_BASE = 6070;

    /** Fault modes the consensus nodes understand. */
    private static final List<String> MODES = Arrays.asList("none", "tamper");

    /** A failed round resolves in about six seconds on the demonstration network. */
    private static final long DEFAULT_ROUND_WAIT_MS = 12000;

    private static final String HOST = System.getProperty("fabric.vm-ip",
            System.getenv("FABRIC_VM_IP") != null ? System.getenv("FABRIC_VM_IP") : "localhost");

    /**
     * The last demonstration round this process ran. The nodes keep their trace and
     * the ledger keeps the block, so after a page reload only the payload would be
     * lost - it is known here and nowhere else. The workflow strip reads it back
     * through the summary, which is also why a reload cannot show a stale round:
     * the memo carries the round instance it belonged to.
     */
    private volatile Map<String, Object> lastRound;

    @Autowired
    private FabricService fabricService;

    private final RestTemplate http;
    private final ObjectMapper mapper = new ObjectMapper();

    public ByzantineDemoService() {
        SimpleClientHttpRequestFactory factory = new SimpleClientHttpRequestFactory();
        factory.setConnectTimeout(2000);
        factory.setReadTimeout(4000);
        this.http = new RestTemplate(factory);
    }

    /** Admin base URL of one consensus node. */
    public String adminUrl(int node) {
        return "http://" + HOST + ":" + (ADMIN_PORT_BASE + node);
    }

    /** Arm or disarm the given nodes. */
    public Map<String, Object> setFault(List<Integer> nodes, String mode) {
        Map<String, Object> result = new HashMap<>();
        if (!MODES.contains(mode)) {
            result.put("mode", mode);
            result.put("rejected", "unknown mode, expected none or tamper");
            result.put("results", new ArrayList<Map<String, Object>>());
            return result;
        }
        if (!"none".equals(mode) && nodes != null && nodes.contains(0)) {
            result.put("mode", mode);
            result.put("rejected", "node 0 is the commander of every round: the demonstration arms deviating "
                    + "deputies only, and no injected strategy makes a commander behave badly");
            result.put("results", new ArrayList<Map<String, Object>>());
            return result;
        }
        List<Integer> targets = (nodes == null || nodes.isEmpty()) ? allNodes() : nodes;
        if (!"none".equals(mode)) {
            //an empty list means everybody, and the commander is never everybody here
            List<Integer> deputies = new ArrayList<>();
            for (Integer node : targets) {
                if (node != 0) {
                    deputies.add(node);
                }
            }
            targets = deputies;
        }
        List<Map<String, Object>> perNode = new ArrayList<>();
        for (Integer node : targets) {
            Map<String, Object> entry = new HashMap<>();
            entry.put("node", node);
            try {
                String body = http.postForObject(adminUrl(node) + "/fault?mode=" + mode, null, String.class);
                entry.put("response", parse(body));
                entry.put("ok", true);
            } catch (Exception e) {
                entry.put("ok", false);
                entry.put("error", e.getMessage());
            }
            perNode.add(entry);
        }
        result.put("mode", mode);
        result.put("results", perNode);
        //the strategy the nodes act on depends on who their partners are, so the
        //hostile set is pushed to every node as part of arming or disarming them
        String peersError = pushHostilePeers();
        if (peersError != null) {
            result.put("peersError", peersError);
        }
        return result;
    }

    /**
     * Tell every node which nodes the demonstration has armed. The malicious
     * strategy of the protocol analysis is collusion-aware: a Byzantine signer
     * offers a conflicting message only to a relay that is Byzantine too, and a
     * Byzantine relay forwards a partner's message untouched instead of forging it.
     * Both rules need the same view of the hostile set, so it is pushed after every
     * change. Returns null when all five accepted, otherwise the first failure.
     */
    private String pushHostilePeers() {
        StringBuilder hostile = new StringBuilder();
        for (Map<String, Object> state : nodeStates()) {
            if (!"none".equals(String.valueOf(state.get("mode")))) {
                if (hostile.length() > 0) {
                    hostile.append(',');
                }
                hostile.append(state.get("node"));
            }
        }
        for (Integer node : allNodes()) {
            try {
                http.postForObject(adminUrl(node) + "/faultpeers?ids=" + hostile, null, String.class);
            } catch (Exception e) {
                return "node " + node + ": " + e.getMessage();
            }
        }
        return null;
    }

    /** Status of every consensus node, including which ones are Byzantine. */
    public List<Map<String, Object>> nodeStates() {
        List<Map<String, Object>> states = new ArrayList<>();
        for (Integer node : allNodes()) {
            Map<String, Object> entry = new HashMap<>();
            entry.put("node", node);
            entry.put("admin", adminUrl(node));
            try {
                Map<String, Object> status = parse(http.getForObject(adminUrl(node) + "/status", String.class));
                entry.putAll(status);
                entry.put("reachable", true);
            } catch (Exception e) {
                entry.put("reachable", false);
                entry.put("error", e.getMessage());
            }
            states.add(entry);
        }
        return states;
    }

    /** Trace steps recorded by one node, for the round timeline. */
    public Map<String, Object> trace(int node, int since, int limit) {
        Map<String, Object> out = new HashMap<>();
        try {
            String body = http.getForObject(adminUrl(node) + "/trace?since=" + since + "&limit=" + limit, String.class);
            out = parse(body);
        } catch (Exception e) {
            out.put("node", node);
            out.put("error", e.getMessage());
        }
        return out;
    }

    /** Trace steps of every node, indexed by node id. */
    public Map<String, Object> traceAll(int since, int limit) {
        Map<String, Object> all = new HashMap<>();
        for (Integer node : allNodes()) {
            all.put(String.valueOf(node), trace(node, since, limit));
        }
        return all;
    }

    /**
     * Run one demonstration round: remember the height, submit a single small
     * transaction (one production record, so the batch stays inside the QDS
     * signature window) and report whether the network agreed on it.
     */
    public Map<String, Object> runRound(long waitMillis) throws Exception {
        long wait = waitMillis > 0 ? waitMillis : DEFAULT_ROUND_WAIT_MS;
        long before = fabricService.getBlockHeight();
        String foodID = "BYZ-" + System.currentTimeMillis();

        //The record is an ordinary production entry, so the demonstration round is
        //the same transaction a user submits from the form: proposal, endorsement,
        //ordering, ledger. The page shows these fields as step 1 of that flow.
        String[] record = {
                foodID, "Organic free-range eggs", "30 eggs per carton",
                "2026-10-20", "2027-10-20", "LOT-2026-1020",
                "QSID-8842", "Green Valley Farm", "12.50", "Zhejiang"
        };

        long started = System.currentTimeMillis();
        String submit = fabricService.submitTransactionBounded(wait, "addProInfo", record);
        long waited = System.currentTimeMillis() - started;

        long after = fabricService.getBlockHeight();
        if (after == before) {
            // the commit event can trail the orderer write by a moment
            for (int i = 0; i < 6 && after == before; i++) {
                Thread.sleep(500);
                after = fabricService.getBlockHeight();
            }
        }

        Map<String, Object> out = new HashMap<>();
        out.put("foodID", foodID);
        Map<String, String> payload = new LinkedHashMap<>();
        payload.put("foodID", record[0]);
        payload.put("foodName", record[1]);
        payload.put("foodSpec", record[2]);
        payload.put("foodMFGDate", record[3]);
        payload.put("foodEXPDate", record[4]);
        payload.put("foodLOT", record[5]);
        payload.put("foodQSID", record[6]);
        payload.put("foodMFRSName", record[7]);
        payload.put("foodProPrice", record[8]);
        payload.put("foodProPlace", record[9]);
        out.put("payload", payload);
        out.put("heightBefore", before);
        out.put("heightAfter", after);
        out.put("blockWritten", after > before);
        out.put("blockNumber", after > before ? after - 1 : null);
        out.put("agreed", after > before);
        out.put("waitedMs", waited);
        out.put("submitMessage", submit);
        List<Map<String, Object>> states = nodeStates();
        out.put("nodes", states);

        Map<String, Object> memo = new LinkedHashMap<>();
        memo.put("foodID", foodID);
        memo.put("payload", payload);
        memo.put("heightBefore", before);
        memo.put("heightAfter", after);
        memo.put("blockNumber", out.get("blockNumber"));
        memo.put("agreed", after > before);
        for (Map<String, Object> n : states) {
            if (n.get("node") instanceof Number && ((Number) n.get("node")).intValue() == 0) {
                memo.put("roundInstance", n.get("roundInstance"));
            }
        }
        lastRound = memo;
        return out;
    }

    /**
     * Aggregate view for the result panel: per-node counters plus the round
     * verdict derived from the ledger itself.
     */
    public Map<String, Object> consensusSummary(int traceSince, int traceLimit) throws Exception {
        List<Map<String, Object>> nodes = nodeStates();
        int malicious = 0;
        int inconsistencies = 0;
        int timeouts = 0;
        int failures = 0;
        int retransmissions = 0;
        boolean allAligned = true;
        Long window = null;

        for (Map<String, Object> n : nodes) {
            if (!"none".equals(String.valueOf(n.get("mode")))) {
                malicious++;
            }
            timeouts += asInt(n.get("timeouts"));
            inconsistencies += asInt(n.get("consistencyTotal"));
            failures += asInt(n.get("safeFailures"));
            retransmissions += asInt(n.get("retransmissions"));
            long count = asInt(n.get("consensusCount"));
            if (window == null) {
                window = count;
            } else if (window != count) {
                allAligned = false;
            }
        }

        Map<String, Object> out = new HashMap<>();
        out.put("height", fabricService.getBlockHeight());
        out.put("nodeCount", NODE_COUNT);
        // QDS gives authenticated broadcast, so the fault bound follows the
        // signed agreement condition n >= 2f+1: five nodes tolerate two
        // Byzantine nodes. The commit rule matches that bound - the primary
        // writes the batch once n-f reports name it - which is why one or two
        // malicious nodes are outvoted while three leave no quorum.
        int faultTolerance = (NODE_COUNT - 1) / 2;
        out.put("faultTolerance", faultTolerance);
        out.put("quorum", NODE_COUNT - faultTolerance);
        out.put("commitRule", "quorum of " + (NODE_COUNT - faultTolerance) + " - the primary commits once n-f reports name the proposed batch");
        out.put("faultModelNote", "QBA signs every broadcast with QDS, so the bound is n >= 2f+1: with five nodes f = 2. The primary commits when " + (NODE_COUNT - faultTolerance) + " of " + NODE_COUNT + " reports agree on the proposed batch, so one or two Byzantine nodes are outvoted and the block is still written; three of them leave the honest nodes with two reports, and after five attempts the batch is discarded as untrustworthy.");
        out.put("maliciousNodes", malicious);
        out.put("timeouts", timeouts);
        out.put("safeFailures", failures);
        out.put("retransmissions", retransmissions);
        out.put("consistencyViolations", inconsistencies);
        out.put("keyWindowsAligned", allAligned);
        out.put("nodes", nodes);
        out.put("lastRound", lastRound);
        out.put("trace", traceAll(traceSince, traceLimit));
        return out;
    }

    private List<Integer> allNodes() {
        List<Integer> nodes = new ArrayList<>();
        for (int i = 0; i < NODE_COUNT; i++) {
            nodes.add(i);
        }
        return nodes;
    }

    private int asInt(Object value) {
        return value instanceof Number ? ((Number) value).intValue() : 0;
    }

    @SuppressWarnings("unchecked")
    private Map<String, Object> parse(String json) throws Exception {
        return mapper.readValue(json, Map.class);
    }
}
