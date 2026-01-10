<template>
	<div class="PaperName" style="display: flex; align-items: center;">
		  <div style="flex: 1; display: flex; align-items: center; justify-content: flex-start; font-size: 20px;">
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Fold /></el-icon>
				<span style="margin-left: 4px;">TRACEABILITY</span>
			</span>
		    <el-divider direction="vertical" />
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Checked /></el-icon>
				<span style="margin-left: 4px;">Submit</span>
			</span>
		    <el-divider direction="vertical" border-style="dashed" />
		    <span style="display: inline-flex; align-items: center;">
				<el-icon><Grape /></el-icon>
				<span style="margin-left: 4px;">Complete Information</span>
			</span>
		  </div>
		  <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><Setting /></el-icon>
		  <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><UserFilled /></el-icon>	  
		  		
	</div>
		<el-divider style="margin: 8px 0;" />
	<div class="submit">
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
			  Submit Complete Information
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
	  
		<el-card class="form-card">
			<div class="table-title">
			  <!-- <el-icon><KnifeFork /></el-icon> -->
			  <span>Submit Product Information</span>
			</div>
		
		  <el-form :model="completeInfo" label-width="150px" class="form-content" label-position="left" ref="productForm">
			<el-row :gutter="150">
			  <!-- Product Information -->
			  <el-col :span="12" v-if="showProductID"><el-form-item label="Traceability Number">
				<el-input v-model="completeInfo.foodID" style="width: 800px;" placeholder="autoGeneration" readonly /></el-form-item></el-col>
			
			  <el-col :span="12"><el-form-item
			   prop="foodName" 
			   label="Product Name"       
			  :rules="[
				{
				  required: true,
				  message: 'Please input product name',
				  trigger: 'blur',
				}]">
				<el-input v-model="completeInfo.foodName" style="width: 800px;" />
			  </el-form-item></el-col>
			  
			  <el-col :span="12"><el-form-item label="Specification"><el-input v-model="completeInfo.foodSpec" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Manufacture Date"><el-input v-model="completeInfo.foodMFGDate" style="width: 800px;" placeholder="xxxx-xx-xx"/></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Expiration Date"><el-input v-model="completeInfo.foodEXPDate" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Product LOT"><el-input v-model="completeInfo.foodLOT" style="width: 800px;" placeholder="LOT_"/></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Product QSID"><el-input v-model="completeInfo.foodQSID" style="width: 800px;" placeholder="QS_"/></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Manufacturer"><el-input v-model="completeInfo.foodMFRSName" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Product Price"><el-input v-model="completeInfo.foodProPrice" style="width: 800px;" placeholder="CNY"/></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Production Place"><el-input v-model="completeInfo.foodProPlace" style="width: 800px;" /></el-form-item></el-col>
		</el-row>
<!-- 			<el-form-item class="center-btn">
			  <el-button type="primary" plain @click="submitCompleteInfo" style="margin-left: 480px;">Submit</el-button>
			</el-form-item> -->
		  </el-form>
		</el-card>
		
		<el-card class="form-card">
			<div class="table-title">
			  <!-- <el-icon><Brush /></el-icon> -->
			  <span>Submit Ingredient Information</span>
			</div>
		 
		  <el-form :model="completeInfo" label-width="150px" class="form-content" label-position="left" ref="ingredientForm">
			<el-row :gutter="150">
			  <!-- Product Information -->

			  <!-- Ingredient Information -->
			  <el-col :span="12">
				<el-form-item 
					prop="ingredientsStr"
					label="Ingredient List"
					:rules="[
					{
						required: true,
						message: 'Please input ingredients',
						trigger: 'blur',
					}]">
					<el-input v-model="completeInfo.ingredientsStr" style="width: 800px;" placeholder="XXX, XXX" />
				</el-form-item>
			  </el-col>
			  
		</el-row>
<!-- 			<el-form-item class="center-btn">
			  <el-button type="primary" plain @click="submitCompleteInfo" style="margin-left: 480px;">Submit</el-button>
			</el-form-item> -->
		  </el-form>
		</el-card>
		
		<el-card class="form-card">
			<div class="table-title">
			  <!-- <el-icon><Box /></el-icon> -->
			  <span>Submit Logistics Information</span>
			</div>
		 
		  <el-form :model="completeInfo" label-width="150px" class="form-content" label-position="left" ref="logisticsForm">
			<el-row :gutter="150">

			  <!-- Ingredient Information -->
