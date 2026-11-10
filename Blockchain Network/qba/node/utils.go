package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"sync"
	"os"
	"bufio"
	//"io/ioutil"
	"strings"
	//"encoding/hex"
	"strconv"
	"log"
)

const(
    KeyCacheWindowLength=100
)

// the execute op num now in state
type ExecuteOpNum struct {
	num    int
	locker *sync.RWMutex
}

func NewExecuteOpNum() *ExecuteOpNum {
	return &ExecuteOpNum{
		num:    0,
		locker: new(sync.RWMutex),
	}
}

func (n *ExecuteOpNum) Get() int {
	return n.num
}

func (n *ExecuteOpNum) Inc() {
	n.num = n.num + 1
}

func (n *ExecuteOpNum) Dec() {
	n.Lock()
	n.num = n.num - 1
	n.UnLock()
}

func (n *ExecuteOpNum) Lock() {
	n.locker.Lock()
}

func (n *ExecuteOpNum) UnLock() {
	n.locker.Unlock()
}

func (n *Node) GetPrimary() message.Identify {
	return message.Identify(0)
}

//key cache
type KeyCache struct {
    cacheData   [][]byte    
    cacheStart  int         
    cacheEnd    int        
    filename    string      
    W           int         
    N           int         
    W0          int         
    totalLines  int         
}

func (kc *KeyCache) CalculateKeyLines(c, q, N int) (kx, ky, kz int) {
    base := (c-1)*N + q
    kx = base*3 - 2
    ky = base*3 - 1  
    kz = base*3
    return kx, ky, kz
}

//read key from cache
func (kc *KeyCache) ReadKeyFromCache(lineNum int) []byte {
    //check whether linenum is in cache range
    if lineNum < kc.cacheStart || lineNum > kc.cacheEnd {
        // not in cache range,need to load key
        if !kc.LoadCacheForLine(lineNum) {
            return nil
        }
    }
    
    // compute cache index
    cacheIndex := lineNum - kc.cacheStart
    if cacheIndex < 0 || cacheIndex >= len(kc.cacheData) {
        log.Printf("cache index error: line num %d, cache range [%d, %d]", lineNum, kc.cacheStart, kc.cacheEnd)
        return nil
    }
    
    return kc.cacheData[cacheIndex]
}

//load cache from target line
func (kc *KeyCache) LoadCacheForLine(targetLine int) bool {
    // compute cache start line
    newStart := kc.calculateOptimalCacheStart(targetLine)
    
    //compute cache size(consider file boundry)
    actualCacheSize := kc.W
    if newStart + kc.W - 1 > kc.totalLines && kc.totalLines > 0 {
        actualCacheSize = kc.totalLines - newStart + 1
    }
    
    // read key from file
    keys, err := kc.readKeysFromFile(newStart, actualCacheSize)
    if err != nil {
        log.Printf("load cache error: %v", err)
        return false
    }
    //renew cache
    kc.cacheData = keys
    kc.cacheStart = newStart
    kc.cacheEnd = newStart + len(keys) - 1
    
    log.Printf("cache renew: range [%d, %d], target line %d", kc.cacheStart, kc.cacheEnd, targetLine)
    return true
}

//compute cache start
func (kc *KeyCache) calculateOptimalCacheStart(targetLine int) int {
    alignedStart := ((targetLine - 1) / kc.W) * kc.W + 1
    return alignedStart
}

//load key of target range from file
func (kc *KeyCache) readKeysFromFile(startLine, count int) ([][]byte, error) {
    file, err := os.Open(kc.filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    currentLine := 0
    var result [][]byte
    
    for scanner.Scan() {
        currentLine++
        
        // skip the line beyond target line
        if currentLine < startLine {
            continue
        }
        
        // read target line
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            result = append(result, nil)
            continue
        }
        
        bits := strings.Fields(line)
        if len(bits) != message.KeyLength {
            log.Printf("line %d format error", currentLine)
            result = append(result, nil)
            continue
        }
        
        key := make([]byte, message.KeyLength)
        for i, bitStr := range bits {
            bit, err := strconv.ParseUint(bitStr, 10, 8)
            if err != nil || bit > 1 {
                log.Printf("line %d have invalid binary byte: %s", currentLine, bitStr)
                key[i] = 0
                continue
            }
            key[i] = byte(bit)
        }
        result = append(result, key)
        
        // stop after read enough key
        if len(result) >= count {
            break
        }
    }
    
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    
    return result, nil
}
//---------------------------------------

