package com.example.springsdk.service;

import com.example.springsdk.config.FabricConfig;
import lombok.extern.slf4j.Slf4j;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Collection;
import java.util.concurrent.CompletableFuture;

@Slf4j
@Service
public class FabricService {

    @Autowired
    private Channel channel;

    @Autowired
    private HFClient hfClient;


    public String addProInfo(String foodID, String foodName, String foodSpec, 
                           String foodMFGDate, String foodEXPDate, String foodLOT,
                           String foodQSID, String foodMFRSName, String foodProPrice, 
                           String foodProPlace) throws Exception {
        
        String[] args = {foodID, foodName, foodSpec, foodMFGDate, foodEXPDate, 
                        foodLOT, foodQSID, foodMFRSName, foodProPrice, foodProPlace};
        
        return invokeChaincode("addProInfo", args);
    }


    public String addIngInfo(String foodID, String... ingredients) throws Exception {
        String[] args = new String[ingredients.length + 1];
        args[0] = foodID;
        System.arraycopy(ingredients, 0, args, 1, ingredients.length);
       
        System.out.println("final args: " + java.util.Arrays.toString(args));
           
        return invokeChaincode("addIngInfo", args);
    }


    public String addLogInfo(String foodID, String logDepartureTm, String logArrivalTm,
                           String logMission, String logDeparturePl, String logDest,
                           String logToSeller, String logStorageTm, String logMOT,
                           String logCopName, String logCost) throws Exception {
        
        String[] args = {foodID, logDepartureTm, logArrivalTm, logMission, logDeparturePl,
                        logDest, logToSeller, logStorageTm, logMOT, logCopName, logCost};
        
        return invokeChaincode("addLogInfo", args);
    }


    public String getFoodInfo(String foodID) throws Exception {
        return queryChaincode("getFoodInfo", foodID);
    }


    public String getProInfo(String foodID) throws Exception {
        return queryChaincode("getProInfo", foodID);
    }


    public String getIngInfo(String foodID) throws Exception {
        return queryChaincode("getIngInfo", foodID);
    }


    public String getLogInfo(String foodID) throws Exception {
        return queryChaincode("getLogInfo", foodID);
    }


    public String getLogInfoList(String foodID) throws Exception {
        return queryChaincode("getLogInfo_l", foodID);
    }


