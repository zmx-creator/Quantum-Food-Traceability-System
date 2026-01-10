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

	buffer         *message.Buffer          

	QBAMainMessage   []message.QBAMessage
	QBAMessageBackup   []message.QBAForwardMessage
	ConsensusCount  int

	QBAKey_X     []byte       
	QBAKey_Y     []byte     
	QBAKey_Z     []byte     
	
	//QBANodeList  []message.Identify

	requestRecv             chan *message.Request    
	QBAPrepareRecv   chan *message.QBAPrepareMessage
	QBAMessageRecv          chan *message.QBAMessage
	QBAForwardRecv chan *message.QBAForwardMessage    
	QBACommitRecv   chan *message.QBACommitMessage 

	QBAMessageSendNotify      chan bool   
	executeNotify          chan bool      

	supports             map[string]consensus.ConsenterSupport
}

func NewNode(cfg *cmd.SharedConfig, support consensus.ConsenterSupport) *Node {
	node := &Node{
		// config
		cfg:	  cfg,
		// http server
		server:   server.NewServer(cfg),
		// information about node
		id:       cfg.Id,
		table:	  cfg.Table,
		faultNum: cfg.FaultNum,
		// the message buffer to store msg
		buffer: message.NewBuffer(), 
		ConsensusCount:1,
		QBAKey_X:make([]byte,message.KeyLength), 
	        QBAKey_Y:make([]byte,message.KeyLength),  
          	QBAKey_Z:make([]byte,message.KeyLength),
          	
          	
          	//QBANodeList:make([]message.Identify,0), 
		QBAMainMessage:make([]message.QBAMessage,0), 
		QBAMessageBackup:make([]message.QBAForwardMessage,0), 
		requestRecv:make(chan *message.Request),           
	        QBAPrepareRecv:make(chan *message.QBAPrepareMessage),
	        QBAMessageRecv:make(chan *message.QBAMessage),
	        QBAForwardRecv:make(chan *message.QBAForwardMessage),    
	        QBACommitRecv:make(chan *message.QBACommitMessage), 
	    	// chan for notify pre-prepare send thread
		QBAMessageSendNotify: make(chan bool),
		// chan for notify execute op and reply thread
		executeNotify:        make(chan bool),
		supports: 			  make(map[string]consensus.ConsenterSupport),
	}
	log.Printf("[Node] the node id:%d, fault number:%d\n", node.id, node.faultNum)
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
	n.server.RegisterChan(n.requestRecv,n.QBAPrepareRecv, n.QBAMessageRecv, n.QBAForwardRecv, n.QBACommitRecv)
	go n.server.Run()
	go n.requestRecvThread()
	go n.QBAMessageSendOneThread()
	go n.QBAForwardThread()
	go n.QBAVerifyAndCommitThread()
	go n.QBACommitRecvThread()
	go n.QBAMessageSendTwoThread()
	go n.executeAndReplyThread()
}
