package com.example.springsdk.controller;

import com.example.springsdk.service.ByzantineDemoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * REST layer for the Byzantine demonstration page.
 *
 * Same {success, data, message} envelope as FabricController so the frontend can
 * treat every endpoint the same way.
 */
@RestController
@RequestMapping("/api/fabric/demo")
@CrossOrigin(origins = "*")
public class ByzantineDemoController {

    @Autowired
    private ByzantineDemoService demoService;

    private ResponseEntity<Map<String, Object>> ok(String message, Object data) {
        Map<String, Object> response = new HashMap<>();
        response.put("success", true);
        response.put("data", data);
        response.put("message", message);
        return ResponseEntity.ok(response);
    }

    private ResponseEntity<Map<String, Object>> fail(String message) {
        Map<String, Object> response = new HashMap<>();
        response.put("success", false);
        response.put("data", null);
        response.put("message", message);
        return ResponseEntity.ok(response);
    }

    /** Arm or disarm nodes: {"nodes":[1,2],"mode":"tamper"}. An empty list means all. */
    @PostMapping("/fault")
    public ResponseEntity<Map<String, Object>> setFault(@RequestBody Map<String, Object> body) {
        try {
            String mode = String.valueOf(body.getOrDefault("mode", "none"));
            List<Integer> nodes = new ArrayList<>();
            Object raw = body.get("nodes");
            if (raw instanceof List) {
                for (Object item : (List<?>) raw) {
                    nodes.add(Integer.parseInt(String.valueOf(item)));
                }
            }
            Map<String, Object> result = demoService.setFault(nodes, mode);
            if (result.get("rejected") != null) {
                return fail(String.valueOf(result.get("rejected")));
            }
            return ok("fault mode " + mode + " applied", result);
        } catch (Exception e) {
            return fail("Failed to set fault mode: " + e.getMessage());
        }
    }

    /** Topology and current mode of every consensus node. */
    @GetMapping("/nodes")
    public ResponseEntity<Map<String, Object>> nodes() {
        try {
            Map<String, Object> data = new HashMap<>();
            data.put("nodeCount", ByzantineDemoService.NODE_COUNT);
            data.put("faultTolerance", (ByzantineDemoService.NODE_COUNT - 1) / 2);
            data.put("nodes", demoService.nodeStates());
            return ok("Consensus node states retrieved successfully", data);
        } catch (Exception e) {
            return fail("Failed to read node states: " + e.getMessage());
        }
    }

    /** One demonstration round: submit a small transaction and report the outcome. */
    @PostMapping("/run")
    public ResponseEntity<Map<String, Object>> run(@RequestParam(value = "waitMs", defaultValue = "0") long waitMs) {
        try {
            return ok("Demonstration round finished", demoService.runRound(waitMs));
        } catch (Exception e) {
            return fail("Demonstration round failed: " + e.getMessage());
        }
    }

    /** Result panel: per-node counters, key-window alignment and the QDS trace. */
    @GetMapping("/consensus")
    public ResponseEntity<Map<String, Object>> consensus(
            @RequestParam(value = "since", defaultValue = "0") int since,
            @RequestParam(value = "limit", defaultValue = "60") int limit) {
        try {
            return ok("Consensus summary retrieved successfully", demoService.consensusSummary(since, limit));
        } catch (Exception e) {
            return fail("Failed to read consensus summary: " + e.getMessage());
        }
    }

    /** QDS step trace of one node. */
    @GetMapping("/trace")
    public ResponseEntity<Map<String, Object>> trace(
            @RequestParam("node") int node,
            @RequestParam(value = "since", defaultValue = "0") int since,
            @RequestParam(value = "limit", defaultValue = "100") int limit) {
        try {
            if (node < 0 || node >= ByzantineDemoService.NODE_COUNT) {
                return fail("Node index must be between 0 and " + (ByzantineDemoService.NODE_COUNT - 1));
            }
            return ok("Node trace retrieved successfully", demoService.trace(node, since, limit));
        } catch (Exception e) {
            return fail("Failed to read node trace: " + e.getMessage());
        }
    }
}
