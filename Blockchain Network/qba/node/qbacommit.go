package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"encoding/json"
	"bytes"
)

func (n *Node) QBACommitRecvThread() {
    for {
        select {
            case msg := <-n.QBACommitRecv:
            //n.buffer.BufferQBACommitMsg(msg)
            switch {
            
            //case I:
            //msg.SenderID == n.GetPrimary()-->Is Primary Node:no
            //msg.SenderID == n.GetPrimary()-->BroadCast Phase:I
            //msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"-->Is Forwarder Node:yes
            case msg.SenderID == n.GetPrimary() && (msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"):  //case I
                n.handleBroadcastPhaseIVerify(msg)
            
            //case III    
            case msg.SenderID != n.GetPrimary() && (msg.CommitType == "verify succeed" || msg.CommitType == "verify failed"):  //case I
                n.handleBroadcastPhaseIIVerify(msg)
            
            //case II: 
            case msg.SenderID == n.GetPrimary() && msg.CommitType == "finished QDS"||msg.CommitType == "finish backup transaction data": 
                n.handlePhaseIQDSFinish(msg)
            
            //case IV    
            case msg.SenderID != n.GetPrimary() && msg.CommitType == "finished QDS and circle all node forwarded" :
                n.handlePhaseIIQDSFinish(msg)
            
            case msg.CommitType == "Ready To Check QDS Signature"://case infer
                n.handleSignCheckPrepare(msg)
                
            case msg.CommitType == "Check QDS Signature" ://case execute
                n.handleSignCheck(msg)
                
            case msg.CommitType == "all signature valid" ://case execute
                n.handleExecute(msg)
            }
        }
    }
}

func (n *Node) handleBroadcastPhaseIVerify(msg *message.QBACommitMessage) {//case I
    //forwarder also do a verify
    verifyresult := message.QDSVerify(msg.Digest,n.QBAMainMessage[0].QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z)

    //only if forwarder and verifier all successfully verify
    //then this QDS verify is successful
    //then node can proceed with subsequent operations
    if verifyresult == "verify succeed" && msg.CommitType == "verify succeed" {
        //notice main node
        //that this node has finished QDS
        content, _, err := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "finished QDS",msg.QdsNum)
        
        //package error check
        if err != nil {
            log.Printf("[BroadCast Phase]generate QBA commit message error")
            return
        }
        
        //send commit message to primary node
        go SendPost(content, n.table[msg.SenderID]+server.QBACommitEntry)
        
        //label:finished QDS
        log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]Node %v finish QDS:%v",n.id,msg.SenderID,n.id, n.id,msg.QdsNum)
    }
    //case:verify failed
    //else{
    //}
}

