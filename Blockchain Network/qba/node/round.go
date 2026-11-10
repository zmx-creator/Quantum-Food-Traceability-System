package node

// Bounded liveness for one consensus instance.
//
// Without this file a round that dies in an unhandled branch leaves the primary
// with IsConsensusRunning set forever: qbasend1 refuses to open another round,
// so every later request queues up and the ordering service stops producing
// blocks. Round gives a round a deadline, a bounded number of retransmissions
// and a terminal safe-failure state that frees the primary for the next round.
//
// Two invariants matter for QDS security:
//
//   - a retry retransmits the byte-identical signed batch. It never re-signs,
//     so a key line is never used for two different messages.
//   - abandoning a round advances ConsensusCount, so the next round consumes
//     fresh key lines instead of replaying the exhausted window.

import (
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"sync"
	"time"
)

// Round states.
const (
	RoundIdle      = "idle"
	RoundRunning   = "running"
	RoundCommitted = "committed"
	RoundFailed    = "failed"
)

// Round is the watchdog state of a single consensus instance. The clock is
// supplied by the caller so the deadline behaviour is deterministic.
type Round struct {
	lock     sync.Mutex
	timeout  time.Duration
	maxRetry int
	state    string
	instance int       // ConsensusCount the round belongs to
	started  time.Time // when the round began
	deadline time.Time
	attempts int // retransmissions performed for the current round
	timeouts int // cumulative timeouts observed by this node
	failures int // cumulative rounds abandoned by this node
}

// NewRound returns an idle round that waits timeout for a terminal outcome and
// retransmits at most maxRetry times.
func NewRound(timeout time.Duration, maxRetry int) *Round {
	return &Round{timeout: timeout, maxRetry: maxRetry, state: RoundIdle}
}

// Begin starts tracking instance. The deadline is armed immediately.
func (r *Round) Begin(instance int, now time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.state = RoundRunning
	r.instance = instance
	r.started = now
	r.deadline = now.Add(r.timeout)
	r.attempts = 0
}

// Expired reports whether a running round has passed its deadline.
func (r *Round) Expired(now time.Time) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.state == RoundRunning && !now.Before(r.deadline)
}

// RetryAllowed reports whether another retransmission is left for this round.
func (r *Round) RetryAllowed() bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.state == RoundRunning && r.attempts < r.maxRetry
}

// Retry re-arms the deadline and counts one retransmission.
func (r *Round) Retry(now time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.state != RoundRunning {
		return
	}
	r.attempts++
	r.timeouts++
	r.deadline = now.Add(r.timeout)
}

// Commit records that the round wrote its block.
func (r *Round) Commit(now time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.state == RoundRunning {
		r.state = RoundCommitted
	}
}

// Fail records a safe failure: no block is written and the primary is released.
func (r *Round) Fail(now time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.state == RoundRunning {
		r.state = RoundFailed
		r.failures++
	}
}

// State returns the state of the current or most recent round.
func (r *Round) State() string {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.state
}

// Instance returns the ConsensusCount the current or most recent round used.
func (r *Round) Instance() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.instance
}

// Cumulative returns the counters the demonstration page reports.
func (r *Round) Cumulative() (retries int, timeouts int, failures int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.attempts, r.timeouts, r.failures
}

// Elapsed reports how long the current round has been running.
func (r *Round) Elapsed(now time.Time) float64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.state != RoundRunning {
		return 0
	}
	return now.Sub(r.started).Seconds()
}

// Liveness control on the node -------------------------------------------------

// roundWatchdog is the primary liveness guard. Only the primary opens rounds,
// so only the primary has to release them.
func (n *Node) roundWatchdog() {
	if n.id != n.GetPrimary() {
		return
	}
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()
	for now := range tick.C {
		if !n.round.Expired(now) {
			continue
		}
		if n.round.RetryAllowed() {
			n.round.Retry(now)
			retries, timeouts, _ := n.round.Cumulative()
			log.Printf("[Round][node:%v] no quorum for instance %v, retransmitting the signed batch unchanged (attempt %v of %v, timeouts %v)",
				n.id, n.round.Instance(), retries, n.maxRetry, timeouts)
			n.retransmitRound()
			n.record("round", "round timed out, retransmitting the same signed bytes", -1, -1, -1, 0, nil)
			continue
		}
		n.abandonRound(fmt.Sprintf("no consensus after %v retransmissions", n.maxRetry))
	}
}

