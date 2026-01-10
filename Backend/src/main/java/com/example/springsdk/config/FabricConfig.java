package com.example.springsdk.config;

import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.PrivateKey;
import java.util.Properties;
import java.util.Set;

@Configuration
public class FabricConfig {

    // IP of ubuntu
    //private static final String VM_IP = "192.168.225.128"; //modify according to the actual IP
    private static final String VM_IP = "localhost"; // modify according to the actual IP
        // //network config
        // private static final String ORDERER_URL = "grpc://" + VM_IP + ":7050";
        // private static final String PEER_URL = "grpc://" + VM_IP + ":7051";
    //network config
    // Orderer config - 5 orderer nodes
    private static final String ORDERER0_URL = "grpc://" + VM_IP + ":6050"; // orderer0.example.com:6050
    private static final String ORDERER1_URL = "grpc://" + VM_IP + ":6051"; // orderer1.example.com:6051
    private static final String ORDERER2_URL = "grpc://" + VM_IP + ":6052"; // orderer2.example.com:6052
    private static final String ORDERER3_URL = "grpc://" + VM_IP + ":6053"; // orderer3.example.com:6053
    private static final String ORDERER4_URL = "grpc://" + VM_IP + ":6054";

    private static final String PEER_URL = "grpc://" + VM_IP + ":7051"; // peer0.orga.com:7051

    //Because static authentication is used when authenticating the certificate in adminUser(), the pre-generated certificate file is used directly - in the crypto-config folder
    private static final String CA_URL = "http://" + VM_IP + ":7054";

    
     // organization config
     private static final String MSP_ID = "Org1MSP";
     private static final String CHANNEL_NAME = "mychannel";
    private static final String CHAINCODE_NAME = "source-app";
    private static final String CHAINCODE_VERSION = "1.0";
//    // organization config
//    private static final String MSP_ID = "OrgAMSP";
//    private static final String CHANNEL_NAME = "mychannel";
//    private static final String CHAINCODE_NAME = "source-app";
//    private static final String CHAINCODE_VERSION = "1.0";
//

     //path of crypto-config
     private static final String CRYPTO_CONFIG_PATH = "src/main/resources/crypto-config";
     private static final String ADMIN_CERT_PATH = CRYPTO_CONFIG_PATH + "/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem";
     private static final String ADMIN_KEY_DIR = CRYPTO_CONFIG_PATH + "/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/";

//    // path of crypto-config
//    private static final String CRYPTO_CONFIG_PATH = "src/main/resources/crypto-config";
//    private static final String ADMIN_CERT_PATH = CRYPTO_CONFIG_PATH + "/peerOrganizations/orga.com/users/Admin@orga.com/msp/signcerts/Admin@orga.com-cert.pem";
//    private static final String ADMIN_KEY_DIR = CRYPTO_CONFIG_PATH + "/peerOrganizations/orga.com/users/Admin@orga.com/msp/keystore/";

    static {
        System.out.println("=== Loading FabricConfig class ===");
        System.out.println("VM_IP: " + VM_IP);
        System.out.println("MSP_ID: " + MSP_ID);
        System.out.println("CHANNEL_NAME: " + CHANNEL_NAME);
        System.out.println("ADMIN_CERT_PATH: " + ADMIN_CERT_PATH);
        System.out.println("ADMIN_KEY_DIR: " + ADMIN_KEY_DIR);
    }

    @Bean
    public HFClient hfClient(FabricUser adminUser) throws Exception {
        System.out.println("=== Creating HFClient ===");
        HFClient client = HFClient.createNewInstance();
        client.setCryptoSuite(CryptoSuite.Factory.getCryptoSuite());
        client.setUserContext(adminUser);
        System.out.println("HFClient created successfully");
        return client;
    }

