<template>
  <div class="left-menu-wrapper fixed-left-menu">
    <el-menu
      active-text-color="#ffd04b"
	  style="margin-top: 40px;"
      background-color="#2d3a4b"
      class="el-menu-vertical-demo custom-menu"
      default-active="/Home/Info"
      text-color="#fff"
      router
    >
      <el-menu-item index="title" style="color: #c2cbd2;" disabled>
        <el-icon><Fold /></el-icon>
        <span>TRACEABILITY</span>
      </el-menu-item>
	  
      <el-menu-item index="/Home/Info">
        <el-icon><House /></el-icon>
        <span>Blockchain</span>
      </el-menu-item>
	  
      <el-sub-menu index="/Home/Submit">	  
        <template #title>
          <el-icon><Pointer /></el-icon>
          <span>Submit</span>
        </template>
        <el-menu-item index="/Home/submitPro"><el-icon><Grape /></el-icon>Product</el-menu-item>
      </el-sub-menu>
	  
      <el-menu-item index="/Home/Transactions">
        <el-icon><Memo /></el-icon>
        <span>Transactions</span>
      </el-menu-item>
	  
	  <el-sub-menu index="/Home/Track">
	    <template #title>
	      <el-icon><Connection /></el-icon>
	      <span>Traceability</span>
	    </template>
	    <el-menu-item index="/Home/FoodTrack"><el-icon><KnifeFork /></el-icon>Product</el-menu-item>
	  </el-sub-menu>
    </el-menu>
  </div>
</template>

<script setup>
</script>

<style scoped>
.left-menu-wrapper {
  height: 100vh;
  background: #2d3a4b;
  box-shadow: 2px 0 6px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}
.fixed-left-menu {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1000;
  width: 200px;
}
.custom-menu {
  border-right: none;
  height: 100%;
  width: 100%;
  background: #2d3a4b;
}
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 100%;
  min-height: 100vh;
}
.el-menu-item, .el-sub-menu__title {
  font-size: 16px !important;
  font-weight: 600;
  height: 56px;
  display: flex;
  align-items: center;
  border-radius: 4px;
  margin: 4px 0;
  transition: background 0.2s;
}
.el-menu-item .el-icon, .el-sub-menu__title .el-icon {
  margin-right: 12px;
}
.el-menu-item.is-active, .el-sub-menu__title.is-active {
  background-color: #3c4e65 !important;
  color: #ffc873 !important;
}
.el-menu-item-group__title {
  font-weight: bold;
  color: #3c4e65 ;
  font-size: 15px;
  padding-left: 12px;
}
.el-menu .el-menu-item-group .el-menu-item {
  font-size: 15px !important;
  font-weight: 400;
  padding-left: 36px !important;
  height: 44px;
}
.el-sub-menu__title span {
  font-size: 16px !important;
  font-weight: 600;
}
.el-sub-menu .el-menu-item {
  font-size: 14px !important;
  font-weight: 400;
  padding-left: 50px !important;
  height: 44px;
}
</style>