    public java.util.List<java.util.Map<String, Object>> getAllFoodProduction() throws Exception {
        java.util.List<java.util.Map<String, Object>> foodList = new java.util.ArrayList<>();
        

        java.util.List<java.util.Map<String, Object>> allTxs = getAllTransactions();
        log.info("Total transactions found: {}", allTxs.size());
        
        for (java.util.Map<String, Object> tx : allTxs) {
            if (tx.get("chaincodeName") != null && tx.get("chaincodeName").toString().contains("source-app")) {
                log.info("Found source-app chaincode transaction: {}", tx.get("transactionID"));
                java.util.List<String> payloads = (java.util.List<String>) tx.get("payloads");
                if (payloads != null) {
                    log.info("Transaction payload count: {}", payloads.size());
                    for (String payload : payloads) {
                        log.info("Processing payload: {}", payload.substring(0, Math.min(100, payload.length())));
                        try {
                            // Parse payload
                            String payloadStr = payload;
                            if (payload.startsWith("[Base64]")) {
                                String base64Data = payload.replace("[Base64] ", "");
                                byte[] decodedBytes = java.util.Base64.getDecoder().decode(base64Data);
                                payloadStr = new String(decodedBytes, java.nio.charset.StandardCharsets.UTF_8);
                                log.info("Base64 decoded: {}", payloadStr.substring(0, Math.min(100, payloadStr.length())));
                            }
                            
                            
                            if (payloadStr.startsWith("{") || payloadStr.startsWith("[")) {
                                com.fasterxml.jackson.databind.ObjectMapper mapper = new com.fasterxml.jackson.databind.ObjectMapper();
                                com.fasterxml.jackson.databind.JsonNode jsonNode = mapper.readTree(payloadStr);
                                
                                
                                if (jsonNode.has("FoodProInfo") || jsonNode.has("FoodName")) {
                                    log.info("found food production information");
                                    java.util.Map<String, Object> foodInfo = new java.util.HashMap<>();
                                    
                                    if (jsonNode.has("FoodID")) {
                                        foodInfo.put("foodID", jsonNode.get("FoodID").asText());
                                    }
                                    
                                    if (jsonNode.has("FoodProInfo")) {
                                        com.fasterxml.jackson.databind.JsonNode proInfo = jsonNode.get("FoodProInfo");
                                        foodInfo.put("foodName", proInfo.get("FoodName").asText(""));
                                        foodInfo.put("foodSpec", proInfo.get("FoodSpec").asText(""));
                                        foodInfo.put("foodMFGDate", proInfo.get("FoodMFGDate").asText(""));
                                        foodInfo.put("foodEXPDate", proInfo.get("FoodEXPDate").asText(""));
                                        foodInfo.put("foodLOT", proInfo.get("FoodLOT").asText(""));
                                        foodInfo.put("foodQSID", proInfo.get("FoodQSID").asText(""));
                                        foodInfo.put("foodMFRSName", proInfo.get("FoodMFRSName").asText(""));
                                        foodInfo.put("foodProPrice", proInfo.get("FoodProPrice").asText(""));
                                        foodInfo.put("foodProPlace", proInfo.get("FoodProPlace").asText(""));
                                    } else if (jsonNode.has("FoodName")) {
                                        
                                        foodInfo.put("foodName", jsonNode.get("FoodName").asText(""));
                                        foodInfo.put("foodSpec", jsonNode.get("FoodSpec").asText(""));
                                        foodInfo.put("foodMFGDate", jsonNode.get("FoodMFGDate").asText(""));
                                        foodInfo.put("foodEXPDate", jsonNode.get("FoodEXPDate").asText(""));
                                        foodInfo.put("foodLOT", jsonNode.get("FoodLOT").asText(""));
                                        foodInfo.put("foodQSID", jsonNode.get("FoodQSID").asText(""));
                                        foodInfo.put("foodMFRSName", jsonNode.get("FoodMFRSName").asText(""));
                                        foodInfo.put("foodProPrice", jsonNode.get("FoodProPrice").asText(""));
                                        foodInfo.put("foodProPlace", jsonNode.get("FoodProPlace").asText(""));
                                    }
                                    
                                    
                                    if (tx.get("timestamp") != null) {
                                        foodInfo.put("timestamp", tx.get("timestamp"));
                                    }
                                    
                                    
                                    if (foodInfo.get("foodName") != null && !foodInfo.get("foodName").toString().isEmpty()) {
                                        foodList.add(foodInfo);
                                        log.info("add food information: {}", foodInfo.get("foodName"));
                                    } else {
                                        log.warn("food name is empty, skip");
                                    }
                                } else {
                                    log.info("no food production information");
                                }
                            } else {   
                                log.info("payload is not JSON format");
                            }
                        } catch (Exception e) {
                            log.warn("parse payload failed: {}", e.getMessage());
                        }
                    }
                } else {
                    log.info("transaction has no payload");
                }
                
              
                if (payloads != null) {
                    for (String payload : payloads) {
                        if (payload.startsWith("Input[0]: addProInfo")) {
                            log.info("found addProInfo transaction, extract data from input parameters");
                            java.util.Map<String, Object> foodInfo = new java.util.HashMap<>();
                            
                        
                            for (String p : payloads) {
                                if (p.startsWith("Input[")) {
                                    String[] parts = p.split(": ", 2);
                                    if (parts.length == 2) {
                                        String index = parts[0].replace("Input[", "").replace("]", "");
                                        String value = parts[1];
                                        
                                        switch (index) {
                                            case "0": // pass
                                                break;
                                            case "1": // FoodID
                                                foodInfo.put("foodID", value);
                                                break;
                                            case "2": // FoodName
                                                foodInfo.put("foodName", value);
                                                break;
                                            case "3": // FoodSpec
                                                foodInfo.put("foodSpec", value);
                                                break;
                                            case "4": // FoodMFGDate
                                                foodInfo.put("foodMFGDate", value);
                                                break;
                                            case "5": // FoodEXPDate
                                                foodInfo.put("foodEXPDate", value);
                                                break;
                                            case "6": // FoodLOT
                                                foodInfo.put("foodLOT", value);
                                                break;
                                            case "7": // FoodQSID
                                                foodInfo.put("foodQSID", value);
                                                break;
                                            case "8": // FoodMFRSName
                                                foodInfo.put("foodMFRSName", value);
                                                break;
                                            case "9": // FoodProPrice
                                                foodInfo.put("foodProPrice", value);
                                                break;
                                            case "10": // FoodProPlace
                                                foodInfo.put("foodProPlace", value);
                                                break;
                                        }
                                    }
                                }
                            }
                            
                            
                            if (tx.get("timestamp") != null) {
                                foodInfo.put("timestamp", tx.get("timestamp"));
                            }
                            
                       
                            if (foodInfo.get("foodName") != null && !foodInfo.get("foodName").toString().isEmpty()) {
                                foodList.add(foodInfo);
                                log.info("add food information from input parameters: {}", foodInfo.get("foodName"));
                            }
                            break; 
                        }
                    }
                }
            } else {
                log.debug("skip non-source-app chaincode transaction: {}", tx.get("chaincodeName"));
            }
        }
        
        log.info("total food information found: {}", foodList.size());
        return foodList;
    }