    @Bean
    public FabricUser adminUser() throws Exception {
        try {
            System.out.println("=== Loading Admin user ===");
            System.out.println("Certificate Path: " + ADMIN_CERT_PATH);
            System.out.println("Private Key Directory: " + ADMIN_KEY_DIR);
            
            // Check if the certificate exists
            File certFile = new File(ADMIN_CERT_PATH);
            if (!certFile.exists()) {
                throw new RuntimeException("Certificate does not exist: " + ADMIN_CERT_PATH);
            }
//            System.out.println("Certificate exists，size: " + certFile.length() + " bytes");
            
            // Read certificate
            String cert = new String(Files.readAllBytes(Paths.get(ADMIN_CERT_PATH)), "UTF-8");
//            System.out.println("First 100 characters of certificate content: " + cert.substring(0, Math.min(100, cert.length())));
            
            // Read private key
            File keyDir = new File(ADMIN_KEY_DIR);
            if (!keyDir.exists()) {
                throw new RuntimeException("Failed to find private key directory: " + ADMIN_KEY_DIR);
            }
            File[] keyFiles = keyDir.listFiles((dir, name) -> name.endsWith("_sk"));
            if (keyFiles == null || keyFiles.length == 0) {
                throw new RuntimeException("Failed to find private key file");
            }
//            System.out.println("Successfully found private key file: " + keyFiles[0].getName());
            
            String keyPem = new String(Files.readAllBytes(keyFiles[0].toPath()), "UTF-8");
            System.out.println("First 100 characters of private key: " + keyPem.substring(0, Math.min(100, keyPem.length())));
            
            PrivateKey privateKey = FabricUser.getPrivateKeyFromString(keyPem);
//            System.out.println("<UNK>: " + privateKey);
            System.out.println("Successfully resolved private key");

             // //return new FabricUser("Admin@org1.example.com", MSP_ID, cert, privateKey);
            // return new FabricUser("Admin@orga.com", MSP_ID, cert, privateKey);
            FabricUser user = new FabricUser("Admin@orga.com", MSP_ID, cert, privateKey);
            System.out.println("Successfully load Admin user: " + user.getName() + ", MSP: " + user.getMspId());
            return user;
        } catch (Exception e) {
            System.err.println("Failed to load Admin user: " + e.getMessage());
            e.printStackTrace();
            throw e;
        }
    }

    @Bean
    public Channel channel(HFClient client) throws Exception {
        try {
            System.out.println("=== Creating Channel ===");
              Peer peer = client.newPeer("peer0.org1.example.com", PEER_URL);
            System.out.println("Successfully created Peer node: peer0.org1.example.com");
            // Orderer orderer = client.newOrderer("orderer.example.com", ORDERER_URL);           

//            Peer peer = client.newPeer("peer0.orga.com", PEER_URL);
//            System.out.println("Successfully created Peer: peer0.orga.com");
            

            //5 Orderer nodes
            Orderer orderer0 = client.newOrderer("orderer0.example.com", ORDERER0_URL);
            Orderer orderer1 = client.newOrderer("orderer1.example.com", ORDERER1_URL);
            Orderer orderer2 = client.newOrderer("orderer2.example.com", ORDERER2_URL);
            Orderer orderer3 = client.newOrderer("orderer3.example.com", ORDERER3_URL);
            Orderer orderer4 = client.newOrderer("orderer4.example.com", ORDERER3_URL);
            System.out.println("Successfully created Orderer nodes");
            
            System.out.println("Trying to connect the channel: " + CHANNEL_NAME);
            Channel channel = client.newChannel(CHANNEL_NAME);
            channel.addPeer(peer);
            
            // channel.addOrderer(orderer);
            // Add all orderer nodes
            channel.addOrderer(orderer0);
            channel.addOrderer(orderer1);
            channel.addOrderer(orderer2);
            channel.addOrderer(orderer3);
            channel.addOrderer(orderer4);
            
            System.out.println("Satrt to initialize channel...");
            channel.initialize();
            System.out.println("Successfully connected channel: " + CHANNEL_NAME);
            
            return channel;
        } catch (Exception e) {
            System.err.println("Failed to create channel: " + e.getMessage());
            e.printStackTrace();
            throw e;
        }
    }

    public static String getChaincodeName() {
        return CHAINCODE_NAME;
    }
    public static String getChaincodeVersion() {
        return CHAINCODE_VERSION;
    }
    public static String getChannelName() {
        return CHANNEL_NAME;
    }
}

// implementation of FabricUser
class FabricUser implements User {
    private String name;
    private String mspId;
    private String cert;
    private PrivateKey privateKey;
    public FabricUser(String name, String mspId, String cert, PrivateKey privateKey) {
        this.name = name;
        this.mspId = mspId;
        this.cert = cert;
        this.privateKey = privateKey;
    }
    @Override
    public String getName() { return name; }
    @Override
    public Set<String> getRoles() { return null; }
    @Override
    public String getAccount() { return null; }
    @Override
    public String getAffiliation() { return null; }
    @Override
    public Enrollment getEnrollment() {
        return new Enrollment() {
            @Override
            public PrivateKey getKey() { return privateKey; }
            @Override
            public String getCert() { return cert; }
        };
    }
    @Override
    public String getMspId() { return mspId; }

    //Tool method: Get PrivateKey from PEM string
    public static PrivateKey getPrivateKeyFromString(String keyPem) throws Exception {
        String privateKeyPEM = keyPem.replace("-----BEGIN PRIVATE KEY-----", "")
                .replace("-----END PRIVATE KEY-----", "")
                .replaceAll("\\s", "");
        byte[] encoded = java.util.Base64.getDecoder().decode(privateKeyPEM);
        java.security.spec.PKCS8EncodedKeySpec keySpec = new java.security.spec.PKCS8EncodedKeySpec(encoded);
        java.security.KeyFactory kf = java.security.KeyFactory.getInstance("EC");
        return kf.generatePrivate(keySpec);
    }
} 
