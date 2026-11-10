package node

import (
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"encoding/json"
	"bytes"
	"time"
)

// QDS outcome helpers shared by the two broadcast-phase verifiers.
//
// failedQDS implements section I.B step 4 of the protocol description: a QDS is
// successful only when the verifier and the forwarder both declare it verified,
// so either negative verdict makes the instance fail.
func failedQDS(local string, peer string) bool {
    return local != "verify succeed" || peer != "verify succeed"
}

// negativeVerdict copies a commit message and marks the copy as a failed QDS
// verdict. The copy keeps the original ID so the role bookkeeping still sees
// this verifier slot as used; isValidVerdict keeps it out of every count.
func negativeVerdict(msg *message.QBACommitMessage) *message.QBACommitMessage {
    recorded := *msg
    recorded.CommitType = "verify failed"
    return &recorded
}

// finishQDSInstance closes one QDS instance from the forwarder side: it tells
// the signer whether every verifier slot has been used and, when it has, clears
// the buffer for the next instance. The failed path uses the same call, which is
// what stops a forged signature from hanging the round forever.
func (n *Node) finishQDSInstance(msg *message.QBACommitMessage, finishCheck bool) {
    commitType := "finished QDS but not all node forwarded"
    if finishCheck {
        commitType = "finished QDS and all node forwarded"
    }
    content, _, err := message.NewQBACommitMessage(msg.Digest, n.id, msg.SenderID, n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, commitType, msg.QdsNum)
    if err != nil {
        log.Printf("[BroadCast]generate QBA commit message error")
        return
    }
    go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
    log.Printf("finished QDS:%v", msg.QdsNum)
    log.Printf("[BroadCast][cur_node:%v,signer:%v,forward:%v]Node %v finish this QDS", n.id, msg.SenderID, n.id, n.id)
    if finishCheck {
        n.buffer.ClearQBACommitMsg()
    }
}
//case I:
//in broadcast I,forwarder node receive the commit message from verifier node

//case:finish this QDS,but not all optional node has performed verifier node

//case:finish this QDS,and all optional node has performed verifier node

//the hallmark of this case:
//Is Primary Node:no
//BroadCast Phase:I
//Is Forwarder Node:yes


//case II:
//in broadcast I,signer node receive the commit notice from forwarder node

//case:finish this QDS,but for this forwarder node,not all optional node has performed verifier node

//case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
//but not all optional node has performed forwarder node 

//case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
//and all optional node has performed forwarder node 

//the hallmark of this case:
//Is Primary Node:yes
//BroadCast Phase:I
//Is Forwarder Node:no


//case III:
//in broadcast I,forwarder node receive the commit message from verifier node

//case:finish this QDS,but not all optional node has performed verifier node

//case:finish this QDS,and all optional node has performed verifier node

//the hallmark of this case:
//Is Primary Node:no
//BroadCast Phase:II
//Is Forwarder Node:yes


//case IV:
//in broadcast II,signer node receive the commit notice from forwarder node

//case:finish this QDS,but for this forwarder node,not all optional node has performed verifier node

//case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
//but not all optional node has performed forwarder node 

//case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
//and all optional node has performed forwarder node 

//the hallmark of this case:
//Is Primary Node:no
//BroadCast Phase:II
//Is Forwarder Node:no


//case V:
//in broadcast II,primary node notice node to send QBA message

//case:not all node has performed signer

//case: all node has performed signer

//the hallmark of this case:
//Is Primary Node:yes
//BroadCast Phase:II
//Is Forwarder Node:no