    public java.util.List<java.util.Map<String, Object>> getAllFoodIngredients() throws Exception {
        java.util.List<java.util.Map<String, Object>> ingredientList = new java.util.ArrayList<>();
        
 
        java.util.List<java.util.Map<String, Object>> allTxs = getAllTransactions();
        
        for (java.util.Map<String, Object> tx : allTxs) {
            if (tx.get("chaincodeName") != null && tx.get("chaincodeName").toString().contains("source-app")) {
                java.util.List<String> payloads = (java.util.List<String>) tx.get("payloads");
                if (payloads != null) {
                    for (String payload : payloads) {
                        try {
                          
                            String payloadStr = payload;
                            if (payload.startsWith("[Base64]")) {
                                String base64Data = payload.replace("[Base64] ", "");
                                byte[] decodedBytes = java.util.Base64.getDecoder().decode(base64Data);
                                payloadStr = new String(decodedBytes, java.nio.charset.StandardCharsets.UTF_8);
                            }
                            
                         
                            if (payloadStr.startsWith("{") || payloadStr.startsWith("[")) {
                                com.fasterxml.jackson.databind.ObjectMapper mapper = new com.fasterxml.jackson.databind.ObjectMapper();
                                com.fasterxml.jackson.databind.JsonNode jsonNode = mapper.readTree(payloadStr);
                                
                              
                                if (jsonNode.has("FoodIngInfo")) {
                                    java.util.Map<String, Object> ingredientInfo = new java.util.HashMap<>();
                                    
                                    if (jsonNode.has("FoodID")) {
                                        ingredientInfo.put("foodID", jsonNode.get("FoodID").asText());
                                    }
                                    
                                    com.fasterxml.jackson.databind.JsonNode ingInfoArray = jsonNode.get("FoodIngInfo");
                                    java.util.List<java.util.Map<String, String>> ingredients = new java.util.ArrayList<>();
                                    if (ingInfoArray.isArray()) {
                                        for (com.fasterxml.jackson.databind.JsonNode ing : ingInfoArray) {
                                            java.util.Map<String, String> ingObj = new java.util.HashMap<>();
                                            ingObj.put("IngID", ing.get("IngID").asText(""));
                                            ingObj.put("IngName", ing.get("IngName").asText(""));
                                            ingredients.add(ingObj);
                                        }
                                    }
                                    ingredientInfo.put("ingredients", ingredients);
                                
                                    if (tx.get("timestamp") != null) {
                                        ingredientInfo.put("timestamp", tx.get("timestamp"));
                                    }
                                    
                              
                                    if (!ingredients.isEmpty()) {
                                        ingredientList.add(ingredientInfo);
                                        log.info("add ingredient information from JSON: FoodID={}, ingredients={}", 
                                            ingredientInfo.get("foodID"), ingredients);
                                    }
                                }
                            }
                        } catch (Exception e) {
                            log.warn("parse ingredient payload failed: {}", e.getMessage());
                        }
                    }
                }
                
              
                if (payloads != null) {
                    for (String payload : payloads) {
                        if (payload.startsWith("Input[0]: addIngInfo")) {
                            log.info("found addIngInfo transaction, extract data from input parameters");
                            java.util.Map<String, Object> ingredientInfo = new java.util.HashMap<>();
                            java.util.List<java.util.Map<String, String>> ingredients = new java.util.ArrayList<>();
                            
                         
                            for (String p : payloads) {
                                if (p.startsWith("Input[")) {
                                    String[] parts = p.split(": ", 2);
                                    if (parts.length == 2) {
                                        String index = parts[0].replace("Input[", "").replace("]", "");
                                        String value = parts[1];
                                        
                                        switch (index) {
                                            case "0": //
                                                break;
                                            case "1": // FoodID
                                                ingredientInfo.put("foodID", value);
                                                break;
                                            default: 
                                                if (Integer.parseInt(index) % 2 == 0) {
                                                    
                                                    if (!value.isEmpty()) {
                                                        java.util.Map<String, String> ingObj = new java.util.HashMap<>();
                                                        ingObj.put("IngID", index);
                                                        ingObj.put("IngName", value);
                                                        ingredients.add(ingObj);
                                                    }
                                                }
                                                break;
                                        }
                                    }
                                }
                            }
                            
                            ingredientInfo.put("ingredients", ingredients);
                            
                      
                            if (tx.get("timestamp") != null) {
                                ingredientInfo.put("timestamp", tx.get("timestamp"));
                            }
                            
                         
                            if (!ingredients.isEmpty()) {
                                ingredientList.add(ingredientInfo);
                                log.info("add ingredient information from input parameters: FoodID={}, ingredients={}", 
                                    ingredientInfo.get("foodID"), ingredients.size());
                            }
                            break; 
                        }
                    }
                }
            } else {
                log.debug("skip non-source-app chaincode transaction: {}", tx.get("chaincodeName"));
            }
        }
        
        log.info("total ingredient information found: {}", ingredientList.size());
        return ingredientList;
    }


