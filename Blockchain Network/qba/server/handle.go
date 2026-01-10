package server

import (
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"log"
	"net/http"
)

func (s *HttpServer) HttpRequest(w http.ResponseWriter, r *http.Request) {
	var msg message.Request
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.requestRecv <- &msg
}

func (s *HttpServer) QBAPrepare(w http.ResponseWriter, r *http.Request) {
	var msg message.QBAPrepareMessage  
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.QBAPrepareRecv <- &msg
}

func (s *HttpServer) QBAMessage(w http.ResponseWriter, r *http.Request) {
	var msg message.QBAMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.QBAMessageRecv <- &msg
}

func (s *HttpServer) QBAForward(w http.ResponseWriter, r *http.Request) {
	var msg message.QBAForwardMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.QBAForwardRecv <- &msg
}

func (s *HttpServer) QBACommit(w http.ResponseWriter, r *http.Request) {
	var msg message.QBACommitMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.QBACommitRecv <- &msg
}
