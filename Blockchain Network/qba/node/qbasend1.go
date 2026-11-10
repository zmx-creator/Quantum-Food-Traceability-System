package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"time"
)

// send pre-prepare thread by request notify or timer
func (n *Node) QBAMessageSendOneThread() {
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
	if !n.IsConsensusRunning {
	    n.IsConsensusRunning=true
	}else if n.IsConsensusRunning {
		return
	}
	// the round is now live: arm the liveness watchdog for this instance
	n.round.Begin(n.ConsensusCount, time.Now())

	// batch request to discard network traffic
	batch := n.buffer.BatchRequest()
	if len(batch) < 1 {
		return
	}
	
	// use QDS Count to read different QDS private key
	QDSNum := n.ConsensusCount*36-35
	log.Printf("QDS time:%v",QDSNum)
	
	//fetch QDS private key of signer
	kx, ky, kz := n.Signer_KeyCache.CalculateKeyLines(n.ConsensusCount, QDSNum,n.Signer_KeyCache.N)
	
	//if the key cache is insufficient,the keys will be read from file
	//first clear the currend key cache
	if kx>n.Signer_KeyCache.cacheEnd {
	    n.Signer_KeyCache.cacheData=make([][]byte,3*KeyCacheWindowLength*(len(n.table)*(len(n.table)-3)+2)) 
	}
	if kx>n.Forwarder_KeyCache.cacheEnd {
	    n.Forwarder_KeyCache.cacheData=make([][]byte,3*KeyCacheWindowLength*(len(n.table)*(len(n.table)-3)+2)) 
	}
	if kx>n.Verifier_KeyCache.cacheEnd {
	    n.Verifier_KeyCache.cacheData=make([][]byte,3*KeyCacheWindowLength*(len(n.table)*(len(n.table)-3)+2)) 
	}
	if len(n.Signer_KeyCache.cacheData[0])==0{
	    result:=n.Signer_KeyCache.LoadCacheForLine(kx)
	    if result==false{
	        log.Printf("read file error")
	        return
	    }
	}
	
	//then read key from file
	if len(n.Forwarder_KeyCache.cacheData[0])==0{
	    result:=n.Forwarder_KeyCache.LoadCacheForLine(kx)
	    if result==false{
	        log.Printf("read file error")
	        return
	    }
	}
	if len(n.Verifier_KeyCache.cacheData[0])==0{
	    result:=n.Verifier_KeyCache.LoadCacheForLine(kx)
	    if result==false{
	        log.Printf("read file error")
	        return
	    }
	}
	log.Printf("Check:key index:%v %v %v",kx,ky,kz)
	n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
	n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
	n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	

	//fetch irreducible polynomial factor
	factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",QDSNum ,true)
	
	//use message digest to send in QBA broadcast round
	//Digest need to convert to byte array and modify to a proper length which is smaller than const 'MessageLength'
	transaction := message.Message{Requests:batch}
	messageDigest,_ := message.Digest(transaction)
	
	//get QDS signature
	sig := message.Sign(messageDigest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
	
	//log:Start BroadCast Phase I
        log.Printf("[BroadCast Phast I]Start BroadCast Phase I")
	
	//packing the QBA Message to Send in QBA broadcast round
        content, msg, _, err := message.NewQBAMessage(n.id,batch,sig,QDSNum)  
                
        //pack error check
        if err != nil {
            log.Printf("[BroadCast Phast I]generate New QBA message error")
            return
        }
        
        //arm the watchdog with the exact bytes on the wire, so a retry retransmits
        //them unchanged rather than producing a second signature over new content
        n.roundContent = content
        n.roundTarget = message.Identify(1)


        //QDS time
       
        
        //output check:node signer fetch the QDS private key of signer
        log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
        log.Printf("irreducible polynomial factor loaded (value not logged)")
        
        //output check:node primary send QBA Message to start QBA broadcast round
        log.Printf("[BroadCast Phast I][cur_node:%v,signer:%v,forward:%v]Send QBA Message To Node %v",n.id,n.id,message.Identify(1),message.Identify(1))
        
        //send QBA Message
        go SendPost(content, n.table[message.Identify(1)] + server.QBAMessageEntry) 

        n.record("sign", "signer produced the QDS signature for the batch", int(n.id), 1, -1, QDSNum, msg.Digest)
        
        //save the QBA Message which already send,to perpare for next send
        //clear the saved QBA message that may be the QBA message which saved in last QBA Consensus
        n.QBAMainMessage=make([]message.QBAMessage,0)
        
        //saved QBA Message
        n.QBAMainMessage=append(n.QBAMainMessage,*msg)
    }	
}


