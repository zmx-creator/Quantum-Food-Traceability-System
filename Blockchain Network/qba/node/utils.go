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

//combine two byte array
func (n *Node)CombineMessage(a, b []byte) []byte {
    result := make([]byte, len(a)+len(b))
    copy(result, a)
    copy(result[len(a):], b)
    return result
}

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

//in the circular collection phase
//use this function to target next node to send QBA message
//the logic for fing the next node:
//node 1-->node 2-->......-->node (n-2)-->node 1
//target next node in this circular manner
//and in this demo,presuppose node(n-1) be verifier node
func (n *Node)CircleNodeIndex()message.Identify{
    if n.id!=message.Identify(len(n.table)-2){
        return n.id+1
    }else{
        return message.Identify(1)
    }
}

//the inverse process of funcion CircleNodeIndex()
func (n *Node)CircleNodeIndex_Backward()message.Identify{
    if n.id!=message.Identify(1){
        return n.id-1
    }else{
        return message.Identify(len(n.table)-2)
    }
}

//read the bits of line lineNum from file filename
//the use of skipFirstbool:
//the first bit of irreducible polynomial is 1 and its useless
//so skip this bit when reading files
//even if you have deleted this bit,you just need to let skipFirstBit be false
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

        if skipFirstBit && len(bits) == 65 {
            bits = bits[1:] // discard first bit
        }

        // verify length
        if len(bits) != message.KeyLength {
            log.Printf("The format of line %d incorrect: require %d bits，but there are actually %d bits", lineNum, message.KeyLength, len(bits))
            return nil
        }

        // changt to byte array
        result := make([]byte, message.KeyLength)
        for i, bitStr := range bits {
            bit, err := strconv.ParseUint(bitStr, 10, 8)
            if err != nil || bit > 1 {
                log.Printf("line %d contain invalid bit: %s", lineNum, bitStr)
                return nil
            }
            result[i] = byte(bit)
        }
        log.Printf("read result:%v",result)
        return result
    }

    if err := scanner.Err(); err != nil {
        log.Printf("scan error: %v", err)
    }
    log.Printf("the file doesn't have line %d (max line: %d)", lineNum, currentLine)
    return nil
}