    public java.util.List<java.util.Map<String, Object>> getAllFoodLogistics() throws Exception {
        java.util.List<java.util.Map<String, Object>> logisticsList = new java.util.ArrayList<>();
        
     
        java.util.List<java.util.Map<String, Object>> allTxs = getAllTransactions();
        
        for (java.util.Map<String, Object> tx : allTxs) {
            if (tx.get("chaincodeName") != null && tx.get("chaincodeName").toString().contains("source-app")) {
                java.util.List<String> payloads = (java.util.List<String>) tx.get("payloads");
                if (payloads != null) {
                    for (String payload : payloads) {
                        try {
                           
                            String payloadStr = payload;
                            if (payload.startsWith("[Base64]")) {
                                String base64Data = payload.replace("[Base64] ", "");
                                byte[] decodedBytes = java.util.Base64.getDecoder().decode(base64Data);
                                payloadStr = new String(decodedBytes, java.nio.charset.StandardCharsets.UTF_8);
                            }
                            
                      
                            if (payloadStr.startsWith("{") || payloadStr.startsWith("[")) {
                                com.fasterxml.jackson.databind.ObjectMapper mapper = new com.fasterxml.jackson.databind.ObjectMapper();
                                com.fasterxml.jackson.databind.JsonNode jsonNode = mapper.readTree(payloadStr);
                                
                             
                                if (jsonNode.has("FoodLogInfo") || jsonNode.has("LogMission")) {
                                    java.util.Map<String, Object> logisticsInfo = new java.util.HashMap<>();
                                    
                                    if (jsonNode.has("FoodID")) {
                                        logisticsInfo.put("foodID", jsonNode.get("FoodID").asText());
                                    }
                                    
                                    if (jsonNode.has("FoodLogInfo")) {
                                        com.fasterxml.jackson.databind.JsonNode logInfo = jsonNode.get("FoodLogInfo");
                                        logisticsInfo.put("logDepartureTm", logInfo.get("LogDepartureTm").asText(""));
                                        logisticsInfo.put("logArrivalTm", logInfo.get("LogArrivalTm").asText(""));
                                        logisticsInfo.put("logMission", logInfo.get("LogMission").asText(""));
                                        logisticsInfo.put("logDeparturePl", logInfo.get("LogDeparturePl").asText(""));
                                        logisticsInfo.put("logDest", logInfo.get("LogDest").asText(""));
                                        logisticsInfo.put("logToSeller", logInfo.get("LogToSeller").asText(""));
                                        logisticsInfo.put("logStorageTm", logInfo.get("LogStorageTm").asText(""));
                                        logisticsInfo.put("logMOT", logInfo.get("LogMOT").asText(""));
                                        logisticsInfo.put("logCopName", logInfo.get("LogCopName").asText(""));
                                        logisticsInfo.put("logCost", logInfo.get("LogCost").asText(""));
                                    } else if (jsonNode.has("LogMission")) {
                                
                                        logisticsInfo.put("logDepartureTm", jsonNode.get("LogDepartureTm").asText(""));
                                        logisticsInfo.put("logArrivalTm", jsonNode.get("LogArrivalTm").asText(""));
                                        logisticsInfo.put("logMission", jsonNode.get("LogMission").asText(""));
                                        logisticsInfo.put("logDeparturePl", jsonNode.get("LogDeparturePl").asText(""));
                                        logisticsInfo.put("logDest", jsonNode.get("LogDest").asText(""));
                                        logisticsInfo.put("logToSeller", jsonNode.get("LogToSeller").asText(""));
                                        logisticsInfo.put("logStorageTm", jsonNode.get("LogStorageTm").asText(""));
                                        logisticsInfo.put("logMOT", jsonNode.get("LogMOT").asText(""));
                                        logisticsInfo.put("logCopName", jsonNode.get("LogCopName").asText(""));
                                        logisticsInfo.put("logCost", jsonNode.get("LogCost").asText(""));
                                    }
                                    
                                   
                                    if (tx.get("timestamp") != null) {
                                        logisticsInfo.put("timestamp", tx.get("timestamp"));
                                    }
                                    
                       
                                    if (logisticsInfo.get("logMission") != null && 
                                        !logisticsInfo.get("logMission").toString().isEmpty()) {
                                        logisticsList.add(logisticsInfo);
                                    }
                                }
                            }
                        } catch (Exception e) {
                            log.warn("parse logistics payload failed: {}", e.getMessage());
                        }
                    }
                }
                
             
                if (payloads != null) {
                    for (String payload : payloads) {
                        if (payload.startsWith("Input[0]: addLogInfo")) {
                            log.info("found addLogInfo transaction, extract data from input parameters");
                            java.util.Map<String, Object> logisticsInfo = new java.util.HashMap<>();
                            
                         
                            for (String p : payloads) {
                                if (p.startsWith("Input[")) {
                                    String[] parts = p.split(": ", 2);
                                    if (parts.length == 2) {
                                        String index = parts[0].replace("Input[", "").replace("]", "");
                                        String value = parts[1];
                                        
                                        switch (index) {
                                            case "0": 
                                                break;
                                            case "1": // FoodID
                                                logisticsInfo.put("foodID", value);
                                                break;
                                            case "2": // LogDepartureTm
                                                logisticsInfo.put("logDepartureTm", value);
                                                break;
                                            case "3": // LogArrivalTm
                                                logisticsInfo.put("logArrivalTm", value);
                                                break;
                                            case "4": // LogMission
                                                logisticsInfo.put("logMission", value);
                                                break;
                                            case "5": // LogDeparturePl
                                                logisticsInfo.put("logDeparturePl", value);
                                                break;
                                            case "6": // LogDest
                                                logisticsInfo.put("logDest", value);
                                                break;
                                            case "7": // LogToSeller
                                                logisticsInfo.put("logToSeller", value);
                                                break;
                                            case "8": // LogStorageTm
                                                logisticsInfo.put("logStorageTm", value);
                                                break;
                                            case "9": // LogMOT
                                                logisticsInfo.put("logMOT", value);
                                                break;
                                            case "10": // LogCopName
                                                logisticsInfo.put("logCopName", value);
                                                break;
                                            case "11": // LogCost
                                                logisticsInfo.put("logCost", value);
                                                break;
                                        }
                                    }
                                }
                            }
                            
            
                            if (tx.get("timestamp") != null) {
                                logisticsInfo.put("timestamp", tx.get("timestamp"));
                            }
                  
                            if (logisticsInfo.get("foodID") != null && !logisticsInfo.get("foodID").toString().isEmpty()) {
                                logisticsList.add(logisticsInfo);
                                log.info("add logistics information from input parameters: FoodID={}", logisticsInfo.get("foodID"));
                            }
                            break; 
                        }
                    }
                }
            } else {
                log.debug("skip non-source-app chaincode transaction: {}", tx.get("chaincodeName"));
            }
        }
        
        log.info("total logistics information found: {}", logisticsList.size());
        return logisticsList;
    }


