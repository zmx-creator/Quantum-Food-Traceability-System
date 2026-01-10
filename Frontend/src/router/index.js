import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/Home.vue'
import Header from '../components/Header.vue'
import LeftMenu from '../components/LeftMenu.vue'
import Home_board from '../components/Home_board.vue'
import Home_Peer from '../components/Home_Peer.vue'
import Home_Blocks from '../components/Home_Blocks.vue'
import BlockChart from '../components/BlockChart.vue'
import PieChart from '../components/PieChart.vue'
import Info from '../components/Info.vue'
import Transactions from '../components/Transactions.vue'
import submitPro from '../components/submitPro.vue'
import FoodTrack from '../components/FoodTrack.vue'


const router = createRouter({
  history: createWebHistory(),
  routes: [
	
    { path: '/', redirect: '/Home' },
    { 	path: '/Home', 
		name: 'Home', 
		component: Home ,
		redirect:'/Home/Info',
		children:[
			{path:'Info',component:Info},
			{path:'submitPro',component:submitPro},
			{path:'Transactions',component:Transactions},
			{path:'FoodTrack',component:FoodTrack},
		]
	
	},
    { path: '/FoodTrack', redirect: '/Home/FoodTrack' },
    { path: '/submitPro', redirect: '/Home/submitPro' },
    // { path: '/login', name: 'Login', component: Login }
	
	//
  ]
})

export default router
