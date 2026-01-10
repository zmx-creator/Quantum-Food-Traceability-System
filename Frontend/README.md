# Frontend

## Contents
```
 Frontend/                    
   ├── public/                 
   ├── src/
   │   ├── api/                 
   │   ├── assets/         
   │   ├── components/        
   │   ├── router/           
   │   ├── views/           
   │   ├── App.vue            
   │   └── main.js           
   ├── index.html            
   ├── package.json          
   ├── vue.config.js          
   └── README.md    
            
```                            

## Requirements

- Node.js >= 14.0.0
- npm >= 6.0.0 
This system was developed by using Vue with node.js language.

## User guide

The system has four functions:

- Display blockchain information
- Submit product information
- Query transaction records
- Query food traceability

To run our system, you need to start the Fabric network and run springboot(Backend folder) first. Then run:
```
# Change to the specified directory
cd ./Frontend

# Clean and compile
npm run serve
```
If running for the first time, you may need to install dependencies first
```
    cd ./Frontend
    npm install
```

If your system is not running on localhost, then in **Frontend/src/api/fabric.js** file you may need to slightly modify the code.