func (n *Node) handleBroadcastPhaseIIVerify(msg *message.QBACommitMessage) {//case III  
    //forwarder also do a verify
    verifyresult := message.QDSVerify(msg.Digest[:len(msg.Digest)-message.KeyLength*2],n.QBAMessageBackup[len(n.QBAMessageBackup)-1].QDSSignature,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,msg.QBAKey_X,msg.QBAKey_Y,msg.QBAKey_Z)

    //only if forwarder and verifier all successfully verify
    //then this QDS verify is successful
    //then node can proceed with subsequent operations
    if verifyresult == "verify succeed" && msg.CommitType == "verify succeed" {
        
        //label:finished QDS
        log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]Node %v finish QDS:%v",n.id,msg.SenderID,n.id, n.id,msg.QdsNum)
        
        //choose the next node of circle
        TargetID:= n.CircleNodeIndex()
        
        //case:finish this QDS,but not all node in circle has done QDS
        if n.id!=msg.SenderID {
            //QDS counter add 1,mean that start next QDS
            qdsnum := msg.QdsNum+1
    
            
            //fetch QDS private key of signer
            log.Printf("QDS time:%v",qdsnum)
            n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-2),false)
            n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-1),false)
            n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum),false)
	
	    //fetch irreducible polynomial factor
            factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
	    
	    //combine message
	    //format:prev_message-prev_sign-message_N-sign_N 
	    tmpmessage :=msg.Digest
	    tmpmessage = n.CombineMessage(tmpmessage,n.QBAMainMessage[0].Digest)
	    tmpmessage = n.CombineMessage(tmpmessage,n.QBAMainMessage[0].QDSSignature.Digest)
	    tmpmessage = n.CombineMessage(tmpmessage,n.QBAMainMessage[0].QDSSignature.Factor)
	    
	    //get QDS signature
            sig := message.Sign(tmpmessage,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
            tmpmessage = n.CombineMessage(tmpmessage,sig.Digest)
            tmpmessage = n.CombineMessage(tmpmessage,sig.Factor)
            
            //output check:node signer fetch the QDS private key of signer
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key X:%v",n.id,n.id,TargetID,n.id,n.QBAKey_X)
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Y:%v",n.id,n.id,TargetID,n.id,n.QBAKey_Y)
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Z:%v",n.id,n.id,TargetID,n.id,n.QBAKey_Z)
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the factor of irreducible polynomial:%v",n.id,n.id,TargetID,n.id,factor)
            
            //pack new QBA message
            //use the pre-saved QBA message to pack new QBA message
            SendQBAMessage := message.QBAMessage{
                Digest:tmpmessage,  
                SenderID:msg.SenderID, 
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
        
            //send QBA message
            go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
            
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
            
            //clean previous QBA message
            n.QBAMessageBackup=n.QBAMessageBackup[:len(n.QBAMessageBackup)-1]
        }else {//case:finish this QDS,but not all node in circle has done QDS
            //package commit message
            //notice primary node that all node in circle has done QDS
            content, _, err := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "finished QDS and circle all node forwarded",msg.QdsNum)
            
            //package error check
            if err != nil {
                log.Printf("[Circular Collection Phase]generate QBA commit message error")
                return
            }
            
            //send commit message to sender
            go SendPost(content, n.table[n.GetPrimary()]+server.QBACommitEntry)
            
            //log:finish this circle send
            log.Printf("[Circular Collection Phase][cur_node:%v,signer:%v,forward:%v]Node %v finish this circle QDS",n.id,msg.SenderID,n.id,n.id)
        }
    }
    //case:verify failed
    //else{
    //}
}

func (n *Node) handlePhaseIQDSFinish(msg *message.QBACommitMessage) {//case II
    //the process of case II is the primary node's procession
    if n.id != n.GetPrimary() {
        return
    }
    
    //save the commit message
    //mean that this node has performed forwarder node 
    n.buffer.BufferQBACommitMsg(msg)
        
    //TargetID-->choose optional node to perform the new forwarder
    //FinishCheck-->and identify whether all optional node has performed forwarder
    TargetID, FinishCheck := n.FindTargetID(msg.SenderID)
    var qdsnum int
    //case:finish this QDS,and for this forwarder node,all optional node has performed verifier node
    //but not all optional node has performed forwarder node 
    if !FinishCheck{
        //QDS counter add 1,mean that start next QDS
        if TargetID!=message.Identify(len(n.table)-1){
            qdsnum = msg.QdsNum+1
            log.Printf("QDS time:%v",qdsnum)
            n.QBAKey_X = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-2),false)
            n.QBAKey_Y = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum-1),false)
            n.QBAKey_Z = n.ReadNodeKey("/QBAKey/key_signer.txt",(3*qdsnum),false)
            
            //fetch irreducible polynomial factor
            factor := n.ReadNodeKey("/QBAKey/irreducible_coefficients.txt",qdsnum,true)
            
            //get QDS signature
            sig := message.Sign(n.QBAMainMessage[0].Digest,factor,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z) 
            
            log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key X:%v",n.id,n.id,TargetID,n.id,n.QBAKey_X)
            log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Y:%v",n.id,n.id,TargetID,n.id,n.QBAKey_Y)
            log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the QDS key Z:%v",n.id,n.id,TargetID,n.id,n.QBAKey_Z)
            log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]signer node %v fetch the factor of irreducible polynomial:%v",n.id,n.id,TargetID,n.id,factor)
            
            //use the pre-saved QBA message to pack new QBA message
            copy(n.QBAMainMessage[0].QDSSignature.Digest,sig.Digest)
            copy(n.QBAMainMessage[0].QDSSignature.Factor,sig.Factor)
            n.QBAMainMessage[0].QdsNum=qdsnum
        }else{
            qdsnum = msg.QdsNum
        }
        
        content,err := json.Marshal(n.QBAMainMessage[0])
        //pack error check
        if err != nil {
            log.Printf("[BroadCast Phase]generate QBA forward message error")
            return
        }
        
        //send QBA message
        go SendPost(content, n.table[TargetID]+server.QBAMessageEntry)
        
        //log:send QBA message
        if TargetID!=message.Identify(len(n.table)-1){
            log.Printf("[BroadCast Phase][cur_node:%v,signer:%v,forward:%v]Send QBA Message to Node %v",n.id,msg.SenderID,n.id, TargetID)
        }
        
    }else{//case:finish this QDS,and all optional node has performed verifier node
    
        //log:finish broadcast I and start broadCast Phase II
        log.Printf("[BroadCast Phase]finished broadcast I")
        log.Printf("[BroadCast Phase]Start BroadCast Phase II")
        
        BackupMessage:=n.buffer.QBACommitBuffer[message.Identify(len(n.table)-1)]
        
        //clear QBA commit message buffer to prepare for next phase
        n.buffer.ClearQBACommitMsg()
        n.buffer.BufferQBACommitMsg(BackupMessage)
        
        //choose an ID to notice it to send QBA message and start QDS
        TargetID, _ := n.FindTargetID(msg.ID)
        content, _, err := message.NewQBAPrepareMessage(n.id,msg.QdsNum)
        if err != nil {
            log.Printf("[BroadCast Phase]generate QBA prepare message error")
            return
        }
        go SendPost(content, n.table[TargetID]+server.QBAPrepareEntry)
        log.Printf("Notice Node %v to send QBA Message",TargetID)
    }
}

