#!/bin/bash

# Change to the specified directory
cd ./Backend

# Clean and compile
echo "Compiling the Spring Boot project"
mvn clean compile

# Start the network
echo "Starting the network"
mvn spring-boot:run