func (n *Node) QBACommitRecvThread() {
    for {
        select {
            case msg := <-n.QBACommitRecv:
            switch {
            
            //case I:
            //msg.SenderID == n.GetPrimary()-->Is Primary Node:no
            //msg.SenderID == n.GetPrimary()-->BroadCast Phase:I
            //msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"-->Is Forwarder Node:yes
            case msg.SenderID == n.GetPrimary() && (msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"):  //case I
                n.handleBroadcastPhaseIVetify(msg)
             
            //case III:
            //msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"msg.SenderID == n.GetPrimary()-->Is Primary Node:no
            //msg.SenderID != n.GetPrimary()-->BroadCast Phase:II
            //msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"-->Is Forwarder Node:yes  
            case msg.SenderID != n.GetPrimary() && (msg.CommitType == "verify succeed" || msg.CommitType == "verify failed")://case III
                n.handleBroadcastPhaseIIVetify(msg)
            
            //case II: 
            case msg.SenderID == n.GetPrimary() && (msg.CommitType == "finished QDS but not all node forwarded" || msg.CommitType == "finished QDS and all node forwarded" ): 
                n.handlePhaseIQDSFinish(msg)
            
            //case IV    
            case msg.SenderID != n.GetPrimary() && (msg.CommitType == "finished QDS but not all node forwarded" || msg.CommitType == "finished QDS and all node forwarded" ):
                n.handlePhaseIIQDSFinish(msg)
            
            //case V
            case msg.CommitType == "finished send QBA message" :
                n.handlePhaseIIFinish(msg)
            
            
            case msg.CommitType == "Ready To Infer"://case infer
                n.handleInfer(msg)
                
            case msg.CommitType == "Successfully infer the result message" ://case execute
                n.handleExecute(msg)
                
            case msg.CommitType == "Consensus Failed" ://case execute
                n.handleExecute(msg)

            //a peer claims this node produced two different values in one round.
            //Nothing in the state machine depends on it: the honest path cannot
            //receive it twice, and the demonstration has to show it
            case msg.CommitType == ConsistencyComplaint :
                n.handleConsistencyComplaint(msg)
            }
        }
    }
}

func (n *Node) handleBroadcastPhaseIVetify(msg *message.QBACommitMessage) {//case I
    //forwarder also do a verify
    verifyresult := message.QDSVerify(msg.Digest,n.QBAMainMessage[0].QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z)

    //only if forwarder and verifier all successfully verify
    //then this QDS verify is successful
    //then node can proceed with subsequent operations
    
    //A QDS instance that fails verification is neither a dead end nor a reason to
    //kill the round: it is recorded as a negative verdict so the round still
    //completes and reaches the inference phase, while Majority and the commit
    //arithmetic skip it. The forged message is therefore visible to the
    //demonstration, it just never wins the vote.
    if failedQDS(verifyresult, msg.CommitType) {
        log.Printf("[BroadCast][cur_node:%v]QDS %v verify failed (local:%v, peer:%v): recorded as a negative verdict, the chain continues", n.id, msg.QdsNum, verifyresult, msg.CommitType)
        n.record("verify", "QDS verification failed: kept out of the judgment list", int(msg.SenderID), int(n.id), int(msg.ID), msg.QdsNum, msg.Digest)
        n.buffer.BufferQBACommitMsg(negativeVerdict(msg))
        _, FinishCheck := n.FindTargetID(msg.SenderID)
        n.finishQDSInstance(msg, FinishCheck)
        return
    }

    if verifyresult == "verify succeed" && msg.CommitType == "verify succeed" {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        _, FinishCheck := n.FindTargetID(n.id)
        //case:finish this QDS,but not all optional node has performed verifier node
        if !FinishCheck {
        //package commit message
        //notice primary node that not all node has performed verifier node
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id, msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z,"finished QDS but not all node forwarded",msg.QdsNum)
            log.Printf("finished QDS:%v",msg.QdsNum)
            //package error check
            if err != nil {
                log.Printf("[BroadCast Phast I]generate QBA forward message error")
                return
            }
        
            //send commit message to primary node
            go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
        
            //log:send commit message to primary node
            log.Printf("[BroadCast Phast I][cur_node:%v,signer:%v,forward:%v]Node %v finish this QDS",n.id,msg.SenderID,n.id, n.id)
        } else {//case:finish this QDS,but not all optional node has performed verifier node
            //package commit message
            //notice primary node that all node has performed verifier node
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "finished QDS and all node forwarded",msg.QdsNum)
            log.Printf("finished QDS:%v",msg.QdsNum)
            //package error check
            if err != nil {
                log.Printf("[BroadCast Phast I]generate QBA commit message error")
                return
            }
        
            //send commit message to primary node
            go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
        
            //clear QBA commit message buffer to prepare for next phase
            n.buffer.ClearQBACommitMsg()
            log.Printf("[BroadCast Phast I][cur_node:%v,signer:%v,forward:%v]Node %v finish this QDS",n.id,msg.SenderID,n.id, n.id)
        }
    }
}

