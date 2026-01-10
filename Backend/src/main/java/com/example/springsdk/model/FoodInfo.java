package com.example.springsdk.model;

import lombok.Data;
import java.util.List;

@Data
public class FoodInfo {
    private String foodID;                    
    private ProInfo foodProInfo;              
    private List<IngInfo> foodIngInfo;        
    private LogInfo foodLogInfo;              
}

@Data
class ProInfo {
    private String foodName;                  
    private String foodSpec;                  
    private String foodMFGDate;               
    private String foodEXPDate;               
    private String foodLOT;                   
    private String foodQSID;                  
    private String foodMFRSName;              
    private String foodProPrice;              
    private String foodProPlace;              
}

@Data
class IngInfo {
    private String ingID;                     
    private String ingName;                   
}

@Data
class LogInfo {
    private String logDepartureTm;            
    private String logArrivalTm;              
    private String logMission;                
    private String logDeparturePl;            
    private String logDest;                   
    private String logToSeller;               
    private String logStorageTm;              
    private String logMOT;                    
    private String logCopName;                
    private String logCost;                   
} 