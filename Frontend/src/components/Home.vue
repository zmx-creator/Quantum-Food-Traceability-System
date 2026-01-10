<template>
	
	<div class='home'>
		<el-container>
<!-- 			<el-header>
				<Header></Header>
			</el-header>
			 -->
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

	// import Header from './Header.vue'
	import LeftMenu from './LeftMenu.vue'
	
	export default {
	  components: {
	    // Header,
	    LeftMenu
	  }
	}

</script>

<style>
	.home {
/* 		height: 100%;
		width: 100%; */
		background-color: #f4f4f5;
		background-image: 
					url("../assets/apple_line2.png");
		/* background-size: 50%; */
		background-size: 40%;
		background-repeat: no-repeat;
		background-position: right bottom ;
		/* background-position: 1400px 550px; */
	}
</style>