    private String invokeChaincode(String function, String... args) throws Exception {
        try {
            log.info("start calling chaincode, function: {}, args: {}", function, java.util.Arrays.toString(args));
            
            ChaincodeID chaincodeID = ChaincodeID.newBuilder()
                    .setName(FabricConfig.getChaincodeName())
                    .setVersion(FabricConfig.getChaincodeVersion())
                    .build();
            
            log.info("chaincode ID: {}", chaincodeID.getName() + ":" + chaincodeID.getVersion());

            TransactionProposalRequest request = hfClient.newTransactionProposalRequest();
            request.setChaincodeID(chaincodeID);
            request.setFcn(function);
            request.setArgs(args);

            log.info("send transaction proposal...");
            Collection<ProposalResponse> responses = channel.sendTransactionProposal(request);
            
            log.info("received {} transaction proposal responses", responses.size());
            

            for (ProposalResponse response : responses) {
                log.info("transaction proposal response status: {}, message: {}", response.getStatus(), response.getMessage());
                if (response.getStatus() != ChaincodeResponse.Status.SUCCESS) {
                    String errorMsg = "transaction proposal failed: " + response.getMessage();
                    log.error(errorMsg);
                    throw new RuntimeException(errorMsg);
                }
            }

            log.info("all transaction proposal responses successful, send transaction to orderer...");
            CompletableFuture<BlockEvent.TransactionEvent> future = channel.sendTransaction(responses);
            BlockEvent.TransactionEvent event = future.get();
            
            log.info("transaction successfully submitted, transaction ID: {}", event.getTransactionID());
            return "transaction successfully submitted, transaction ID: " + event.getTransactionID();
            
        } catch (Exception e) {
            log.error("call chaincode failed, function: {}, error: {}", function, e.getMessage(), e);
            throw e;
        }
    }


    private String queryChaincode(String function, String... args) throws Exception {
        try {
            log.info("start querying chaincode, function: {}, args: {}", function, java.util.Arrays.toString(args));
            
            ChaincodeID chaincodeID = ChaincodeID.newBuilder()
                    .setName(FabricConfig.getChaincodeName())
                    .setVersion(FabricConfig.getChaincodeVersion())
                    .build();

            QueryByChaincodeRequest request = hfClient.newQueryProposalRequest();
            request.setChaincodeID(chaincodeID);
            request.setFcn(function);
            request.setArgs(args);

            log.info("send query request...");
            Collection<ProposalResponse> responses = channel.queryByChaincode(request);
            
            log.info("received {} query responses", responses.size());
            
            for (ProposalResponse response : responses) {
                log.info("query response status: {}, message: {}", response.getStatus(), response.getMessage());
                if (response.getStatus() == ChaincodeResponse.Status.SUCCESS) {
                    String result = new String(response.getChaincodeActionResponsePayload());
                    log.info("query successful, result: {}", result);
                    return result;
                } else {
                    String errorMsg = "query failed: " + response.getMessage();
                    log.error(errorMsg);
                    throw new RuntimeException(errorMsg);
                }
            }
            
            log.warn("no valid query response received");
            return null;
            
        } catch (Exception e) {
            log.error("query chaincode failed, function: {}, error: {}", function, e.getMessage(), e);
            throw e;
        }
    }


