package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
)

func (n *Node) QBAForwardThread() {
    for {
        select {
            case msg := <-n.QBAMessageRecv:

            //the primary drives the key schedule: never lag behind its round
            n.realignKeyWindow(msg.QdsNum)
            //the label to define current broadcast phase
            //the label to output
            var BroadCastNotify string
            //determine broadcast phase by sender ID in received message
            //if sender ID is the ID of primary node,then current broadcast phase is phase I
            if msg.SenderID == n.GetPrimary() {
                BroadCastNotify="[BroadCast Phast I]"
            }
            //if sender ID isn't the ID of primary node,then current broadcast phase is phase II
            if msg.SenderID != n.GetPrimary() {
                BroadCastNotify="[BroadCast Phast II]"
            }
            
            //fetch QDS private key of forwarder
                kx, ky, kz := n.Forwarder_KeyCache.CalculateKeyLines(n.ConsensusCount, msg.QdsNum,n.Forwarder_KeyCache.N)
                
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
	        
	        //then read key from file
	        if len(n.Signer_KeyCache.cacheData[0])==0{
	            result:=n.Signer_KeyCache.LoadCacheForLine(kx)
	            if result==false{
	                log.Printf("read file error")
	                return
	            }
	        }
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
                
                n.QBAKey_X = n.Forwarder_KeyCache.ReadKeyFromCache(kx)
                n.QBAKey_Y = n.Forwarder_KeyCache.ReadKeyFromCache(ky)
                n.QBAKey_Z = n.Forwarder_KeyCache.ReadKeyFromCache(kz)
            
            //pack the forward message to forward
            content1, _, err1 := message.NewQBAForwardMessage(msg.Digest,n.id,msg.SenderID,msg.QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QdsNum)
            
            //pack error check
            if err1 != nil {
                log.Printf(BroadCastNotify+"generate QBA forward message error")
                return
            }
            //choose a node to forward message
            TargetID,_:=n.FindTargetID(msg.SenderID)
            //log:receive message from signer node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Receive QBA Message From Node %v",n.id,msg.SenderID,n.id,TargetID,msg.SenderID)
            
            //output check:node forwarder fetch the QDS private key of forwarder
            log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
            
            //send the QBA forward message to chosen node
            go SendPost(content1,n.table[TargetID]+server.QBAForwardEntry)

            n.record("forward", "forwarder relayed the signature to the verifier", int(msg.SenderID), int(n.id), int(TargetID), msg.QdsNum, msg.Digest)
            
            //log:forwarder node forward message to verifier node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Forward QBA Message To Node %v",n.id,msg.SenderID,n.id,TargetID,TargetID)      
            
            //in broadcast phase I,when forwarder node receive QBA message
            //forwarder node will save the message from primary node
            //and put the message in the infer message buffer
            if msg.SenderID==n.GetPrimary() { 
                n.QBAMainMessage=make([]message.QBAMessage,0)
                n.QBAMainMessage=append(n.QBAMainMessage,*msg)
                
                //use the message send from signer node to pack the vote message 
                //put the vote message in the vote message buffer as the infer list
                _, msg3, err3 := message.NewQBAVoteMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,false)
                
                //pack error check
                if err3 != nil {
                    log.Printf(BroadCastNotify+"generate QBA forward message error")
                    return
                }
                
                //consistency check on the entry this node records for itself: two
                //different values under one signer are only possible if that signer
                //equivocated
                //put the vote message in the vote message buffer
                n.buffer.BufferQBAVoteMsg(msg3)
                
                //log:node put the QBA message in the infer list
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Join Forward Message:%v To Infer List",n.id,msg.SenderID,n.id,TargetID,n.ConvertBInaryByteToString(msg.Digest[:8])) 
            } 
            
            //in broadcast phase II
            //when forwarder node receive QBA message
            //this node put the QBA message in QBA commit message buffer as infer list
            if msg.SenderID != n.GetPrimary() { 
                //pack QBA vote message
                _, msg2, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"ready to verify",msg.QdsNum)
                
                //pack error check
                if err2 != nil {
                    log.Printf(BroadCastNotify+"generate QBA forward message error")
                    return
                }
                copy(n.QBAMainMessage[0].QDSSignature.Digest,msg.QDSSignature.Digest)
                copy(n.QBAMainMessage[0].QDSSignature.Factor,msg.QDSSignature.Factor)
                //put the vote message in QBA commit message buffer
                n.buffer.BufferQBACommitMsg(msg2)
            }      
        } 
    }
}


//func (n *Node) checkPrepareMsg(msg *message.Prepare) bool {
//	if n.view != msg.View {
//		return false
//	}
//	if !n.sequence.CheckBound(msg.Sequence) {
//		return false
//	}
//	return true
//}