//case IV
func (n *Node) handlePhaseIIQDSFinish(msg *message.QBACommitMessage) {
    //the process of case II is the primary node's procession
    if n.id != n.GetPrimary() {
        return
    }
    if n.id == n.GetPrimary() {
        //save the message in QBA commit message buffer
        //mean that this node is already performed verifier
        //in subsequent QDS,when forwarder node choose node to forward message
        //then the node will avoid th forward message to this node
        n.buffer.BufferQBACommitMsg(msg)
        
        //TargetID-->choose optional node to perform the new forwarder
        //FinishCheck-->and identify whether all optional node has performed forwarder
        TargetID, FinishCheck := n.FindTargetID(n.id)
        
        //case:not all optional node has performed signer
        if !FinishCheck {
            //notice next node to send QBA message
            content, _, err := message.NewQBAPrepareMessage(n.id,msg.QdsNum)
            if err != nil {
                log.Printf("[Circular Collection Phase]generate QBA prepare message error")
                return
            }
            go SendPost(content, n.table[TargetID]+server.QBAPrepareEntry)
        } else {
            //notice all node to prepare for checking signature
            content1, _, err1 := message.NewQBACommitMessage(msg.Digest, n.id,msg.SenderID ,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "Ready To Check QDS Signature",msg.QdsNum)
            if err1 != nil {
                log.Printf("[Circular Collection Phase]generate QBA commit message error")
                return
            }
            for k,v := range n.table{
                if k == n.id||k==message.Identify(len(n.table)-1) {
                    continue
                }
                go SendPost(content1,v+server.QBACommitEntry)
            }  
            log.Printf("[Circular Collection Phase]Ready to check all signature")
            n.buffer.ClearQBACommitMsg()
        }
    } 
}

func (n *Node) handleSignCheckPrepare(msg *message.QBACommitMessage) {
    if n.id != n.GetPrimary() {
        //combine message and signature
        //the message has containing all signature
        //send the signature to verifier node
        SavedMessage:=n.CombineMessage(n.QBAMessageBackup[0].Digest,n.QBAMessageBackup[0].QDSSignature.Digest)
        content1, _, err1 := message.NewQBACommitMessage(SavedMessage, n.id,msg.SenderID ,n.QBAKey_X, n.QBAKey_Y, n.QBAKey_Z, "Check QDS Signature",msg.QdsNum)
        if err1 != nil {
            log.Printf("[Circular Collection Phase]generate QBA commit message error")
            return
        }
        go SendPost(content1, n.table[message.Identify(len(n.table)-1)]+server.QBACommitEntry)
        log.Printf("Problem Check")
    } 
}

