'use strict';

const G = require("./source-app-open.js");

const logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger("Test");

let bc, contx;
let index = 0;

module.exports.init = async (blockchain, context, args) => {
    bc = blockchain;
    contx = context;
    index = 0;
};

module.exports.run = async() => {
    if (G.foodIDs.length === 0) {
        throw new Error("No food IDs available. Run open test first.");
    }
    
    let foodID = G.foodIDs[index % G.foodIDs.length]; 
    index++;
    

    let txArgs = {
        chaincodeFunction: "addIngInfo",
        chaincodeArguments: [
            foodID,
            `ING${Date.now()}1`, 
            'apple',             
            `ING${Date.now()}2`, 
            'suger'              
        ]
    };

    return bc.invokeSmartContract(contx, G.contractID, G.contractVer, txArgs, 30000);
};

module.exports.end = async() => {
};
