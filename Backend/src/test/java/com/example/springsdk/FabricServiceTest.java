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
public class FabricServiceTest {

    @Autowired
    private FabricService fabricService;

    // private static final String TEST_FOOD_ID = "TEST001";

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
    public void testTransactionSubmission() {
        try {
            System.out.println("\n" + StringUtils.repeat("=", 60));
            System.out.println("Transaction submission function test");
            System.out.println(StringUtils.repeat("=", 60));
            

            String testId = "TX_TEST_" + System.currentTimeMillis();
            System.out.println("Test ID: " + testId);
            

            System.out.println("\nStep 1: Test production information submission");
            String proResult = fabricService.addProInfo(
                testId,
                "APPLE_TEST",
                "500g",
                "2025-06-23",
                "2025-07-23",
                "LOT_TEST" + testId,
                "QS_TEST" + testId,
                "MingCampany_TEST",
                "CNY 5.00",
                "Liaoning"
            );
            System.out.println("Production information submission result: " + proResult);
            assertNotNull(proResult);
            assertTrue(proResult.contains("Transaction successfully submitted"));
            
    
            System.out.println("Waiting for transaction confirmation...");
            Thread.sleep(3000);
            

            System.out.println("\nStep 2: Verify production information");
            String proInfo = fabricService.getProInfo(testId);
            System.out.println("Queried production information: " + proInfo);
            assertNotNull(proInfo);
            

            System.out.println("\nStep 3: Test ingredient information submission");
            String ingResult = fabricService.addIngInfo(
                testId,
                "1", "apple",
                "2", "sugar",
                "3", "lemon acid"
            );
            System.out.println("Ingredient information submission result: " + ingResult);
            assertNotNull(ingResult);
            assertTrue(ingResult.contains("Transaction successfully submitted"));
            
 
            System.out.println("Waiting for transaction confirmation...");
            Thread.sleep(3000);
            

            System.out.println("\nStep 4: Verify ingredient information");
            String ingInfo = fabricService.getIngInfo(testId);
            System.out.println("Queried ingredient information: " + ingInfo);
            assertNotNull(ingInfo);
            

            System.out.println("\nStep 5: Test logistics information submission");
            String logResult = fabricService.addLogInfo(
                testId,
                "08:00 AM",
                "10:00 AM",
                "transport",
                "Liaoning",
                "Nanning",
                "Nanning Supermarket",
                "2 hours",
                "Truck",
                "SF Logistics",
                "CNY 50.00"
            );
            System.out.println("Logistics information submission result: " + logResult);
            assertNotNull(logResult);
            assertTrue(logResult.contains("Transaction successfully submitted"));
            
 
            System.out.println("Waiting for transaction confirmation...");
            Thread.sleep(3000);
            

            System.out.println("\nStep 6: Verify logistics information");
            String logInfo = fabricService.getLogInfo(testId);
            System.out.println("Queried logistics information: " + logInfo);
            assertNotNull(logInfo);
            
  
            System.out.println("\nStep 7: Verify complete food information");
            String foodInfo = fabricService.getFoodInfo(testId);
            System.out.println("Complete food information: " + foodInfo);
            assertNotNull(foodInfo);
            

            System.out.println("\nStep 8: Test logistics information list");
            String logListInfo = fabricService.getLogInfoList(testId);
            System.out.println("Logistics information list: " + logListInfo);
            assertNotNull(logListInfo);
            
            System.out.println("\n" + StringUtils.repeat("=", 60));
            System.out.println("Transaction submission function test completed! All transactions were successfully submitted and verified!");
            System.out.println(StringUtils.repeat("=", 60) + "\n");
            
        } catch (Exception e) {
            System.err.println("Transaction submission function test failed: " + e.getMessage());
            e.printStackTrace();
            fail("Transaction submission function test failed: " + e.getMessage());
        }
    }


}
