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

If your system is not running on localhost, then in ** /Backend/src/main/java/com/example/springsdk/config/FabricConfig.java ** file you may need to slightly modify the code.

```FabricConfig.java
private static final String VM_IP = "192.168.61.128";
private static final String ORDERER_URL = "grpc://" + VM_IP + ":7050";
private static final String PEER_URL = "grpc://" + VM_IP + ":7051";
```


