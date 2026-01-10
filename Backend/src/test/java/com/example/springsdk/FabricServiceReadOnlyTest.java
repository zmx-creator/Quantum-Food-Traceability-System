package com.example.springsdk;

import com.example.springsdk.service.FabricService;
import org.apache.commons.lang3.StringUtils;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@ActiveProfiles("test")
public class FabricServiceReadOnlyTest {

    @Autowired
    private FabricService fabricService;

    private static final String TEST_FOOD_ID = "TEST001";
    private static final String WORKFLOW_FOOD_ID = "WORKFLOW001";

    @Test
    public void testGetBlockchainInfo() {
        try {
            String info = fabricService.getBlockchainInfo();
            System.out.println("Blockchain information: " + info);
            assertNotNull(info);
            assertTrue(info.contains("Latest block height"));
        } catch (Exception e) {
            System.err.println("Failed to get blockchain information: " + e.getMessage());
            fail("Blockchain information test failed");
        }
    }

    @Test
    public void testGetFoodInfo() {
        try {
            String result = fabricService.getFoodInfo(TEST_FOOD_ID);
            System.out.println("Food information query result: " + result);
            assertNotNull(result);
        } catch (Exception e) {
            System.err.println("Failed to get food information: " + e.getMessage());
            fail("Food information test failed");
        }
    }

    @Test
    public void testGetProInfo() {
        try {
            String result = fabricService.getProInfo(TEST_FOOD_ID);
            System.out.println("Production information query result: " + result);
            assertNotNull(result);
        } catch (Exception e) {
            System.err.println("Failed to get production information: " + e.getMessage());
            fail("Production information test failed");
        }
    }

    @Test
    public void testGetIngInfo() {
        try {
            String result = fabricService.getIngInfo(TEST_FOOD_ID);
            System.out.println("Ingredient information query result: " + result);
            assertNotNull(result);
        } catch (Exception e) {
            System.err.println("Failed to get ingredient information: " + e.getMessage());
            fail("Ingredient information test failed");
        }
    }

    @Test
    public void testGetLogInfo() {
        try {
            String result = fabricService.getLogInfo(TEST_FOOD_ID);
            System.out.println("Logistics information query result: " + result);
            assertNotNull(result);
        } catch (Exception e) {
            System.err.println("Failed to get logistics information: " + e.getMessage());
            fail("Logistics information test failed");
        }
    }

    @Test
    public void testGetLogInfoList() {
        try {
            String result = fabricService.getLogInfoList(TEST_FOOD_ID);
            System.out.println("Logistics information list query result: " + result);
            assertNotNull(result);
        } catch (Exception e) {
            System.err.println("Failed to get logistics information list: " + e.getMessage());
            fail("Logistics information list test failed");
        }
    }

    @Test
    public void testReadOnlyWorkflow() {
        try {
            System.out.println("\n" + StringUtils.repeat("=", 50));
            System.out.println("Starting read-only test");
            System.out.println(StringUtils.repeat("=", 50));
            
            // 1. Get blockchain information
            System.out.println("\nGetting blockchain information");
            String blockchainInfo = fabricService.getBlockchainInfo();
            System.out.println("Blockchain information: " + blockchainInfo);
            
            // 2. Query complete food information
            System.out.println("\nQuerying complete food information");
            String foodInfo = fabricService.getFoodInfo(WORKFLOW_FOOD_ID);
            System.out.println("Complete food information: " + foodInfo);
            
            // 3. Query various types of information separately
            System.out.println("\nQuerying production information");
            String proInfo = fabricService.getProInfo(WORKFLOW_FOOD_ID);
            System.out.println("Production information: " + proInfo);
            
            System.out.println("\nQuerying ingredient information");
            String ingInfo = fabricService.getIngInfo(WORKFLOW_FOOD_ID);
            System.out.println("Ingredient information: " + ingInfo);
            
            System.out.println("\nQuerying logistics information");
            String logInfo = fabricService.getLogInfo(WORKFLOW_FOOD_ID);
            System.out.println("Logistics information: " + logInfo);
            
            System.out.println("\n" + StringUtils.repeat("=", 50));
            System.out.println("DONE!");
            System.out.println(StringUtils.repeat("=", 50) + "\n");
            
            // Verify results
            assertNotNull(blockchainInfo);
            assertNotNull(foodInfo);
            assertNotNull(proInfo);
            assertNotNull(ingInfo);
            assertNotNull(logInfo);
            
        } catch (Exception e) {
            System.err.println("Read-only test failed: " + e.getMessage());
            e.printStackTrace();
            fail("Read-only test failed");
        }
    }
}
