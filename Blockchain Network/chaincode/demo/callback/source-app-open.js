'use strict';

const logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger("Test");

const contractID = "source-app";
const contractVer = "1.0";

let bc, contx;
let foodIDs = []; 

module.exports.init = async (blockchain, context, args) => {
    bc = blockchain;
    contx = context;
};

module.exports.run = async() => {
   
    let foodID = "food_" + Math.random().toString(36).substr(2, 9);
    foodIDs.push(foodID);
    
  
    let txArgs = {
        chaincodeFunction: "addProInfo",
        chaincodeArguments: [
            foodID,                    // FoodID
            `TESTFOOD_${Date.now()}`,  // FoodName
            '500g',                // FoodSpec
            '2025-09-28',             // FoodMFGDate
            '2026-09-28',             // FoodEXPDate
            `BATCH${Date.now()}`,     // FoodLOT
            'QS123456789',            // FoodQSID
            'Supermarket',        // FoodMFRSName
            '25.50',                  // FoodProPrice
            'Beijing'                  // FoodProPlace
        ]
    };

    return bc.invokeSmartContract(contx, contractID, contractVer, txArgs, 30000);
};

module.exports.end = async() => { 
};


module.exports.contractID = contractID;
module.exports.contractVer = contractVer;
module.exports.foodIDs = foodIDs;
