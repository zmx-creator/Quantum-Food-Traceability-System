package node

import (
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/qba/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
)

var GNode *Node = nil

type Node struct {
	cfg    *cmd.SharedConfig          
	server *server.HttpServer       

	id       message.Identify         
	table    map[message.Identify]string        
	faultNum uint
	fault    *FaultState
	round    *Round
	trace    *Trace
	maxRetry int

	// consistency check bookkeeping, and how often one node may be asked per
	// round to send a message again
	consistency      *Consistency
	consistencyRetry int

	roundContent []byte
	roundTarget message.Identify                      

	buffer         *message.Buffer          
	
	IsConsensusRunning  bool

	QBAMainMessage   []message.QBAMessage
	ConsensusCount  int

	QBAKey_X     []byte       
	QBAKey_Y     []byte     
	QBAKey_Z     []byte     
	
	Signer_KeyCache  KeyCache
	Forwarder_KeyCache  KeyCache
	Verifier_KeyCache  KeyCache

	requestRecv             chan *message.Request    
	QBAPrepareRecv   chan *message.QBAPrepareMessage
	QBAMessageRecv          chan *message.QBAMessage
	QBAForwardRecv chan *message.QBAForwardMessage    
	QBACommitRecv   chan *message.QBACommitMessage 
	QBAVoteRecv      chan *message.QBAVoteMessage 

	QBAMessageSendNotify      chan bool   
	executeNotify          chan bool      

	supports             map[string]consensus.ConsenterSupport
}

func NewNode(cfg *cmd.SharedConfig, support consensus.ConsenterSupport) *Node {
    node := &Node{
        // config
        cfg:cfg,
        // http server
        server:   server.NewServer(cfg),
        // information about node
        id:       cfg.Id,
        table:	  cfg.Table,
        faultNum: cfg.FaultNum,
        fault: NewFaultState(),
        round: NewRound(cfg.RoundTimeout, cfg.MaxRetry),
        trace: NewTrace(),
        maxRetry: cfg.MaxRetry,
        consistency: NewConsistency(),
        consistencyRetry: cfg.ConsistencyRetry,
        // the message buffer to store msg
        buffer: message.NewBuffer(), 
        IsConsensusRunning:false,
        ConsensusCount:1,
        QBAKey_X:make([]byte,8), 
        QBAKey_Y:make([]byte,8),  
        QBAKey_Z:make([]byte,8),
          	
        Signer_KeyCache:KeyCache {
        cacheData:make([][]byte,3*KeyCacheWindowLength*36),   
            cacheStart:1,         
            cacheEnd:3*KeyCacheWindowLength*36,     
            filename:"/QBAKey/key_signer.txt",   
            W:3*KeyCacheWindowLength*36,
            N:36,      
            W0:KeyCacheWindowLength,        
            totalLines:500000,         
        },
        Forwarder_KeyCache:KeyCache {
            cacheData:make([][]byte,3*KeyCacheWindowLength*36),   
            cacheStart:1,         
            cacheEnd:3*KeyCacheWindowLength*36,     
            filename:"/QBAKey/key_forwarder.txt",   
            W:3*KeyCacheWindowLength*36,
            N:36,      
            W0:KeyCacheWindowLength,        
            totalLines:500000,         
        },
        Verifier_KeyCache:KeyCache {
            cacheData:make([][]byte,3*KeyCacheWindowLength*36),   
            cacheStart:1,         
            cacheEnd:3*KeyCacheWindowLength*36,     
            filename:"/QBAKey/key_verifier.txt",   
            W:3*KeyCacheWindowLength*36,
            N:36,      
            W0:KeyCacheWindowLength,        
            totalLines:500000,         
        },
          	
          	//QBANodeList:make([]message.Identify,0), 
		QBAMainMessage:make([]message.QBAMessage,0), 
		requestRecv:make(chan *message.Request),           
	        QBAPrepareRecv:make(chan *message.QBAPrepareMessage),
	        QBAMessageRecv:make(chan *message.QBAMessage),
	        QBAForwardRecv:make(chan *message.QBAForwardMessage),    
	        QBACommitRecv:make(chan *message.QBACommitMessage), 
	        QBAVoteRecv:make(chan *message.QBAVoteMessage), 
	    	// chan for notify pre-prepare send thread
		QBAMessageSendNotify: make(chan bool),
		// chan for notify execute op and reply thread
		executeNotify:        make(chan bool),
		supports: 			  make(map[string]consensus.ConsenterSupport),
	}
	log.Printf("[Node] the node id:%d, fault number:%d, commit quorum:%d\n", node.id, node.faultNum, len(node.table)-int(node.faultNum))
	node.RegisterChain(support)
	return node
}

func (n *Node) RegisterChain(support consensus.ConsenterSupport) {
	if _, ok := n.supports[support.ChainID()]; ok {
		return
	}
	log.Printf("[Node] Register the chain(%s)", support.ChainID())
	n.supports[support.ChainID()] = support
}



func (n *Node) Run() {
	// first register chan for server
	n.server.RegisterChan(n.requestRecv,n.QBAPrepareRecv, n.QBAMessageRecv, n.QBAForwardRecv, n.QBACommitRecv, n.QBAVoteRecv)
	// demonstration endpoints (no-op unless a fault mode is armed)
	n.server.RegisterAdmin(n.HttpFaultControl, n.HttpFaultPeers, n.HttpNodeStatus, n.HttpNodeTrace)
	go n.server.Run()
	go n.requestRecvThread()
	go n.QBAMessageSendOneThread()
	go n.QBAForwardThread()
	go n.QBAVerifyAndCommitThread()
	go n.QBACommitRecvThread()
	go n.QBAMessageSendTwoThread()
	go n.executeAndReplyThread()
	go n.roundWatchdog()
}