func (n *Node) handleSignCheck(msg *message.QBACommitMessage) {
    if n.id == message.Identify(len(n.table)-1) {
        n.buffer.BufferQBACommitMsg(msg)
        var commitresult string
        if len(n.buffer.QBACommitBuffer)==(len(n.table)-2){
            for k,v:=range n.buffer.QBACommitBuffer{
                TmpMessage:=v.Digest
                for i:=0;i<int((len(TmpMessage))/(message.MessageLength+message.KeyLength));i++{
                    //in this message format
                    //extract to signatures at a time
                    sig_Digest1:=TmpMessage[i*(message.MessageLength+4*message.KeyLength)+message.MessageLength:i*(message.MessageLength+4*message.KeyLength)+message.        MessageLength+2*message.KeyLength]
                    sig_Digest2:=TmpMessage[i*(message.MessageLength+4*message.KeyLength)+message.MessageLength+2*message.KeyLength:i*(message.MessageLength+4*message.KeyLength)+message.MessageLength+4*message.KeyLength]
                    verify_message1:=TmpMessage[i*(message.MessageLength+4*message.KeyLength):i*(message.MessageLength+4*message.KeyLength)+message.MessageLength]
                    verify_message2:=TmpMessage[:i*(message.MessageLength+4*message.KeyLength)+message.MessageLength+2*message.KeyLength]
                    
                    //ErrolLabel=false:sign not found
                    ErrorLabel:=make([]bool,2)
                    ErrorLabel[0]=false
                    ErrorLabel[1]=false
                    for _,savedQBAMsg := range n.QBAMessageBackup{
                        //compare with the saved signature
                        if bytes.Equal(savedQBAMsg.QDSSignature.Digest,sig_Digest1[:64])&&bytes.Equal(savedQBAMsg.QDSSignature.Factor,sig_Digest1[64:128]){
                            ErrorLabel[0]=true
                            if ErrorLabel[0]==true{
                                //fetch QDS private key of verifier
                                tmp_QBAKey_X := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum-2),false)
                                tmp_QBAKey_Y := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum-1),false)
                                tmp_QBAKey_Z := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum),false)
                                sig1:=message.QDSSignature{
                                    Digest:sig_Digest1[:message.KeyLength],
                                    Factor:sig_Digest1[message.KeyLength:2*message.KeyLength],
                                }
                                //check the extracted signature
                                commitresult =message.QDSVerify(verify_message1,sig1,savedQBAMsg.QBAKey_X,savedQBAMsg.QBAKey_Y,savedQBAMsg.QBAKey_Z,tmp_QBAKey_X,tmp_QBAKey_Y,tmp_QBAKey_Z)
                                if commitresult=="verify failed"{
                                    //use commit resule to pack the commit message to send back to forwarder node
                                    content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,commitresult,msg.QdsNum)
             
                                    //pack error check
                                    if err2 != nil {
                                        log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                                        return
                                    }
                                    log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Verify failed:invalid signature:%v",n.id,msg.SenderID,msg.ID,n.id,sig_Digest1)
                                    go SendPost(content2,n.table[k]+server.QBACommitEntry)
                                    return
                                }
                            }
                        }
                        if bytes.Equal(savedQBAMsg.QDSSignature.Digest,sig_Digest2[:64])&&bytes.Equal(savedQBAMsg.QDSSignature.Factor,sig_Digest2[64:128]){
                            //skip auxiliary binary sequence
                            ErrorLabel[1]=true
                            if ErrorLabel[1]==true{
                                //fetch QDS private key of verifier
                                tmp_QBAKey_X := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum-2),false)
                                tmp_QBAKey_Y := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum-1),false)
                                tmp_QBAKey_Z := n.ReadNodeKey("/QBAKey/key_verifier.txt",(3*savedQBAMsg.QdsNum),false)
                                sig2:=message.QDSSignature{
                                    Digest:sig_Digest2[:message.KeyLength],
                                    Factor:sig_Digest2[message.KeyLength:2*message.KeyLength],
                                }
                                tmpcommitresult :=message.QDSVerify(verify_message2,sig2,savedQBAMsg.QBAKey_X,savedQBAMsg.QBAKey_Y,savedQBAMsg.QBAKey_Z,tmp_QBAKey_X,tmp_QBAKey_Y,tmp_QBAKey_Z)
                                if tmpcommitresult=="verify failed"{
                                    //use commit resule to pack the commit message to send back to forwarder node
                                    content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,tmpcommitresult,msg.QdsNum)
             
                                    //pack error check
                                    if err2 != nil {
                                        log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                                        return
                                    }
                                    log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Verify failed:invalid signature:%v",n.id,msg.SenderID,msg.ID,n.id,sig_Digest2)
                                    go SendPost(content2,n.table[k]+server.QBACommitEntry)
                                    return
                                }  
                            }
                        }
                    }
                    //signature not found
                    if ErrorLabel[0]==false{
                        //use commit resule to pack the commit message to send back to forwarder node
                        content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"verify failed:signature not found",msg.QdsNum)
             
                        //pack error check
                        if err2 != nil {
                            log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                            return
                        }
                        log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Verify failed:sign not found:%v",n.id,msg.SenderID,msg.ID,n.id,sig_Digest1)
                        go SendPost(content2,n.table[k]+server.QBACommitEntry)
                        return
                    }
                    //signature not found
                    if ErrorLabel[1]==false{
                        //use commit resule to pack the commit message to send back to forwarder node
                        content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"verify failed:signature not found",msg.QdsNum)
                     
                        //pack error check
                        if err2 != nil {
                            log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                            return
                        }
                        log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Verify failed:sign not found:%v",n.id,msg.SenderID,msg.ID,n.id,sig_Digest2)
                        go SendPost(content2,n.table[k]+server.QBACommitEntry)
                        return
                    }    
                }
            } 
            for k,_ := range n.buffer.QBACommitBuffer{
                //use commit resule to pack the commit message to send back to forwarder node
                content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"all signature valid",msg.QdsNum)
                     
                //pack error check
                if err2 != nil {
                    log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                    return
                }
                log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]all signature valid",n.id,msg.SenderID,msg.ID,n.id)
                    go SendPost(content2,n.table[k]+server.QBACommitEntry)
            } 
            n.buffer.ClearQBACommitMsg()
            n.QBAMessageBackup=n.QBAMessageBackup[:0]
        }
    }
}