func (n *Node) handleBroadcastPhaseIIVetify(msg *message.QBACommitMessage) {//case III
    //forwarder also do a verify
    verifyresult := message.QDSVerify(msg.Digest,n.QBAMainMessage[0].QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z)

    //only if forwarder and verifier all successfully verify
    //then this QDS verify is successful
    //then node can proceed with subsequent operations
    
    //A QDS instance that fails verification is neither a dead end nor a reason to
    //kill the round: it is recorded as a negative verdict so the round still
    //completes and reaches the inference phase, while Majority and the commit
    //arithmetic skip it. The forged message is therefore visible to the
    //demonstration, it just never wins the vote.
    if failedQDS(verifyresult, msg.CommitType) {
        log.Printf("[BroadCast][cur_node:%v]QDS %v verify failed (local:%v, peer:%v): recorded as a negative verdict, the chain continues", n.id, msg.QdsNum, verifyresult, msg.CommitType)
        n.record("verify", "QDS verification failed: kept out of the judgment list", int(msg.SenderID), int(n.id), int(msg.ID), msg.QdsNum, msg.Digest)
        n.buffer.BufferQBACommitMsg(negativeVerdict(msg))
        _, FinishCheck := n.FindTargetID(msg.SenderID)
        n.finishQDSInstance(msg, FinishCheck)
        return
    }

    if verifyresult == "verify succeed" && msg.CommitType == "verify succeed" {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        _, FinishCheck := n.FindTargetID(msg.SenderID)
        //case:finish this QDS,but not all optional node has performed verifier node
        if !FinishCheck {
        //package commit message
        //notice primary node that not all node has performed verifier node
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id, msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z,"finished QDS but not all node forwarded",msg.QdsNum)
            log.Printf("finished QDS:%v",msg.QdsNum)
            //package error check
            if err != nil {
                log.Printf("[BroadCast Phast II]generate QBA forward message error")
                return
            }
        
            //send commit message to primary node
            go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
        
            //log:send commit message to primary node
            log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Node %v finish this QDS",n.id,msg.SenderID,n.id, n.id)
        }else {//case:finish this QDS,but not all optional node has performed verifier node
            //package commit message
            //notice primary node that all node has performed verifier node
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "finished QDS and all node forwarded",msg.QdsNum)
            log.Printf("finished QDS:%v",msg.QdsNum)
            //package error check
            if err != nil {
                log.Printf("[BroadCast Phast II]generate QBA commit message error")
                return
            }
            
            //log:node received all commit message in QBA commit message buffer
            Output := "["
            count:=0
            for _,v := range n.buffer.QBACommitBuffer{
                Output+=n.ConvertBInaryByteToString(v.Digest[:8])
                count+=1
                if count<len(n.buffer.QBACommitBuffer){
                    Output+=","
                }
            }
            Output += "]"
            log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Receive Message List Is:%v",n.id,msg.SenderID,n.id,Output)
            
            //node infer the result from these commit message 
            voteResult := n.buffer.Vote()   
            
            //log:node infer the result from these commit message 
            log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Choose Message:%v To Join The Infer List",n.id,msg.SenderID,n.id,string(voteResult.Digest[:8]))
            
            //use the infer result to pack vote message and save in vote message buffer
            _, voteMsg, err1 := message.NewQBAVoteMessage(voteResult.Digest, msg.SenderID, msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, voteResult.AllVoteEqual)
            
            //pack error check
            if err1 != nil {
                log.Printf("[BroadCast Phast II]generate QBA vote message error")
                return
            }
            
            //save the vote message, but only after the consistency check: what this
            //signer attests now has to equal what it attested earlier in the same
            //round. A malicious signer signs a different message for every recipient
            //and every one of those signatures verifies, so this comparison is the
            //only thing that exposes the equivocation.
            n.buffer.BufferQBAVoteMsg(voteMsg)
            
            //send commit message to sender
            go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
            
            //clear QBA commit message buffer to prepare for next phase
            n.buffer.ClearQBACommitMsg()
        }
    }
}

