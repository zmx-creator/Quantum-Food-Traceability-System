package message

import (
    //"log"
    //"bytes"
   // "sync"
)

type Buffer struct {
	requestQueue  []*Request

	QBACommitBuffer map[Identify]*QBACommitMessage
        executeQueue  []*QBAMessage
}

func NewBuffer() *Buffer {
	return &Buffer{
		requestQueue:  make([]*Request, 0),

		QBACommitBuffer: make(map[Identify]*QBACommitMessage),
		
		executeQueue:  make([]*QBAMessage,0),

	}
}

// buffer about request
func (b *Buffer) AppendToRequestQueue(req *Request) {
	b.requestQueue = append(b.requestQueue, req)
}

func (b *Buffer) AppendToExecuteQueue(req *QBAMessage) {
	b.executeQueue = append(b.executeQueue, req)
}


func (b *Buffer) BatchRequest() (batch []*Request) {
	batch = make([]*Request, 0)
	for _, r := range b.requestQueue {
		batch = append(batch, r)
	}
	b.requestQueue = make([]*Request, 0)
	return
}

func (b *Buffer) BatchExecute() (batch []*QBAMessage) {
	batch = make([]*QBAMessage, 0)
	for _, r := range b.executeQueue {
		batch = append(batch, r)
	}
	b.executeQueue = make([]*QBAMessage, 0)
	return
}


func (b *Buffer) SizeofRequestQueue() (l int) {
	l = len(b.requestQueue)
	return  l
}

func (b *Buffer) ClearRequestQueue() {
	b.requestQueue = make([]*Request, 0)
}

func (b *Buffer) ClearExecuteQueue() {
	b.executeQueue = make([]*QBAMessage, 0)
}
//-------------------------------------------
// buffer about pre-prepare
func (b *Buffer) BufferQBACommitMsg(msg *QBACommitMessage) {
	b.QBACommitBuffer[msg.ID] = msg
}

func (b *Buffer) ClearQBACommitMsg() {
	b.QBACommitBuffer=make(map[Identify]*QBACommitMessage)
}

func (b *Buffer) IsExistQBACommitMsg(id Identify) bool {
	if len(b.QBACommitBuffer) == 0{
	    return false
	}
	if len(b.QBACommitBuffer) > 0{
	    for k,_:=range b.QBACommitBuffer{
	        if k==id {
	            return true
	        }
            }
	}
	return false
}


