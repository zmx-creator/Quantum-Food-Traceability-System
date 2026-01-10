package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"log"
	"encoding/json"
)

// send pre-prepare thread by request notify or timer
func (n *Node) QBAMessageSendTwoThread() {
    for {
        select {
            case msg:=<-n.QBAPrepareRecv:
                //log.Printf("Receive Notice From Node %v",msg.SenderID)
                if n.id != n.GetPrimary(){
                    //QDS counter add 1,mean that start next QDS
                    qdsnum := msg.QdsNum+1
                    log.Printf("QDS time:%v",qdsnum)
                    //fetch QDS private key of signer
                    n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-2),false)
                    n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-1),false)
                    n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum),false)
	
                    //fetch irreducible polynomial factor
                    factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	
	            
                    digest:=n.CombineMessage(n.QBAMainMessage[0].Digest,n.QBAMainMessage[0].QDSSignature.Digest)
                    digest =n.CombineMessage(digest,n.QBAMainMessage[0].QDSSignature.Factor)
                    //get QDS signature
                    sig := message.Sign(digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
                    digest =n.CombineMessage(digest,sig.Digest)
                    digest =n.CombineMessage(digest,sig.Factor)
                    //choose same node as the forwarder node to send QBA message
                    TargetID:= n.CircleNodeIndex()
        
                    //output check:node signer fetch the QDS private key of signer
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key X:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_X)
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Y:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_Y)
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Z:%v",n.id,n.id,message.Identify(1),n.id,n.QBAKey_Z)
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the factor of irreducible polynomial:%v",n.id,n.id,message.Identify(1),n.id,factor)
                    
                    //use the pre-saved QBA message to pack new QBA message
                    SendQBAMessage := message.QBAMessage{
                        Digest:digest,  
                        SenderID:n.id, 
                        QDSSignature: message.QDSSignature{
                            Digest:sig.Digest,
                            Factor:sig.Factor,
                        },  
                        Message:n.QBAMainMessage[0].Message,
                        QdsNum:qdsnum,
                    }
                    content,err := json.Marshal(SendQBAMessage)
                    //pack error check
                    if err != nil {
                        log.Printf("[Circular Collection Phase]generate QBA forward message error")
                        return
                    }
                    
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]node %v start circle send process",n.id,n.id,TargetID,n.id) 
                    
                    //send QBA message
                    go SendPost(content, n.table[TargetID] + server.QBAMessageEntry) 
                    
                    //log:send QBA message
                    log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]Send QBA Message To Node %v",n.id,n.id,TargetID,TargetID)
                }
            }
        }
}


