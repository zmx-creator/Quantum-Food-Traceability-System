<template>
	<div class="PaperName" style="display: flex; align-items: center;">
		  <div style="flex: 1; display: flex; align-items: center; justify-content: flex-start; font-size: 20px;">
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Fold /></el-icon>
				<span style="margin-left: 4px;">TRACEABILITY</span>
			</span>
		    <el-divider direction="vertical" />
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Promotion /></el-icon>
				<span style="margin-left: 4px;">Information Track</span>
			</span>
		    <el-divider direction="vertical" border-style="dashed" />
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><KnifeFork /></el-icon>
				<span style="margin-left: 4px;">Complete Information</span>
			</span>
		  </div>
		  <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><Setting /></el-icon>
		  <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><UserFilled /></el-icon>	  
		  		
	</div>
	<el-divider />
	<div class="header" style="display: flex; align-items: center; justify-content: space-between; font-size: 24px; margin-bottom: 20px;">
	  <el-text class="header-title" size="large" style="margin-left: 10px;">
	
			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:10px;height:10px;margin-right: 8px;">	
					  
			<img src="../assets/orangeNode.png" alt="icon"
			  style="vertical-align: middle;width:12px;height:12px;margin-right: 8px;">
			<img src="../assets/greenNode.png" alt="icon"
			  style="vertical-align: middle;width:14px;height:14px;margin-right: 8px;">
			<img src="../assets/blueNode.png" alt="icon"
			  style="vertical-align: middle;width:16px;height:16px;margin-right: 8px;">
		  Trace Back Complete Information
		  
			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:18px;height:18px;margin-left: 8px;">
			<img src="../assets/greenNode.png" alt="icon"
				style="vertical-align: middle;width:14px;height:14px;margin-left: 8px;">
			<img src="../assets/orangeNode.png" alt="icon"
				style="vertical-align: middle;width:12px;height:12px;margin-left: 8px;">	
				  
			<img src="../assets/blueNode.png" alt="icon"
				style="vertical-align: middle;width:10px;height:10px;margin-left: 8px;">		
	
		</el-text>
	
	</div>


  <el-card class="search-card">
    <div class="track-search-bar">
      <el-input v-model="searchId" placeholder="Input the traceability number to search:" style="width: 300px; margin-right: 12px;" clearable />
      <el-button type="primary" @click="searchData">Search</el-button>
      <el-button @click="clearSearch">Clear</el-button>
    </div>
  </el-card>


  <el-card class="track-card" v-if="showProductTable">
    <div class="table-title">
      <el-icon><KnifeFork /></el-icon>
      <span>Product Information</span>
    </div>
    <el-table :data="productTableData" border style="width: 100%; margin-top: 16px;"
      :header-cell-style="{textAlign: 'center'}"
      :cell-style="{textAlign: 'center'}"
    >
      <el-table-column prop="foodID" label="Product ID" />
      <el-table-column prop="foodName" label="Product Name" />
      <el-table-column prop="foodSpec" label="Specification" />
      <el-table-column prop="timestamp" label="Timestamp" />
      <el-table-column prop="foodLOT" label="Product LOT" />
      <el-table-column prop="foodQSID" label="Product QSID" />
      <el-table-column prop="foodMFRSName" label="Manufacturer" />
      <el-table-column prop="foodProPrice" label="Product Price" />
      <el-table-column prop="foodProPlace" label="Production Place" />
    </el-table>
  </el-card>

 
  <el-card class="track-card" v-if="showIngredientTable">
    <div class="table-title">
      <el-icon><Brush /></el-icon>
      <span>Ingredient Information</span>
    </div>
    <el-table :data="ingredientTableData" border style="width: 100%; margin-top: 16px;">
      <el-table-column prop="foodID" label="Product ID" width="120" />
      <el-table-column prop="timestamp" label="Timestamp" width="160" />
      <el-table-column prop="ingredients" label="Ingredient List">
        <template #default="scope">
          <span v-if="Array.isArray(scope.row.ingredients)">
            <template v-if="scope.row.ingredients.length && typeof scope.row.ingredients[0] === 'object' && scope.row.ingredients[0].IngName">
              {{ scope.row.ingredients.map(i => mapIngredientName(i.IngName, scope.row.foodID)).join('，') }}
            </template>
            <template v-else>
              {{ scope.row.ingredients.join('，') }}
            </template>
          </span>
          <span v-else>{{ scope.row.ingredients }}</span>
        </template>
      </el-table-column>
    </el-table>
  </el-card>


  <el-card class="track-card" v-if="showLogisticsTable">
    <div class="table-title">
      <el-icon><Box /></el-icon>
      <span>Logistics Information</span>
    </div>
    <el-table :data="logisticsTableData" border style="width: 100%; margin-top: 16px;"
      :header-cell-style="{textAlign: 'center'}"
      :cell-style="{textAlign: 'center'}"
    >
      <el-table-column prop="foodID" label="Product ID" />
      <el-table-column prop="timestamp" label="Timestamp" />
      <el-table-column prop="logDeparturePl" label="Departure Place" />
      <el-table-column prop="logDest" label="Destination" />
      <el-table-column prop="logToSeller" label="Seller" />
      <el-table-column prop="logStorageTm" label="Storage Time" />
      <el-table-column prop="logMOT" label="Mode of Transport" />
      <el-table-column prop="logCopName" label="Logistics Company" />
    </el-table>
  </el-card>


  <el-card class="no-data-card" v-if="!showProductTable && !showIngredientTable && !showLogisticsTable && searchId">
    <el-empty description="No data found for the specified Product ID" />
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAllFoodProduction, getAllFoodIngredients, getAllFoodLogistics } from '../api/fabric'
import { ElMessage } from 'element-plus'