    public String getBlockchainInfo() throws Exception {
        BlockInfo blockInfo = channel.queryBlockByNumber(channel.queryBlockchainInfo().getHeight() - 1);
        return "latest block height: " + blockInfo.getBlockNumber() + 
               ", block hash: " + blockInfo.getPreviousHash();
    }


    public java.util.List<java.util.Map<String, Object>> getAllTransactions() throws Exception {
        java.util.List<java.util.Map<String, Object>> txList = new java.util.ArrayList<>();
        BlockchainInfo blockchainInfo = channel.queryBlockchainInfo();
        long height = blockchainInfo.getHeight();
        for (long i = 0; i < height; i++) {
            BlockInfo blockInfo = channel.queryBlockByNumber(i);
            int blockNumber = (int) blockInfo.getBlockNumber();
            for (BlockInfo.EnvelopeInfo envelopeInfo : blockInfo.getEnvelopeInfos()) {
                java.util.Map<String, Object> tx = new java.util.HashMap<>();
                tx.put("blockNumber", blockNumber);
                tx.put("transactionID", envelopeInfo.getTransactionID());
                tx.put("timestamp", envelopeInfo.getTimestamp());
                tx.put("type", envelopeInfo.getType().toString());
                tx.put("creatorMSP", envelopeInfo.getCreator() != null ? envelopeInfo.getCreator().getMspid() : null);
                tx.put("creator", envelopeInfo.getCreator() != null ? envelopeInfo.getCreator().getId() : null);

             
                if (envelopeInfo.getType() == BlockInfo.EnvelopeType.TRANSACTION_ENVELOPE) {
                    BlockInfo.TransactionEnvelopeInfo txe = (BlockInfo.TransactionEnvelopeInfo) envelopeInfo;
                    java.util.List<String> payloads = new java.util.ArrayList<>();
                    String chaincodeName = null;
                    for (BlockInfo.TransactionEnvelopeInfo.TransactionActionInfo actionInfo : txe.getTransactionActionInfos()) {
                     
                        int inputCount = actionInfo.getChaincodeInputArgsCount();
                        log.info("chaincode input parameters count: {}", inputCount);
                        
                        for (int j = 0; j < inputCount; j++) {
                            byte[] chaincodeInput = actionInfo.getChaincodeInputArgs(j);
                            if (chaincodeInput != null && chaincodeInput.length > 0) {
                                try {
                                    String inputStr = new String(chaincodeInput, java.nio.charset.StandardCharsets.UTF_8);
                                    log.info("chaincode input parameters[{}]: {}", j, inputStr);
                                    payloads.add("Input[" + j + "]: " + inputStr);
                                } catch (Exception e) {
                                    log.warn("parse chaincode input parameters[{}] failed: {}", j, e.getMessage());
                                }
                            }
                        }
                        
                       
                        byte[] payload = actionInfo.getProposalResponsePayload();
                        if (payload != null) {
                            String payloadStr;
                            try {
                                payloadStr = new String(payload, java.nio.charset.StandardCharsets.UTF_8);
                                log.info("original payload length: {}, content: {}", payload.length, payloadStr);
                                
                             
                                boolean isValidJson = false;
                                boolean containsChaincodeFunction = false;
                                
                            
                                try {
                                    new com.fasterxml.jackson.databind.ObjectMapper().readTree(payloadStr);
                                    isValidJson = true;
                                } catch (Exception e) {
                                  
                                }
                                
                             
                                if (payloadStr.contains("addProInfo") || payloadStr.contains("addIngInfo") || 
                                    payloadStr.contains("addLogInfo") || payloadStr.contains("getProInfo") ||
                                    payloadStr.contains("getIngInfo") || payloadStr.contains("getLogInfo")) {
                                    containsChaincodeFunction = true;
                                }
                                
                             
                                int nonPrintable = 0;
                                for (char c : payloadStr.toCharArray()) {
                                    if (c < 0x09 || (c > 0x0D && c < 0x20) || c > 0x7E) nonPrintable++;
                                }
                                
                              
                                if (isValidJson || containsChaincodeFunction || 
                                    payloadStr.length() == 0 || ((double)nonPrintable / payloadStr.length()) <= 0.5) {
                                    
                                    log.info("keep original payload string: {}", payloadStr);
                                    if (!payloadStr.isEmpty()) {
                                        payloads.add("Response: " + payloadStr);
                                    }
                                } else {
                                   
                                    payloadStr = "[Base64] " + java.util.Base64.getEncoder().encodeToString(payload);
                                    log.info("convert to Base64: {}", payloadStr.substring(0, Math.min(100, payloadStr.length())));
                                    payloads.add("Response: " + payloadStr);
                                }
                            } catch (Exception e) {
                                payloadStr = "[Base64] " + java.util.Base64.getEncoder().encodeToString(payload);
                                log.info("Exception handling, convert to Base64: {}", payloadStr.substring(0, Math.min(100, payloadStr.length())));
                                payloads.add("Response: " + payloadStr);
                            }
                        } else {
                            log.info("payload is null");
                        }
                        if (actionInfo.getChaincodeIDName() != null && chaincodeName == null) {
                            chaincodeName = actionInfo.getChaincodeIDName();
                        }
                    }
                    tx.put("payloads", payloads);
                    tx.put("chaincodeName", chaincodeName);
                }
                txList.add(tx);
            }
        }
        return txList;
    }

