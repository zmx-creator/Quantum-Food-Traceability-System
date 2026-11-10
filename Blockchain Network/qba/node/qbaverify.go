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

            //a node that missed an abandon report catches up here, whatever role
            //this round assigns to it
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
            
            //log:verifier node received message from forwarder node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Receive Forward Message From Node %v",n.id,msg.SenderID,msg.ID,n.id,msg.ID)
            
            //fetch QDS private key of verifier
            kx, ky, kz := n.Verifier_KeyCache.CalculateKeyLines(n.ConsensusCount, msg.QdsNum,n.Verifier_KeyCache.N)
            
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
            n.QBAKey_X = n.Verifier_KeyCache.ReadKeyFromCache(kx)
            n.QBAKey_Y = n.Verifier_KeyCache.ReadKeyFromCache(ky)
            n.QBAKey_Z = n.Verifier_KeyCache.ReadKeyFromCache(kz)
            
            //output check:node forwarder fetch the QDS private key of verifier
            log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
            
            //get the verify result
            commitresult :=message.QDSVerify(msg.Digest,msg.QDSSignature,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z)

		//Consistency check (protocol description, broadcast phase II): a signer owns
		//its QDS key lines, so it is free to sign anything and the signature verifies
		//perfectly. What exposes a signer that shows different values to different
		//recipients is this comparison with what the same signer already attested in
		//this round. An inconsistent message is answered as a failed verdict, which is
		//what keeps it out of the judgment list, and its sender is asked to repeat.
		//A conflicting message that arrived through a malicious relay is the collusion
		//the analysis grants: two malicious parties make it verify, and the honest
		//answer is the majority over the three relays, so the value is kept and
		//counted. A conflicting message that arrived through an honest relay is the
		//signer lying to an honest forwarder, which is what this check exists to stop.
		if commitresult == "verify succeed" && !n.fault.IsHostile(int(msg.ID)) &&
			!n.consistentSigner(msg.SenderID, msg.Digest, msg.QdsNum) {
			commitresult = "verify failed"
		}
            
            //use commit resule to pack the commit message to send back to forwarder node
            content1, _, err1 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,commitresult,msg.QdsNum)
            
            //pack error check
            if err1 != nil {
                log.Printf(BroadCastNotify+"generate qba-commit message error")
                return
            }
            
            //send message back to forwarder node
            go SendPost(content1,n.table[msg.ID]+server.QBACommitEntry)

            n.record("verify", commitresult, int(msg.SenderID), int(msg.ID), int(n.id), msg.QdsNum, msg.Digest)
            
            //log:send message back to forwarder node
            log.Printf(BroadCastNotify+"[cur_node:%v,signer:%v,forward:%v,verifier:%v]Finish Verify And Send Back To Node %v",n.id,msg.SenderID,msg.ID,n.id,msg.ID)
        }
    }
}

