package node

import (
	"log"
)

func (n *Node) requestRecvThread() {
    log.Printf("[Node] start recv the request thread")
        for {
            msg := <- n.requestRecv
            n.buffer.AppendToRequestQueue(msg)
            n.QBAMessageSendNotify <- true
     }
}