    public Channel getChannel() {
        return this.channel;
    }


public java.util.Map<String, Object> getNodeLogs(String nodeName) throws Exception {
    java.util.Map<String, Object> nodeLogs = new java.util.HashMap<>();
    
    try {
        log.info("get node logs: {}", nodeName);
        
        java.util.Map<String, String> nodeConfigs = new java.util.HashMap<>();
        nodeConfigs.put("peer0.org1.example.com", "grpc://localhost:7051");
        nodeConfigs.put("orderer0.example.com", "grpc://localhost:6050");
        nodeConfigs.put("orderer1.example.com", "grpc://localhost:6051");
        nodeConfigs.put("orderer2.example.com", "grpc://localhost:6052");
        nodeConfigs.put("orderer3.example.com", "grpc://localhost:6053");
        nodeConfigs.put("orderer4.example.com", "grpc://localhost:6054");
        
        if (!nodeConfigs.containsKey(nodeName)) {
            throw new RuntimeException("node not found: " + nodeName);
        }
        
        String nodeUrl = nodeConfigs.get(nodeName);
        boolean isPeer = nodeName.startsWith("peer");
        boolean isOrderer = nodeName.startsWith("orderer");
        
        nodeLogs.put("nodeName", nodeName);
        nodeLogs.put("nodeType", isPeer ? "Peer" : (isOrderer ? "Orderer" : "Unknown"));
        nodeLogs.put("nodeUrl", nodeUrl);
        
        java.util.Map<String, Object> statusInfo = new java.util.HashMap<>();
        
        try {
            ProcessBuilder pb = new ProcessBuilder("docker", "ps", "--filter", "name=" + nodeName, "--format", "{{.Status}}");
            Process process = pb.start();
            java.io.BufferedReader reader = new java.io.BufferedReader(new java.io.InputStreamReader(process.getInputStream()));
            String status = reader.readLine();
            process.waitFor();
            
            if (status != null && !status.trim().isEmpty()) {
                statusInfo.put("status", "Online");
                statusInfo.put("connection", "Running");
                statusInfo.put("lastSeen", java.time.LocalDateTime.now().toString());
                statusInfo.put("containerStatus", status.trim());
            } else {
                statusInfo.put("status", "Offline");
                statusInfo.put("connection", "Container not running");
                statusInfo.put("lastSeen", "Unknown");
            }
            
        } catch (Exception e) {
            log.warn("check Docker container status failed: {}", e.getMessage());
            statusInfo.put("status", "Unknown");
            statusInfo.put("connection", "Unknown");
            statusInfo.put("lastSeen", "Unknown");
        }
        
        java.util.List<String> rawLogs = new java.util.ArrayList<>();
        
        try {
            java.util.Map<String, String> containerNameMap = new java.util.HashMap<>();
            containerNameMap.put("peer0.org1.example.com", "peer0.org1.example.com");
            containerNameMap.put("orderer0.example.com", "orderer0.example.com");
            containerNameMap.put("orderer1.example.com", "orderer1.example.com");
            containerNameMap.put("orderer2.example.com", "orderer2.example.com");
            containerNameMap.put("orderer3.example.com", "orderer3.example.com");
            containerNameMap.put("orderer4.example.com", "orderer4.example.com");
            
            String containerName = containerNameMap.get(nodeName);
            log.info("Looking up container: {}, mapped to: {}", nodeName, containerName);
            
            if (containerName == null) {
                rawLogs.add("Error: No container mapping found for node " + nodeName);
                rawLogs.add("Supported nodes: " + String.join(", ", containerNameMap.keySet()));
            } else {
                ProcessBuilder checkPb = new ProcessBuilder("docker", "ps", "-a", "--format", "{{.Names}}");
                Process checkProcess = checkPb.start();
                java.io.BufferedReader checkReader = new java.io.BufferedReader(new java.io.InputStreamReader(checkProcess.getInputStream()));
                String actualContainerName = null;
                String line;
                while ((line = checkReader.readLine()) != null) {
                    String container = line.trim();
                    if (container.equals(containerName)) {
                        actualContainerName = container;
                        break;
                    }
                }
                checkProcess.waitFor();
                
                if (actualContainerName == null) {
                    rawLogs.add("Error: Container " + containerName + " does not exist");
                    rawLogs.add("Please check if the Docker container is running");
                    rawLogs.add("Currently available containers:");
                    
                    ProcessBuilder listPb = new ProcessBuilder("docker", "ps", "-a", "--format", "{{.Names}}");
                    Process listProcess = listPb.start();
                    java.io.BufferedReader listReader = new java.io.BufferedReader(new java.io.InputStreamReader(listProcess.getInputStream()));
                    String listLine;
                    while ((listLine = listReader.readLine()) != null) {
                        String container = listLine.trim();
                        if (container.contains("peer") || container.contains("orderer")) {
                            rawLogs.add("  - " + container);
                        }
                    }
                    listProcess.waitFor();
                } else {
                    containerName = actualContainerName;
                    log.info("Confirmed container exists: {}", containerName);
                }
            }
            
            if (containerName == null || containerName.trim().isEmpty()) {
                rawLogs.add("Error: Container not found for " + nodeName);
                rawLogs.add("Hint: Please check if the container name is correct");
            } else {
                log.info("Start reading container logs: {}", containerName);
                ProcessBuilder pb = new ProcessBuilder("docker", "logs", containerName, "--tail", "100");
                Process process = pb.start();
                
                java.io.BufferedReader reader = new java.io.BufferedReader(new java.io.InputStreamReader(process.getInputStream()));
                String line;
                while ((line = reader.readLine()) != null) {
                    rawLogs.add(line);
                }
                
                java.io.BufferedReader errorReader = new java.io.BufferedReader(new java.io.InputStreamReader(process.getErrorStream()));
                String errorLine;
                while ((errorLine = errorReader.readLine()) != null) {
                    rawLogs.add(errorLine);
                }
                
                int exitCode = process.waitFor();
                log.info("Docker logs command exit code: {}", exitCode);
                
                if (rawLogs.isEmpty()) {
                    rawLogs.add("Container " + containerName + " has no log output yet");
                    rawLogs.add("Possible reasons: Container just started or log level is set too high");
                }
            }
            
        } catch (Exception e) {
            log.warn("Failed to read Docker logs: {}", e.getMessage(), e);
            rawLogs.add("Error: Failed to read Docker logs - " + e.getMessage());
            rawLogs.add("Please ensure Docker service is running and current user has permission to execute docker commands");
        }
        
        nodeLogs.put("statusInfo", statusInfo);
        nodeLogs.put("logs", rawLogs);
        nodeLogs.put("logCount", rawLogs.size());
        nodeLogs.put("logSource", "Docker Container");
        
        log.info("successfully get Docker logs for node {}, total {} lines", nodeName, rawLogs.size());
        
    } catch (Exception e) {
        log.error("get node logs failed: {}", e.getMessage(), e);
        nodeLogs.put("error", "get node logs failed: " + e.getMessage());
    }
    
    return nodeLogs;
}


