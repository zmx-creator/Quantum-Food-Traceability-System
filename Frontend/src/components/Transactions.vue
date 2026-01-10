<template>

	<div class="PaperName" style="display: flex; align-items: center;">
		  <div style="flex: 1; display: flex; align-items: center; justify-content: flex-start; font-size: 20px;">
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Fold /></el-icon>
				<span style="margin-left: 4px;">TRACEABILITY</span>
			</span>
		    <el-divider direction="vertical" />
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Memo /></el-icon>
				<span style="margin-left: 4px;">Transaction Information</span>
			</span>

		  </div>
		   <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><Setting /></el-icon>
		   <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><UserFilled /></el-icon>	  
		  
	</div>
	<el-divider style="margin: 8px 0;" />
  
  <div class="transactions">
    <div class="header" style="display: flex; align-items: center; justify-content: space-between; font-size: 24px;">
      <el-text class="header-title" size="large" style="margin-left: 10px;">

			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:10px;height:10px;margin-right: 8px;">	
					  
			<img src="../assets/orangeNode.png" alt="icon"
			  style="vertical-align: middle;width:12px;height:12px;margin-right: 8px;">
			<img src="../assets/greenNode.png" alt="icon"
			  style="vertical-align: middle;width:14px;height:14px;margin-right: 8px;">
			<img src="../assets/blueNode.png" alt="icon"
			  style="vertical-align: middle;width:16px;height:16px;margin-right: 8px;">
		  Blockchain Transaction Records
			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:18px;height:18px;margin-left: 8px;">
			<img src="../assets/greenNode.png" alt="icon"
				style="vertical-align: middle;width:14px;height:14px;margin-left: 8px;">
			<img src="../assets/orangeNode.png" alt="icon"
				style="vertical-align: middle;width:12px;height:12px;margin-left: 8px;">	
				  
			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:10px;height:10px;margin-left: 8px;">		

	
		</el-text>
      <el-button   type="primary" plain @click="refreshData" :disabled="loading" class="refresh-btn" >
        {{ loading ? 'Load...' : 'Refresh' }}
      </el-button>
    </div>
    <el-descriptions :column="2" border style="box-shadow: 4px 2px 8px 0 rgba(0,0,0,0.10);">
      <el-descriptions-item label="Transaction Volume">{{ transactions.length }}</el-descriptions-item>
      <el-descriptions-item label="Blockchain Height">{{ blockHeight }}</el-descriptions-item>
      <el-descriptions-item label="DataHash">{{ latestBlockHash }}</el-descriptions-item>
    </el-descriptions>

    <div v-if="loading" class="loading">
      <p>Loading...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>fail to load : {{ error }}</p>
      <button @click="refreshData">retry</button>
    </div>

    <div v-else-if="transactions.length === 0" class="empty">
      <p>No transaction records available</p>
    </div>

    <div v-else class="transactions-table">
      <el-table :data="transactions" border style="width: 100% ;" >
        <el-table-column prop="transactionID" label="TransactionID" min-width="180">
          <template #default="scope">
            <span class="tx-id">{{ formatTransactionId(scope.row.transactionID) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="blockNumber" label="Block Number" min-width="80" />
        <el-table-column prop="timestamp" label="Timestamp" min-width="160">
          <template #default="scope">
            {{ formatTimestamp(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="creatorMSP" label="Creator MSP" min-width="120">
          <template #default="scope">
            {{ scope.row.creatorMSP || 'N/A' }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="Transaction Type" min-width="120">
          <template #default="scope">
            {{ formatType(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="chaincodeName" label="Chaincode Name" min-width="120">
          <template #default="scope">
            {{ scope.row.chaincodeName || 'N/A' }}
          </template>
        </el-table-column>
        <el-table-column label="Operation" min-width="80">
          <template #default="scope">
            <el-button size="small" type="primary" plain @click="showTransactionDetails(scope.row)">Detials</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>


    <div v-if="selectedTransaction" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Transaction details</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="detail-item">
            <label>TransactionID:</label>
            <span>{{ selectedTransaction.transactionID }}</span>
          </div>
          <div class="detail-item">
            <label>Block Number:</label>
            <span>{{ selectedTransaction.blockNumber }}</span>
          </div>
          <div class="detail-item">
            <label>Timestamp:</label>
            <span>{{ formatTimestamp(selectedTransaction.timestamp) }}</span>
          </div>
          <div class="detail-item">
            <label>CreatorMSP:</label>
            <span>{{ selectedTransaction.creatorMSP || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <label>CreatorID:</label>
            <span>{{ selectedTransaction.creator || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <label>Type:</label>
            <span>{{ formatType(selectedTransaction.type) }}</span>
          </div>
          <div class="detail-item">
            <label>Chaincode Name:</label>
            <span>{{ selectedTransaction.chaincodeName || 'N/A' }}</span>
          </div>
          <div v-if="selectedTransaction.payloads && selectedTransaction.payloads.length > 0" class="detail-item">
            <label>Payload Data:</label>
            <div class="payloads">
              <div v-for="(payload, index) in selectedTransaction.payloads" :key="index" class="payload-item">
                <span class="payload-index">{{ index + 1 }}.</span>
                <span class="payload-content">{{ payload }}</span>
              </div>
            </div>
          </div>
<!--          <div class="detail-item">
            <label>ProductID:</label>
            <span>{{ getProductIDFromPayload(selectedTransaction) }}</span>
          </div> -->
        </div>
      </div>
    </div>
	
		
	</div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAllTransactions, getBlockchainInfo, getBlockchainDetails } from '@/api/fabric'

const transactions = ref([])
const blockchainInfo = ref('')
const loading = ref(false)
const error = ref('')
const selectedTransaction = ref(null)
const blockHeight = ref('')
const latestBlockHash = ref('')


const fetchTransactions = async () => {
  loading.value = true
  error.value = ''
  try {
    const [txResponse, infoResponse, blockDetailRes] = await Promise.all([
      getAllTransactions(),
      getBlockchainInfo(),
      getBlockchainDetails()
    ])
    if (txResponse.data.success) {
      transactions.value = txResponse.data.data || []
    } else {
      error.value = txResponse.data.message || 'failed to get'
    }
    if (infoResponse.data.success) {
      blockchainInfo.value = infoResponse.data.data || ''
    }
    if (blockDetailRes.data.success && blockDetailRes.data.data) {
      blockHeight.value = blockDetailRes.data.data.totalBlocks || ''

      const blocks = blockDetailRes.data.data.blocks || []
      if (blocks.length > 0) {
        latestBlockHash.value = blocks[blocks.length - 1].dataHash || ''
      } else {
        latestBlockHash.value = ''
      }
    }
  } catch (err) {
    console.error('failed to get:', err)
    error.value = err.response?.data?.message || err.message || 'network error'
  } finally {
    loading.value = false
  }
}


const refreshData = () => {
  fetchTransactions()
}


const formatTransactionId = (id) => {
  if (!id) return 'N/A'
  if (id.length <= 16) return id
  return `${id.substring(0, 8)}...${id.substring(id.length - 8)}`
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}


const formatType = (type) => {
  if (!type) return 'N/A'
  return type
}


const showTransactionDetails = (transaction) => {
  selectedTransaction.value = transaction
}


const closeModal = () => {
  selectedTransaction.value = null
}




onMounted(() => {
  fetchTransactions()
})
</script>

<style scoped>
.transactions {
  /* padding: 20px; */
/*  width: 100%;     
  margin: 0 40px; */
/*  max-width: unset; */
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 50px;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 15px;
}

.header h2 {
  margin: 0;
  color: #333;
}

.stats {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-item .label {
  font-weight: bold;
  color: #666;
}

.stat-item .value {
  color: #2196F3;
  font-weight: bold;
}

.refresh-btn {
  /* background-color: #2196F3;
  color: white;
  border: none;
  border-radius: 4px; */
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
}

.refresh-btn:hover:not(:disabled) {
  /* background-color: #1976D2; */
}

.refresh-btn:disabled {
  /* background-color: #ccc; */
  cursor: not-allowed;
}

.loading, .error, .empty {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #f44336;
}

.error button {
  margin-top: 10px;
  padding: 8px 16px;
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.transactions-table {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.08);
  padding: 24px 18px 18px 18px;
  margin-top: 18px;
  overflow-x: auto;
}

.transactions-table .el-table {
  border-radius: 8px;
  overflow: hidden;
}

.transactions-table .el-table__header th {
  font-weight: bold;
  font-size: 20px;
  text-align: center;
  background: #bdc3c9 !important;
  color: #333;
}


.transactions-table .el-table .el-table__header-wrapper .el-table__header th {
  background: #f5f7faf5f7fa !important;
}


.transactions-table .el-table th.el-table__cell {
  background: #f5f7faf5f7fa !important;
}


:deep(.el-table__header th) {
  background: #f5f7fa !important;
}

.transactions-table .el-table__cell {
  font-size: 18px;
  text-align: center;
  vertical-align: middle;
  padding: 18px 8px;
}

.transactions-table .el-table__row {
  transition: background 0.2s;
}
.transactions-table .el-table__row:hover {
  background: #f0f7ff !important;
}


.transactions-table .el-pagination {
  display: flex;
  justify-content: center;
  margin: 24px 0 0 0;
}

.tx-id {
  font-family: monospace;
  font-size: 12px;
}


.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 600px;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.detail-item {
  margin-bottom: 15px;
}

.detail-item label {
  font-weight: bold;
  color: #666;
  display: block;
  margin-bottom: 5px;
}

.detail-item span {
  color: #333;
  word-break: break-all;
}

.payloads {
  margin-top: 10px;
}

.payload-item {
  display: flex;
  margin-bottom: 8px;
  padding: 8px;
  background-color: #f8f9fa;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
}

.payload-index {
  color: #666;
  margin-right: 8px;
  min-width: 20px;
}

.payload-content {
  color: #333;
  word-break: break-all;
}

/* @media (max-width: 768px) {
  .header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .stats {
    flex-direction: column;
    gap: 10px;
  }
  
  .tx-id {
    font-size: 12px;
  }

  .transactions-table .el-table {
    font-size: 15px;
  }
} */
</style> 