//when signer node or forwarder node need to choose node to forward QBA message
//then node use this function to choose optional node
func (n *Node) FindTargetID(id message.Identify) (message.Identify,bool){
    //create a node list to collect all optional node
    var AvailableNode []message.Identify
    
    
    for k,_:=range n.table{
        //'k!=n.id':don't send message to self
        //'k!= n.GetPrimary()':don't send message to primary node
        //'k!=id':don't send message to signer\
        //'!n.buffer.IsExistQBACommitMsg(k)':
        //use QBA commit message buffer to represent node that already send
        //by identifying which nodes have send message from the message queue
        //then identify which node don't need to send message to
        if k!=n.id && k!= n.GetPrimary() && !n.buffer.IsExistQBACommitMsg(k) && k!=id{
            AvailableNode=append(AvailableNode,k)
        }
    }
    
    //the label that whether all optional node have sent message
    FinishCheck:=false
    
    //if there are still optional node,then set the label to false
    if(len(AvailableNode)>0){
        FinishCheck=false
    }
    
    //if there are no optional node
    //then set the label to true and set any ID output
    if(len(AvailableNode)==0){
        FinishCheck=true
        AvailableNode=append(AvailableNode,message.Identify(0))
    }
    
    return AvailableNode[0],FinishCheck
}

func (n *Node) FindTargetIDFromVoteResult(id message.Identify) (message.Identify,bool){
	var AvailableNode []message.Identify
	for k,_:=range n.table{
		if k!=n.id && k!= n.GetPrimary() && !n.buffer.IsExistQBAVoteMsg(k) && k!=id{
			AvailableNode=append(AvailableNode,k)
		}
	}
	FinishCheck:=false
	if(len(AvailableNode)>0){
		FinishCheck=false
	}
	if(len(AvailableNode)==0){
		FinishCheck=true
		AvailableNode=append(AvailableNode,message.Identify(0))
	}
	return AvailableNode[0],FinishCheck
}


func StringCalc(a string, b string) string {
	return ""
}

func (n *Node)ConvertBInaryByteToString(num []byte)string{
    result := ""
    for i:=0;i<len(num);i++{
        if num[i]==byte(0){
            result+="0"
        }
        if num[i]==byte(1){
            result+="1"
        }
    }
    return result
}

func (n *Node) ReadNodeKey(filename string, lineNum int, skipFirstBit bool) []byte {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("open file error: %v", err)
        return nil
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    currentLine := 0

    for scanner.Scan() {
        currentLine++
        if currentLine != lineNum {
            continue
        }

        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            return nil
        }

        bits := strings.Fields(line) 


        if len(bits) != message.KeyLength {
            log.Printf("line %d format error: need %d bytes, actually got %d bytes", lineNum, message.KeyLength, len(bits))
            return nil
        }

        // convert to byte array
        result := make([]byte, message.KeyLength)
        for i, bitStr := range bits {
            bit, err := strconv.ParseUint(bitStr, 10, 8)
            if err != nil || bit > 1 {
                log.Printf("line %d have invalid binary byte: %s", lineNum, bitStr)
                return nil
            }
            result[i] = byte(bit)
        }
        log.Printf("read result:%v",result)
        return result
    }

    if err := scanner.Err(); err != nil {
        log.Printf("sacn error: %v", err)
    }
    log.Printf("file dont have %d lines (max line: %d)", lineNum, currentLine)
    return nil
}