const searchId = ref('')
const productTableData = ref([])
const ingredientTableData = ref([])
const logisticsTableData = ref([])


const allProductData = ref([])
const allIngredientData = ref([])
const allLogisticsData = ref([])


const showProductTable = ref(false)
const showIngredientTable = ref(false)
const showLogisticsTable = ref(false)


function mapIngredientName(ingName, foodID) {

  if (ingName && ingName.startsWith('ING_')) {

    const ingredientName = getIngredientNameFromHistory(ingName, foodID)
    if (ingredientName) {
      return ingredientName
    }
    

    return 'ingredient:' + ingName.split('_').pop()
  }
  

  return ingName
}


function getIngredientNameFromHistory(ingID, foodID) {

  

  const ingredientHistory = JSON.parse(localStorage.getItem('ingredientHistory') || '{}')
  
  if (ingredientHistory[foodID] && ingredientHistory[foodID][ingID]) {
    return ingredientHistory[foodID][ingID]
  }
  
  return null
}


function saveIngredientHistory(foodID, ingredientPairs) {
  const history = JSON.parse(localStorage.getItem('ingredientHistory') || '{}')
  
  if (!history[foodID]) {
    history[foodID] = {}
  }
  

  for (let i = 0; i < ingredientPairs.length; i += 2) {
    const ingID = ingredientPairs[i]
    const ingName = ingredientPairs[i + 1]
    history[foodID][ingID] = ingName
  }
  
  localStorage.setItem('ingredientHistory', JSON.stringify(history))
}


onMounted(() => {
  fetchAllData()
})


function fetchAllData() {

  getAllFoodProduction().then(res => {
    if (res.data && res.data.success) {
      const foods = res.data.data || []
      allProductData.value = foods.map(food => ({
        foodID: food.foodID || '',
        foodName: food.foodName || '',
        foodSpec: food.foodSpec || '',
        timestamp: food.timestamp ? new Date(food.timestamp).toLocaleString('zh-CN') : '',
        foodLOT: food.foodLOT || '',
        foodQSID: food.foodQSID || '',
        foodMFRSName: food.foodMFRSName || '',
        foodProPrice: food.foodProPrice || '',
        foodProPlace: food.foodProPlace || ''
      }))
      console.log('get data:', allProductData.value)
    }
  }).catch(err => {
    console.error('get data:', err)
  })


  getAllFoodIngredients().then(res => {
    if (res.data && res.data.success) {
      const ingredients = res.data.data || []
      allIngredientData.value = ingredients.map(item => ({
        foodID: item.foodID || '',
        timestamp: item.timestamp ? new Date(item.timestamp).toLocaleString('zh-CN') : '',
        ingredients: item.ingredients || []
      }))
      console.log('get ingredient:', allIngredientData.value)
    }
  }).catch(err => {
    console.error('get ingredient:', err)
  })


  getAllFoodLogistics().then(res => {
    if (res.data && res.data.success) {
      const logistics = res.data.data || []
      allLogisticsData.value = logistics.map(item => ({
        foodID: item.foodID || '',
        timestamp: item.timestamp ? new Date(item.timestamp).toLocaleString('zh-CN') : '',
        logDeparturePl: item.logDeparturePl || '',
        logDest: item.logDest || '',
        logToSeller: item.logToSeller || '',
        logStorageTm: item.logStorageTm || '',
        logMOT: item.logMOT || '',
        logCopName: item.logCopName || ''
      }))
      console.log('get log:', allLogisticsData.value)
    }
  }).catch(err => {
    console.error('get log:', err)
  })
}


function searchData() {
  if (!searchId.value.trim()) {
    ElMessage.warning('Please enter a Product ID to search')
    return
  }
  
  const searchTerm = searchId.value.toLowerCase()
  

  const filteredProductData = allProductData.value.filter(food => 
    food.foodID && food.foodID.toLowerCase().includes(searchTerm)
  )
  productTableData.value = filteredProductData
  showProductTable.value = filteredProductData.length > 0


  const filteredIngredientData = allIngredientData.value.filter(item => 
    item.foodID && item.foodID.toLowerCase().includes(searchTerm)
  )
  ingredientTableData.value = filteredIngredientData
  showIngredientTable.value = filteredIngredientData.length > 0


  const filteredLogisticsData = allLogisticsData.value.filter(item => 
    item.foodID && item.foodID.toLowerCase().includes(searchTerm)
  )
  logisticsTableData.value = filteredLogisticsData
  showLogisticsTable.value = filteredLogisticsData.length > 0

  console.log('result:', {
    product: filteredProductData,
    ingredient: filteredIngredientData,
    logistics: filteredLogisticsData
  })
  

  console.log('log status:', {
    allLogisticsData: allLogisticsData.value,
    filteredLogisticsData: filteredLogisticsData,
    showLogisticsTable: showLogisticsTable.value,
    searchTerm: searchTerm
  })
}


function clearSearch() {
  searchId.value = ''
  showProductTable.value = false
  showIngredientTable.value = false
  showLogisticsTable.value = false
  productTableData.value = []
  ingredientTableData.value = []
  logisticsTableData.value = []
}


defineExpose({
  saveIngredientHistory
})
</script>

<style scoped>
.search-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  border-radius: 8px;
  padding: 16px;
}

.track-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  border-radius: 8px;
  padding: 0 0 8px 0;
}

.track-search-bar {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.table-title {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  color: #409EFF;
  margin: 16px 0 8px 16px;
}

.table-title .el-icon {
  margin-right: 8px;
  font-size: 20px;
}

.no-data-card {
  margin-top: 20px;
  text-align: center;
  padding: 40px;
}

::v-deep .el-table th,
::v-deep .el-table td {
  text-align: center !important;
}
</style> 
