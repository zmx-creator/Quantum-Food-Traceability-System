<template>
  <div class="chart-section">
    <h3>Block Transaction Chart</h3>
    <div class="chart-panel">
      <div id="blockChart" style="width:100%;height:200px;"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { getBlockchainDetails } from '../api/fabric'

const blocks = ref([])

onMounted(async () => {
  try {
    const res = await getBlockchainDetails()
    if (res.data.success && res.data.data) {
      blocks.value = (res.data.data.blocks || []).map(b => ({
        number: b.blockNumber,
        tx: b.transactionCount
      }))
    }
  } catch (e) {
    console.error('Failed to fetch block information:', e)
    blocks.value = []
  }

  // Echarts chart initialization
  setTimeout(() => {
    const chartDom = document.getElementById('blockChart')
    if (chartDom) {
      const blockChart = echarts.init(chartDom)
      blockChart.setOption({
        title: {
          text: 'Block Transaction Statistics',
          left: 'center',
          textStyle: {
            fontSize: 14,
            color: '#333'
          }
        },
        tooltip: {
          trigger: 'axis',
          formatter: function(params) {
            return `Block ${params[0].name}<br/>Transactions: ${params[0].value}`
          }
        },
        xAxis: { 
          type: 'category', 
          data: blocks.value.map(b => 'Block ' + b.number),
          axisLabel: {
            rotate: 45,
            fontSize: 10
          }
        },
        yAxis: { 
          type: 'value',
          name: 'Transaction Count'
        },
        series: [{ 
          data: blocks.value.map(b => b.tx), 
          type: 'line',
          smooth: true,
          lineStyle: {
            color: '#2196F3',
            width: 3
          },
          itemStyle: {
            color: '#2196F3'
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgba(33, 150, 243, 0.3)'
              }, {
                offset: 1, color: 'rgba(33, 150, 243, 0.1)'
              }]
            }
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
