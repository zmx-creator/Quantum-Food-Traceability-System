<template>
  <el-card class="chaincode-info-card">
    <el-table :data="chaincodes" border style="width: 100%" size="small">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="Chaincode Name" />
      <el-table-column prop="version" label="Version" width="120" />
      <el-table-column prop="channel" label="Channel" width="120" />
      <el-table-column prop="txCount" label="Tx Count" width="120" />
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBlockchainDetails, getAllTransactions } from '../api/fabric'

const chaincodes = ref([])

onMounted(async () => {

  const [blockRes, txRes] = await Promise.all([
    getBlockchainDetails(),
    getAllTransactions()
  ])
  let chaincodeName = 'source-app'
  let version = '1.0'
  let channel = 'mychannel'
  let txCount = 0
  if (txRes.data.success && Array.isArray(txRes.data.data)) {

    txCount = txRes.data.data.filter(t => t.chaincodeName === chaincodeName).length
  }
  chaincodes.value = [{
    id: 1,
    name: chaincodeName,
    version,
    channel,
    txCount
  }]
})
</script>

<style scoped>
.chaincode-info-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  border-radius: 8px;
  padding: 0 0 8px 0;
}
</style> 
