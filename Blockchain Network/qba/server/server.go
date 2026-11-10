package server

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"log"
	"net/http"
	"strconv"
)

const (
	RequestEntry= "/request"
	QBAPrepareEntry= "/qbapreprepare"
	QBAMessageEntry= "/qbamessage"
	QBAForwardEntry= "/qbaforward"
	QBACommitEntry= "/qbacommit"
	QBAVoteEntry= "/qbavote"
	// Admin endpoints used by the malicious-node demonstration page
	FaultEntry  = "/fault"
	PeersEntry  = "/faultpeers"
	StatusEntry = "/status"
	TraceEntry  = "/trace"
)

// HTTP request listener
type HttpServer struct {
	port   int
	server *http.Server
	requestRecv       chan *message.Request
	QBAPrepareRecv    chan *message.QBAPrepareMessage
	QBAMessageRecv    chan *message.QBAMessage
	QBAForwardRecv    chan *message.QBAForwardMessage
	QBACommitRecv     chan *message.QBACommitMessage
	QBAVoteRecv       chan *message.QBAVoteMessage

	// admin handlers supplied by the node through RegisterAdmin
	faultHandle  func(http.ResponseWriter, *http.Request)
	peersHandle  func(http.ResponseWriter, *http.Request)
	statusHandle func(http.ResponseWriter, *http.Request)
	traceHandle  func(http.ResponseWriter, *http.Request)
}

func NewServer(cfg *cmd.SharedConfig) *HttpServer {
	httpServer := &HttpServer{
		port:   cfg.Port,
		server: nil,
	}
	// set server
	return httpServer
}

// config server: to register the handle chan
func (s *HttpServer) RegisterChan(r    chan *message.Request,p chan *message.QBAPrepareMessage,m chan *message.QBAMessage,f chan *message.QBAForwardMessage, c chan *message.QBACommitMessage,v chan *message.QBAVoteMessage) {
	log.Printf("[Server] register the chan for listen func")
	s.requestRecv           = r
	s.QBAMessageRecv          = m
	s.QBAPrepareRecv          = p
	s.QBAForwardRecv   = f
	s.QBACommitRecv    = c
	s.QBAVoteRecv      = v
}

// RegisterAdmin installs the demonstration-only fault control and status handlers.
// They are passed in rather than imported because the node package already
// depends on this one.
func (s *HttpServer) RegisterAdmin(fault, peers, status, trace func(http.ResponseWriter, *http.Request)) {
	s.faultHandle = fault
	s.peersHandle = peers
	s.statusHandle = status
	s.traceHandle = trace
}

func (s *HttpServer) Run() {
	// register server service and run
	log.Printf("[Node] start the listen server")
	s.registerServer()
}

func (s *HttpServer) registerServer() {
	log.Printf("[Server] set listen port:%d\n", s.port)

	httpRegister := map[string]func(http.ResponseWriter, *http.Request){
		RequestEntry:    s.HttpRequest,
		QBAMessageEntry: s.QBAMessage,
		QBAPrepareEntry: s.QBAPrepare,
		QBAForwardEntry:    s.QBAForward,
		QBACommitEntry:     s.QBACommit,
		QBAVoteEntry: s.QBAVote,
	}

	if s.faultHandle != nil {
		httpRegister[FaultEntry] = s.faultHandle
	}
	if s.peersHandle != nil {
		httpRegister[PeersEntry] = s.peersHandle
	}
	if s.statusHandle != nil {
		httpRegister[StatusEntry] = s.statusHandle
	}
	if s.traceHandle != nil {
		httpRegister[TraceEntry] = s.traceHandle
	}

	mux := http.NewServeMux()
	for k, v := range httpRegister {
		log.Printf("[Server] register the func for %s", k)
		mux.HandleFunc(k, v)
	}

	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: mux,
	}

	if err := s.server.ListenAndServe(); err != nil {
		log.Printf("[Server Error] %s", err)
		return
	}
}
