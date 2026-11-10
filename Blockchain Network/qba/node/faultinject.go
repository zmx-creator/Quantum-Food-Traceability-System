package node

// Runtime fault injection used by the malicious-node demonstration page.
//
// A consensus node can be switched between two behaviours while the orderer
// process keeps running, through POST /fault on the QBA HTTP listener:
//
//	none    participate honestly (the default, and the only mode the
//	        performance experiments ever use)
//	tamper  the malicious strategy the protocol analysis assumes, and nothing
//	        beyond it: a malicious forwarder relays what it received untouched,
//	        because unforgeability of QDS holds whenever two parties of the
//	        instance are honest, so the broadcasting list of the honest nodes
//	        stays m1; a malicious signer offers a conflicting message only to a
//	        relay that is malicious as well, the collusion of the five-party
//	        case, where two of the three parties make it verify and the honest
//	        answer is the majority over the three relays; and it reports a
//	        decision for another message to the commander
//
// Both halves model the disloyal deputy of section III of the protocol
// description: it modifies what it relays, it sends m41 to one honest node and
// m42 to another, and it reports a decision for a message other than the one
// proposed.
//
// The honest path never pays for this. Every hook returns as soon as it sees
// mode "none", and no state outside this file is consulted by consensus.

import (
	"sort"
	"strconv"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"log"
	"net/http"
	"sync"
	"time"
)

// Fault modes.
const (
	FaultNone   = "none"
	FaultTamper = "tamper"
)

// FaultState holds the fault configuration of one node. It is written by the
// HTTP admin handler and read by the consensus goroutines.
type FaultState struct {
	lock       sync.Mutex
	mode       string
	injections int
	hostile    map[int]bool
}

// NewFaultState returns state for a node that participates honestly.
func NewFaultState() *FaultState {
	return &FaultState{mode: FaultNone, hostile: map[int]bool{}}
}

// Mode returns the current fault mode.
func (f *FaultState) Mode() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.mode
}

// Swap installs a new mode, resets the injection counter and returns the
// mode that was in effect before.
func (f *FaultState) Swap(mode string) string {
	f.lock.Lock()
	defer f.lock.Unlock()
	previous := f.mode
	f.mode = mode
	f.injections = 0
	return previous
}

// Count records one injection and returns the running total.
func (f *FaultState) Count() int {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.injections++
	return f.injections
}

// Injected returns the number of injections since the last mode change.
func (f *FaultState) Injected() int {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.injections
}

// SetHostile records which nodes the demonstration has set to malicious. The
// strategy the protocol description analyses depends on knowing where a collusion
// is possible: two malicious parties inside one QDS instance can make a
// conflicting message verify, one cannot.
func (f *FaultState) SetHostile(ids []int) {
	f.lock.Lock()
	defer f.lock.Unlock()
	set := make(map[int]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	f.hostile = set
}

// IsHostile reports whether the demonstration armed that node.
func (f *FaultState) IsHostile(id int) bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.hostile[id]
}

// Hostile lists the armed nodes in order, for the admin response.
func (f *FaultState) Hostile() []int {
	f.lock.Lock()
	defer f.lock.Unlock()
	out := make([]int, 0, len(f.hostile))
	for id := range f.hostile {
		out = append(out, id)
	}
	sort.Ints(out)
	return out
}

// SetFault validates mode and applies it, returning the previous mode.
func (n *Node) SetFault(mode string) (string, error) {
	switch mode {
	case FaultNone, FaultTamper:
	default:
		return n.fault.Mode(), fmt.Errorf("unknown fault mode %q: want none or tamper", mode)
	}
	return n.fault.Swap(mode), nil
}

