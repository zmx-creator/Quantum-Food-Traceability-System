<template>
  <div class="summary card-shadow">
    <div class="summary-item">
      <div class="icon block"></div>
      <div class="count">
		<span style="display: inline-flex; align-items: center;">
			<img src="../assets/greenNode.png" alt="icon" 
				style="width:20px;height:20px;margin-right: 16px;">
			{{ blockCount }}
		</span>	
	  </div>
      <div class="label">BLOCKS</div>
    </div>
    <div class="summary-item">
      <div class="icon tx"></div>
      <div class="count">
			<span style="display: inline-flex; align-items: center;">
				<img src="../assets/blueNode.png" alt="icon" 
				style="width:20px;height:20px;margin-right: 16px;">
				{{ txCount }}
			</span>	
		</div>
      <div class="label">TRANSACTIONS</div>
    </div>
    <div class="summary-item">
      <div class="icon node"></div>
      <div class="count">
			<span style="display: inline-flex; align-items: center;">
				<img src="../assets/orangeNode.png" alt="icon" 
				style="width:20px;height:20px;margin-right: 16px;">
				{{ nodeCount }}
			</span>	
		</div>
      <div class="label">NODES</div>
    </div>
    <div class="summary-item">
      <div class="icon chaincode"></div>
      <div class="count">
			<span style="display: inline-flex; align-items: center;">
				<img src="../assets/redNode.png" alt="icon" 
				style="width:20px;height:20px;margin-right: 16px;">
				{{ chaincodeCount }}
			</span>		  
	  </div>
      <div class="label">CHAINCODES</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBlockchainDetails, getAllTransactions, getPeers, getOrderers } from '../api/fabric'

const blockCount = ref(0)
const txCount = ref(0)
const nodeCount = ref(0)
const chaincodeCount = ref(1) 

onMounted(async () => {

  try {
    const res = await getBlockchainDetails()
    if (res.data.success && res.data.data) {
      blockCount.value = res.data.data.totalBlocks
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
  } catch (e) {
    nodeCount.value = 'failed to get'
  }
})
</script>

<style scoped>
.summary { 
  display: flex; 
  justify-content: space-around; 
  margin-bottom: 20px; 
 
}
.card-shadow {
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,0.10);
  background: #fff;
  padding: 30px 0;
}
.summary-item { 
  text-align: center; 
}
.summary-item .count { 
  font-size: 2em; 
  font-weight: bold; 
}
</style>
