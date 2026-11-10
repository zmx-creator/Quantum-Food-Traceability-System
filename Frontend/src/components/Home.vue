<template>
	
	<div class='home'>
		<el-container>
			<el-container class="body">
			
				<el-aside width="200px" style="height: 100vh; padding: 0;">
					<LeftMenu></LeftMenu>						
				</el-aside>
				
				<el-main style="min-height: 100vh;">
					<router-view ></router-view>
				</el-main>
			</el-container>
		</el-container>
		

	</div>
</template>

<script>

	import LeftMenu from './LeftMenu.vue'
	
	export default {
	  components: {
	    LeftMenu
	  }
	}

</script>

<style>
	.home {
		background-color: #f4f4f5;
		background-image: 
					url("../assets/apple_line2.png");
		background-size: 40%;
		background-repeat: no-repeat;
		background-position: right bottom ;
	}
</style>
