<template>
  <div class="block-section">
    <h3>Block Information</h3>
    <div class="timeline-scroll">
      <el-timeline>
        <el-timeline-item
          v-for="(block, idx) in blocks"
          :key="block.number"
          :timestamp="block.time || '---'"
          placement="top"
          color="#2196F3"
        >
          <div :class="['block-title', { 'odd-bg': idx % 2 === 0 }]">Block {{ block.number }}</div>
          <div :class="['block-info', { 'odd-bg': idx % 2 === 0 }]">
            <div >Channel Name: {{ block.channel }}</div>
            <div>Datahash: {{ block.hash }}</div>
            <div >Number of Tx: {{ block.tx }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElTimeline, ElTimelineItem } from 'element-plus'
import { getBlockchainDetails, getAllTransactions, getPeers, getOrderers } from '../api/fabric'

const blockCount = ref(0)
const txCount = ref(0)
const nodeCount = ref(0)
const chaincodeCount = ref(1) 

const peers = ref([])
const blocks = ref([])

onMounted(async () => {

  try {
    const res = await getBlockchainDetails()
    if (res.data.success && res.data.data) {
      blockCount.value = res.data.data.totalBlocks
      blocks.value = (res.data.data.blocks || []).map(b => ({
        number: b.blockNumber,
        channel: 'mychannel',
        hash: b.dataHash || '',
        tx: b.transactionCount,
        time: '' 
      }))
    }
  } catch (e) {
    blockCount.value = 'failed to get'
  }


  try {
    const res = await getAllTransactions()
    if (res.data.success && Array.isArray(res.data.data)) {
      txCount.value = res.data.data.length
    }
  } catch (e) {
    txCount.value = 'failed to get'
  }


  try {
    const [peerRes, ordererRes] = await Promise.all([getPeers(), getOrderers()])
    let peerList = Array.isArray(peerRes.data.data) ? peerRes.data.data : []
    let ordererList = Array.isArray(ordererRes.data.data) ? ordererRes.data.data : []
    nodeCount.value = peerList.length + ordererList.length
    peers.value = [...peerList.map(p => p.name), ...ordererList.map(o => o.name)]
  } catch (e) {
    nodeCount.value = 'failed to get'
  }
})
</script>

<style scoped>
.block-section {
  margin-bottom: 20px;
  background-color: #f4f4f5;
 /* border: 2px solid #e0e0e0;; */
}

.block-section h3 {
  margin-bottom: 10px;
  color: #333;
  
}

.timeline-scroll {
  max-height: 300px;
  overflow-y: auto;
  padding-right: 4px;
/*   border: 2px solid #f4f4f5;;
   border-radius: 10px;
   box-shadow: 0 2px 8px 0 rgba(0,0,0,0.10);
   background-color: #ffffff; */
}

.block-title {
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 8px;
  padding: 6px 0 0 0;
  background-color: #617a9e;
}

.block-info {
  font-size: 0.9em;
  color: #555;
  padding: 8px 12px 8px 12px;
  border-radius: 6px;
  margin-bottom: 8px;
  text-align: left;
}

.block-info > div {
  margin-bottom: 4px;
}

/* .odd-bg {
  background: #e0e0e0 !important;
} */
</style>
