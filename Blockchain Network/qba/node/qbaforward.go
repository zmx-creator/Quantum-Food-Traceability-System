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
            //---------------------------------------------------------------------------------------------------------------------------------------------------
            //the label to define current broadcast phase
            //the label to output
            var BroadCastNotify string
            //determine broadcast phase by sender ID in received message
            //if sender ID is the ID of primary node,then current broadcast phase is phase I
            if msg.SenderID == n.GetPrimary() {
                BroadCastNotify="[BroadCast Phase]"
            }
            //if sender ID isn't the ID of primary node,then current broadcast phase is phase II
            if msg.SenderID != n.GetPrimary() {
                BroadCastNotify="[Circular Collection Phase]"
            }
            //---------------------------------------------------------------------------------------------------------------------------------------------------
            
            //main process
            if n.id!=message.Identify(len(n.table)-1){
                log.Printf("Check:Sender ID:%v",msg.SenderID)
                //fetch QDS private key of forwarder
                n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_forwarder.txt",(3*msg.QdsNum-2),false)
                n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_forwarder.txt",(3*msg.QdsNum-1),false)
                n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_forwarder.txt",(3*msg.QdsNum),false)
            
                //pack the forward message to forward
                content1, msg1, err1 := message.NewQBAForwardMessage(msg.Digest,n.id,msg.SenderID,msg.QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QdsNum)
            
                //pack error check
                if err1 != nil {
                    log.Printf(BroadCastNotify+"generate QBA forward message error")
                    return
                }
            
                //forward message to RELIABLE NODE:pre-set node 4 be reliable node
                TargetID:=message.Identify(len(n.table)-1)
            
                //log:receive message from signer node
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Receive QBA Message From Node %v",n.id,msg.SenderID,n.id,TargetID,msg.SenderID)
            
                //output check:node forwarder fetch the QDS private key of forwarder
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]forwarder node %v fetch the QDS key X:%v",n.id,n.CircleNodeIndex_Backward(),n.id,TargetID,n.id,n.QBAKey_X)
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]forwarder node %v fetch the QDS key Y:%v",n.id,n.CircleNodeIndex_Backward(),n.id,TargetID,n.id,n.QBAKey_Y)
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]forwarder node %v fetch the QDS key Z:%v",n.id,n.CircleNodeIndex_Backward(),n.id,TargetID,n.id,n.QBAKey_Z)
                //send the QBA forward message to chosen node
                go SendPost(content1,n.table[TargetID]+server.QBAForwardEntry)
            
                //log:forwarder node forward message to verifier node
                log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Forward QBA Message To Node %v",n.id,msg.SenderID,n.id,TargetID,TargetID)      
            
                if msg.SenderID == n.GetPrimary() {
                    n.QBAMainMessage=n.QBAMainMessage[:0]
                    n.QBAMainMessage=append(n.QBAMainMessage,*msg)
                }
                //if sender ID isn't the ID of primary node,then current broadcast phase is phase II
                if msg.SenderID != n.GetPrimary() {
                    n.QBAMessageBackup=append(n.QBAMessageBackup,*msg1)
                }
            }
            
            //verify node just backup transaction data
            //to write block when consensus succeed
            if n.id==message.Identify(len(n.table)-1){
                n.QBAMainMessage=n.QBAMainMessage[:0]
                n.QBAMainMessage=append(n.QBAMainMessage,*msg)
                //use commit resule to pack the commit message to send back to forwarder node
                content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"finish backup transaction data",msg.QdsNum)
             
                //pack error check
                if err2 != nil {
                    log.Printf(BroadCastNotify+"generate qba-commit message error")
                    return
                }
                //log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v]finish backup transaction data",n.id,msg.SenderID)
                go SendPost(content2,n.table[msg.SenderID]+server.QBACommitEntry)
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
