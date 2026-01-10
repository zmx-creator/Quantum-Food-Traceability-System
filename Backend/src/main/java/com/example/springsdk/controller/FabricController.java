package com.example.springsdk.controller;

import com.example.springsdk.service.FabricService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;
import javax.xml.bind.DatatypeConverter;

@Slf4j
@RestController
@RequestMapping("/api/fabric")
@CrossOrigin(origins = "*")
public class FabricController {

    @Autowired
    private FabricService fabricService;

    /**
     * Get blockchain information
     */
    @GetMapping("/info")
    public ResponseEntity<Map<String, Object>> getBlockchainInfo() {
        Map<String, Object> response = new HashMap<>();
        try {
            String info = fabricService.getBlockchainInfo();
            response.put("success", true);
            response.put("data", info);
            response.put("message", "Got blockchain information successfullu");
        } catch (Exception e) {
            log.error("Failed to get blockchain information", e);
            response.put("success", false);
            response.put("message", "Failed to get blockchain information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Add food production information
     */
    @PostMapping("/food/production")
    public ResponseEntity<Map<String, Object>> addProInfo(@RequestBody Map<String, Object> requestBody) {
        
        Map<String, Object> response = new HashMap<>();
        try {
            String foodID = (String) requestBody.get("foodID");
            String foodName = (String) requestBody.get("foodName");
            String foodSpec = (String) requestBody.get("foodSpec");
            String foodMFGDate = (String) requestBody.get("foodMFGDate");
            String foodEXPDate = (String) requestBody.get("foodEXPDate");
            String foodLOT = (String) requestBody.get("foodLOT");
            String foodQSID = (String) requestBody.get("foodQSID");
            String foodMFRSName = (String) requestBody.get("foodMFRSName");
            String foodProPrice = (String) requestBody.get("foodProPrice");
            String foodProPlace = (String) requestBody.get("foodProPlace");
            
            if (foodID == null || foodName == null) {
                response.put("success", false);
                response.put("message", "Missing required parameters：foodID or foodName");
                return ResponseEntity.ok(response);
            }
            
            String result = fabricService.addProInfo(foodID, foodName, foodSpec, foodMFGDate, 
                                                   foodEXPDate, foodLOT, foodQSID, foodMFRSName, 
                                                   foodProPrice, foodProPlace);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Added food production information successfully");
        } catch (Exception e) {
            log.error("Failed to add food production information", e);
            response.put("success", false);
            response.put("message", "Failed to add food production information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Add food ingredient information
     */
    @PostMapping("/food/ingredients")
    public ResponseEntity<Map<String, Object>> addIngInfo(@RequestBody Map<String, Object> requestBody) {

        Map<String, Object> response = new HashMap<>();
        try {
            String foodID = (String) requestBody.get("foodID");
            Object ingredientPairsObj = requestBody.get("ingredientPairs");
            
            if (foodID == null || ingredientPairsObj == null) {
                response.put("success", false);
                response.put("message", "Missing required parameters：foodID or ingredientPairs");
                return ResponseEntity.ok(response);
            }
            
            // Convert the ingradient array to the format expected by the chain code
            java.util.List<String> ingredientPairsList;
            if (ingredientPairsObj instanceof java.util.List) {
                ingredientPairsList = (java.util.List<String>) ingredientPairsObj;
            } else if (ingredientPairsObj instanceof String[]) {
                ingredientPairsList = java.util.Arrays.asList((String[]) ingredientPairsObj);
            } else {
                response.put("success", false);
                response.put("message", "ingredientPairs formatting error");
                return ResponseEntity.ok(response);
            }

            // System.out.println("ingredientPairsObj = " + ingredientPairsObj);
            // System.out.println("ingredientPairsList = " + ingredientPairsList);
            
            // Check if the number of parameters is even -- pairs of IDs and names
            if (ingredientPairsList.size() % 2 != 0) {
                response.put("success", false);
                response.put("message", "Ingredient information must be in pairs of IDs and names");
                return ResponseEntity.ok(response);
            }
            
            // Converted to an array of strings passed to the chain code
            String[] ingredientPairs = ingredientPairsList.toArray(new String[0]);
            
            String result = fabricService.addIngInfo(foodID, ingredientPairs);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Added food ingredient information successfully");
        } catch (Exception e) {
            log.error("Failed to add food ingredient information", e);
            response.put("success", false);
            response.put("message", "Failed to add food ingredient information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Add logistics information
     */
    @PostMapping("/food/logistics")
    public ResponseEntity<Map<String, Object>> addLogInfo(@RequestBody Map<String, Object> requestBody) {
        
        Map<String, Object> response = new HashMap<>();
        try {
            String foodID = (String) requestBody.get("foodID");
            String logDepartureTm = (String) requestBody.get("logDepartureTm");
            String logArrivalTm = (String) requestBody.get("logArrivalTm");
            String logMission = (String) requestBody.get("logMission");
            String logDeparturePl = (String) requestBody.get("logDeparturePl");
            String logDest = (String) requestBody.get("logDest");
            String logToSeller = (String) requestBody.get("logToSeller");
            String logStorageTm = (String) requestBody.get("logStorageTm");
            String logMOT = (String) requestBody.get("logMOT");
            String logCopName = (String) requestBody.get("logCopName");
            String logCost = (String) requestBody.get("logCost");
            
            if (foodID == null) {
                response.put("success", false);
                response.put("message", "Missing required parameters：foodID");
                return ResponseEntity.ok(response);
            }
            
            String result = fabricService.addLogInfo(foodID, logDepartureTm, logArrivalTm,
                                                   logMission, logDeparturePl, logDest,
                                                   logToSeller, logStorageTm, logMOT,
                                                   logCopName, logCost);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Added logistics information successfully");
        } catch (Exception e) {
            log.error("Failed to add logistics information", e);
            response.put("success", false);
            response.put("message", "Failed to add logistics information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get complete information about food
     */
    @GetMapping("/food/{foodID}")
    public ResponseEntity<Map<String, Object>> getFoodInfo(@PathVariable String foodID) {
        Map<String, Object> response = new HashMap<>();
        try {
            String result = fabricService.getFoodInfo(foodID);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Got food information successfully");
        } catch (Exception e) {
            log.error("failed to get food information", e);
            response.put("success", false);
            response.put("message", "failed to get food information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get food production information by foodID
     */
    @GetMapping("/food/{foodID}/production")
    public ResponseEntity<Map<String, Object>> getProInfo(@PathVariable String foodID) {
        Map<String, Object> response = new HashMap<>();
        try {
            String result = fabricService.getProInfo(foodID);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Got production information successfully");
        } catch (Exception e) {
            log.error("Failed to get production information", e);
            response.put("success", false);
            response.put("message", "Failed to get production information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get food ingredient information by foodID
     */
    @GetMapping("/food/{foodID}/ingredients")
    public ResponseEntity<Map<String, Object>> getIngInfo(@PathVariable String foodID) {
        Map<String, Object> response = new HashMap<>();
        try {
            String result = fabricService.getIngInfo(foodID);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Got production ingredient successfully");
        } catch (Exception e) {
            log.error("Failed to get ingredient information", e);
            response.put("success", false);
            response.put("message", "Failed to get ingredient information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get food logistics information by foodID
     */
    @GetMapping("/food/{foodID}/logistics")
    public ResponseEntity<Map<String, Object>> getLogInfo(@PathVariable String foodID) {
        Map<String, Object> response = new HashMap<>();
        try {
            String result = fabricService.getLogInfo(foodID);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Got production logistics successfully");
        } catch (Exception e) {
            log.error("Failed to get logistics information", e);
            response.put("success", false);
            response.put("message", "Failed to get logistics information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get food logistics information -- in list form
     */
    @GetMapping("/food/{foodID}/logistics/list")
    public ResponseEntity<Map<String, Object>> getLogInfoList(@PathVariable String foodID) {
        Map<String, Object> response = new HashMap<>();
        try {
            String result = fabricService.getLogInfoList(foodID);
            response.put("success", true);
            response.put("data", result);
            response.put("message", "Got production logistics list successfully");
        } catch (Exception e) {
            log.error("Failed to get logistics information list", e);
            response.put("success", false);
            response.put("message", "Failed to get logistics information list: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get all blockchain transaction details
     */
    @GetMapping("/blockchain/transactions")
    public ResponseEntity<Map<String, Object>> getAllTransactions() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, Object>> txList = fabricService.getAllTransactions();
            response.put("success", true);
            response.put("data", txList);
            response.put("message", "Got all blockchain transaction details successfully");
        } catch (Exception e) {
            log.error("Failed to get all blockchain transaction details", e);
            response.put("success", false);
            response.put("message", "Failed to get all blockchain transaction details: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     *  Get Peer node information
     */
    @GetMapping("/peers")
    public ResponseEntity<Map<String, Object>> getPeers() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, String>> peers = new java.util.ArrayList<>();
            for (org.hyperledger.fabric.sdk.Peer peer : fabricService.getChannel().getPeers()) {
                java.util.Map<String, String> info = new java.util.HashMap<>();
                info.put("name", peer.getName());
                info.put("url", peer.getUrl());
                peers.add(info);
            }
            response.put("success", true);
            response.put("data", peers);
            response.put("message", "Get peer node information success");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Get peer node information failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get Orderer nodes information
     */
    @GetMapping("/orderers")
    public ResponseEntity<Map<String, Object>> getOrderers() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, String>> orderers = new java.util.ArrayList<>();
            for (org.hyperledger.fabric.sdk.Orderer orderer : fabricService.getChannel().getOrderers()) {
                java.util.Map<String, String> info = new java.util.HashMap<>();
                info.put("name", orderer.getName());
                info.put("url", orderer.getUrl());
                orderers.add(info);
            }
            response.put("success", true);
            response.put("data", orderers);
            response.put("message", "Get orderer node information success");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Get orderer node information failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get node logs information
     */
    @GetMapping("/nodes/{nodeName}/logs")
    public ResponseEntity<Map<String, Object>> getNodeLogs(@PathVariable String nodeName) {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.Map<String, Object> nodeLogs = fabricService.getNodeLogs(nodeName);
            response.put("success", true);
            response.put("data", nodeLogs);
            response.put("message", "Got node logs information successfully");
        } catch (Exception e) {
            log.error("Failed to get node logs information", e);
            response.put("success", false);
            response.put("message", "Failed to get node logs information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get blockchain details
     */
    @GetMapping("/blockchain/details")
    public ResponseEntity<Map<String, Object>> getBlockchainDetails() {
        Map<String, Object> response = new HashMap<>();
        try {
            org.hyperledger.fabric.sdk.Channel channel = fabricService.getChannel();
            org.hyperledger.fabric.sdk.BlockchainInfo blockchainInfo = channel.queryBlockchainInfo();
            long height = blockchainInfo.getHeight();
            java.util.List<java.util.Map<String, Object>> blocks = new java.util.ArrayList<>();
            for (long i = 0; i < height; i++) {
                org.hyperledger.fabric.sdk.BlockInfo blockInfo = channel.queryBlockByNumber(i);
                java.util.Map<String, Object> block = new java.util.HashMap<>();
                block.put("blockNumber", blockInfo.getBlockNumber());
                byte[] dataHash = blockInfo.getDataHash();
                String dataHashHex = dataHash != null ? DatatypeConverter.printHexBinary(dataHash) : null;
                block.put("dataHash", dataHashHex);
                block.put("previousHash", blockInfo.getPreviousHash() != null ? blockInfo.getPreviousHash().toString() : null);
                block.put("transactionCount", blockInfo.getTransactionCount());
                blocks.add(block);
            }
            java.util.Map<String, Object> data = new java.util.HashMap<>();
            data.put("totalBlocks", height);
            data.put("blocks", blocks);
            response.put("success", true);
            response.put("data", data);
            response.put("message", "Get blockchain details success");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Get blockchain details failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get all food production information --for production traceability
     */
    @GetMapping("/food/production/all")
    public ResponseEntity<Map<String, Object>> getAllFoodProduction() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, Object>> foodList = fabricService.getAllFoodProduction();
            response.put("success", true);
            response.put("data", foodList);
            response.put("message", "Get all food production information successfully");
        } catch (Exception e) {
            log.error("Failed to get all food production information", e);
            response.put("success", false);
            response.put("message", "Failed to get all food production information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Obtain all food ingredient information -- for ingredient traceability
     */
    @GetMapping("/food/ingredients/all")
    public ResponseEntity<Map<String, Object>> getAllFoodIngredients() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, Object>> ingredientList = fabricService.getAllFoodIngredients();
            response.put("success", true);
            response.put("data", ingredientList);
            response.put("message", "Get all food ingredient information successfully");
        } catch (Exception e) {
            log.error("Failed to get all food ingredient information", e);
            response.put("success", false);
            response.put("message", "Failed to get all food ingredient information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Access to all food logistics information -- for logistics traceability
     */
    @GetMapping("/food/logistics/all")
    public ResponseEntity<Map<String, Object>> getAllFoodLogistics() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, Object>> logisticsList = fabricService.getAllFoodLogistics();
            response.put("success", true);
            response.put("data", logisticsList);
            response.put("message", "Get all food logistics information successfully");
        } catch (Exception e) {
            log.error("Failed to get all food logistics information", e);
            response.put("success", false);
            response.put("message", "Failed to get all food logistics information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Get blockchain details -- including number of blocks, number of transactions per block
     */
    @GetMapping("/blockchain/details/full")
    public ResponseEntity<Map<String, Object>> getBlockchainDetailsFull() {
        Map<String, Object> response = new HashMap<>();
        try {
            org.hyperledger.fabric.sdk.Channel channel = fabricService.getChannel();
            org.hyperledger.fabric.sdk.BlockchainInfo blockchainInfo = channel.queryBlockchainInfo();
            long height = blockchainInfo.getHeight();
            java.util.List<java.util.Map<String, Object>> blocks = new java.util.ArrayList<>();
            for (long i = 0; i < height; i++) {
                org.hyperledger.fabric.sdk.BlockInfo blockInfo = channel.queryBlockByNumber(i);
                java.util.Map<String, Object> block = new java.util.HashMap<>();
                block.put("blockNumber", blockInfo.getBlockNumber());
                byte[] dataHash = blockInfo.getDataHash();
                String dataHashHex = dataHash != null ? DatatypeConverter.printHexBinary(dataHash) : null;
                block.put("dataHash", dataHashHex);
                block.put("previousHash", blockInfo.getPreviousHash() != null ? blockInfo.getPreviousHash().toString() : null);
                block.put("transactionCount", blockInfo.getTransactionCount());
                blocks.add(block);
            }
            java.util.Map<String, Object> data = new java.util.HashMap<>();
            data.put("totalBlocks", height);
            data.put("blocks", blocks);
            response.put("success", true);
            response.put("data", data);
            response.put("message", "Get blockchain details success");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Get blockchain details failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Test API - Verify Blockchain Connectivity
     */
    @GetMapping("/test/blockchain")
    public ResponseEntity<Map<String, Object>> testBlockchain() {
        Map<String, Object> response = new HashMap<>();
        try {
            // Test Blockchain Connectivity
            org.hyperledger.fabric.sdk.Channel channel = fabricService.getChannel();
            org.hyperledger.fabric.sdk.BlockchainInfo blockchainInfo = channel.queryBlockchainInfo();
            
            // Get all transactions
            java.util.List<java.util.Map<String, Object>> allTxs = fabricService.getAllTransactions();
            
            java.util.Map<String, Object> data = new java.util.HashMap<>();
            data.put("blockchainHeight", blockchainInfo.getHeight());
            data.put("totalTransactions", allTxs.size());
            data.put("sourceAppTransactions", allTxs.stream()
                .filter(tx -> tx.get("chaincodeName") != null && 
                             tx.get("chaincodeName").toString().contains("source-app"))
                .count());
            
            response.put("success", true);
            response.put("data", data);
            response.put("message", "Blockchain connection test successful");
        } catch (Exception e) {
            log.error("Blockchain connection test failed", e);
            response.put("success", false);
            response.put("message", "Blockchain connection test failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * View details of all transactions
     */
    @GetMapping("/test/transactions")
    public ResponseEntity<Map<String, Object>> getAllTransactionsDetail() {
        Map<String, Object> response = new HashMap<>();
        try {
            java.util.List<java.util.Map<String, Object>> allTxs = fabricService.getAllTransactions();
            
            java.util.List<java.util.Map<String, Object>> txDetails = new java.util.ArrayList<>();
            for (java.util.Map<String, Object> tx : allTxs) {
                java.util.Map<String, Object> detail = new java.util.HashMap<>();
                detail.put("transactionID", tx.get("transactionID"));
                detail.put("blockNumber", tx.get("blockNumber"));
                detail.put("timestamp", tx.get("timestamp"));
                detail.put("chaincodeName", tx.get("chaincodeName"));
                detail.put("type", tx.get("type"));
                detail.put("payloads", tx.get("payloads"));
                txDetails.add(detail);
            }
            
            response.put("success", true);
            response.put("data", txDetails);
            response.put("message", "Get all transaction details successfully");
        } catch (Exception e) {
            log.error("Failed to get transaction details", e);
            response.put("success", false);
            response.put("message", "Failed to get transaction details: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Check chain code status
     */
    @GetMapping("/chaincode/status")
    public ResponseEntity<Map<String, Object>> checkChaincodeStatus() {
        Map<String, Object> response = new HashMap<>();
        try {
            String status = fabricService.checkChaincodeStatus();
            response.put("success", true);
            response.put("data", status);
            response.put("message", "Chaincode status check completed");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Chaincode status check failed: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }

    /**
     * Getting channel information
     */
    @GetMapping("/channel/info")
    public ResponseEntity<Map<String, Object>> getChannelInfo() {
        Map<String, Object> response = new HashMap<>();
        try {
            String info = fabricService.getChannelInfo();
            response.put("success", true);
            response.put("data", info);
            response.put("message", "Channel information retrieved successfully");
        } catch (Exception e) {
            response.put("success", false);
            response.put("message", "Failed to get channel information: " + e.getMessage());
        }
        return ResponseEntity.ok(response);
    }
} 
