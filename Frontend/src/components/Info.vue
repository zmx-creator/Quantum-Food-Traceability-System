<template>
	<div class="body">
			<div class="PaperName" style="display: flex; align-items: center;">
				  <div style="flex: 1; display: flex; align-items: center; justify-content: flex-start; font-size: 20px;">
					<span style="display: inline-flex; align-items: center;">
						<el-icon><Fold /></el-icon>
						<span style="margin-left: 4px;">TRACEABILITY</span>
					</span>
					<el-divider direction="vertical" />
					<span style="display: inline-flex; align-items: center;">
						<el-icon><View /></el-icon>
						<span style="margin-left: 4px;">Blockchain Information</span>
					</span>
					<el-divider direction="vertical" border-style="dashed" />

				  </div>
				   <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><Setting /></el-icon>
				   <el-icon style="margin-left: auto;font-size:22px ;margin-right: 40px;"><UserFilled /></el-icon>	  
				  
			</div>

			<el-divider style="margin: 8px 0;" />
		  <div class="dashboard">
			<div class="top-stats">
			  <Home_board />
			</div>
			<div class="main-content">
			  <div class="left">
				<Home_Peer />
				<Home_Blocks />
			  </div>
			  <div class="right">
				<PieChart />
				<ChannelInfo />
				<ChaincodeInfo />
				<!-- <Home_Blocks /> -->
			  </div>
			</div>
		  </div>		
	</div>

</template>

<script>
import Home_board from './Home_board.vue'
import Home_Peer from './Home_Peer.vue'
import Home_Blocks from './Home_Blocks.vue'
import PieChart from './PieChart.vue'
import ChannelInfo from  './ChannelInfo.vue'
import ChaincodeInfo from  './ChaincodeInfo.vue'
import { HomeFilled } from '@element-plus/icons-vue'

export default {
  components: {
    Home_board,
    Home_Peer,
    Home_Blocks,
    PieChart,
    ChaincodeInfo,
    ChannelInfo,
    HomeFilled
  }
}
</script>

<style scoped>
/* .PaperName{
	background-color: #617a9e;
} */
.body{
	background-color: #f4f4f5;
	height: 100%;
}
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 10px;
  
}
.top-stats {

}
.main-content {
  display: flex;
  gap: 20px;
}
.left {
  flex: 1;
  gap: 30px;
}
.right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
}
</style>