// equivocate signs a message other than the one this node already put in front
// of the group, and returns the wire bytes for the current instance. A signer
// owns its QDS key lines, so it can sign anything and the signature verifies:
// QDS stops forgery by an outsider, nothing but the consistency check stops a
// signer from equivocating. The alternative is offered only to a relay that is
// itself malicious, the collusion of the protocol analysis, and one flip is used
// for the whole round so the colluding relay carries a single m2 to everybody.
// The message the node keeps for itself is never touched.
func (n *Node) equivocate(base *message.QBAMessage, factor []byte, qdsNum int, forwardID int) (content []byte, ok bool) {
	if n.fault.Mode() != FaultTamper || len(base.Digest) == 0 {
		return nil, false
	}
	//The five-party analysis grants a conflicting message only where the relay is
	//in on it, because two malicious parties are then the majority of the QDS
	//instance. Signed for an honest relay it would be caught by the consistency
	//check of the round, so a malicious signer keeps to the agreed message there.
	if !n.fault.IsHostile(forwardID) {
		return nil, false
	}
	lying := *base
	lying.Digest = make([]byte, len(base.Digest))
	copy(lying.Digest, base.Digest)
	count := n.fault.Count()
	//one stable flip per signer: the colluding relay carries this single
	//alternative to every verifier it serves, and because the flip depends on the
	//signer the two disloyal deputies put out two different values, the m2 and
	//the m3 of the protocol figure
	lying.Digest[0] ^= byte(0x10 + uint(n.id))
	lying.QDSSignature = message.Sign(lying.Digest, factor, n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z)
	lying.SenderID = n.id
	lying.QdsNum = qdsNum
	content, err := json.Marshal(&lying)
	if err != nil {
		log.Printf("[Fault][node:%v] cannot marshal the equivocated message: %v", n.id, err)
		return nil, false
	}
	log.Printf("[Fault][node:%v] tamper: signed a message other than the one already sent (injections=%v)", n.id, count)
	n.record("fault", "signed a different message for this recipient", int(n.id), -1, -1, qdsNum, lying.Digest)
	return content, true
}

// lie returns a decision value other than the one the node really computed.
//
// The commit report carries only the digest of a decision, so one bit flip is
// enough to make the claim disagree with the honest value: the primary counts
// the report, sees that it does not name the proposed batch and never lets it
// join the quorum. The flip is the one the node applies when it signs a
// conflicting message, so a deviating deputy reports the alternative it put in
// front of the group during this round, the m2 and m3 of the protocol figure.
// The argument is never modified.
func (n *Node) lie(decision []byte) []byte {
	out := append(make([]byte, 0, len(decision)), decision...)
	if len(out) == 0 {
		return out
	}
	count := n.fault.Count()
	out[0] ^= byte(0x10 + uint(n.id))
	log.Printf("[Fault][node:%v] tamper: reporting a decision for a different message (injections=%v)", n.id, count)
	return out
}

// HttpFaultControl serves POST /fault?mode=tamper|none.
func (n *Node) HttpFaultControl(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = FaultNone
	}
	previous, err := n.SetFault(mode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"success\":false,\"message\":\"%v\"}", err)
		return
	}
	log.Printf("[Fault][node:%v] fault mode changed %v -> %v", n.id, previous, mode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"success\":true,\"node\":%v,\"previous\":\"%v\",\"mode\":\"%v\",\"injections\":0}",
		n.id, previous, mode)
}

// HttpFaultPeers serves POST /faultpeers?ids=1,2: the collusion map the strategy
// above needs. The demonstration backend pushes it after every fault change, so a
// node never acts on a stale view of who its partners are.
func (n *Node) HttpFaultPeers(w http.ResponseWriter, r *http.Request) {
	ids := []int{}
	for _, part := range strings.Split(r.URL.Query().Get("ids"), ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.Atoi(part)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "ids must be node numbers"})
			return
		}
		ids = append(ids, id)
	}
	n.fault.SetHostile(ids)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "node": n.id, "hostile": n.fault.Hostile()})
}

// HttpNodeStatus serves GET /status, the state the demonstration page polls.
// The consensus fields are a best-effort snapshot: the current implementation
// does not guard them.
func (n *Node) HttpNodeStatus(w http.ResponseWriter, r *http.Request) {
	retries, timeouts, failures := n.round.Cumulative()
	consRound, consTotal := n.consistency.Counters()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"node\":%v,\"primary\":%v,\"mode\":\"%v\",\"injections\":%v,"+
		"\"consensusCount\":%v,\"consensusRunning\":%v,\"roundState\":\"%v\",\"roundInstance\":%v,"+
		"\"retransmissions\":%v,\"timeouts\":%v,\"safeFailures\":%v,\"roundElapsed\":%.2f,"+
		"\"consistencyViolations\":%v,\"consistencyTotal\":%v}",
		n.id, n.GetPrimary(), n.fault.Mode(), n.fault.Injected(), n.ConsensusCount, n.IsConsensusRunning,
		n.round.State(), n.round.Instance(), retries, timeouts, failures, n.round.Elapsed(time.Now()),
		consRound, consTotal)
}