func (n *Node) handlePhaseIQDSFinish(msg *message.QBACommitMessage) {//case II
    //the process of case II is the primary node's procession
    if n.id != n.GetPrimary() {
        return
    }
    
    //case:finish this QDS,but for this forwarder node,not all optional node has performed verifier node
    if msg.CommitType == "finished QDS but not all node forwarded" {
        //QDS counter add 1,mean that start next QDS
        qdsnum := msg.QdsNum+1
        
        //fetch QDS private key of signer
        
        log.Printf("QDS time:%v",qdsnum)
        kx, ky, kz := n.Forwarder_KeyCache.CalculateKeyLines(n.ConsensusCount, qdsnum,n.Signer_KeyCache.N)
        n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
        n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
        n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	
	//fetch irreducible polynomial factor
	factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	//get QDS signature
	sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
    
        //choose same node as the forwarder node to send QBA message
        TargetID := msg.ID
        
        //output check:node signer fetch the QDS private key of signer
        log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
        log.Printf("irreducible polynomial factor loaded (value not logged)")
        //use the pre-saved QBA message to pack new QBA message
        copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
        copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
        n.QBAMainMessage[0].QdsNum=qdsnum
        
        content,err := json.Marshal(n.QBAMainMessage[0])
        //pack error check
        if err != nil {
            log.Printf("[BroadCast Phast I]generate QBA forward message error")
            return
        }
        
        //send QBA message
        go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
        log.Printf("[BroadCast Phast I][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
    }
    if msg.CommitType == "finished QDS and all node forwarded" {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        
        
        //TargetID-->choose optional node to perform the new forwarder
        //FinishCheck-->and identify whether all optional node has performed forwarder
        TargetID, FinishCheck := n.FindTargetID(n.id)
        
        //case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
        //but not all optional node has performed forwarder node 
        if !FinishCheck {
            //QDS counter add 1,mean that start next QDS
            qdsnum := msg.QdsNum+1
    
            
            //fetch QDS private key of signer
            log.Printf("QDS time:%v",qdsnum)
            kx, ky, kz := n.Forwarder_KeyCache.CalculateKeyLines(n.ConsensusCount, qdsnum,n.Signer_KeyCache.N)
            n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
            n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
            n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	
	    //fetch irreducible polynomial factor
	    factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	
	    //get QDS signature
	    sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
        
            //output check:node signer fetch the QDS private key of signer
            log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
            log.Printf("irreducible polynomial factor loaded (value not logged)")
        
            //use the pre-saved QBA message to pack new QBA message
            copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
            copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
            n.QBAMainMessage[0].QdsNum=qdsnum
        
            content,err := json.Marshal(n.QBAMainMessage[0])
            //pack error check
            if err != nil {
                log.Printf("[BroadCast Phast I]generate QBA forward message error")
                return
            }
        
            //send QBA message
            go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
            log.Printf("[BroadCast Phast I][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
        }else{//case:finish this QDS,and all optional node has performed verifier node
            //log:finish broadcast I and start broadCast Phase II
            log.Printf("[BroadCast Phast I]finished broadcast I")
            log.Printf("[BroadCast Phast I]Start BroadCast Phase II")
            
            //clear QBA commit message buffer to prepare for next phase
            n.buffer.ClearQBACommitMsg()
            
            //choose an ID to notice it to send QBA message and start QDS
            TargetID, _ := n.FindTargetID(msg.ID)
            content, _, err := message.NewQBAPrepareMessage(n.id,msg.QdsNum)
            if err != nil {
                log.Printf("[BroadCast Phast I]generate QBA prepare message error")
                return
            }
            go SendPost(content, n.table[TargetID]+server.QBAPrepareEntry)
            log.Printf("Notice Node %v to send QBA Message",TargetID)
        }
    }
}

//case IV
func (n *Node) handlePhaseIIQDSFinish(msg *message.QBACommitMessage) {
    //the process of case II is the primary node's procession
    if n.id == n.GetPrimary() {
        return
    }
    
    //case:finish this QDS,but for this forwarder node,not all optional node has performed verifier node
    if msg.CommitType == "finished QDS but not all node forwarded" {
        //QDS counter add 1,mean that start next QDS
        qdsnum := msg.QdsNum+1
        
        //fetch QDS private key of signer
        log.Printf("QDS time:%v",qdsnum)
        kx, ky, kz := n.Forwarder_KeyCache.CalculateKeyLines(n.ConsensusCount, qdsnum,n.Signer_KeyCache.N)
        n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
        n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
        n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	
	//fetch irreducible polynomial factor
	factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	
	//get QDS signature
	sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
    
        //choose same node as the forwarder node to send QBA message
        TargetID := msg.ID
        
        //output check:node signer fetch the QDS private key of signer
        log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
        log.Printf("irreducible polynomial factor loaded (value not logged)")
        
        //use the pre-saved QBA message to pack new QBA message
        copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
        copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
        n.QBAMainMessage[0].QdsNum=qdsnum
        
        content,err := json.Marshal(n.QBAMainMessage[0])
        //pack error check
        if err != nil {
            log.Printf("[BroadCast Phast II]generate QBA forward message error")
            return
        }
        
        //demonstration hook: this copy is the signer's own, so it is signed as usual
        //unless the relay chosen for it is deviating too - the collusion of the protocol
        //analysis, and the only hop where a conflicting message can pass the QDS check.
        if lying, ok := n.equivocate(&n.QBAMainMessage[0], factor, qdsnum, int(TargetID)); ok {
            content = lying
        }
        //send QBA message
        go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
        log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
    }
    if msg.CommitType == "finished QDS and all node forwarded" {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        
        
        //TargetID-->choose optional node to perform the new forwarder
        //FinishCheck-->and identify whether all optional node has performed forwarder
        TargetID, FinishCheck := n.FindTargetID(n.id)
        
        //case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
        //but not all optional node has performed forwarder node 
        if !FinishCheck {
            //QDS counter add 1,mean that start next QDS
            qdsnum := msg.QdsNum+1
            
            log.Printf("QDS time:%v",qdsnum)
            //fetch QDS private key of signer
            kx, ky, kz := n.Forwarder_KeyCache.CalculateKeyLines(n.ConsensusCount, qdsnum,n.Signer_KeyCache.N)
            n.QBAKey_X = n.Signer_KeyCache.ReadKeyFromCache(kx)
            n.QBAKey_Y = n.Signer_KeyCache.ReadKeyFromCache(ky)
            n.QBAKey_Z = n.Signer_KeyCache.ReadKeyFromCache(kz)
	
	    //fetch irreducible polynomial factor
	    factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	
	    //get QDS signature
	    sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
        
            //output check:node signer fetch the QDS private key of signer
            log.Printf("QDS key lines %v %v %v for this signature (key values are not logged)", kx, ky, kz)
            log.Printf("irreducible polynomial factor loaded (value not logged)")
        
            //use the pre-saved QBA message to pack new QBA message
            copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
            copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
            n.QBAMainMessage[0].QdsNum=qdsnum
        
            content,err := json.Marshal(n.QBAMainMessage[0])
            //pack error check
            if err != nil {
                log.Printf("[BroadCast Phast II]generate QBA forward message error")
                return
            }
        
            //demonstration hook: this copy is the signer's own, so it is signed as usual
            //unless the relay chosen for it is deviating too - the collusion of the protocol
            //analysis, and the only hop where a conflicting message can pass the QDS check.
            if lying, ok := n.equivocate(&n.QBAMainMessage[0], factor, qdsnum, int(TargetID)); ok {
                content = lying
            }
            //send QBA message
            go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
            log.Printf("[BroadCast Phast II][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
        }else{//case:finish this QDS,and all optional node has performed verifier node
            //pack commit message
            //notice primary node to change target node to notice it to send QBA Message
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID, n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "finished send QBA message",msg.QdsNum)
            
            //pack error check
            if err != nil {
                log.Printf("[BroadCast Phast II]generate QBA commit message error")
                return
            }  
            
            //send
            go SendPost(content, n.table[n.GetPrimary()]+server.QBACommitEntry)
            n.buffer.ClearQBACommitMsg()
        }
    }
}

//case V
func (n *Node) handlePhaseIIFinish(msg *message.QBACommitMessage) {
    if n.id == n.GetPrimary() {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        //TargetID-->choose optional node to perform the new forwarder
        //FinishCheck-->and identify whether all optional node has performed forwarder
        TargetID, FinishCheck := n.FindTargetID(n.id)
        
        //case:not all node has performed signer
        if !FinishCheck {
            //notice next node to send QBA message
            content, _, err := message.NewQBAPrepareMessage(n.id,msg.QdsNum)
            if err != nil {
                log.Printf("[BroadCast Phast II]generate QBA prepare message error")
                return
            }
            go SendPost(content, n.table[TargetID]+server.QBAPrepareEntry)
        } else {
            //notice all node to infer final infer result 
            content1, _, err1 := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID ,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "Ready To Infer",msg.QdsNum)
            if err1 != nil {
                log.Printf("[BroadCast Phast II]generate QBA commit message error")
                return
            }
            for k,v := range n.table{
                if k == n.id {
                    continue
                }
                go SendPost(content1,v+server.QBACommitEntry)
            }  
            n.buffer.ClearQBACommitMsg()
            n.buffer.ClearQBAVoteMsg()
        }
    } 
}
func (n *Node) handleInfer(msg *message.QBACommitMessage) {//case infer
    if n.id == n.GetPrimary() {
        return
    }
    //The judgment list is indexed by the node that forwarded the message, which
    //is the list the protocol description writes down for one deputy as
    //{m1, m2, m3, m41}. Entries are printed as forwarder:value so the
    //demonstration page can show what each node actually received.
    InferList:="["
    count :=0
    for k,v := range n.buffer.QBAVoteBuffer{
        InferList+=fmt.Sprintf("%v:%v", k, digestPrefix(v.Result.Digest))
        count++
        if count < len(n.buffer.QBAVoteBuffer){
            InferList+=","
        }
    }
    InferList += "]"
    log.Printf("[Infer Phase][cur_node:%v,signer:%v,forward:%v]Node %v's Infer List:%v",n.id,msg.SenderID,n.id,n.id,InferList)
    voteResult:=n.buffer.Infer()
    if voteResult == nil {
        //Nothing survived verification: the node has no decision to report, and
        //the primary fails the round through the commit quorum instead.
        log.Printf("[Infer Phase][cur_node:%v]judgment list is empty, no decision to report",n.id)
        n.record("infer", "judgment list "+InferList+" is empty: no decision", int(msg.SenderID), int(msg.SenderID), int(n.id), msg.QdsNum, nil)
        return
    }
    log.Printf("[Infer Phase][cur_node:%v,signer:%v,forward:%v]Node %v Choose Message:%v To Be The Infer Result",n.id,msg.SenderID,n.id,n.id,n.ConvertBInaryByteToString(voteResult.Digest[:8]))
    var Result string
    if bytes.Equal(voteResult.Digest,n.QBAMainMessage[0].Digest) && voteResult.AllVoteEqual == false{
         Result = "Successfully infer the result message"
    }else{
         //A tie at the top of the vote count is the abstention symbol of the
         //protocol: this node cannot decide, and it does not block the others.
         Result = "Consensus Failed"
    }
    decision := voteResult.Digest
    n.record("infer", "judgment list "+InferList+" -> majority "+digestPrefix(decision), int(msg.SenderID), int(msg.SenderID), int(n.id), msg.QdsNum, decision)
    if n.fault.Mode() != FaultNone {
        //A deviating deputy does not abstain and does not follow the majority: it
        //reports a decision for a different message, the m41/m42/m43 of the
        //protocol description. It is a lie about its own vote, so no signature
        //has to be forged for it, and an honest quorum never counts it.
        decision = n.lie(decision)
        Result = "Successfully infer the result message"
        n.record("infer", "byzantine report: claims it decided "+digestPrefix(decision)+" instead of the proposed batch", int(msg.SenderID), int(msg.SenderID), int(n.id), msg.QdsNum, decision)
    }
    content,_,err:=message.NewQBACommitMessage(decision, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, Result,msg.QdsNum)
    if err != nil {
        log.Printf("[Infer Phase]generate QBA commit message error")
        return
    }  
    go SendPost(content,n.table[n.GetPrimary()]+server.QBACommitEntry)
}


func (n *Node) handleExecute(msg *message.QBACommitMessage) {//case last
    if n.id == n.GetPrimary() {
        //Commit rule of this implementation.
        //
        //A report counts when it carries a countable verdict AND the decision it
        //names is the very batch this primary proposed. Arrival alone proves
        //nothing: a malicious node answers too, it answers with a different
        //message. The previous rule waited for len(table)-1 reports, which
        //measured liveness and not agreement.
        //
        //With QDS-authenticated broadcast the fault bound is f=(n-1)/2, so the
        //quorum is n-f=3 for five nodes: one or two malicious nodes are outvoted
        //and the block is still written, three of them leave the honest side with
        //two reports, no quorum, and the batch is discarded once the retries
        //allowed by the round watchdog are used up.
        //
        //Liveness follows the bounded policy of round.go: a batch that does not
        //gather its quorum is retransmitted a fixed number of times and then
        //discarded, and the round reports the discard instead of blocking.
        if n.round.State() != RoundRunning || len(n.QBAMainMessage) == 0 {
            return
        }
        n.buffer.BufferQBACommitMsg(msg)
        proposed := n.QBAMainMessage[0].Digest
        agreeing, answered := n.buffer.CountAgreeing(proposed)
        agreeing++ //the primary decided for the batch it proposed
        quorum := len(n.table) - int(n.faultNum)
        retries, _, _ := n.round.Cumulative()
        log.Printf("[Commit Phase][cur_node:%v]report from node %v: %v -> %v/%v agree with the proposed batch (quorum %v, answered %v, retransmissions %v)",
            n.id, msg.ID, msg.CommitType, agreeing, len(n.table), quorum, answered, retries)
        n.record("commit", fmt.Sprintf("node %v reported %v: %v/%v agree with the proposed batch, quorum %v",
            msg.ID, digestPrefix(msg.Digest), agreeing, len(n.table), quorum), int(n.id), -1, int(msg.ID), msg.QdsNum, msg.Digest)
        if agreeing < quorum {
            log.Printf("[Commit Phase][cur_node:%v]no quorum yet (%v/%v), waiting", n.id, agreeing, quorum)
            return
        }
        content,_,err:=message.NewQBACommitMessage(proposed, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "Successfully infer the result message",msg.QdsNum)
        if err != nil {
            log.Printf("[Commit Phase]generate QBA commit message error")
            return
        }
        for k,v := range n.table{
            if k==n.id{
                continue
            }
            go SendPost(content,v+server.QBACommitEntry)
        }
        log.Printf("[Infer Phase][cur_node:%v,signer:%v]Choose Message:%v To Be Consensus Result",n.id,n.id,string(n.QBAMainMessage[0].Digest[:8]))
        log.Printf("[Commit Phase][cur_node:%v]quorum reached %v/%v: writing the block", n.id, agreeing, len(n.table))
        log.Printf("Finish Check")
        n.record("commit", fmt.Sprintf("quorum %v/%v reached: the batch is written into a block", agreeing, len(n.table)), int(n.id), -1, -1, msg.QdsNum, proposed)
        //settle the round here so a late duplicate report can never write a
        //second block for the same instance
        n.round.Commit(time.Now())
        n.buffer.AppendToExecuteQueue(&n.QBAMainMessage[0])
        n.executeNotify<-true
        n.buffer.ClearQBAVoteMsg()
        n.buffer.ClearQBACommitMsg()
    }
    if n.id != n.GetPrimary() {
            if msg.CommitType == "Consensus Failed"{
                log.Printf("[Infer Phase][cur_node:%v]Consensus Failed",n.id)
                //advance this node's key window together with the rest of the group,
                //otherwise the node stays on an exhausted window and every later
                //round fails verification
                n.abandonAsFollower("the primary abandoned the round")
            }
            if msg.CommitType == "Successfully infer the result message"{
            log.Printf("[Infer Phase][cur_node:%v,signer:%v]Choose Message:%v To Be Consensus Result",n.id,n.id,string(n.QBAMainMessage[0].Digest[:8]))
            log.Printf("[Infer Phase][cur_node:%v,signer:%v]Finish Consensus",n.id,n.id)
            n.buffer.AppendToExecuteQueue(&n.QBAMainMessage[0])
            n.executeNotify<-true
            n.buffer.ClearQBAVoteMsg()
            n.buffer.ClearQBACommitMsg()
            }
    }
}