func (n *Node) handleExecute(msg *message.QBACommitMessage) {
    //main node:execute to write block
    if n.id == n.GetPrimary() || n.id == message.Identify(len(n.table)-1) {
        n.buffer.BufferQBACommitMsg(msg)
        if len(n.buffer.QBACommitBuffer)==(len(n.table)-2){
            log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Consensus succeed",n.id,msg.SenderID,msg.ID,n.id)
            n.buffer.AppendToExecuteQueue(&n.QBAMainMessage[0])
            n.executeNotify<-true
            n.buffer.ClearQBACommitMsg()
        }
    }
    //main node:execute to write block
    if n.id != n.GetPrimary() && n.id != message.Identify(len(n.table)-1){
        var SavedMessage [][]byte
        for i:=0;i<(len(n.table)-2);i++{
            SavedMessage=append(SavedMessage,n.QBAMessageBackup[0].Digest[i*(message.MessageLength+4*message.KeyLength):i*(message.MessageLength+4*message.KeyLength)+message.MessageLength])
        }
        maxCount:=0
        for i:=0;i<(len(n.table)-2);i++{
            if bytes.Equal(SavedMessage[i],n.QBAMainMessage[0].Digest){
                maxCount++
            }
        }
        if maxCount>= int((len(n.table)-2)/2)+1{
            content2, _, err2 := message.NewQBACommitMessage(msg.Digest,n.id,msg.SenderID,n.QBAKey_X,n.QBAKey_Y,n.QBAKey_Z,"all signature valid",msg.QdsNum)
             
            //pack error check
            if err2 != nil {
                log.Printf("[Verify and Infer Phase]generate qba-commit message error")
                return
            }
            log.Printf("[Verify and Infer Phase][cur_node:%v,signer:%v,forward:%v,verifier:%v]Consensus succeed",n.id,msg.SenderID,msg.ID,n.id)
            go SendPost(content2,n.table[n.GetPrimary()]+server.QBACommitEntry)
            go SendPost(content2,n.table[message.Identify(len(n.table)-1)]+server.QBACommitEntry)
            n.QBAMessageBackup=n.QBAMessageBackup[:0]
            n.buffer.AppendToExecuteQueue(&n.QBAMainMessage[0])
            n.executeNotify<-true
        }
    }
}

