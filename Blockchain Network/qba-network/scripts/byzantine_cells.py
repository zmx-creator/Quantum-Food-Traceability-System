#!/usr/bin/env python3
"""Check the recorded trace against the broadcast tables the protocol states.

The five party case analysed in the protocol description has a loyal commander and two
deviating deputies, R3 and R4. For that case the protocol names the message every
relay puts in front of every receiver, so the tables are expected values rather
than a description of behaviour:

    MR1_S      every relay delivers m1, the batch the commander signed
    MR2_S->Ri  for a loyal signer (R1, R2) every relay delivers m1
    MR2_S->R3  the loyal relays R1 and R2 deliver m1 while the deviating relay R4
               delivers m2, and it delivers the same m2 to both receivers
    MR2_S->R4  the loyal relays R1 and R2 deliver m1 while the deviating relay R3
               delivers m3 to both receivers

A row of the figure holds three entries because the relay keeps for itself the
message it received from the signer, which is the value it then hands to the two
verifiers. The invariant behind that row is therefore checked here as well: one
relay carries one message through a sub round, never two.

Run on the network host after one round with nodes 3 and 4 malicious, which is what the
acceptance script does. Exit status is the number of cells that did not match.
"""

import json
import re
import sys
import urllib.request

PORTS = {0: 6070, 1: 6071, 2: 6072, 3: 6073, 4: 6074}
DEVIATING = {3: 4, 4: 3}          # signer -> the relay that is in on it
RELAYS = (1, 2, 3, 4)
NAMES = {}
FAILED = []


def trace(node):
    url = "http://localhost:%d/trace?since=0&limit=1500" % PORTS[node]
    return json.loads(urllib.request.urlopen(url, timeout=10).read().decode())["steps"]


def flip(digest, signer):
    """The message a deviating signer puts out: its own batch digest, one bit off."""
    raw = bytearray(bytes.fromhex(digest))
    raw[0] ^= 0x10 + signer
    return raw.hex()


def label(digest):
    return NAMES.get(digest, digest)


def check(name, ok, detail):
    if ok:
        print("  PASS  %s (%s)" % (name, detail))
    else:
        print("  FAIL  %s (%s)" % (name, detail))
        FAILED.append(name)


def gather():
    """The trace of every node, each step tagged with the node that recorded it."""
    steps = []
    for node in sorted(PORTS):
        for s in trace(node):
            s["node"] = node
            steps.append(s)
    return steps


def newest_round(steps):
    """The newest round with a commander signature, and the batch it signed (m1)."""
    signs = [s for s in steps if s["Phase"] == "sign" and s["Signer"] == 0]
    if not signs:
        return -1, ""
    rnd = max(s["Round"] for s in signs)
    return rnd, [s["Digest"] for s in signs if s["Round"] == rnd][0]


def deliveries(steps, rnd):
    """Every delivery of that round, recorded once, by the relay that made it."""
    cells = {}
    for s in steps:
        if s["Phase"] == "forward" and s["Round"] == rnd:
            cells[(s["Signer"], s["Forward"], s["Verifier"])] = s["Digest"]
    return cells


def tally(step):
    """The running count of reports naming the proposed batch, read off the log line."""
    m = re.search(r"(\d+)/\d+ agree", str(step["Detail"]))
    return int(m.group(1)) if m else -1


def last_attempt(steps):
    """Reports of the final attempt only: a smaller tally opens the next attempt."""
    runs = [[]]
    previous = -1
    for s in steps:
        value = tally(s)
        if value < previous and runs[-1]:
            runs.append([])
        runs[-1].append(s)
        previous = value
    return runs[-1]


