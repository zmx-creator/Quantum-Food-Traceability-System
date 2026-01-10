package message

import (
	"encoding/json"
	cb "github.com/hyperledger/fabric/protos/common"
)

type TimeStamp uint64 
type Identify uint64  
type View Identify    
type Sequence int64   

const TYPENORMAL = "normal"
const TYPECONFIG = "config"
const prepare = "prepare"

// Operation
type Operation struct {
    Envelope  *cb.Envelope
    ChannelID string
    ConfigSeq uint64
    Type      string
}

// Result
type Result struct {
}

// Request
type Request struct {
    Op        Operation `json:"operation"`
    TimeStamp TimeStamp `json:"timestamp"`
    ID        Identify  `json:"clientID"`
}

// Message
type Message struct {
    Requests []*Request `json:"requests"`
}

//QBAPrepareMessage:
//in Circularcollection phase
//main node send this message to notice secondary node except verifier node
//to start circuler send process and send QBA message
type QBAPrepareMessage struct {
    SenderID    Identify `json:"senderid"`  
    Message  string  `json:"message"` 
    QdsNum   int   `json:"qdsnum"`  
}

//QDS signature
type QDSSignature struct {
    Digest     []byte  `json:"digest"`  
    Factor     []byte   `json:"factor"`   
}

type QBAMessage struct {
	SenderID    Identify `json:"senderid"`  
	Message    Message  `json:"message"`
        QDSSignature   QDSSignature `json:"qdssignature"`  
        QdsNum   int   `json:"qdsnum"`  
	Digest   []byte  `json:"digest"`       
}

type QBAForwardMessage struct {
	ID         Identify   `json:"forwarderID"` 
	SenderID   Identify  `json:"senderID"` 
        Digest     []byte   `json:"digest"`  
        QDSSignature   QDSSignature  `json:"qdssignature"` 
	QBAKey_X   []byte   `json:"QBAKey_X"`  
	QBAKey_Y   []byte   `json:"QBAKey_Y"`  
	QBAKey_Z   []byte   `json:"QBAKey_Z"`  
	QdsNum   int   `json:"qdsnum"` 
}

type QBACommitMessage struct {
	ID         Identify   `json:"forwarderID"` 
	SenderID   Identify  `json:"senderID"` 
        Digest     []byte  `json:"digest"`  
	QBAKey_X   []byte   `json:"QBAKey_X"`  
	QBAKey_Y   []byte   `json:"QBAKey_Y"`  
	QBAKey_Z   []byte   `json:"QBAKey_Z"`  
	CommitType string   `json:"QBACommit"`  
	QdsNum   int   `json:"qdsnum"`  
}


func NewQBAMessage(id Identify,batch []*Request,sig QDSSignature,count int) ([]byte, *QBAMessage, []byte, error) {
	message := Message{Requests: batch}
	d, err := Digest(message)
	if err != nil {
		return []byte{}, nil, make([]byte,0), nil
	}
	qbaMessage := &QBAMessage{
		Digest:   d,  
	      SenderID:   id, 
              QDSSignature: QDSSignature{
                  Digest:sig.Digest,
                  Factor:sig.Factor,
              },  
              Message:message,
              QdsNum:count,
	}
	ret, err := json.Marshal(qbaMessage)
	if err != nil {
		return []byte{}, nil, make([]byte,0), nil
	}
	return ret, qbaMessage, d, nil
}


func NewQBAForwardMessage(d []byte,id Identify,sender Identify,sig QDSSignature,x []byte,y []byte,z []byte,count int) ([]byte, *QBAForwardMessage, error) {
	qbaForwardMessage := &QBAForwardMessage{
		ID:id,
		SenderID:sender,
		Digest:d,
		QDSSignature: QDSSignature{
                  Digest:sig.Digest,
                  Factor:sig.Factor,
              },  
		QBAKey_X:x,
	    QBAKey_Y:y,
	    QBAKey_Z:z,
	    QdsNum:count,
	}
	ret, err := json.Marshal(qbaForwardMessage)
	if err != nil {
		return []byte{}, nil, nil
	}
	return ret, qbaForwardMessage, nil
}

func NewQBACommitMessage(d []byte,id Identify,sender Identify,x []byte,y []byte,z []byte,commit string,count int) ([]byte, *QBACommitMessage, error) {
	qbaCommitMessage := &QBACommitMessage{
		ID:id,
		SenderID:sender,
		Digest:d,
		QBAKey_X:x,
	        QBAKey_Y:y,
	        QBAKey_Z:z,
		CommitType:commit,
		QdsNum:count,
	}
	ret, err := json.Marshal(qbaCommitMessage)
	if err != nil {
		return []byte{}, nil, nil
	}
	return ret, qbaCommitMessage, nil
}
func NewQBAPrepareMessage(id Identify,count int) ([]byte, *QBAPrepareMessage, error) {
	
    qbaPrepareMessage := &QBAPrepareMessage{
		SenderID:id,
		Message:prepare,
		QdsNum:count,
	}
	ret, err := json.Marshal(qbaPrepareMessage)
	if err != nil {
		return []byte{}, nil, nil
	}
	return ret, qbaPrepareMessage, nil
}
