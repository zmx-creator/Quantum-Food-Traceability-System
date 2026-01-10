import axios from 'axios'


export function getBlockchainDetails() {
  return axios.get('http://localhost:8080/api/fabric/blockchain/details')
}


export function getAllTransactions() {
  return axios.get('http://localhost:8080/api/fabric/blockchain/transactions')
}


export function getPeers() {
  return axios.get('http://localhost:8080/api/fabric/peers')
}


export function getOrderers() {
  return axios.get('http://localhost:8080/api/fabric/orderers')
}

export function getNodeLogs(nodeName) {
  return axios.get(`http://localhost:8080/api/fabric/nodes/${nodeName}/logs`)
}


export function getBlockchainInfo() {
  return axios.get('http://localhost:8080/api/fabric/info')
}



export function addProInfo(params) {
  return axios.post('http://localhost:8080/api/fabric/food/production', params)
}


export function addIngInfo(foodID, ingredients) {
  const data = {
    foodID: foodID,
    ingredientPairs: ingredients
  }
  return axios.post('http://localhost:8080/api/fabric/food/ingredients', data)
}


export function addLogInfo(params) {
  return axios.post('http://localhost:8080/api/fabric/food/logistics', params)
}


export function getFoodInfo(foodID) {
  return axios.get(`http://localhost:8080/api/fabric/food/${foodID}`)
}


export function getProInfo(foodID) {
  return axios.get(`http://localhost:8080/api/fabric/food/${foodID}/production`)
}


export function getIngInfo(foodID) {
  return axios.get(`http://localhost:8080/api/fabric/food/${foodID}/ingredients`)
}


export function getLogInfo(foodID) {
  return axios.get(`http://localhost:8080/api/fabric/food/${foodID}/logistics`)
}


export function getLogInfoList(foodID) {
  return axios.get(`http://localhost:8080/api/fabric/food/${foodID}/logistics/list`)
}


export function getAllFoodProduction() {
  return axios.get('http://localhost:8080/api/fabric/food/production/all')
}


export function getAllFoodIngredients() {
  return axios.get('http://localhost:8080/api/fabric/food/ingredients/all')
}


export function getAllFoodLogistics() {
  return axios.get('http://localhost:8080/api/fabric/food/logistics/all')
} 