<!-- 			  <el-col :span="12"><el-form-item label="Ingredient List"><el-input v-model="completeInfo.ingredientsStr" style="width: 600px;" placeholder="XXX, XXX" /></el-form-item></el-col>
			  
	 -->		  <!-- Logistics Information -->
			  <el-col :span="12">
				  <el-form-item 
				  prop="logDepartureTm"
				  label="Departure Time"
					:rules="[
					{
						required: true,
						message: 'Please input Departure Time',
						trigger: 'blur',
					}]"
				  >
					  <el-input v-model="completeInfo.logDepartureTm" style="width: 800px;" placeholder="xxxx-xx-xx xx:xx AM/PM"/></el-form-item>
			  </el-col>
			  <el-col :span="12"><el-form-item label="Arrival Time"><el-input v-model="completeInfo.logArrivalTm" style="width: 800px;" placeholder="xxxx-xx-xx xx:xx AM/PM"/></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Mission"><el-input v-model="completeInfo.logMission" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Departure Place"><el-input v-model="completeInfo.logDeparturePl" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Destination"><el-input v-model="completeInfo.logDest" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Seller"><el-input v-model="completeInfo.logToSeller" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Storage Time"><el-input v-model="completeInfo.logStorageTm" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Mode of Transport"><el-input v-model="completeInfo.logMOT" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Logistics Company"><el-input v-model="completeInfo.logCopName" style="width: 800px;" /></el-form-item></el-col>
			  <el-col :span="12"><el-form-item label="Cost"><el-input v-model="completeInfo.logCost" style="width: 800px;" placeholder="CNY" /></el-form-item></el-col>
			</el-row>
			<el-form-item class="center-btn">
			  <el-button type="primary" plain @click="submitCompleteInfo" style="margin-left: 580px;">Submit</el-button>
			</el-form-item>
		  </el-form>
		</el-card>
    <div v-if="msg" class="msg">{{ msg }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { addProInfo, addIngInfo, addLogInfo } from '../api/fabric'
import { ElMessage } from 'element-plus'

const completeInfo = ref({
  // Product Information
  foodID: '', foodName: '', foodSpec: '', foodMFGDate: '', foodEXPDate: '', foodLOT: '', foodQSID: '', foodMFRSName: '', foodProPrice: '', foodProPlace: '',
  // Ingredient Information
  ingredientsStr: '',
  // Logistics Information
  logDepartureTm: '', logArrivalTm: '', logMission: '', logDeparturePl: '', logDest: '', logToSeller: '', logStorageTm: '', logMOT: '', logCopName: '', logCost: ''
})
const msg = ref('')
const showProductID = ref(false)

const productForm = ref(null)
const ingredientForm = ref(null)
const logisticsForm = ref(null)

async function submitCompleteInfo() {

  let valid = true
  await productForm.value.validate((valid1) => { if (!valid1) valid = false })
  await ingredientForm.value.validate((valid2) => { if (!valid2) valid = false })
  await logisticsForm.value.validate((valid3) => { if (!valid3) valid = false })
  if (!valid) {
    ElMessage.error('Please fill out the complete form！')
    return
  }
  try {

    if (!completeInfo.value.foodID) {
      const now = new Date();
      const pad = n => n.toString().padStart(2, '0');
      const timeStr = `${now.getFullYear()}${pad(now.getMonth()+1)}${pad(now.getDate())}${pad(now.getHours())}${pad(now.getMinutes())}${pad(now.getSeconds())}`;
      completeInfo.value.foodID = `PRODUCTTRACK${timeStr}`;
    }
    showProductID.value = true
    // Submit Product Information
    const proInfo = {
      foodID: completeInfo.value.foodID,
      foodName: completeInfo.value.foodName,
      foodSpec: completeInfo.value.foodSpec,
      foodMFGDate: completeInfo.value.foodMFGDate,
      foodEXPDate: completeInfo.value.foodEXPDate,
      foodLOT: completeInfo.value.foodLOT,
      foodQSID: completeInfo.value.foodQSID,
      foodMFRSName: completeInfo.value.foodMFRSName,
      foodProPrice: completeInfo.value.foodProPrice,
      foodProPlace: completeInfo.value.foodProPlace
    }
    await addProInfo(proInfo)

    // Submit Ingredient Information
    if (completeInfo.value.ingredientsStr.trim()) {
      const ings = completeInfo.value.ingredientsStr.split(',').map(i => i.trim()).filter(i => i)
      const ingredientPairs = []
      ings.forEach((ingredient, index) => {
        ingredientPairs.push(`ING_${Date.now()}_${index}`) 
        ingredientPairs.push(ingredient) 
      })
      

      saveIngredientHistory(completeInfo.value.foodID, ingredientPairs)
      
      await addIngInfo(completeInfo.value.foodID, ingredientPairs)
    }

    // Submit Logistics Information
    const logInfo = {
      foodID: completeInfo.value.foodID,
      logDepartureTm: completeInfo.value.logDepartureTm,
      logArrivalTm: completeInfo.value.logArrivalTm,
      logMission: completeInfo.value.logMission,
      logDeparturePl: completeInfo.value.logDeparturePl,
      logDest: completeInfo.value.logDest,
      logToSeller: completeInfo.value.logToSeller,
      logStorageTm: completeInfo.value.logStorageTm,
      logMOT: completeInfo.value.logMOT,
      logCopName: completeInfo.value.logCopName,
      logCost: completeInfo.value.logCost
    }
    await addLogInfo(logInfo)

    ElMessage.success('Complete information submitted successfully')
  } catch (e) {
    ElMessage.error('Failed to submit information')
  }
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
  console.log('save history:', history)
}
</script>

<style scoped>
.submit {
  padding: 20px;
  width: 100%;
  margin: 0 auto;
}
.form-card {
 margin: 20px 40px;
  /* width: 100%; */
  border-radius: 8px;
  box-shadow: 0 2px 8px #e0e2e4;
  padding: 20px 30px 10px 30px;
}

.table-title {
  display: flex;
  align-items: left;
  margin-bottom: 16px;
  font-size: 18px;
  font-weight: bold;
  color: #409EFF;
  /* margin: 0px 0 0px 16px; */
}

.form-title {
  font-size: 1.2em;
  font-weight: bold;
  margin-bottom: 18px;
  color: #409EFF;
}
.form-content {
  max-width: 900px;
}


.form-content .el-form-item__label {
  text-align: left !important;
  width: 100%;
  justify-content: flex-start !important;
}

.center-btn {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
.el-input {
  width: 800px !important;
}
.msg { margin-top: 20px; color: green; }
</style> 
