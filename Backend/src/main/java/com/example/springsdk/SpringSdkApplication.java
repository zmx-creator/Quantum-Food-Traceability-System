package com.example.springsdk;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class SpringSdkApplication {

    public static void main(String[] args) {
        SpringApplication.run(SpringSdkApplication.class, args);
        System.out.println("=== Hyperledger Fabric Java SDK run successfully ===");
        System.out.println("App address: http://localhost:8080");
        System.out.println("API file: http://localhost:8080/api/fabric/info");
        System.out.println("================================================");
    }
} 