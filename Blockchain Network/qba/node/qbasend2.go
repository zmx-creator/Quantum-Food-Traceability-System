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

            //a node that missed an abandon report catches up here, whatever role
            //this round assigns to it
            n.realignKeyWindow(msg.QdsNum)
                if n.id != n.GetPrimary(){
                    //QDS counter add 1,mean that start next QDS
                    qdsnum := msg.QdsNum+1
                    log.Printf("QDS time:%v",qdsnum)
	    
                    //fetch QDS private key of signer
                    kx, ky, kz := n.Signer_KeyCache.CalculateKeyLines(n.ConsensusCount, qdsnum,n.Signer_KeyCache.N)
                    n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
                    n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
                    n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	
                    //fetch irreducible polynomial factor
                    factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	
                    //get QDS signature
                    sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
    
                    //choose same node as the forwarder node to send QBA message
                    TargetID, _ := n.FindTargetID(n.id)

        
                    //output check:node signer fetch the QDS private key of signer
                    log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
                    log.Printf("irreducible polynomial factor loaded (value not logged)")
        
                    //use the pre-saved QBA message to pack new QBA message
                    copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
                    copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
                    n.QBAMainMessage[0].SenderID = n.id
                    n.QBAMainMessage[0].QdsNum=qdsnum
                    content,err := json.Marshal(n.QBAMainMessage[0])

                    //demonstration hook: a malicious signer offers a different message
                    //only when the relay it is handing this copy to is malicious too,
                    //the collusion of the protocol analysis. With an honest relay it
                    //keeps to the agreed message, which the consistency check of the
                    //round would otherwise expose, so that relay's row stays the honest
                    //value exactly as the protocol table shows it.
                    if lying, ok := n.equivocate(&n.QBAMainMessage[0], factor, qdsnum, int(TargetID)); ok {
                        content = lying
                    }
                    //pack error check
                    if err != nil {
                        log.Printf("[BroadCast Phast II]generate QBA forward message error")
                        return
                    }
        
                    
                    //just use the varible msg
                    log.Printf("receive notice from node %v",msg.SenderID)
                    
                    //send QBA message
                    go SendPost(content, n.table[TargetID] + server.QBAMessageEntry) 
                    
                    //log:send QBA message
                    log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Send QBA Message To Node %v",n.id,n.id,TargetID,TargetID)
                }
            }
        }
}


