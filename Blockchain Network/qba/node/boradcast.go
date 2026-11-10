package node

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"github.com/hyperledger/fabric/orderer/consensus/qba/server"
	"log"
	"net/http"
	"time"
)

func (n *Node) SendPrimary(msg *message.Request) {
	content, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error to marshal json")
		return
	}
	go SendPost(content, n.table[n.GetPrimary()] + server.RequestEntry)
}

func (n *Node) BroadCast(content []byte, handle string) {
	for k, v := range n.table {
		// do not send to my self
		if k == n.id {
			continue
		}
		go SendPost(content, v + handle)
	}
}

// postTimeout bounds every node-to-node POST. Without it a peer that accepts
// the connection but never answers keeps the request goroutine alive forever.
const postTimeout = 3 * time.Second

var postClient = &http.Client{Timeout: postTimeout}

// SendPost delivers content to url. Failures are logged: the round watchdog is
// what decides whether a lost message matters.
func SendPost(content []byte, url string) {
	buff := bytes.NewBuffer(content)
	resp, err := postClient.Post(url, "application/json", buff)
	if err != nil {
		log.Printf("[Send] send to %s error: %s", url, err)
		return
	}
	resp.Body.Close()
}

