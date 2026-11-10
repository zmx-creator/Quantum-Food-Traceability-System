<template>
  <div class="chart-section">
    <h3>MSP Distribution</h3>
    <div class="chart-panel">
      <div id="pieChart" style="width:100%;height:200px;"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { getPeers, getOrderers } from '../api/fabric'

const mspData = ref([])

onMounted(async () => {
  try {
    const [peerRes, ordererRes] = await Promise.all([getPeers(), getOrderers()])
    let peerList = Array.isArray(peerRes.data.data) ? peerRes.data.data : []
    let ordererList = Array.isArray(ordererRes.data.data) ? ordererRes.data.data : []
    

    const mspCount = {}
    peerList.forEach(peer => {
      const msp = peer.name.split('.')[0] || 'Unknown'
      mspCount[msp] = (mspCount[msp] || 0) + 1
    })
    ordererList.forEach(orderer => {
      const msp = orderer.name.split('.')[0] || 'Unknown'
      mspCount[msp] = (mspCount[msp] || 0) + 1
    })
    
    mspData.value = Object.entries(mspCount).map(([name, value]) => ({ name, value }))
  } catch (e) {
    console.error('Failed to get MSP information', e)

    mspData.value = [
      { value: 2, name: 'OrdererMSP' },
      { value: 6, name: 'Org1MSP' }
    ]
  }


  setTimeout(() => {
    const chartDom = document.getElementById('pieChart')
    if (chartDom) {
      const pieChart = echarts.init(chartDom)
      pieChart.setOption({
        title: {
          // text: 'MSP nodes distribution',
          left: 'center',
          textStyle: {
            fontSize: 14,
            color: '#333'
          }
        },
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: { 
          bottom: 0,
          textStyle: {
            fontSize: 12
          }
        },
        series: [{
          name: 'MSP distribution',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['50%', '45%'],
          data: mspData.value,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          label: {
            show: true,
            formatter: '{b}: {c}'
          },
          labelLine: {
            show: true
          }
        }]
      })
    }
  }, 200)
})
</script>

<style scoped>
.chart-section {
  margin-bottom: 20px;
}

.chart-section h3 {
  margin-bottom: 10px;
  color: #333;
}

.chart-panel {
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
</style> 
