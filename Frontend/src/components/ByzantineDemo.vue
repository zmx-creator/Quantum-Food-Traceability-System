<template>
  <div class="byz-page">
    <div class="head">
      <h2>QBA consensus with quantum digital signatures</h2>
      <div class="facts-bar">
        <span><b>{{ nodeCount }}</b> consensus nodes</span>
        <span>tolerates <b>f = {{ faultTolerance }}</b> malicious</span>
        <span>writes a block when <b>{{ quorum }} of {{ nodeCount }}</b> agree on the transaction</span>
      </div>
    </div>

    <div v-if="!syncOk" class="alert">
      Cannot reach the backend at {{ apiBase }}. The page is showing its last known state.
    </div>

    <!-- ---------------- nodes ---------------- -->
    <div class="panel">
      <div class="ptitle">Consensus nodes</div>
      <div class="topology">
        <div v-for="n in nodes" :key="n.node" class="node" :class="nodeClass(n)">
          <div class="ntag" v-if="n.node === 0">PRIMARY</div>
          <div class="nname">Node {{ n.node }}</div>
          <div class="nstate">{{ isHostile(n) ? 'MALICIOUS' : 'HONEST' }}</div>
          <div class="nrole">{{ roleOf(n.node) ? 'this round: ' + roleOf(n.node) : '' }}</div>
        </div>
      </div>
      <div v-if="hostileRoleWarning" class="note">{{ hostileRoleWarning }}</div>
    </div>

    <!-- ---------------- control + outcome ---------------- -->
    <div class="grid2">
      <div class="panel">
        <div class="ptitle">Control</div>
        <label class="lab">Fault mode</label>
        <el-select v-model="mode" :disabled="busy" size="large" class="w100">
          <el-option label="Malicious - the strategy the protocol describes" value="tamper" />
          <el-option label="none - every node behaves honestly" value="none" />
        </el-select>

        <label class="lab">Malicious nodes</label>
        <el-checkbox-group v-model="selected" :disabled="busy" class="picker">
          <el-checkbox v-for="i in deputyIndexes" :key="i" :value="i">Node {{ i }}</el-checkbox>
        </el-checkbox-group>
        <p class="dim">Node 0 is the commander of every round, so it is not one of the choices: the
        demonstration injects deviating deputies.</p>
        <div class="quick">
          <span class="qlab">quick</span>
          <el-button v-for="k in [0, 1, 2, 3]" :key="k" :disabled="busy" @click="pickCount(k)">
            {{ k }} node{{ k === 1 ? '' : 's' }}
          </el-button>
        </div>

        <div v-if="pendingHint" class="hint" :class="{ warn: hintWarn }">{{ pendingHint }}</div>

        <div class="btns">
          <el-button type="primary" size="large" :loading="busy" @click="runRound">Run round</el-button>
          <el-button size="large" :disabled="busy" @click="applySelection">Apply</el-button>
          <el-button size="large" :disabled="busy" @click="resetFaults">Reset demo</el-button>
          <el-button size="large" :disabled="busy" @click="refresh">Refresh</el-button>
        </div>
      </div>

      <div class="panel">
        <div class="ptitle">Outcome of this round</div>
        <div class="verdict" :class="verdictClass">{{ verdict.text }}</div>
        <div class="txcard">
          <div class="txcap">the transaction this round put to the orderers</div>
          <template v-if="txInfo.foodID">
            <div class="txid">{{ txInfo.foodID }}</div>
            <div class="txfields">
              <span>{{ txInfo.foodName }}</span><span>{{ txInfo.foodSpec }}</span>
              <span>lot {{ txInfo.foodLOT }}</span><span>{{ txInfo.foodMFRSName }}</span>
              <span>{{ txInfo.foodProPrice }}</span><span>{{ txInfo.foodProPlace }}</span>
            </div>
            <div class="txdigest">signed digest = <code>{{ proposedValue || '-' }}</code>
              <span class="dim">- the batch every honest node names m1</span></div>
          </template>
          <div v-else class="dim">No transaction on screen yet - press Run round.</div>
        </div>

        <div class="attempts" v-if="frames.length">
          <span class="acap">attempts on this batch</span>
          <span v-if="!attempts.retries" class="achip ok">quorum reached on the first attempt</span>
          <template v-else>
            <span v-for="k in attempts.retries" :key="k" class="achip fail">no quorum</span>
            <span class="achip" :class="attempts.discarded ? 'fail' : 'ok'">
              {{ attempts.discarded ? 'retries used up' : 'quorum reached' }}
            </span>
            <span class="adetail">A timeout retransmits the very same signed bytes: the batch and its
              signatures do not change. When the retries are used up the batch is discarded as
              untrustworthy and no block is written.</span>
          </template>
        </div>
        <div class="ledgerline">{{ ledgerText }}</div>
        <div v-if="trackLink" class="stepfoot">
          The record reads back from the ledger:
          <router-link class="tlink" :to="trackLink">open this FoodID in Traceability</router-link>
        </div>
      </div>
    </div>

    <!-- ---------------- QDS timeline ---------------- -->
    <div class="panel">
      <div class="ptitle">QDS timeline - step through this round</div>
      <p class="expl">
        The signatures of this round in the order they were performed, then the judgment list of each deputy,
        then what each deputy reported. A round finishes in a fraction of a second, so this panel
        replays the trace the nodes recorded. Click any cell to jump to it.
      </p>
      <div class="tlbar">
        <el-button type="primary" size="large" :disabled="!frames.length" @click="togglePlay">
          {{ playing ? 'Pause' : 'Play' }}
        </el-button>
        <el-button size="large" :disabled="!frames.length" @click="stepBy(-1)">Back</el-button>
        <el-button size="large" :disabled="!frames.length" @click="stepBy(1)">Next</el-button>
        <el-button size="large" :disabled="!frames.length" @click="restartTimeline">Restart</el-button>
        <span class="qlab">speed</span>
        <el-select v-model="speedMs" size="large" class="wspeed">
          <el-option label="slow" :value="1000" />
          <el-option label="normal" :value="520" />
          <el-option label="fast" :value="240" />
        </el-select>
        <el-switch v-model="followLive" active-text="Follow live" class="sw" />
      </div>
      <div v-if="!frames.length" class="dim">No round on screen yet - press Run round and the timeline fills in.</div>
      <div v-else class="lanes">
        <div v-for="(ln, li) in lanes" :key="li" class="lane">
          <div class="lcap">
            <div class="lt">{{ ln.title }}</div>
            <div class="ls">{{ ln.sub }}</div>
          </div>
          <div class="lbody">
            <div v-for="(grp, gi) in ln.items" :key="gi" class="lgroup">
              <div v-if="grp.label" class="glab">{{ grp.label }}</div>
              <div class="gcells">
                <button v-for="x in grp.items" :key="x.i" class="cell"
                        :class="[x.f.cellClass, valueTone(x.f.name), { at: x.i === cursor }]"
                        :title="x.f.title" @click="goTo(x.i)">
                  {{ x.f.short }}<i v-if="x.f.dot" class="dot"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="legend">
        <span><i class="k v1"></i>m1 - the batch the commander proposed</span>
        <span><i class="k v2"></i>m2 - another message put out in this round</span>
        <span><i class="k v3"></i>m3</span>
        <span><i class="k v4"></i>m4 - the third message of the round</span>
        <span><i class="k dotk"></i>a malicious action took place here</span>
        <span class="legnote">The colour says <b>which message</b>, not whether it verified: a coloured
        cell is a signature that passed the QDS check, and the card below states the verdict of the
        signature now on screen.</span>
      </div>

      <div class="maprow">
        <div class="mapbox">
          <svg viewBox="0 0 400 372" class="pentagon">
            <defs>
              <marker id="arrOk" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#1f7a35" />
              </marker>
              <marker id="arrBad" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#b3222c" />
              </marker>
              <marker id="arrWarn" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#c97f10" />
              </marker>
              <marker id="arrV2" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#c97f10" />
              </marker>
              <marker id="arrV3" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#6b52c4" />
              </marker>
              <marker id="arrV4" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto">
                <path d="M0,0 L6,3 L0,6 z" fill="#b0398f" />
              </marker>
            </defs>
            <line v-for="l in links" :key="l.k" :x1="l.x1" :y1="l.y1" :x2="l.x2" :y2="l.y2"
                  :class="l.msg ? 'lk lk' + valueTone(l.msg) : (l.tone ? 'lk lk' + l.tone : 'lk')" />
            <template v-if="current && current.kind === 'qds'">
              <line :x1="pos(current.signer).x" :y1="pos(current.signer).y" :x2="pos(current.forward).x" :y2="pos(current.forward).y"
                    :class="'flow ' + (current.name ? 'f' + valueTone(current.name) : 'f' + current.tone)"
                    :marker-end="current.name ? valueArrow(current.name) : arrowFor(current.tone)" />
              <line :x1="pos(current.forward).x" :y1="pos(current.forward).y" :x2="pos(current.verifier).x" :y2="pos(current.verifier).y"
                    :class="'flow ' + (current.name ? 'f' + valueTone(current.name) : 'f' + current.tone)"
                    :marker-end="current.name ? valueArrow(current.name) : arrowFor(current.tone)" />
            </template>
            <g v-for="n in nodeIndexes" :key="'p' + n">
              <circle :cx="pos(n).x" :cy="pos(n).y" r="31" :class="['ndot', isHostileNode(n) ? 'ndHostile' : 'ndOk', n === 0 ? 'ndPrimary' : '']" />
              <text :x="pos(n).x" :y="pos(n).y - 1" class="ndlab">{{ n === 0 ? 'S' : 'R' + n }}</text>
              <text :x="pos(n).x" :y="pos(n).y + 17" class="ndsub">node {{ n }}</text>
            </g>
          </svg>
          <div class="mapcap">
            The five consensus nodes: S is the commander, R1-R4 the deputies. The thick arrow is the signature
            now on screen - from the signer, through the relay, to the verifier. A link keeps the colour of the
            message it last carried, so the links carrying m2 or m3 stand out from the m1 web.
          </div>
        </div>
        <div class="mapbox">
          <div class="tables">
            <div v-for="tb in tables" :key="tb.key" class="subtable" :class="{ wide: tb.key === 'd1' }">
              <div class="stab">{{ tb.title }} <span class="stabsub">{{ tb.sub }}</span></div>
              <table class="mtable">
                <thead>
                  <tr>
                    <th class="axis">F \ V</th>
                    <th v-for="c in tb.cols" :key="'c' + c" :class="{ hostileHead: isHostileNode(c) }">R{{ c }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in tb.rows" :key="'r' + r">
                    <th :class="{ hostileHead: isHostileNode(r) }">R{{ r }}</th>
                    <td v-for="c in tb.cols" :key="'cc' + c" :class="{ cellNow: isNowCell(tb, r, c) }">
                      <span v-for="(v, k) in (tb.cells[r + ':' + c] || [])" :key="k" class="mchip"
                            :class="['m' + v.tone, valueTone(v.name)]">
                        {{ v.name }}<template v-if="v.times > 1"> x{{ v.times }}</template>
                      </span>
                    </td>
                  </tr>
                  <tr class="mfoot">
                    <th>majority</th>
                    <td v-for="c in tb.cols" :key="'mf' + c">
                      <span class="mchip" :class="['m' + tb.maj[c].tone, valueTone(tb.maj[c].text)]">{{ tb.maj[c].text }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="!tables.length" class="dim">No signature replayed yet - the tables fill in as the round plays.</div>
            <div v-if="reportsLine.length" class="reports">
              <span class="rlab">reports the deputies send to the commander:</span>
              <span class="rdim">read "R3: m1 &gt; m2" as: R3 decided m1 but reported m2</span>
              <span v-for="(r, i) in reportsLine" :key="i" class="mchip" :class="'m' + r.tone">{{ r.text }}</span>
            </div>
          </div>
          <div class="mapcap">
            Rows are the relay F, columns the receiver V, as in the protocol figure: a cell is the message that
            relay delivered there, the diagonal is what the relay keeps for itself, and the majority line is the
            receiver's own list.
          </div>
        </div>
      </div>

<div v-if="current" class="frame" :class="'tone-' + current.tone">
        <div class="fhead">
          <span class="fphase">{{ current.phase }}</span>
          <span class="fno">frame {{ cursor + 1 }} of {{ frames.length }}</span>
        </div>

        <div v-if="current.kind === 'qds'" class="qdsrow">
          <div class="role" :class="{ hostile: isHostileNode(current.signer) }">
            <div class="rlab">signer</div><div class="rname">Node {{ current.signer }}</div>
          </div>
          <div class="ar">&rarr;</div>
          <div class="role" :class="{ hostile: isHostileNode(current.forward) }">
            <div class="rlab">forwarder</div><div class="rname">Node {{ current.forward }}</div>
          </div>
          <div class="ar">&rarr;</div>
          <div class="role" :class="{ hostile: isHostileNode(current.verifier) }">
            <div class="rlab">verifier</div><div class="rname">Node {{ current.verifier }}</div>
          </div>
          <div class="msg">
            message <code>{{ current.digest || '-' }}</code> = <b class="bval" :class="valueTone(current.name)">{{ current.name || 'm?' }}</b>
            <span class="dim">{{ current.digestNote }}</span>
          </div>
        </div>

        <div v-else-if="current.kind === 'judge'" class="jrow">
          <div class="rname">Node {{ current.node }} takes the majority of its judgment list</div>
          <div class="chips">
            <span v-for="(e, i) in current.entries" :key="i" class="chip"
                  :class="e.value === proposedValue ? 'chipOk' : 'chipBad'">
              m{{ e.label }}=<code>{{ e.value }}</code>
            </span>
            <span v-if="!current.entries.length" class="dim">empty - nothing was accepted</span>
          </div>
          <div class="msg">majority <code>{{ current.decision || '-' }}</code></div>
        </div>

        <div v-else class="jrow">
          <div class="rname">{{ current.headline }}</div>
          <div class="msg">{{ current.text }}</div>
        </div>

        <div class="badge" :class="'b-' + current.tone">{{ current.badge }}</div>
        <div v-for="(note, i) in current.notes" :key="i" class="fnote" :class="note.cls">{{ note.text }}</div>
      </div>

      <el-button link class="toggle" @click="showRaw = !showRaw">
        {{ showRaw ? 'Hide' : 'Show' }} raw QDS trace ({{ traceSteps.length }} rows)
      </el-button>
      <div v-if="showRaw" class="rawbox">
        <table class="raw">
          <thead>
            <tr><th>time</th><th>node</th><th>phase</th><th>qds #</th><th>signer</th><th>fwd</th><th>ver</th><th>value</th><th>detail</th></tr>
          </thead>
          <tbody>
            <tr v-for="(s, i) in traceSteps" :key="i" :class="rowClass(s)">
              <td>{{ s.Time }}</td><td>{{ s.node }}</td><td>{{ s.Phase }}</td><td>{{ s.QdsNum }}</td>
              <td :class="{ hostileCell: isHostileNode(s.Signer) }">{{ s.Signer }}</td>
              <td :class="{ hostileCell: isHostileNode(s.Forward) }">{{ s.Forward }}</td>
              <td :class="{ hostileCell: isHostileNode(s.Verifier) }">{{ s.Verifier }}</td>
              <td><code>{{ s.Digest }}</code></td><td class="det">{{ detailText(s.Detail) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { ElNotification } from 'element-plus'
import { getDemoNodes, getDemoConsensus, setDemoFault, runDemoRound, API_BASE, http } from '../api/fabric'
import greenNode from '../assets/greenNode.png'
import redNode from '../assets/redNode.png'

const nodeCount = ref(5)
const faultTolerance = ref(2)
const quorum = ref(3)
const commitRule = ref("quorum of 3 - the primary commits once n-f reports name the proposed batch")
const nodes = ref([])
const summary = ref({ maliciousNodes: 0, timeouts: 0, retransmissions: 0, safeFailures: 0, keyWindowsAligned: true })
const trace_raw = ref([])
//trace is the part of the fetched buffer a round panel may read: nothing at or
//below the reset floor.
const trace = computed(() => trace_raw.value.filter(s => Number(s.Round || 0) > clearedRound.value))
//The analysis reads every fetched step; only the timeline table is
//truncated, because one retried round produces hundreds of QDS steps.
const traceSteps = computed(() => trace.value.slice(-150))
const result = ref(null)
const mode = ref('tamper')
//selected is what the operator asked for; nodes.value is what the network is.
//They used to share a variable, so every poll overwrote a fresh checkbox tick
//with the honest state of the network and the selection vanished a second after
//it was made. The network seeds the selection once, on the first read.
const selected = ref([])
const intentLoaded = ref(false)
const busy = ref(false)
//syncOk makes a dead poll loop visible: without it a page that cannot reach the
//backend keeps showing the state it last read, which looks like a button that
//does nothing.
const syncOk = ref(true)
const apiBase = API_BASE
let timer = null
//clearedRound is the round the operator has just reset away: every panel that
//describes a round reads the newest round above that floor, so Reset really
//empties the page instead of keeping the previous demonstration on screen.
const clearedRound = ref(-1)
//A node answers /trace with the OLDEST steps of its ring buffer, capped at
//"limit", so a small window would hand the page last weeks rounds and hide the
//round that is on screen. The window therefore covers the whole buffer.
const TRACE_WINDOW = 1500
const pollLimit = ref(TRACE_WINDOW)

const nodeIndexes = computed(() => Array.from({ length: nodeCount.value }, (_, i) => i))
//The commander cannot be armed. Nothing in the injected strategy makes a commander
//behave badly, so offering node 0 would let the page claim a Byzantine node that the
//network does not treat as one.
const deputyIndexes = computed(() => nodeIndexes.value.filter(i => i !== 0))
const hostileSet = computed(() => new Set(nodes.value.filter(isHostile).map(n => n.node)))

function isHostile(n) {
  return n.mode && n.mode !== 'none'
}
function isHostileNode(id) {
  return id >= 0 && hostileSet.value.has(id)
}
function modeClass(n) {
  return (!n.mode || n.mode === 'none') ? 'mode-ok' : 'mode-tamper'
}

// roleOf reports the role the latest traced round assigned to a node, so the
// page can show whether a Byzantine node ended up holding a decisive role.
// roleMap reports the role each node holds in the newest QDS instance of the
// newest round. Roles rotate on every signature inside a round, so reporting
// them per round would always say "signer/forwarder/verifier" and hide the
// thing this demonstration is about: which role the Byzantine node currently
// occupies.
const roleMap = computed(() => {
  const map = {}
  if (latestRound.value < 0) return map
  let newest = -1
  let latestQds = -1
  for (const s of trace.value) {
    if (s.Round > newest || (s.Round === newest && s.QdsNum > latestQds)) {
      newest = s.Round
      latestQds = s.QdsNum
    }
  }
  if (newest < 0) return map
  const put = (id, role) => {
    if (id < 0) return
    if (!map[id]) map[id] = role
    else if (!map[id].includes(role)) map[id] = map[id] + '/' + role
  }
  for (const s of trace.value) {
    if (s.Round !== newest || s.QdsNum !== latestQds) continue
    put(s.Signer, 'signer')
    put(s.Forward, 'forwarder')
    put(s.Verifier, 'verifier')
  }
  return map
})

function roleOf(id) {
  return roleMap.value[id] || ''
}
const hostileRoleWarning = computed(() => {
  const hit = [...hostileSet.value].filter(id => roleMap.value[id]).sort((a, b) => a - b)
  if (!hit.length) return ''
  const plural = hit.length > 1
  return 'Malicious node' + (plural ? 's ' : ' ') + hit.join(', ')
    + (plural ? ' hold' : ' holds') + ' a role in the current round ('
    + hit.map(id => id + ': ' + roleMap.value[id]).join(', ') + ').'
})

// latestRound is the newest round that produced a decision step. The bookkeeping
// lines of a discarded round are stamped with the window that was opened after
// the abandonment, so they are skipped here; otherwise the round whose outcome
// is on screen would look empty.
const latestRound = computed(() => {
  let m = -1
  for (const s of trace.value) {
    if (s.Phase === 'round') continue
    if ((s.Round || 0) > m) m = s.Round || 0
  }
  return m
})

// parseList reads the judgment list published by an inference step:
// "[2:a3f1c0de,4:11b7ee02] -> majority a3f1c0de".
function parseList(detail) {
  const m = /\[(.*?)\]/.exec(detail || '')
  if (!m || !m[1]) return []
  return m[1].split(',').map(part => {
    const i = part.indexOf(':')
    return i < 0 ? null : { fwd: part.slice(0, i).trim(), value: part.slice(i + 1).trim() }
  }).filter(Boolean)
}

// proposedValue is the value the honest majority agrees on, so it is the batch
// the primary signed: the protocol's m1.
const proposedValue = computed(() => {
  //the batch the commander signed is the reference value of the round: reading it
  //off the primary signature step keeps the labels honest even in a round that is
  //discarded, where the judgment lists hold no majority at all
  for (const s of trace.value) {
    if (s.Phase === 'sign' && String(s.node) === '0' && s.Round === latestRound.value && s.Digest) {
      return String(s.Digest)
    }
  }
  const counts = {}
  for (const s of trace.value) {
    if (s.Round !== latestRound.value || s.Phase !== 'infer') continue
    for (const e of parseList(s.Detail)) counts[e.value] = (counts[e.value] || 0) + 1
  }
  let best = ''
  let bestCount = 0
  for (const value of Object.keys(counts)) {
    if (counts[value] > bestCount) { bestCount = counts[value]; best = value }
  }
  return best
})

const networkVerdict = computed(() => {
  const lines = commitTally.value
  if (lines.some(l => String(l).indexOf('block written') >= 0)) {
    return { agreed: true, text: 'CONSENSUS REACHED - the transaction was written to the ledger' }
  }
  //A round that never reached the quorum leaves no "written" line at all, so its
  //own terminal state is the evidence: the primary sits in the failed state and
  //the batch was discarded.
  const primary = (summary.value.nodes || []).find(x => String(x.node) === '0')
  const failed = primary && String(primary.roundState) === 'failed'
  if (failed || lines.some(l => String(l).indexOf('discarded') >= 0)) {
    return { agreed: false, text: 'CONSENSUS NOT REACHED - the transaction was discarded' }
  }
  return null
})

// commitTally shows the arithmetic the primary performed. The report of a
// discarded round is stamped with the key window opened after the abandonment,
// so the commit reports of the round on screen win whenever there are any:
// otherwise the previous round's discard would be read as part of this one.
const commitTally = computed(() => {
  const commits = []
  const discards = []
  if (latestRound.value < 0) return []
  for (const s of trace.value) {
    if (String(s.node) !== '0') continue
    if (s.Phase === 'commit' && s.Round === latestRound.value) commits.push(String(s.Detail))
    else if (s.Phase === 'round' && String(s.Detail).indexOf('quorum') >= 0) discards.push(String(s.Detail))
  }
  return (commits.length ? commits : discards).slice(-12)
})

function rowClass(s) {
  if (s.Phase === 'fault') return 'row-fault'
  if (s.Phase === 'round') return 'row-round'
  if (s.Phase === 'verify' && s.Detail && s.Detail.indexOf('failed') >= 0) return 'row-bad'
  if (s.Phase === 'commit') return 'row-ok'
  return ''
}

const armedSet = computed(() => new Set(nodes.value.filter(isHostile).map(x => x.node)))
const wantedSet = computed(() => new Set(mode.value === 'none' ? [] : selected.value))
//selectionPending says whether the screen and the network disagree, which is the
//one thing the operator needs to know before pressing anything.
const selectionPending = computed(() => {
  const want = wantedSet.value
  const now = armedSet.value
  if (want.size !== now.size) return true
  for (const v of want) { if (!now.has(v)) return true }
  return false
})
//The nodes record the protocol vocabulary in their trace; the raw dump is the one place
//that text reaches the screen, so the terminology of the demonstration is applied here,
//display only - nothing on the wire or in the judgment changes.
function detailText(d) {
  const t = (d === undefined || d === null) ? '' : String(d)
  return t.replace('byzantine report', 'malicious report').replace('the QDS signature', 'the QDS')
}

//The demonstration speaks of malicious nodes; the protocol analysis and the wire
//values keep the fault-model name, so the labels are mapped here and nowhere else.
const modeNames = { none: 'honest', tamper: 'malicious' }
const modeName = computed(() => modeNames[mode.value] || mode.value)

const pendingHint = computed(() => {
  if (!selected.value.length) {
    return armedSet.value.size
      ? armedSet.value.size + ' node(s) are still malicious on the network - press Reset demo to disarm them.'
      : ''
  }
  if (mode.value === 'none') {
    return 'Fault mode is "none", so the selected nodes are armed as honest.'
  }
  if (!selectionPending.value) {
    return selected.value.length + ' node(s) are armed with the ' + modeName.value + ' strategy (read back from the nodes).'
  }
  return selected.value.length + ' node(s) selected with the ' + modeName.value + ' strategy but not applied - press Apply selection, or Run demonstration round to arm and run one round.'
})

async function applySelection() {
  busy.value = true
  try {
    const n = await applyFaults()
    await refresh()
    ElNotification({
      title: 'Faults applied',
      message: n + ' node state(s) written; the cards now read back ' + armedSet.value.size + ' malicious node(s).',
      type: 'success'
    })
  } catch (e) {
    ElNotification({ title: 'Applying faults failed', message: String(e.message || e), type: 'error', duration: 0 })
  } finally {
    busy.value = false
    await refresh()
  }
}

function pickCount(k) {
  //only the four deputies can be deviating, so the quick choices count from R1 up
  const order = deputyIndexes.value
  selected.value = order.slice(0, k).sort((a, b) => a - b)
}

async function refresh() {
  try {
    const [n, c] = await Promise.all([getDemoNodes(), getDemoConsensus(0, pollLimit.value)])
    if (n.data.success) {
      nodeCount.value = n.data.data.nodeCount
      faultTolerance.value = n.data.data.faultTolerance
      nodes.value = n.data.data.nodes
      if (!intentLoaded.value) {
        selected.value = nodes.value.filter(isHostile).map(x => x.node)
        intentLoaded.value = true
      }
    }
    if (c.data.success) {
      const d = c.data.data
      quorum.value = d.quorum
      if (d.commitRule) commitRule.value = d.commitRule
      summary.value = d
      const merged = []
      for (const key of Object.keys(d.trace || {})) {
        for (const s of (d.trace[key].steps || [])) merged.push({ ...s, node: key })
      }
      merged.sort((a, b) => (a.Time < b.Time ? -1 : a.Time > b.Time ? 1 : a.node - b.node))
      trace_raw.value = merged
    }
    syncOk.value = true
  } catch (e) {
    syncOk.value = false
    console.error('demo status unavailable', e)
  }
}

// callFault turns a refused request, an unreachable node or a partially applied
// change into an exception. Every caller reports it, so a control that had no
// effect on the network can never look like a control that was never pressed.
async function callFault(nodes, wanted) {
  const r = await setDemoFault(nodes, wanted)
  const data = (r && r.data) || {}
  if (!data.success) {
    throw new Error(data.message || 'the backend rejected the request')
  }
  const results = (data.data && data.data.results) || []
  const failed = results.filter(x => !x.ok).map(x => x.node)
  if (failed.length) {
    throw new Error('node ' + failed.join(', ') + ' did not accept mode "' + wanted + '"')
  }
  if (!results.length) {
    throw new Error('the backend applied the change to no node at all')
  }
  return results.length
}

// arm the chosen nodes and force everybody else back to honest, so what the
// page shows always matches what the network is doing
async function applyFaults() {
  let applied = 0
  const others = deputyIndexes.value.filter(i => !selected.value.includes(i))
  if (others.length) applied += await callFault(others, 'none')
  //the selected nodes get the mode that is on screen, including none: skipping them
  //when the mode was set back to none would leave them armed while the page said
  //they were honest
  if (selected.value.length) applied += await callFault(selected.value, mode.value)
  return applied
}

// ---------------------------------------------------------------------------
// the workflow strip
// ---------------------------------------------------------------------------
//What the page presents as the transaction of the round. Nothing here is staged:
//the fields are the ones the round really submitted through addProInfo, the digest
//is the batch the orderers signed, and the record is read back from the ledger.
const ledgerRecord = ref(null)


async function readBackRecord(foodID) {
  ledgerRecord.value = null
  if (!foodID) return
  try {
    const r = await http.get('/food/' + encodeURIComponent(foodID) + '/production')
    const d = r.data && r.data.success ? r.data.data : null
    if (d && (typeof d !== 'object' || Object.keys(d).length)) ledgerRecord.value = d
  } catch (e) {
    ledgerRecord.value = null
  }
}


//The strip follows the round that is on screen. A freshly opened tab has no run
//result of its own, so it falls back to the round the backend last ran; a reset
//hides that memo the same way it hides the trace of the discarded round.
const wf = computed(() => {
  if (result.value) return result.value
  const memo = summary.value.lastRound
  if (!memo) return null
  return Number(memo.roundInstance || 0) <= clearedRound.value ? null : memo
})


watch(() => (wf.value && wf.value.agreed ? wf.value.foodID : ''), id => readBackRecord(id))


const txInfo = computed(() => {
  const r = wf.value
  return (r && r.payload) || {}
})


//Liveness of the round, read off what the primary recorded. Every timeout line is
//one retransmission of the same signed batch - a retry is not a new transaction -
//and the discard line is what the protocol does once the retries are used up.
const attempts = computed(() => {
  let retries = 0
  for (const s of trace.value) {
    if (Number(s.node) !== 0 || s.Phase !== 'round') continue
    if (latestRound.value >= 0 && Number(s.Round) !== latestRound.value) continue
    if (String(s.Detail).indexOf('timed out') >= 0) retries++
  }
  //The line announcing the abandonment is stamped after the round has already been
  //released, so it can fall outside the round on screen. The verdict is the
  //authority here: retries that did not end in agreement are retries used up.
  return { retries, discarded: retries > 0 && verdict.value.agreed === false }
})


const trackLink = computed(() => {
  const r = wf.value
  if (!r || !r.agreed || !r.foodID) return null
  return { path: '/Home/FoodTrack', query: { foodID: r.foodID } }
})


async function runRound() {
  busy.value = true
  stopPlay()
  cursor.value = -1
  try {
    try {
      await applyFaults()
    } catch (e) {
      ElNotification({
        title: 'Faults were not applied',
        message: String(e.message || e) + ' - the round is skipped so the result cannot be misread.',
        type: 'error',
        duration: 0
      })
      return
    }
    await refresh()
    const r = await runDemoRound(15000)
    result.value = r.data.success ? r.data.data : null
    if (result.value) sessionStorage.setItem('byzantineLastRound', JSON.stringify(result.value))
    if (!r.data.success) {
      ElNotification({ title: 'Demonstration round failed', message: r.data.message, type: 'error', duration: 0 })
    } else if (!r.data.data.agreed) {
      ElNotification({
        title: 'Consensus not reached',
        message: 'The retries were used up without agreement, so the transaction is '
               + 'discarded as untrustworthy and no block is written.',
        type: 'error',
        duration: 0
      })
    } else {
      ElNotification({ title: 'Consensus reached', message: 'Block ' + r.data.data.blockNumber + ' written.', type: 'success' })
    }
  } catch (e) {
    ElNotification({ title: 'Request failed', message: String(e), type: 'error', duration: 0 })
  } finally {
    busy.value = false
    await refresh()
    afterRound()
    schedulePoll()
  }
}

async function resetFaults() {
  busy.value = true
  try {
    selected.value = []
    mode.value = 'none'
    await callFault([], 'none')
    //the demonstration screen belongs to the reset as well: the outcome, the
    //QDS counters, the broadcast list and the timeline all describe the round
    //that has just been discarded
    result.value = null
    ledgerRecord.value = null
    sessionStorage.removeItem('byzantineLastRound')
    stopPlay()
    cursor.value = -1
    clearedRound.value = Math.max(clearedRound.value, latestRound.value)
    //read the network back before claiming success: the nodes, not the page,
    //decide whether a node is Byzantine
    await refresh()
    const left = nodes.value.filter(isHostile).map(x => x.node)
    if (!syncOk.value) {
      ElNotification({
        title: 'Reset could not be verified',
        message: 'The backend is unreachable at ' + apiBase + ', so the page cannot confirm anything. '
               + 'Check the Spring Boot application on port 8080.',
        type: 'error',
        duration: 0
      })
    } else if (left.length) {
      ElNotification({
        title: 'Reset incomplete',
        message: 'node ' + left.join(', ') + ' still reports a fault mode after the request.',
        type: 'error',
        duration: 0
      })
    } else {
      ElNotification({
        title: 'Nodes reset',
        message: 'All ' + nodeCount.value + ' nodes report mode none (read back from the nodes), and the round display is cleared.',
        type: 'success'
      })
    }
  } catch (e) {
    ElNotification({ title: 'Reset failed', message: String(e.message || e), type: 'error', duration: 0 })
  } finally {
    busy.value = false
    await refresh()
  }
}

onMounted(async () => {
  const seen = sessionStorage.getItem('byzantineLastRound')
  if (seen) {
    try { result.value = JSON.parse(seen) } catch (e) { sessionStorage.removeItem('byzantineLastRound') }
  }
  await refresh()
  //poll through a running round as well: watching the cards move is the
  //demonstration, and a frozen page reads as a dead button. The checkbox
  //selection is still the user's while busy.
  schedulePoll()
})
onUnmounted(() => { if (timer) clearTimeout(timer) })

//One self-scheduling poll instead of a fixed interval: a live round is ticked
//faster so the timeline grows while it runs, an idle page is ticked slowly.
function schedulePoll() {
  if (timer) clearTimeout(timer)
  const fast = followLive.value || busy.value
  timer = setTimeout(async () => {
    await refresh()
    schedulePoll()
  }, fast ? 700 : 2000)
}
// ---------------------------------------------------------------------------
// Presentation layer: everything below is scoped to the round on screen, so the
// page never mixes a per-round result with counters accumulated since start.
// ---------------------------------------------------------------------------
const showRaw = ref(false)

const roundSteps = computed(() =>
  trace.value.filter(s => s.Round === latestRound.value && s.Phase !== 'round'))

// Rejected entries: a value a peer delivered that this node could not accept.
// Two different mechanisms reject two different things, and the distinction is
// the point of the demonstration, so both are reported with their reason.
const rejectedBy = computed(() => {
  const byNode = {}
  const push = (node, accused, value, reason) => {
    if (node === undefined || node === null || !value) return
    const list = byNode[node] || (byNode[node] = [])
    const found = list.find(r => r.accused === accused && r.value === value && r.reason === reason)
    if (found) { found.times++ ; return }
    list.push({ accused, value, reason, times: 1 })
  }
  //one round repeats the same rejection dozens of times (a signer takes part in
  //several QDS instances), so the view keeps the distinct cases and counts them
  for (const s of roundSteps.value) {
    const detail = String(s.Detail || '')
    if (s.Phase === 'consistency' && detail.indexOf('inconsistent message rejected') >= 0) {
      push(s.node, s.Verifier, s.Digest, 'inconsistent with its earlier message')
    } else if (s.Phase === 'verify' && detail.indexOf('kept out of the judgment list') >= 0) {
      push(s.node, s.Signer, s.Digest, 'QDS forged')
    }
  }
  return byNode
})

const verdict = computed(() => {
  if (result.value) {
    return result.value.agreed
      ? { agreed: true, text: 'CONSENSUS REACHED - the transaction was written to the ledger' }
      : { agreed: false, text: 'CONSENSUS NOT REACHED - the transaction was discarded' }
  }
  //The fallback describes the round the network last ran, which is what a freshly
  //opened page should answer. Once the operator has pressed Reset there is no round
  //on screen at all, and a stale verdict line would contradict the empty panels.
  if (networkVerdict.value && !clearedRound.value) return networkVerdict.value
  return { agreed: null, text: 'No round on screen yet' }
})
const verdictClass = computed(() =>
  verdict.value.agreed === true ? 'ok' : verdict.value.agreed === false ? 'bad' : 'idle')

const ledgerText = computed(() => {
  if (verdict.value.agreed === null) return '-'
  const height = result.value ? result.value.heightAfter : summary.value.height
  if (verdict.value.agreed) {
    const block = result.value && result.value.blockNumber !== undefined && result.value.blockNumber !== null
      ? 'block #' + result.value.blockNumber + ' - ' : ''
    return block + 'height ' + height
  }
  return 'no block written - height ' + height + ' unchanged'
})

const hintWarn = computed(() =>
  (mode.value === 'none' && selected.value.length > 0) || selected.value.length > faultTolerance.value)

function nodeClass(n) {
  return { hostile: isHostile(n), primary: n.node === 0 }
}

// ---------------------------------------------------------------------------
// QDS timeline. One round is over in a fraction of a second, so the panel
// replays the trace the nodes recorded: one frame per QDS signature, then one
// frame per judgment list, then what each deputy reported. "Follow
// live" pins the playhead to the newest frame while a round is running.
// ---------------------------------------------------------------------------
const playing = ref(false)
const speedMs = ref(520)
const followLive = ref(false)
const cursor = ref(-1)
let playTimer = null

function firstRole(steps, field) {
  for (const s of steps) {
    const v = Number(s[field])
    if (v >= 0) return v
  }
  return -1
}

function buildFrames() {
  const steps = trace.value.filter(s => s.Round === latestRound.value && s.Phase !== 'round')
  if (!steps.length) return []
  const groups = new Map()
  const forgings = []
  const judges = new Map()
  const dishonest = []
  const reports = []
  for (const s of steps) {
    if (s.Phase === 'commit') {
      if (String(s.node) === '0' && String(s.Detail).indexOf('block written') < 0) reports.push(s)
      continue
    }
    if (s.Phase === 'infer') {
      if (String(s.Detail).indexOf('byzantine report') >= 0) dishonest.push(s)
      else judges.set(String(s.node), s)
      continue
    }
    const q = Number(s.QdsNum) || 0
    if (q <= 0) {
      //the forging of a relayed signature is recorded without a signature number
      if (s.Phase === 'fault') forgings.push(s)
      continue
    }
    if (!groups.has(q)) groups.set(q, [])
    groups.get(q).push(s)
  }

  const proposed = proposedValue.value
  const out = []
  const relaySlots = {}
  //m1 is the proposed batch; every other value of the round is named by the
  //order in which it first appears, the m2 and m3 of the protocol figure
  const names = {}
  let nextName = 2
  const keys = Array.from(groups.keys()).sort((a, b) => a - b)
  //Number the alternatives the way the protocol figure does: the honest batch is
  //m1 and the messages the deputies put out take m2, m3... in the order of their
  //signer, so the sub-round of R3 is the one carrying m2 and R4 the one on m3.
  {
    const bySigner = new Map()
    for (const q of keys) {
      const st = groups.get(q)
      const relay = st.find(s => s.Phase === 'forward')
      const sg = relay ? Number(relay.Signer) : firstRole(st, 'Signer')
      const wv = st.find(s => s.Digest)
      const dg = wv ? String(wv.Digest) : ''
      if (dg !== '' && proposed !== '' && dg !== proposed && !bySigner.has(dg)) bySigner.set(dg, sg)
    }
    Array.from(bySigner.entries())
      .sort((a, b) => a[1] - b[1])
      .forEach(entry => { names[entry[0]] = 'm' + nextName++ })
  }
  keys.forEach((q, i) => {
    const st = groups.get(q)
    const relay = st.find(s => s.Phase === 'forward')
    const signer = relay ? Number(relay.Signer) : firstRole(st, 'Signer')
    const forward = relay ? Number(relay.Forward) : firstRole(st, 'Forward')
    const verifier = relay ? Number(relay.Verifier) : firstRole(st, 'Verifier')
    const withValue = st.find(s => s.Digest)
    const digest = withValue ? String(withValue.Digest) : ''
    const inconsistent = st.some(s => s.Phase === 'consistency')
    const ok = st.some(s => s.Phase === 'verify' && String(s.Detail).indexOf('verify succeed') >= 0)
    const failed = st.some(s => s.Phase === 'verify' && String(s.Detail).indexOf('verify succeed') < 0)
    const notes = []
    let dot = false
    for (const s of st) {
      if (s.Phase === 'fault') {
        notes.push({ text: 'Node ' + (Number(s.Signer) >= 0 ? Number(s.Signer) : s.node) + ' - ' + s.Detail, cls: 'sBad' })
        dot = true
      } else if (s.Phase === 'consistency') {
        notes.push({ text: String(s.Detail), cls: 'sWarn' })
      }
    }
    let tone = 'idle'
    let badge = 'Relayed - the verdicts on this signature have not come back yet'
    let cellClass = 'cIdle'
    if (inconsistent) {
      tone = 'warn'; cellClass = 'cWarn'
      badge = 'REJECTED BY THE CONSISTENCY CHECK - this signer had already shown the node another message'
    } else if (ok) {
      tone = 'ok'; cellClass = 'cOk'
      badge = 'VERIFIED - the forwarder and the verifier both accept this signature'
    } else if (failed) {
      tone = 'bad'; cellClass = 'cBad'
      badge = 'REJECTED BY QDS - the signature does not verify, the message is kept out of the judgment list'
    }
    const other = digest !== '' && proposed !== '' && digest !== proposed
    if (other) dot = true
    let name = ''
    if (digest !== '') {
      if (!other) name = 'm1'
      else {
        if (!names[digest]) names[digest] = 'm' + nextName++
        name = names[digest]
      }
    }
    out.push({
      kind: 'qds',
      signer, forward, verifier, digest, name,
      digestNote: proposed === '' ? ''
        : (other ? ' - not the proposed batch: this receiver was shown a different message'
                 : ' - the batch the commander proposed'),
      tone, badge, cellClass, dot, notes,
      phase: signer === 0
        ? 'Broadcast phase I - the commander signs the batch, a deputy relays it'
        : 'Broadcast phase II - deputy ' + signer + ' signs what it received and distributes it',
      short: String(i + 1),
      title: 'signature ' + (i + 1) + ' of ' + keys.length + ': node ' + signer + ' signs, node '
        + forward + ' relays, node ' + verifier + ' verifies'
    })
    if (forward >= 0) {
      if (!relaySlots[forward]) relaySlots[forward] = []
      relaySlots[forward].push(out.length - 1)
    }
  })

  //A forged relay is recorded without the signature number, but every Byzantine
  //node forges exactly once per signature it relays, so the forgeries map onto
  //the signatures where that node was the forwarder, in the order they happened.
  forgings.sort((a, b) => (String(a.Time) < String(b.Time) ? -1 : 1))
  const used = {}
  for (const f of forgings) {
    const node = String(f.node)
    const slots = relaySlots[node] || []
    const k = used[node] || 0
    used[node] = k + 1
    if (k >= slots.length) continue
    const frame = out[slots[k]]
    frame.notes.push({ text: 'Node ' + node + ' forged the QDS while relaying this message', cls: 'sBad' })
    frame.dot = true
    if (frame.tone === 'idle') {
      frame.tone = 'bad'
      frame.cellClass = 'cBad'
      frame.badge = 'REJECTED BY QDS - the forged signature does not verify'
    }
  }

  const judgeNodes = Array.from(judges.keys()).sort((a, b) => Number(a) - Number(b))
  for (const node of judgeNodes) {
    const detail = String(judges.get(node).Detail || '')
    const m = /->\s*majority\s*([0-9a-f]*)/.exec(detail)
    const entries = parseList(detail).map(e => ({
      fwd: e.fwd,
      value: e.value,
      label: e.value === proposed ? '1' : String(e.fwd) + String(node)
    }))
    const lies = dishonest.filter(x => String(x.node) === node)
    const empty = detail.indexOf('is empty') >= 0
    out.push({
      kind: 'judge',
      node: Number(node),
      entries,
      decision: m ? m[1] : '',
      notes: lies.map(x => ({ text: String(x.Detail), cls: 'sBad' })),
      tone: lies.length ? 'bad' : (empty ? 'warn' : 'ok'),
      badge: lies.length
        ? 'DISHONEST REPORT - this node sends the commander a decision about another message'
        : (empty ? 'NO DECISION - nothing reached this deputy intact'
                 : 'REPORT SENT - this deputy reports its majority to the commander'),
      cellClass: lies.length ? 'cBad' : 'cJudge',
      dot: lies.length > 0,
      phase: 'Inference - each deputy takes the majority of its own judgment list',
      short: 'J' + node,
      title: 'judgment list of node ' + node
    })
  }

  reports.sort((a, b) => (String(a.Time) < String(b.Time) ? -1 : 1))
  //the tally inside one attempt only grows, so a report counting fewer
  //agreements than the one before it opens the next attempt
  let attempt = 1
  let lastTally = -1
  let attempts = 1
  reports.forEach((s, i) => {
    const text = String(s.Detail)
    const quorum = text.indexOf('reached') >= 0
    let rname = ''
    if (s.Digest) {
      const rd = String(s.Digest)
      if (proposed !== '' && rd === proposed) rname = 'm1'
      else {
        if (!names[rd]) names[rd] = 'm' + nextName++
        rname = names[rd]
      }
    }
    const at = text.indexOf(' agree')
    let tally = -1
    if (at >= 0) {
      const slash = text.lastIndexOf('/', at)
      const space = text.lastIndexOf(' ', slash)
      tally = Number(text.slice(space + 1, slash))
      if (lastTally >= 0 && tally <= lastTally) attempt++
      lastTally = tally
      attempts = attempt
    }
    const who = (/^node (\d+)/.exec(text) || [])[1]
    const reporter = who === undefined ? String(s.Verifier) : who
    const reported = 'it decided ' + (rname || 'nothing') + ' - '
      + (rname === 'm1' ? 'the transaction the commander proposed' : 'a message of this round')
    const enough = rname === 'm1'
    out.push({
      kind: 'report',
      name: rname,
      digest: s.Digest ? String(s.Digest) : '',
      headline: (attempt > 1 ? 'Attempt ' + attempt + ' - ' : '')
        + (quorum ? 'the commander has enough agreement to write the block'
                  : 'node ' + reporter + ' reports its decision'),
      text: reported,
      notes: [],
      tone: quorum ? 'ok' : (enough ? 'ok' : 'info'),
      badge: quorum
        ? 'QUORUM REACHED - the transaction is written into a block'
        : (enough ? 'IT DECIDED THE PROPOSED TRANSACTION' : 'IT DECIDED ANOTHER MESSAGE - not counted'),
      cellClass: 'cReport',
      dot: false,
      attempt,
      phase: 'Commit - each deputy reports its decision to the commander',
      short: 'C' + (i + 1),
      title: reported
    })
  })

  const v = verdict.value
  if (v.agreed !== null && (reports.length || dishonest.length)) {
    out.push({
      kind: 'report',
      headline: 'Outcome of the round',
      text: ledgerText.value,
      notes: [],
      tone: v.agreed ? 'ok' : 'bad',
      badge: v.agreed
        ? (hostileSet.value.size
            ? 'BLOCK WRITTEN - the honest nodes were enough'
            : 'BLOCK WRITTEN - all five nodes decided the same transaction')
        : 'NOT WRITTEN - the transaction was discarded after the retries',
      cellClass: v.agreed ? 'cOk' : 'cBad',
      dot: false,
      attempt: -1,
      phase: 'Decision',
      short: v.agreed ? 'OK' : 'X',
      title: v.text
    })
  }
  return out
}

const frames = computed(() => buildFrames())

//The round has two broadcast phases and the page shows them as two separate
//lanes, because that is the distinction the protocol rests on: in the first one
//the commander signs and a deputy relays, in the second every deputy re-signs
//what it received. Phase II is further split per signing deputy.
const lanes = computed(() => {
  const numbered = frames.value.map((f, i) => ({ f, i }))
  const p1 = numbered.filter(x => x.f.kind === 'qds' && x.f.signer === 0)
  const p2 = numbered.filter(x => x.f.kind === 'qds' && x.f.signer !== 0)
  const groups = []
  for (const x of p2) {
    const last = groups[groups.length - 1]
    if (last && last.signer === x.f.signer) last.items.push(x)
    else groups.push({ signer: x.f.signer, label: 'node ' + x.f.signer + ' signs', items: [x] })
  }
  const judge = numbered.filter(x => x.f.kind === 'judge')
  const commit = numbered.filter(x => x.f.kind === 'report')
  const out = []
  if (p1.length) {
    out.push({
      title: 'Broadcast phase I - ' + p1.length + ' signatures',
      sub: 'the commander signs the batch, one deputy relays it to the other two',
      items: [{ signer: 0, label: 'node 0 is the signer', items: p1 }]
    })
  }
  if (groups.length) {
    out.push({
      title: 'Broadcast phase II - ' + p2.length + ' signatures',
      sub: 'each deputy re-signs what it received and distributes it',
      items: groups
    })
  }
  if (judge.length) {
    out.push({
      title: 'Inference - ' + judge.length + ' judgment lists',
      sub: 'each deputy takes the majority of its own judgment list',
      items: [{ signer: -1, label: '', items: judge }]
    })
  }
  if (commit.length) {
    const commitGroups = []
    for (const x of commit) {
      const last = commitGroups[commitGroups.length - 1]
      if (last && last.key === x.f.attempt) last.items.push(x)
      else commitGroups.push({ key: x.f.attempt, label: '', items: [x] })
    }
    const retried = commitGroups.some(g => g.key > 1)
    for (const g of commitGroups) {
      g.label = !retried ? '' : (g.key < 0 ? 'the round ends' : 'attempt ' + g.key)
    }
    if (commitGroups.length > 6) {
      const folded = [{ key: 0, label: 'the reports gathered over the retries of this round', items: [] }]
      for (const g of commitGroups) folded[0].items.push(...g.items)
      commitGroups.length = 0
      commitGroups.push(...folded)
    }
    out.push({
      title: 'Commit - ' + commit.length + ' reports',
      sub: 'each deputy answers the commander, and the round is settled',
      items: commitGroups
    })
  }
  return out
})
const current = computed(() => (cursor.value >= 0 && cursor.value < frames.value.length ? frames.value[cursor.value] : null))

// The map is drawn the way the protocol figure draws it: the commander at the
// top of the pentagon, the four deputies at the other vertices, and every QDS
// signature shown as the path signer -> relay -> verifier. Both the map and the
// broadcast list grow with the playhead, so replaying the round also replays the
// figure.
const PENTA = { 0: { x: 200, y: 42 }, 1: { x: 52, y: 150 }, 4: { x: 348, y: 150 }, 2: { x: 112, y: 324 }, 3: { x: 288, y: 324 } }

function pos(id) {
  return PENTA[id] || { x: 200, y: 190 }
}
function linkKey(a, b) {
  return Math.min(a, b) + '-' + Math.max(a, b)
}
//Each message of the round keeps one colour wherever it is printed: the cells, the
//broadcast tables, the links of the map and the replay card. m1 is the proposed
//batch and stays green; m2 and up get their own hue so a conflicting message is
//readable at a glance instead of looking like another verified m1.
const VALUE_HUES = { m2: 'v2', m3: 'v3' }


function valueTone(name) {
  const head = String(name || '').split(' ')[0]
  if (head === 'm1') return 'v1'
  if (VALUE_HUES[head]) return VALUE_HUES[head]
  if (head.indexOf('m') === 0) return 'v4'
  return 'v0'
}


function valueArrow(name) {
  const tone = valueTone(name)
  return 'url(#arr' + tone.charAt(0).toUpperCase() + tone.slice(1) + ')'
}


function arrowFor(tone) {
  return tone === 'ok' ? 'url(#arrOk)' : (tone === 'bad' ? 'url(#arrBad)' : 'url(#arrWarn)')
}

const visibleQds = computed(() => frames.value.slice(0, cursor.value + 1).filter(f => f.kind === 'qds'))

const links = computed(() => {
  const ids = nodeIndexes.value
  const out = []
  for (let a = 0; a < ids.length; a++) {
    for (let b = a + 1; b < ids.length; b++) {
      out.push({
        k: linkKey(ids[a], ids[b]), tone: '',
        x1: pos(ids[a]).x, y1: pos(ids[a]).y, x2: pos(ids[b]).x, y2: pos(ids[b]).y
      })
    }
  }
  const byKey = {}
  for (const l of out) byKey[l.k] = l
  for (const f of visibleQds.value) {
    const tone = f.tone === 'idle' ? '' : f.tone
    const mark = (x, y) => {
      const l = byKey[linkKey(x, y)]
      if (!l) return
      l.tone = tone
      if (f.name) l.msg = f.name
    }
    mark(f.signer, f.forward)
    mark(f.forward, f.verifier)
  }
  return out
})

//The broadcast lists of the protocol figure, split by depth. Depth 1 holds the
//signatures of the batch the commander signed; depth 2 is one sub-round per
//deputy that re-signed what it received, named MR2_S->Ri. Rows are the relay and
//columns the receiver, exactly as the figure prints F \ V, so a cell answers the
//one question the protocol cares about: which message did this relay show here.
const tables = computed(() => {
  const depth1 = []
  const bySigner = new Map()
  for (const f of visibleQds.value) {
    if (f.signer === 0) {
      depth1.push(f)
      continue
    }
    if (!bySigner.has(f.signer)) bySigner.set(f.signer, [])
    bySigner.get(f.signer).push(f)
  }
  const out = []
  if (depth1.length) {
    out.push(makeTable('d1', 'Depth 1 - MR1_S', 'the batch S signed, relayed by each deputy', depth1))
  }
  for (const s of Array.from(bySigner.keys()).sort((a, b) => a - b)) {
    out.push(makeTable('d2' + s, 'Depth 2 - MR2_S-' + '>' + 'R' + s,
      'R' + s + ' re-signs what it received and distributes it', bySigner.get(s)))
  }
  return out
})

function makeTable(key, title, sub, list) {
  const cells = {}
  const rows = []
  const cols = []
  for (const f of list) {
    if (f.forward < 0 || f.verifier < 0) continue
    if (rows.indexOf(f.forward) < 0) rows.push(f.forward)
    if (cols.indexOf(f.verifier) < 0) cols.push(f.verifier)
    const k = f.forward + ':' + f.verifier
    const bucket = cells[k] || (cells[k] = [])
    const found = bucket.find(v => v.name === f.name && v.tone === f.tone)
    if (found) found.times++
    else bucket.push({ name: f.name || 'm?', times: 1, tone: f.tone })
  }
  rows.sort((a, b) => a - b)
  cols.sort((a, b) => a - b)
  //The protocol table prints one more entry per row: what the relay itself keeps,
  //the message it received from the signer and then handed to the two verifiers.
  //It sits on the diagonal, so a row reads m1, m1, m1 or m2, m2, m2 exactly as the
  //figure prints it, and every honest column ends up with three entries.
  for (const r of rows) {
    const own = list.find(f => f.forward === r && f.tone === 'ok' && f.name)
    if (own) cells[r + ':' + r] = [{ name: own.name, times: 1, tone: 'ok' }]
  }
  //each receiver decides out of its own column and only accepted messages are
  //countable: an empty column or a tie is the abstention of the protocol
  const maj = {}
  for (const c of cols) {
    const votes = {}
    const add = (nm) => { if (nm) votes[nm] = (votes[nm] || 0) + 1 }
    for (const r of rows) {
      for (const v of (cells[r + ':' + c] || [])) {
        if (v.tone !== 'ok') continue
        for (let i = 0; i < v.times; i++) add(v.name)
      }
    }
    let best = ''
    let bestCount = 0
    let tie = false
    for (const nm of Object.keys(votes)) {
      if (votes[nm] > bestCount) { bestCount = votes[nm]; best = nm; tie = false }
      else if (votes[nm] === bestCount) tie = true
    }
    maj[c] = best === '' ? { text: 'abstain', tone: 'warn' }
      : (tie ? { text: best + ' (tie)', tone: 'warn' } : { text: best, tone: 'ok' })
  }
  return { key, title, sub, rows, cols, cells, maj }
}

function isNowCell(tb, relay, receiver) {
  const f = current.value
  if (!f || f.kind !== 'qds') return false
  if (f.forward !== relay || f.verifier !== receiver) return false
  return tb.key === 'd1' ? f.signer === 0 : String(f.signer) === tb.key.slice(2)
}

//Every distinct value of the round keeps the short name the frames gave it, so
//the tables, the replay card and the report line all speak the same language.
const valueLabels = computed(() => {
  const map = {}
  for (const f of frames.value) {
    //the values a node only ever *reports* are named as well, otherwise a
    //dishonest report would have to be shown as a raw digest
    if (f.digest && f.name) map[String(f.digest)] = f.name
  }
  return map
})

//What each node decided out of the messages it was shown, and whether the report
//it sent the commander matches that decision.
const judgeVisible = computed(() => frames.value.slice(0, cursor.value + 1).filter(f => f.kind === 'judge'))

function majorityOf(node) {
  if (latestRound.value < 0) return null
  if (Number(node) === 0) {
    const v = verdict.value
    if (v.agreed === null) return null
    return { text: 'm1', note: 'the batch it proposed', tone: v.agreed ? 'ok' : 'warn' }
  }
  const j = judgeVisible.value.find(x => x.node === Number(node))
  if (!j) return null
  if (!j.entries.length) return { text: '-', note: 'no decision to report', tone: 'warn' }
  const label = valueLabels.value[j.decision] || (j.decision || '-')
  if (j.notes.length) {
    //the trace says what the node claimed instead: show the lie, not only that
    //there was one, so the line reads like the protocol figure annotation
    const claim = claimedValue(j.notes[0].text)
    const claimName = claim ? (valueLabels.value[claim] || claim) : ''
    return {
      text: claimName ? label + ' > ' + claimName : label,
      note: 'decided > reported, in the order the values appeared',
      tone: 'bad',
      lied: true
    }
  }
  return { text: label, note: 'reported to the commander', tone: 'ok' }
}

//the digest an honest report would have carried is left out of the trace line,
//so the value a liar claims is read straight off the sentence that records it
function claimedValue(text) {
  const s = String(text || '')
  const at = s.indexOf('decided ')
  if (at < 0) return ''
  return s.slice(at + 8, at + 20).trim()
}

const reportsLine = computed(() => {
  const out = []
  for (const n of nodeIndexes.value) {
    const m = majorityOf(n)
    if (m) out.push({ node: n, text: (Number(n) === 0 ? 'S' : 'R' + n) + ': ' + m.text, tone: m.tone })
  }
  return out
})

function stopPlay() {
  playing.value = false
  if (playTimer) { clearInterval(playTimer); playTimer = null }
}

function startPlay() {
  stopPlay()
  if (!frames.value.length) return
  if (cursor.value >= frames.value.length - 1) cursor.value = -1
  playing.value = true
  playTimer = setInterval(() => {
    if (cursor.value >= frames.value.length - 1) { stopPlay(); return }
    cursor.value++
  }, speedMs.value)
}

function togglePlay() {
  if (playing.value) stopPlay()
  else startPlay()
}

function stepBy(delta) {
  stopPlay()
  const next = cursor.value + delta
  if (next < -1) return
  cursor.value = Math.min(next, frames.value.length - 1)
}

function restartTimeline() {
  stopPlay()
  cursor.value = -1
  startPlay()
}

function goTo(i) {
  stopPlay()
  cursor.value = i
}

//From the start of a round to the last frame of it: the replay is entered after
//the round has finished, and follow-live pins the playhead to the newest frame
//while the round is still running.
function afterRound() {
  if (followLive.value) {
    stopPlay()
    cursor.value = frames.value.length - 1
    return
  }
  cursor.value = -1
  startPlay()
}

watch(speedMs, () => { if (playing.value) startPlay() })
watch(followLive, () => { if (followLive.value) stopPlay() })
//Two ways the playhead follows the round on its own: while a live round grows it
//stays on the newest frame, and a page opened onto a finished round shows all of
//it instead of an empty map.
watch(() => frames.value.length, (n) => {
  if (followLive.value && busy.value) cursor.value = n - 1
  else if (n > 0 && cursor.value < 0 && !playing.value && !busy.value) cursor.value = n - 1
})
onUnmounted(stopPlay)

</script>

<style scoped>
/* Sized for screen recording: the smallest text on this page is 15px. */
.byz-page { padding: 14px 18px 40px; color: #22303f; font-size: 17px; line-height: 1.5; }
.head { display: flex; align-items: baseline; justify-content: space-between; flex-wrap: wrap; gap: 10px; margin-bottom: 14px; }
.head h2 { margin: 0; font-size: 28px; letter-spacing: .2px; }
.facts-bar { display: flex; gap: 18px; font-size: 16px; color: #46586b; flex-wrap: wrap; }
.facts-bar b { color: #16222e; }
.alert { background: #fdeaea; border-left: 6px solid #c0303a; color: #8f1d24; padding: 12px 14px; font-size: 16px; border-radius: 4px; margin-bottom: 14px; }
.panel { background: #fff; border: 1px solid #dde4ec; border-radius: 8px; padding: 18px 20px 20px; margin-bottom: 18px; }
.ptitle { font-weight: 700; font-size: 20px; margin-bottom: 14px; color: #16222e; }
.grid2 { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.topology { display: grid; grid-template-columns: repeat(5, 1fr); gap: 14px; }
.node { border: 3px solid #cfd9e4; border-radius: 8px; padding: 16px 10px; text-align: center; background: #f8fbfd; }
.node.hostile { border-color: #c0303a; background: #fff4f4; }
.node.primary { border-style: double; }
.ntag { display: inline-block; background: #16222e; color: #ffd04b; font-size: 13px; font-weight: 700; padding: 2px 9px; border-radius: 10px; margin-bottom: 8px; }
.nname { font-size: 24px; font-weight: 700; }
.nstate { font-size: 19px; font-weight: 700; color: #1f7a35; margin-top: 4px; letter-spacing: .5px; }
.node.hostile .nstate { color: #b3222c; }
.nrole { font-size: 15px; color: #5b6c7d; margin-top: 6px; min-height: 21px; }
.note { margin-top: 12px; padding: 10px 12px; border-left: 5px solid #e6a23c; background: #fdf6ec; font-size: 16px; }
.lab { display: block; font-size: 16px; font-weight: 700; margin: 14px 0 7px; }
.w100 { width: 100%; }
.picker { display: flex; flex-wrap: wrap; gap: 6px 20px; font-size: 17px; }
.quick { margin-top: 12px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.qlab { font-size: 15px; color: #5b6c7d; }
.hint { margin-top: 14px; padding: 10px 12px; border-left: 5px solid #52a35c; background: #f1f8ee; color: #2f5d2a; font-size: 15.5px; border-radius: 0 4px 4px 0; }
.hint.warn { border-left-color: #e6a23c; background: #fdf6ec; color: #8a5a15; }
.btns { margin-top: 18px; display: flex; gap: 10px; flex-wrap: wrap; }
.verdict { font-weight: 800; font-size: 21px; padding: 14px; border-radius: 6px; margin-bottom: 14px; text-align: center; }
.verdict.ok { background: #e4f6e8; color: #1c6b2f; }
.verdict.bad { background: #fdeaea; color: #a91f28; }
.verdict.idle { background: #eef2f6; color: #66788a; }
.dim { color: #8494a3; font-size: 15px; }
.expl { font-size: 16px; color: #46586b; margin: 0 0 16px; }
.sBad { color: #a91f28; }
.sWarn { color: #9a5c0e; }
.sOk { color: #1c6b2f; }
.toggle { margin-top: 12px; font-size: 15px; }
.rawbox { max-height: 320px; overflow: auto; border: 1px solid #eef2f6; margin-top: 10px; }
.raw { width: 100%; border-collapse: collapse; font-size: 15px; }
.raw th { position: sticky; top: 0; background: #f2f6f9; text-align: left; padding: 8px; }
.raw td { padding: 7px 8px; border-bottom: 1px solid #f4f7f9; }
.raw .det { white-space: normal; min-width: 320px; }
.row-fault td, .row-bad td { background: #fdeeee; }
.row-round td { background: #fff5e8; }
.row-ok td { background: #eaf7ed; }
.hostileCell { color: #b3222c; font-weight: 700; }
.chip { display: inline-block; margin: 3px 8px 3px 0; padding: 5px 11px; border-radius: 14px; background: #eaf7ed; border: 1px solid #cfe9d5; font-size: 15.5px; }
.chipBad { background: #fdeeee; border-color: #f2d2d4; color: #a91f28; }
.good { color: #1c6b2f; font-weight: 700; }
.bad { color: #a91f28; font-weight: 700; }
code { font-family: Consolas, monospace; font-size: 15px; }/* ---- QDS timeline ---- */
/* ---- pentagon map and broadcast list ---- */
.maprow { display: grid; grid-template-columns: 340px 1fr; gap: 24px; align-items: start; margin-bottom: 18px; }
.pentagon { width: 100%; height: auto; }
.mapcap { font-size: 14.5px; color: #6a7b8c; margin-top: 8px; line-height: 1.45; }
.lk { stroke: #dce4eb; stroke-width: 2; }
.lkok { stroke: #86c093; stroke-width: 4; }
.lkbad { stroke: #e09aa0; stroke-width: 4; }
.lkwarn { stroke: #f0c894; stroke-width: 4; }
.flow { stroke-width: 7; stroke-linecap: round; opacity: .92; }
.fok { stroke: #1f7a35; }
.fv1 { stroke: #1f7a35; }
.fv2 { stroke: #c97f10; }
.fv3 { stroke: #6b52c4; }
.fv4 { stroke: #b0398f; }
.lkv1 { stroke: #86c093; stroke-width: 4; }
.lkv2 { stroke: #e0a63d; stroke-width: 4; }
.lkv3 { stroke: #9c8ade; stroke-width: 4; }
.lkv4 { stroke: #d170b4; stroke-width: 4; }
.bval.v1 { color: #1c6b2f; }
.bval.v2 { color: #8a5b00; }
.bval.v3 { color: #49329a; }
.bval.v4 { color: #8c1d67; }
.legend .v1 { background: #eaf7ed; border-color: #b9dfc2; }
.legend .v2 { background: #fdf1d8; border-color: #e6bd68; }
.legend .v3 { background: #f1ecfb; border-color: #c2b1e8; }
.legend .v4 { background: #fbe7f3; border-color: #e19ccd; }
.legnote { flex-basis: 100%; color: #6a7b8c; font-size: 14.5px; }
.txcard { margin: 12px 0; padding: 11px 13px; background: #f7fafc; border: 1px solid #dde5ee; border-left: 4px solid #2f6fd0; border-radius: 7px; }
.txcap { font-size: 13.5px; letter-spacing: .04em; text-transform: uppercase; color: #6a7b8c; }
.txid { margin-top: 3px; font-family: Consolas, Menlo, monospace; font-size: 19px; font-weight: 700; color: #1b3a63; }
.txfields { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 7px; }
.txfields span { padding: 3px 9px; background: #fff; border: 1px solid #dbe3ea; border-radius: 12px; font-size: 14.5px; color: #33465a; }
.txdigest { margin-top: 7px; font-size: 15px; color: #33465a; }
.attempts { display: flex; flex-wrap: wrap; align-items: center; gap: 6px; margin-top: 4px; padding-top: 10px; border-top: 1px dashed #dde5ee; }
.acap { margin-right: 3px; font-size: 13.5px; letter-spacing: .04em; text-transform: uppercase; color: #6a7b8c; }
.achip { padding: 3px 10px; border-radius: 12px; font-size: 14px; font-weight: 600; background: #f1f4f8; border: 1px solid #d5dde6; color: #55667a; }
.achip.fail { background: #fdeeee; border-color: #eab7bb; color: #a91f28; }
.achip.ok { background: #e3f4e7; border-color: #a9d6b4; color: #17632a; }
.adetail { flex-basis: 100%; font-size: 14.5px; color: #55667a; }
.ledgerline { margin-top: 10px; font-size: 17px; font-weight: 600; color: #33465a; }
.stepfoot { margin-top: 10px; font-size: 14.5px; color: #45566a; }
.tlink { color: #1f5fbf; font-weight: 600; }
.ptt-note { margin-left: 10px; font-size: 13.5px; font-weight: 400; color: #7d8b99; }
.fbad { stroke: #b3222c; }
.fwarn { stroke: #c97f10; }
.fidle { stroke: #66788a; }
.ndot { fill: #dceafa; stroke: #2f6fb5; stroke-width: 3; }
.ndHostile { fill: #fbdcdc; stroke: #b3222c; }
.ndPrimary { fill: #ffeab8; stroke: #a97a00; }
.ndlab { font-size: 22px; font-weight: 800; fill: #1d2c3a; text-anchor: middle; }
.ndsub { font-size: 12.5px; fill: #46586b; text-anchor: middle; }
.tables { display: grid; grid-template-columns: 1fr 1fr; gap: 14px 18px; }
.subtable.wide { grid-column: 1 / -1; }
.stab { font-size: 16px; font-weight: 800; color: #2c3d4e; margin-bottom: 6px; }
.stabsub { font-size: 14px; font-weight: 400; color: #6a7b8c; margin-left: 6px; }
.reports { grid-column: 1 / -1; border-top: 1px dashed #dde4ec; padding-top: 10px; display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }
.rlab { font-size: 15px; color: #46586b; font-weight: 700; margin-right: 4px; }
.rdim { font-size: 14px; color: #8494a3; margin-right: 6px; }
.mtable { width: 100%; border-collapse: collapse; font-size: 16px; }
.mtable th { background: #f2f6f9; padding: 6px 6px; font-size: 15px; color: #2c3d4e; }
.mtable .axis { font-size: 14px; color: #52647a; }
.mtable td { border: 1px solid #e6ecf2; padding: 5px 4px; height: 40px; vertical-align: middle; text-align: center; }
.mtable .cellSelf { background: #f7f9fb; }
.mtable .cellNow { box-shadow: inset 0 0 0 3px #1f6fd0; }
.mchip { display: inline-block; padding: 4px 9px; border-radius: 12px; font-size: 16px; font-weight: 700; background: #e3f4e7; color: #17632a; border: 1px solid #a9d6b4; margin: 2px 3px 2px 0; }
.mtable .mfoot .mchip { font-size: 17px; }
.mchip.mbad { background: #fdeeee; border-color: #eab7bb; color: #a91f28; text-decoration: line-through; }
.mchip.mwarn { background: #fff5e8; border-color: #f0cfa0; color: #9a5c0e; text-decoration: line-through; }
.mchip.midle, .mchip.m { background: #f1f4f8; border-color: #d5dde6; color: #66788a; }
.mchip.v1 { background: #e3f4e7; border-color: #a9d6b4; color: #17632a; }
.mchip.v2 { background: #fdf1d8; border-color: #e6bd68; color: #7c5405; }
.mchip.v3 { background: #f1ecfb; border-color: #c2b1e8; color: #49329a; }
.mchip.v4 { background: #fbe7f3; border-color: #e19ccd; color: #8c1d67; }
.mchip.v2, .mchip.v3, .mchip.v4 { text-decoration: none; }
.hostileHead { color: #b3222c; }
.mtable .mfoot th { background: #f7f9fb; font-size: 15px; color: #2c3d4e; text-align: right; }
.mtable .mfoot td { background: #fcfdfe; }
.mnote { font-size: 14px; color: #52647a; margin-top: 4px; line-height: 1.3; }
.tlbar { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; margin-bottom: 14px; }
.wspeed { width: 130px; }
.sw { margin-left: auto; }
.lanes { display: flex; flex-direction: column; gap: 12px; margin-bottom: 12px; }
.lane { display: grid; grid-template-columns: 250px 1fr; gap: 16px; align-items: start; border-top: 1px dashed #e2e8ef; padding-top: 10px; }
.lane:first-child { border-top: none; padding-top: 0; }
.lcap .lt { font-size: 17px; font-weight: 800; color: #2c3d4e; }
.lcap .ls { font-size: 14.5px; color: #6a7b8c; margin-top: 2px; }
.lbody { display: flex; flex-wrap: wrap; gap: 18px; }
.lgroup .glab { font-size: 14.5px; color: #5b6c7d; margin-bottom: 5px; font-weight: 700; }
.gcells { display: flex; flex-wrap: wrap; gap: 4px; }
.cell { position: relative; min-width: 32px; height: 32px; border: 1px solid #cfd8e2; border-radius: 5px; background: #f5f8fa; font-size: 14.5px; color: #46586b; cursor: pointer; padding: 0 5px; }
.cell.cOk { background: #eaf7ed; border-color: #b9dfc2; color: #1c6b2f; }
.cell.cBad { background: #fdeeee; border-color: #eab7bb; color: #a91f28; }
.cell.cWarn { background: #fff5e8; border-color: #f0cfa0; color: #9a5c0e; }
.cell.v1 { background: #eaf7ed; border-color: #b9dfc2; color: #1c6b2f; }
.cell.v2 { background: #fdf1d8; border-color: #e6bd68; color: #7c5405; }
.cell.v3 { background: #f1ecfb; border-color: #c2b1e8; color: #49329a; }
.cell.v4 { background: #fbe7f3; border-color: #e19ccd; color: #8c1d67; }
.cell.cJudge { background: #eef3fb; border-color: #c9d9ee; color: #2c5080; }
.cell.cReport { background: #f1f4f8; border-color: #d5dde6; color: #35506b; }
.cell.at { outline: 3px solid #1f6fd0; outline-offset: 1px; font-weight: 800; }
.cell .dot { position: absolute; right: 2px; top: 2px; width: 8px; height: 8px; border-radius: 50%; background: #b3222c; }
.legend { display: flex; flex-wrap: wrap; gap: 8px 20px; font-size: 15px; color: #46586b; margin-bottom: 16px; }
.legend .k { display: inline-block; width: 15px; height: 15px; border-radius: 3px; border: 1px solid #cfd8e2; vertical-align: -2px; margin-right: 6px; }
.legend .cOk { background: #eaf7ed; border-color: #b9dfc2; }
.legend .cBad { background: #fdeeee; border-color: #eab7bb; }
.legend .cWarn { background: #fff5e8; border-color: #f0cfa0; }
.legend .cJudge { background: #eef3fb; border-color: #c9d9ee; }
.legend .cReport { background: #f1f4f8; border-color: #d5dde6; }
.legend .dotk { width: 9px; height: 9px; border-radius: 50%; background: #b3222c; border: none; }
.frame { border: 1px solid #e5ebf1; border-left: 9px solid #9fb0c0; border-radius: 8px; padding: 16px 18px; background: #fbfcfd; }
.frame.tone-ok { border-left-color: #1f7a35; }
.frame.tone-bad { border-left-color: #b3222c; }
.frame.tone-warn { border-left-color: #e6a23c; }
.fhead { display: flex; justify-content: space-between; gap: 16px; font-size: 16.5px; color: #46586b; margin-bottom: 14px; }
.fphase { font-weight: 800; }
.qdsrow { display: flex; align-items: center; gap: 14px; flex-wrap: wrap; }
.role { border: 2px solid #cfd8e2; border-radius: 8px; padding: 10px 18px; min-width: 140px; text-align: center; background: #fff; }
.role.hostile { border-color: #d33; background: #fdeeee; }
.rlab { font-size: 14.5px; color: #6a7b8c; letter-spacing: .6px; text-transform: uppercase; }
.rname { font-size: 22px; font-weight: 800; }
.role.hostile .rname { color: #b3222c; }
.ar { font-size: 28px; color: #9fb0c0; }
.msg { font-size: 18px; }
.jrow .chips { margin: 10px 0; }
.badge { margin-top: 14px; display: inline-block; font-size: 19px; font-weight: 800; padding: 10px 16px; border-radius: 6px; }
.b-ok { background: #e4f6e8; color: #1c6b2f; }
.b-bad { background: #fdeaea; color: #a91f28; }
.b-warn { background: #fff5e8; color: #9a5c0e; }
.b-info { background: #eaf1fb; color: #2c5080; }
.b-idle { background: #eef2f6; color: #66788a; }
.fnote { margin-top: 10px; font-size: 17px; font-weight: 700; border-left: 5px solid #b3222c; padding-left: 10px; color: #a91f28; }
.fnote.sWarn { border-left-color: #e6a23c; color: #9a5c0e; }

</style>

