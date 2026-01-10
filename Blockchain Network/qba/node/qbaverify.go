package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
)

func (n *Node) QBAVerifyAndCommitThread() {
    for {
        select {
            case msg := <-n.QBAForwardRecv:
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
            
            n.QBAMessageBackup=append(n.QBAMessageBackup,*msg)
            var TmpMessage []byte
            
            if msg.SenderID == n.GetPrimary() {
                TmpMessage=n.CombineMessage(TmpMessage,msg.Digest)
            }
            //if sender ID isn't the ID of primary node,then current broadcast phase is phase II
            if msg.SenderID != n.GetPrimary() {
                TmpMessage=n.CombineMessage(TmpMessage,msg.Digest[:len(msg.Digest)-message.KeyLength*2])
            }
            
            //log:verifier node received message from forwarder node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Receive Forward Message From Node %v",n.id,msg.SenderID,msg.ID,n.id,msg.ID)
            
            //fetch QDS private key of verifier
            n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*msg.QdsNum-2),false)
            n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*msg.QdsNum-1),false)
            n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*msg.QdsNum),false)
            
            //output check:node forwarder fetch the QDS private key of verifier
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]verifier node %v fetch the QDS key X:%v",n.id,msg.SenderID,msg.ID,n.id,n.id,n.QBAKey_X)
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]verifier node %v fetch the QDS key Y:%v",n.id,msg.SenderID,msg.ID,n.id,n.id,n.QBAKey_Y)
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]verifier node %v fetch the QDS key Z:%v",n.id,msg.SenderID,msg.ID,n.id,n.id,n.QBAKey_Z)
            
            //get the verify result
            commitresult :=message.QDSVerify(TmpMessage,msg.QDSSignature,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z)
            
            //use commit resule to pack the commit message to send back to forwarder node
            content1, _, err1 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,commitresult,msg.QdsNum)
            
            //pack error check
            if err1 != nil {
                log.Printf(BroadCastNotify+"generate qba-commit message error")
                return
            }
            
            //send message back to forwarder node
            go SendPost(content1,n.table[msg.ID]+server.QBACommitEntry)
            
            //log:send message back to forwarder node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Finish Verify And Send Back To Node %v",n.id,msg.SenderID,msg.ID,n.id,msg.ID)
        }
    }
}

