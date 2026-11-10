import axios from 'axios'

// The demo talks to one Spring Boot application, so its address is kept in a single
// place: point the frontend at another host by setting VUE_APP_API_BASE in .env.
export const API_BASE = process.env.VUE_APP_API_BASE || 'http://localhost:8080'

// Shared instance; every path below is relative to the Fabric REST prefix.
const http = axios.create({
  baseURL: `${API_BASE}/api/fabric`
})

export { http }
export function getBlockchainDetails() {
  return http.get('/blockchain/details')
}


export function getAllTransactions() {
  return http.get('/blockchain/transactions')
}


export function getPeers() {
  return http.get('/peers')
}


export function getOrderers() {
  return http.get('/orderers')
}

export function getNodeLogs(nodeName) {
  return http.get(`/nodes/${nodeName}/logs`)
}


export function getBlockchainInfo() {
  return http.get('/info')
}



export function addProInfo(params) {
  return http.post('/food/production', params)
}


export function addIngInfo(foodID, ingredients) {
  const data = {
    foodID: foodID,
    ingredientPairs: ingredients
  }
  return http.post('/food/ingredients', data)
}


export function addLogInfo(params) {
  return http.post('/food/logistics', params)
}


export function getFoodInfo(foodID) {
  return http.get(`/food/${foodID}`)
}


export function getProInfo(foodID) {
  return http.get(`/food/${foodID}/production`)
}


export function getIngInfo(foodID) {
  return http.get(`/food/${foodID}/ingredients`)
}


export function getLogInfo(foodID) {
  return http.get(`/food/${foodID}/logistics`)
}


export function getLogInfoList(foodID) {
  return http.get(`/food/${foodID}/logistics/list`)
}


export function getAllFoodProduction() {
  return http.get('/food/production/all')
}


export function getAllFoodIngredients() {
  return http.get('/food/ingredients/all')
}


export function getAllFoodLogistics() {
  return http.get('/food/logistics/all')
} 


// ---------------------------------------------------------------------------
// Byzantine demonstration (see Backend ByzantineDemoController)
// ---------------------------------------------------------------------------

export function getDemoNodes() {
  return http.get('/demo/nodes')
}


export function setDemoFault(nodes, mode) {
  return http.post('/demo/fault', { nodes, mode })
}


export function runDemoRound(waitMs = 0) {
  return http.post('/demo/run', null, { params: { waitMs } })
}


export function getDemoConsensus(since = 0, limit = 60) {
  return http.get('/demo/consensus', { params: { since, limit } })
}


export function getDemoTrace(node, since = 0, limit = 200) {
  return http.get('/demo/trace', { params: { node, since, limit } })
}