// retransmitRound resends the byte-identical signed batch. It deliberately does
// not sign anything again: repeating a signature over the same message does not
// reuse a QDS key line for a second message.
func (n *Node) retransmitRound() {
	if len(n.roundContent) == 0 {
		return
	}
	go SendPost(n.roundContent, n.table[n.roundTarget]+server.QBAMessageEntry)
}

// abandonRound ends the current round without writing a block and releases the
// primary so that the next client request opens a fresh round on fresh keys.
// The pending batch is dropped, not re-queued: a client whose round failed has
// to submit again, which keeps the key pool from draining while a malicious
// node stays online.
func (n *Node) abandonRound(reason string) {
	if !n.IsConsensusRunning {
		return
	}
	n.round.Fail(time.Now())
	_, _, failures := n.round.Cumulative()
	//Tell the followers that this instance is over. Every node derives its QDS
	//key lines from its own ConsensusCount, so an abandoned round has to be
	//abandoned by all five or the key windows drift apart.
	if len(n.QBAMainMessage) > 0 {
		if content, _, err := message.NewQBACommitMessage(n.QBAMainMessage[0].Digest, n.id, n.id,
			n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "Consensus Failed", n.ConsensusCount*36-35); err == nil {
			n.BroadCast(content, server.QBACommitEntry)
		}
	}
	n.buffer.ClearQBACommitMsg()
	n.buffer.ClearQBAVoteMsg()
	n.roundContent = nil
	n.IsConsensusRunning = false
	n.ConsensusCount++
	log.Printf("[Round][node:%v] safe failure (total %v): %v; instance %v abandoned, batch discarded as untrustworthy, next round starts on fresh keys",
		n.id, failures, reason, n.round.Instance())
	n.record("round", fmt.Sprintf("no quorum after %v attempts: the transaction is discarded as untrustworthy and no block is written (%v)", n.maxRetry, reason), -1, -1, -1, 0, nil)
}

// abandonAsFollower releases a follower from a round the primary gave up on.
//
// Every node derives its QDS key lines from its own ConsensusCount, so an
// abandoned round has to advance on all five nodes or the key windows drift
// apart and every later round fails verification. The guard on QBAMainMessage
// keeps duplicate Consensus Failed reports from advancing twice.
func (n *Node) abandonAsFollower(reason string) {
	if len(n.QBAMainMessage) == 0 {
		return
	}
	n.QBAMainMessage = nil
	n.buffer.ClearQBACommitMsg()
	n.buffer.ClearQBAVoteMsg()
	n.ConsensusCount++
	log.Printf("[Round][node:%v] released from the abandoned round (%v), key window now %v",
		n.id, reason, n.ConsensusCount)
}

// realignKeyWindow follows the primary's key schedule.
//
// The primary owns the round counter, so any node that missed an abandon report
// a retransmitted report that arrived too late) would otherwise stay on
// an exhausted window forever. QDSNum and ConsensusCount are related by
// QDSNum = ConsensusCount*36 - 35 + i for i in [0,35], so the integer inverse
// below recovers the round index from any message of that round.
func (n *Node) realignKeyWindow(qdsNum int) {
	instance := (qdsNum + 35) / 36
	if instance <= n.ConsensusCount {
		return
	}
	log.Printf("[Round][node:%v] realigning key window %v -> %v to follow the primary",
		n.id, n.ConsensusCount, instance)
	n.ConsensusCount = instance
	n.QBAMainMessage = nil
	n.buffer.ClearQBACommitMsg()
	n.buffer.ClearQBAVoteMsg()
}
