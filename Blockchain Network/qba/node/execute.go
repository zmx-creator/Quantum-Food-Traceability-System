package node

// ready to execute the msg(digest) send to execute queue
func (n *Node) readytoExecute(digest string) {
	// buffer to ExcuteQueue
	n.executeNotify<-true
}
