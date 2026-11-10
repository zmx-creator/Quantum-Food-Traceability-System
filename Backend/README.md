# Backend 

## Contents
```
    Backend/                 
    ├── src/
    │   ├── main/
    │   │   ├── java/com/example/springsdk/
    │   │   │   ├── config/        
    │   │   │   ├── controller/    
    │   │   │   ├── model/           
    │   │   │   ├── service/      
    │   │   │   └── SpringSdkApplication.java      
    │   │   └── resources/
    │   │       ├── application.properties         
    │   │       └── crypto-config/                 
    │   └── test/
    ├── pom.xml                                     
    ├── mvnw 
    ├── README.md                      
    └── mvnw.cmd                                   
```

## Requirements
- JDK 1.8
- Maven 3.x

This system was developed by using Springboot with java language. To run this system, java8 is needed.

## User guide

The system has four functions:
- Connect to Fabric network
- Submit product information
- Query transaction records
- Query food traceability

To run our system, you need to start the Fabric network first, then run:

```
# Change to the specified directory
cd ./Backend

# Clean and compile
mvn clean compile

# Start the network
mvn spring-boot:run

```

If the blockchain network runs on another machine, point the backend at it with the
`fabric.vm-ip` system property or the `FABRIC_VM_IP` environment variable. Both default
to `localhost`, where the deployment publishes the five orderers on 6050-6054 and the
peer on 7051.

```FabricConfig.java
private static final String VM_IP = System.getProperty("fabric.vm-ip", /* FABRIC_VM_IP */ "localhost");
private static final String ORDERER0_URL = "grpc://" + VM_IP + ":6050";   // orderer0 .. orderer4
private static final String PEER_URL = "grpc://" + VM_IP + ":7051";
```