    public String checkChaincodeStatus() throws Exception {
        try {
            log.info("check chaincode status...");
            
    
            ChaincodeID chaincodeID = ChaincodeID.newBuilder()
                    .setName(FabricConfig.getChaincodeName())
                    .setVersion(FabricConfig.getChaincodeVersion())
                    .build();
            
            log.info("chaincode ID: {}", chaincodeID.getName() + ":" + chaincodeID.getVersion());
            
    
            QueryByChaincodeRequest request = hfClient.newQueryProposalRequest();
            request.setChaincodeID(chaincodeID);
            request.setFcn("getFoodInfo");
            request.setArgs("TEST_CHECK");
            
            Collection<ProposalResponse> responses = channel.queryByChaincode(request);
            
            for (ProposalResponse response : responses) {
                log.info("chaincode status check response: {}", response.getMessage());
                if (response.getStatus() == ChaincodeResponse.Status.SUCCESS) {
                    return "chaincode status is normal, can be called normally";
                } else {
                    return "chaincode status is abnormal: " + response.getMessage();
                }
            }
            
            return "cannot get chaincode status information";
            
        } catch (Exception e) {
            log.error("check chaincode status failed: {}", e.getMessage(), e);
            return "check chaincode status failed: " + e.getMessage();
        }
    }


    public String getChannelInfo() throws Exception {
        try {
            BlockchainInfo blockchainInfo = channel.queryBlockchainInfo();
            return String.format("channel name: %s, block height: %d, current block hash: %s", 
                channel.getName(), 
                blockchainInfo.getHeight(),
                blockchainInfo.getCurrentBlockHash());
        } catch (Exception e) {
            log.error("get channel info failed: {}", e.getMessage(), e);
            throw e;
        }
    }
} 