def figure():
    steps = gather()
    rnd, proposed = newest_round(steps)
    if rnd < 0:
        print("  FAIL  a round with the commander signature on record (none)")
        return 1
    NAMES[proposed] = "m1"
    for signer in DEVIATING:
        NAMES[flip(proposed, signer)] = "m%d" % (signer - 1)
    cells = deliveries(steps, rnd)

    def row(signer, relay):
        return [d for k, d in cells.items() if k[0] == signer and k[1] == relay]

    def expected(signer, relay):
        if DEVIATING.get(signer) == relay:
            return flip(proposed, signer)
        return proposed

    def held_by(signer, receiver):
        """What the protocol counts in the list of that receiver: one entry per relay."""
        out = []
        own = row(signer, receiver)
        if own:
            out.append(own[0])
        for relay in RELAYS:
            if relay == signer or relay == receiver:
                continue
            out.extend(d for k, d in cells.items() if k[0] == signer and k[1] == relay and k[2] == receiver)
        return out

    first = [d for k, d in cells.items() if k[0] == 0]
    check("MR1_S puts the proposed batch in every cell of its 12",
          len(first) == 12 and all(d == proposed for d in first),
          "%d cells, %s" % (len(first), ", ".join(sorted(set(label(d) for d in first)))))

    for signer in (1, 2):
        sub = [d for k, d in cells.items() if k[0] == signer]
        relays = sorted(set(k[1] for k in cells if k[0] == signer))
        check("MR2_S->R%d keeps all three loyal relays on the proposed batch" % signer,
              len(sub) == 6 and len(relays) == 3 and all(d == proposed for d in sub),
              "%d cells, relays %s" % (len(sub), relays))

    for signer, partner in sorted(DEVIATING.items()):
        bad = []
        for relay in RELAYS:
            if relay == signer:
                continue
            want = expected(signer, relay)
            got = row(signer, relay)
            if len(got) != 2 or any(d != want for d in got):
                bad.append("R%d=%s" % (relay, ",".join(label(d) for d in got)))
        check("MR2_S->R%d carries %s in the whole row of the deviating relay R%d" % (signer, label(flip(proposed, signer)), partner),
              not bad, bad and "honest rows wrong: %s" % ", ".join(bad) or "6 cells as the figure prints them")

    split = []
    for signer in (0, 1, 2, 3, 4):
        for relay in RELAYS:
            if relay == signer:
                continue
            if len(set(row(signer, relay))) > 1:
                split.append("signer R%d relay R%d" % (signer, relay))
    check("no relay splits a sub round between two messages", not split,
          split and ", ".join(split) or "36 rows of 3 equal entries")

    #The first broadcast round reaches a receiver through all four other deputies,
    #a second round sub round through three, so the expected list length differs.
    wrong = []
    for signer in (0, 1, 2, 3, 4):
        for receiver in (1, 2):
            if receiver == signer:
                continue
            entries = held_by(signer, receiver)
            count = {}
            for d in entries:
                count[d] = count.get(d, 0) + 1
            best = max(count.items(), key=lambda kv: kv[1]) if count else (None, 0)
            if len(entries) != (4 if signer == 0 else 3) or best[0] != proposed:
                wrong.append("signer %d receiver R%d %s" % (signer, receiver, [label(d) for d in entries]))
    check("every honest receiver still decides the proposed batch in all sub rounds", not wrong,
          wrong and "; ".join(wrong) or "four entries in MR1_S and three in each sub round, majority m1")

    return len(FAILED)


def lone(deputy):
    """One deviating deputy and no partner: it reaches nobody but the commander.

    This is the case the tables cannot show, because a conflicting message needs a
    deviating relay to carry it. What the trace has to prove is therefore the two
    claims the demonstration makes for it: nothing but the proposed batch ever
    arrives anywhere, and the value the deputy reports to the commander is counted
    and dismissed while the honest reports still reach the quorum.
    """
    steps = gather()
    rnd, proposed = newest_round(steps)
    if rnd < 0:
        print("  FAIL  a round with the commander signature on record (none)")
        return 1
    NAMES[proposed] = "m1"
    NAMES[flip(proposed, deputy)] = "m2"

    cells = deliveries(steps, rnd)
    astray = [k for k, d in cells.items() if d != proposed]
    check("a lone deviating deputy places no message in any honest list",
          bool(cells) and not astray,
          astray and ("%d of %d cells carry another value" % (len(astray), len(cells)))
          or ("all %d deliveries carry the proposed batch" % len(cells)))

    #What the commander writes down is the running count of reports that name the
    #proposed batch, so an arithmetic check settles it: the report of the deviating
    #deputy is present, and it is not one of the reports the final count is made of.
    claims = [s for s in steps if s["node"] == deputy and s["Phase"] == "infer"
              and "byzantine report" in str(s["Detail"]) and s["Digest"] != proposed]
    counted = [s for s in steps if s["node"] == 0 and s["Phase"] == "commit"
               and s["Round"] == rnd and " reported " in str(s["Detail"])]
    last = last_attempt(counted)
    mine = [s for s in last if str(s["Detail"]).startswith("node %d " % deputy)
            and s["Digest"] != proposed]
    named = len([s for s in last if s["Digest"] == proposed])
    final = max([tally(s) for s in last] or [-1])
    committed = [s for s in steps if s["node"] == 0 and s["Phase"] == "commit"
                 and s["Round"] == rnd and "reached" in str(s["Detail"])]
    check("its dishonest report is counted and adds nothing to the quorum",
          bool(claims) and bool(mine) and bool(committed) and final == named + 1,
          bool(mine) and ("%d/5 named the batch at commit from %d honest reports plus the "
                          "commander itself, while %s reported %s"
                          % (final, named, mine[0]["Detail"].split(" reported")[0],
                             label(mine[0]["Digest"])))
          or "no dishonest report on record")
    return len(FAILED)


if __name__ == "__main__":
    if len(sys.argv) > 2 and sys.argv[1] == "--lone":
        raise SystemExit(lone(int(sys.argv[2])))
    raise SystemExit(figure())