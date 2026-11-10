package node

// Consistency check between the QDS instances of one round.
//
// Protocol description of QBA, section II and section III: in broadcast phase II
// every deputy re-signs what it received and hands it to the other deputies, and
// each honest deputy compares the message it is being shown now with the message
// the same deputy showed it before ("consistency check"). If the two differ the
// deputy is asked to send again, and a deputy that keeps producing different
// values for different recipients is exposed as equivocation.
//
// Why this check is needed at all: a signer holds its own QDS key lines, so it
// can sign any message it likes and the signature verifies perfectly. QDS stops
// forgery by outsiders; nothing but this comparison stops a signer from
// equivocating. The demonstration shows the difference: a forged signature dies
// at the verification step, an equivocating one passes verification and dies
// here.
//
// Two properties of the honest path matter here:
//
//   - it never fires. Honest deputies repeat the same value, so the check is a
//     digest comparison against an entry that is already there, and no message
//     is ever rejected or re-sent.
//   - when it does fire the entry already recorded is kept. A peer is held to
//     what it attested first, the contradictory second attestation is dropped
//     from the judgment list, and the sender is complained to at most
//     QBA_CONSISTENCY_RETRY times per round so a malicious node cannot use the
//     complaint channel to flood the network.

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"sync"
)

// ConsistencyComplaint is the CommitType of the resend request.
const ConsistencyComplaint = "inconsistent message, resend"

// Consistency is the per-round bookkeeping of the check.
type Consistency struct {
	lock       sync.Mutex
	instance   int
	attested   map[message.Identify][]byte
	asked      map[message.Identify]int
	violations int // conflicts observed in the current round
	total      int // conflicts observed since the process started
}

// NewConsistency returns empty bookkeeping.
func NewConsistency() *Consistency {
	return &Consistency{
		attested: make(map[message.Identify][]byte),
		asked:    make(map[message.Identify]int),
	}
}

// beginRound clears the per-round state the first time a new round is seen.
func (c *Consistency) beginRound(instance int) {
	if c.instance == instance {
		return
	}
	c.instance = instance
	c.attested = make(map[message.Identify][]byte)
	c.asked = make(map[message.Identify]int)
	c.violations = 0
}

// attested reports the value this node has already seen from a signer.
func (c *Consistency) attestedValue(instance int, signer message.Identify) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.beginRound(instance)
	old, ok := c.attested[signer]
	return old, ok
}

// note stores the value a signer attested to this node.
func (c *Consistency) note(instance int, signer message.Identify, digest []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.beginRound(instance)
	c.attested[signer] = append(make([]byte, 0, len(digest)), digest...)
}

// conflict counts one violation and returns the running figures.
func (c *Consistency) conflict(instance int) (violations int, total int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.beginRound(instance)
	c.violations++
	c.total++
	return c.violations, c.total
}

// allowComplaint reports whether a resend request may still be sent to this
// signer in this round, and counts the one about to be sent.
func (c *Consistency) allowComplaint(instance int, signer message.Identify, maxRetry int) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.beginRound(instance)
	if c.asked[signer] >= maxRetry {
		return false
	}
	c.asked[signer]++
	return true
}

// Counters returns the violations of the current round and the cumulative total.
func (c *Consistency) Counters() (round int, total int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.violations, c.total
}

// checkConsistency is the one call site helper: it compares a value about to
// join the judgment list against the entry already recorded under the same
// signer. keep is false when the new value contradicts the recorded one, in
// which case the caller leaves the list untouched.
// consistentSigner is the call site helper. It says whether a message that just
// passed QDS verification also passes the consistency check, that is whether it
// repeats what the same signer already attested to this node. The first message
// of a signer is trivially consistent and is remembered; a repeat is consistent;
// a contradiction is counted, its sender is asked to send again, and the caller
// must treat the message as a failed verdict so it cannot join the judgment
// list.
func (n *Node) consistentSigner(signer message.Identify, digest []byte, qdsNum int) bool {
	if signer == n.id {
		return true
	}
	old, seen := n.consistency.attestedValue(n.ConsensusCount, signer)
	if !seen {
		n.consistency.note(n.ConsensusCount, signer, digest)
		return true
	}
	if bytes.Equal(old, digest) {
		return true
	}
	violations, total := n.consistency.conflict(n.ConsensusCount)
	log.Printf("[Consistency][node:%v]node %v presented %v after attesting %v: inconsistent value rejected (round conflict %v, total %v)",
		n.id, signer, digestPrefix(digest), digestPrefix(old), violations, total)
	n.record("consistency", fmt.Sprintf("node %v showed %v now but %v before: inconsistent message rejected",
		signer, digestPrefix(digest), digestPrefix(old)), int(n.id), -1, int(signer), qdsNum, digest)
	n.askForResend(signer, old, digest, qdsNum)
	return false
}

// askForResend implements the "make it send again" step of the protocol, within
// a bounded number of requests per signer per round.
func (n *Node) askForResend(signer message.Identify, old []byte, now []byte, qdsNum int) {
	if signer == n.id || int(signer) >= len(n.table) {
		return
	}
	if !n.consistency.allowComplaint(n.ConsensusCount, signer, n.consistencyRetry) {
		log.Printf("[Consistency][node:%v]resend already requested %v times from node %v, giving up on it for this round",
			n.id, n.consistencyRetry, signer)
		return
	}
	content, _, err := message.NewQBACommitMessage(old, n.id, signer, n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, ConsistencyComplaint, qdsNum)
	if err != nil {
		log.Printf("[Consistency][node:%v]cannot build the resend request: %v", n.id, err)
		return
	}
	log.Printf("[Consistency][node:%v]asking node %v to send its message again (it claimed %v and then %v)",
		n.id, signer, digestPrefix(old), digestPrefix(now))
	n.record("consistency", fmt.Sprintf("requested node %v to send its message again", signer), int(n.id), -1, int(signer), qdsNum, old)
	go SendPost(content, n.table[signer]+server.QBACommitEntry)
}

// handleConsistencyComplaint runs at the node a peer accuses. The honest path
// cannot be accused twice, so all this does is make the accusation observable;
// a node under injection keeps its behaviour by construction, which is
// precisely the case where the protocol repeats and eventually abandons.
func (n *Node) handleConsistencyComplaint(msg *message.QBACommitMessage) {
	accuser := msg.ID
	log.Printf("[Consistency][node:%v]node %v complains that this node produced different values in instance %v",
		n.id, accuser, n.ConsensusCount)
	n.record("consistency", fmt.Sprintf("node %v was accused by node %v of sending inconsistent values", n.id, accuser),
		int(accuser), -1, int(n.id), msg.QdsNum, msg.Digest)
}
