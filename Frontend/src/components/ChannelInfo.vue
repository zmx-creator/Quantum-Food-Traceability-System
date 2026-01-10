<template>
  <el-card class="channel-info-card">
    <el-table :data="[channelInfo]" border style="width: 100%" size="small">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="Channel Name" />
      <el-table-column prop="blocks" label="Blocks" width="100" />
      <el-table-column prop="transactions" label="Transactions" width="120" />
      <!-- <el-table-column prop="timestamp" label="Timestamp" min-width="200" /> -->
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBlockchainDetails, getAllTransactions } from '../api/fabric'

const channelInfo = ref({
  id: '',
  name: '',
  blocks: 0,
  transactions: 0,
  timestamp: ''
})

onMounted(async () => {

  const [blockRes, txRes] = await Promise.all([
    getBlockchainDetails(),
    getAllTransactions()
  ])
  if (blockRes.data.success && blockRes.data.data) {
    channelInfo.value.id = 3 
    channelInfo.value.name = 'mychannel' 
    channelInfo.value.blocks = blockRes.data.data.totalBlocks || 0
    channelInfo.value.timestamp = blockRes.data.data.blocks && blockRes.data.data.blocks.length > 0
      ? blockRes.data.data.blocks[blockRes.data.data.blocks.length - 1].timestamp || ''
      : ''
  }
  if (txRes.data.success && Array.isArray(txRes.data.data)) {
    channelInfo.value.transactions = txRes.data.data.length
  }
})
</script>

<style scoped>
.channel-info-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  border-radius: 8px;
  padding: 0 0 8px 0;
}
</style> 
