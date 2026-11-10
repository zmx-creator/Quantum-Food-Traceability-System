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
		chaincodeFunction: "addLogInfo",
		chaincodeArguments: [
			foodID,                        // FoodID
			'2025-09-28 08:00:00',        // LogDepartureTm
			'2025-09-28 16:00:00',        // LogArrivalTm
			Math.random() > 0.5 ? 'Transportation' : 'Storage', // LogMission
			'Beijing Warehouse',          // LogDeparturePl
			'Shanghai Warehouse',         // LogDest
			'Test Supermarket',           // LogToSeller
			'24 hours',                   // LogStorageTm
			'Cold Chain Transportation',  // LogMOT
			'Test Logistics Company',     // LogCopName
			'150.00'                      // LogCost
		]
	};

    return bc.invokeSmartContract(contx, G.contractID, G.contractVer, txArgs, 30000);
};

module.exports.end = async() => {
};
