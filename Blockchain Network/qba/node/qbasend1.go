package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"time"
)

// send pre-prepare thread by request notify or timer
func (n *Node) QBAMessageSendOneThread() {
    // TODO change timer duration from config
    duration := time.Second
    timer := time.After(duration)
    for {
        select {
            // recv request or time out
            case <-n.QBAMessageSendNotify:
            if n.id == n.GetPrimary(){
                n.QBAMessageSendOneHandleFunc()
            }
        case <-timer:
            timer = nil
            if n.id == n.GetPrimary(){
                n.QBAMessageSendOneHandleFunc()
            }
            timer = time.After(duration)
        }
    }
}

func (n *Node) QBAMessageSendOneHandleFunc(){
    if n.id == n.GetPrimary(){
	if n.buffer.SizeofRequestQueue() < 1 {
		return
	}
	// batch request to discard network traffic
	batch := n.buffer.BatchRequest()
	if len(batch) < 1 {
		return
	}
	
	// use QDS Count to read different QDS private key
	QDSNum := (len(n.table)*((len(n.table)-3))+2)*(n.ConsensusCount-1)+1
	
	//label:Start QDS:QDS time
	log.Printf("QDS time:%v",QDSNum)
	
	//fetch QDS private key of signer
	n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*QDSNum-2),false)
	n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*QDSNum-1),false)
	n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*QDSNum),false)
	
	//test message
	//testmessage := []byte{1, 0, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1}
        
	//fetch irreducible polynomial factor
	factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",QDSNum ,true)
	
	//use message digest to send in QBA broadcast round
	//Digest need to convert to byte array and modify to a proper length which is smaller than const 'MessageLength'
	transaction := message.Message{Requests:batch}
	messageDigest,_ := message.Digest(transaction)
	
	//get QDS signature
	sig := message.Sign(messageDigest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
	
	//log:Start BroadCast Phase I
        log.Printf("[BroadCast Phase]Start BroadCast Phase")
	
	//packing the QBA Message to Send in QBA broadcast round
        content, msg, _, err := message.NewQBAMessage(n.id,batch,sig,QDSNum)  
                
        //pack error check
        if err != nil {
            log.Printf("[BroadCast Phast I]generate New QBA message error")
            return
        }
        
        //output check:node signer fetch the QDS private key of signer
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key X:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_X)
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Y:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_Y)
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Z:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_Z)
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the factor of irreducible polynomial:%v",n.id,n.id,message.Identify(1),n.id,factor)
        
        
        //output check:node primary send QBA Message to start QBA broadcast round
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]Send QBA Message To Node %v",n.id,n.id,message.Identify(1),message.Identify(1))
        
        //send QBA Message
        go SendPost(content, n.table[message.Identify(1)] + server.QBAMessageEntry) 
        
        //save the QBA Message which already send,to perpare for next send
        //clear the saved QBA message that may be the QBA message which saved in last QBA Consensus
        n.QBAMainMessage=make([]message.QBAMessage,0)
        
        //saved QBA Message
        n.QBAMainMessage=append(n.QBAMainMessage,*msg)
    }	
}



