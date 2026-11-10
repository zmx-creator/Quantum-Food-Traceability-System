package node

// Per-round trace for the demonstration page.
//
// The consensus roles (signer, forwarder, verifier) rotate inside a round, and
// every step already writes a log line. The trace keeps the same facts in a
// small ring buffer so the demonstration page can render the QDS round as it
// happens without scraping container logs.
//
// Only public values are recorded: the message digest prefix, the round and QDS
// numbers and the role assignment. QDS private key lines are deliberately not
// stored here and are no longer printed to the log.

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// traceCapacity is how many steps each node keeps. A round is 36 QDS instances
// and a node takes part in most of them, so one round produces on the order of a
// hundred steps per node; a retried round, which is what a round with malicious nodes looks
// like, repeats the whole set once per attempt. 1500 entries therefore hold a
// few rounds of history, which is what the timeline view replays, and the memory
// cost stays well under a megabyte per node.
const traceCapacity = 1500

// TraceStep is one observable step of a consensus round.
type TraceStep struct {
	Seq      int    // monotonically increasing per node
	Time     string // wall clock, RFC3339 with milliseconds
	Phase    string // request sign forward verify vote infer commit fault round
	Round    int    // ConsensusCount when the step happened
	QdsNum   int    // which QDS instance inside the round
	Signer   int    // role assignment for this step, -1 when not yet known
	Forward  int
	Verifier int
	Digest   string // first four bytes of the message digest, hex
	Detail   string // short human readable outcome
}

// Trace is the ring buffer behind GET /trace.
type Trace struct {
	lock  sync.Mutex
	steps []TraceStep
	seq   int
}

// NewTrace returns an empty trace.
func NewTrace() *Trace {
	return &Trace{steps: make([]TraceStep, 0, traceCapacity)}
}

// Append records one step and returns its sequence number.
func (t *Trace) Append(s TraceStep) int {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.seq++
	s.Seq = t.seq
	s.Time = time.Now().Format("15:04:05.000")
	t.steps = append(t.steps, s)
	if len(t.steps) > traceCapacity {
		t.steps = t.steps[len(t.steps)-traceCapacity:]
	}
	return s.Seq
}

// Since returns up to limit steps with a sequence number greater than since.
func (t *Trace) Since(since int, limit int) []TraceStep {
	t.lock.Lock()
	defer t.lock.Unlock()
	if limit <= 0 || limit > traceCapacity {
		limit = traceCapacity
	}
	out := make([]TraceStep, 0, limit)
	for _, s := range t.steps {
		if s.Seq > since {
			out = append(out, s)
			if len(out) == limit {
				break
			}
		}
	}
	return out
}

// Latest returns the most recent step, or false when nothing was recorded.
func (t *Trace) Latest() (TraceStep, bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if len(t.steps) == 0 {
		return TraceStep{}, false
	}
	return t.steps[len(t.steps)-1], true
}

// digestPrefix renders a short, safe fingerprint of a digest. Six bytes are
// enough to tell two food records apart in the demonstration page while never
// exposing anything usable: the digest is a public value anyway, the QDS key
// lines are not.
func digestPrefix(d []byte) string {
	if len(d) < 6 {
		return hex.EncodeToString(d)
	}
	return hex.EncodeToString(d[:6])
}

// record is the one call sites use. Unknown roles are -1 so the page can show a
// role that has not been assigned yet instead of implying node 0 holds it.
func (n *Node) record(phase, detail string, signer, forward, verifier int, qdsNum int, digest []byte) {
	if n.trace == nil {
		return
	}
	n.trace.Append(TraceStep{
		Phase:    phase,
		Round:    n.ConsensusCount,
		QdsNum:   qdsNum,
		Signer:   signer,
		Forward:  forward,
		Verifier: verifier,
		Digest:   digestPrefix(digest),
		Detail:   detail,
	})
}

// HttpNodeTrace serves GET /trace?since=<seq>&limit=<n> for the demonstration page.
func (n *Node) HttpNodeTrace(w http.ResponseWriter, r *http.Request) {
	since := 0
	limit := 100
	if v := r.URL.Query().Get("since"); v != "" {
		fmt.Sscanf(v, "%d", &since)
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &limit)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"node":  n.id,
		"steps": n.trace.Since(since, limit),
	}); err != nil {
		log.Printf("[Trace][node:%v] encode error: %v", n.id, err)
	}
}
