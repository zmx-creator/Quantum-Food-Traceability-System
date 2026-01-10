<template>
  <div class="peer-section">
    <h3 >Peer&Orderer Nodes</h3>
    <table class="peer-table">
      <thead>
        <tr>
          <th>Node Name</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="peer in peers" :key="peer" @click="showNodeDetails(peer)" class="clickable-row">
          <td>{{ peer }}</td>
        </tr>
      </tbody>
    </table>


    <div v-if="selectedNode" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Node Details - {{ selectedNode.nodeName }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="detail-item">
            <label>Node Name:</label>
            <span>{{ selectedNode.nodeName }}</span>
          </div>
          <div class="detail-item">
            <label>Node Type:</label>
            <span>{{ selectedNode.nodeType }}</span>
          </div>
          <div class="detail-item">
            <label>Node Address:</label>
            <span>{{ selectedNode.nodeUrl || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <label>Status:</label>
            <span>{{ selectedNode.statusInfo?.status || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <label>Last Active Time:</label>
            <span>{{ formatTimestamp(selectedNode.statusInfo?.lastSeen) }}</span>
          </div>
          <div class="detail-item">
            <label>Log Count:</label>
            <span>{{ selectedNode.logCount || 0 }}</span>
          </div>
          <div v-if="selectedNode.logs && selectedNode.logs.length > 0" class="detail-item">
            <label>Node Logs ({{ selectedNode.logCount }} lines):</label>
            <div class="logs-container">
              <div v-for="(log, index) in selectedNode.logs" :key="index" class="log-item">
                <div class="log-message">{{ log }}</div>
              </div>
            </div>
          </div>
          <div v-else class="detail-item">
              <label>Node Logs:</label>
            <div class="logs-container">
              <div class="log-item">
                <div class="log-message">No logs output</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPeers, getOrderers, getNodeLogs } from '../api/fabric'

const peers = ref([])
const selectedNode = ref(null)

onMounted(async () => {
  try {
    const [peerRes, ordererRes] = await Promise.all([getPeers(), getOrderers()])
    let peerList = Array.isArray(peerRes.data.data) ? peerRes.data.data : []
    let ordererList = Array.isArray(ordererRes.data.data) ? ordererRes.data.data : []
    peers.value = [...peerList.map(p => p.name), ...ordererList.map(o => o.name)]
  } catch (e) {
    console.error('Failed to get node information:', e)
    peers.value = []
  }
})


const showNodeDetails = async (nodeName) => {
  try {
    const response = await getNodeLogs(nodeName)
    if (response.data.success) {
      selectedNode.value = response.data.data
    } else {
      console.error('Failed to get node logs:', response.data.message)
    }
  } catch (error) {
    console.error('Failed to get node logs:', error)
  }
}


const closeModal = () => {
  selectedNode.value = null
}


const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A'
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN')
  } catch (e) {
    return timestamp
  }
}
</script>

<style scoped>
.peer-section {
  margin-bottom: 20px;
}

.peer-section h3 {
  margin-bottom: 10px;
  color: #333;
}

.peer-table { 
  width: 100%; 
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.peer-table th,
.peer-table td {
  padding: 12px 8px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.peer-table th {
  background-color: #e6e6e6;
  font-weight: bold;
  color: #333;
}

.peer-table tr:hover {
  background-color: #f8f9fa;
}

.clickable-row {
  cursor: pointer;
  transition: background-color 0.2s;
}

.clickable-row:hover {
  background-color: #f0f7ff !important;
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
  max-width: 700px;
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

.logs-container {
  margin-top: 10px;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 10px;
}

.log-item {
  margin-bottom: 12px;
  padding: 8px;
  background-color: #f8f9fa;
  border-radius: 4px;
  border-left: 3px solid #2196F3;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.log-timestamp {
  font-size: 12px;
  color: #666;
}

.log-level {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 3px;
  font-weight: bold;
  text-transform: uppercase;
}

.log-level-info {
  background-color: #e3f2fd;
  color: #1976d2;
}

.log-level-debug {
  background-color: #f3e5f5;
  color: #7b1fa2;
}

.log-level-warn {
  background-color: #fff3e0;
  color: #f57c00;
}

.log-level-error {
  background-color: #ffebee;
  color: #d32f2f;
}

.log-message {
  font-size: 13px;
  color: #333;
  line-height: 1.4;
}
</style>
