package message

import (
    "log"
    "bytes"
   // "sync"
)

type Buffer struct {
	requestQueue  []*Request

	QBACommitBuffer map[Identify]*QBACommitMessage
        

	QBAVoteBuffer map[Identify]*QBAVoteMessage
	
        executeQueue  []*QBAMessage
}

func NewBuffer() *Buffer {
	return &Buffer{
		requestQueue:  make([]*Request, 0),

		QBACommitBuffer: make(map[Identify]*QBACommitMessage),

		QBAVoteBuffer:    make(map[Identify]*QBAVoteMessage),
		
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


//-------------------------------------------

// buffer about pre-prepare
func (b *Buffer) BufferQBAVoteMsg(msg *QBAVoteMessage) {
	b.QBAVoteBuffer[msg.ID] = msg
}

func (b *Buffer) ClearQBAVoteMsg() {
	b.QBAVoteBuffer=make(map[Identify]*QBAVoteMessage)
}

func (b *Buffer) IsExistQBAVoteMsg(id Identify) bool {
	if len(b.QBAVoteBuffer) == 0{
	    return false
	}
	if len(b.QBAVoteBuffer) > 0{
	    for k,_:=range b.QBAVoteBuffer{
	        if k==id {
	            return true
	        }
            }
	}
	return false
}

func (b *Buffer) Consensus(digest []byte)bool{
    var result bool
    ConsensusCount:=0
    for _,v := range b.QBAVoteBuffer {
        if bytes.Equal(v.Result.Digest,digest) {
            ConsensusCount+=1
        }
    }
    if ConsensusCount ==len(b.QBAVoteBuffer){
        for _,v := range b.QBAVoteBuffer {
            if v.Result.AllVoteEqual != true {
                result=true
            }
            if v.Result.AllVoteEqual == true {
                result=false
            }
        }
    }
    if ConsensusCount <len(b.QBAVoteBuffer){
        result=false
    }
    return result
}

func (b *Buffer) FetchQBAVoteMsg(id Identify) (ret *QBAVoteMessage) {
	ret = nil
	if _, ok := b.QBAVoteBuffer[id]; !ok {
		log.Printf("[Buffer] error to find qba-vote msg from id (%d)", id)
		return
	}
	ret = b.QBAVoteBuffer[id]
	return ret
}

func (b *Buffer) Vote() (*VoteResult) {
	result := b.Majority()
	return result
}

func (b *Buffer) Infer() (*VoteResult) {
	result := b.MajorityFromVoteList()
	return result
}



// isValidVerdict reports whether a buffered message may be counted, either in a
// judgment list or in the commit arithmetic.
//
// Protocol description of QBA, section I.B step 4 and section III: a QDS that
// fails verification does not invalidate the round, it simply contributes
// nothing to the judgment list, and the honest parties decide by majority over
// what they kept. A vote that ends in a tie is the abstention symbol of the
// protocol and is not a decision either. Entries are therefore filtered here
// rather than aborting the round when the first mismatch appears.
func isValidVerdict(commitType string) bool {
	return commitType != "verify failed" && commitType != "Consensus Failed"
}

// CountAgreeing returns how many buffered commits answered with a countable
// verdict, and how many of those carry exactly the given decision digest.
//
// The distinction is the whole point of the commit rule: every node answers,
// including a malicious one, so counting answers measures liveness while
// counting matching decisions measures agreement.
func (b *Buffer) CountAgreeing(digest []byte) (agreeing int, answered int) {
	for _, v := range b.QBACommitBuffer {
		if v == nil {
			continue
		}
		answered++
		if isValidVerdict(v.CommitType) && bytes.Equal(v.Digest, digest) {
			agreeing++
		}
	}
	return
}

func (b *Buffer) Majority() *VoteResult {
    maps := make(map[Identify][]byte, len(b.QBACommitBuffer))
    for k, v := range b.QBACommitBuffer {
        if v != nil && isValidVerdict(v.CommitType) {
            maps[k] = v.Digest
        }
    }

    valueCount := make(map[string]int) // Changed to use string keys since []byte isn't comparable
    for _, v := range maps {
        valueCount[string(v)]++ // Convert []byte to string for map key
    }
    if len(valueCount) == 0 {
        return nil
    }

    maxCount := 0
    for _, count := range valueCount {
        if count > maxCount {
            maxCount = count
        }
    }
    
    var maxValues []byte
    AllEqual := false
    for vStr, count := range valueCount {
        if count == maxCount {
            v := []byte(vStr) // Convert back to []byte
            if maxValues == nil {
                maxValues = v
            }
            for v1Str, count1 := range valueCount {
                if count1 == maxCount && !bytes.Equal(v, []byte(v1Str)) {
                    AllEqual = true
                    break
                }
            }
        }
        if AllEqual {
            break
        }
    }
    
    var maxVote Identify
    for _, v := range b.QBACommitBuffer {
        if v != nil && bytes.Equal(v.Digest, maxValues) {
            maxVote = v.SenderID
            break // Assuming we just need one sender ID
        }
    }
    
    result := &VoteResult{
        Digest:       maxValues,
        SenderID:     maxVote,
        AllVoteEqual: AllEqual,
    }
    return result
}

func (b *Buffer) MajorityFromVoteList() *VoteResult {
    // Build the Identify -> Digest mapping
    maps := make(map[Identify][]byte, len(b.QBAVoteBuffer))
    for k, v := range b.QBAVoteBuffer {
        if v != nil {
            maps[k] = v.Result.Digest
        }
    }

    // Count how often each Digest occurs
    valueCount := make(map[string]int) // []byte cannot be used as a map key, so convert to string
    for _, v := range maps {
        valueCount[string(v)]++
    }
    
    if len(valueCount) == 0 {
        return nil
    }

    // Find the highest occurrence count
    maxCount := 0
    for _, count := range valueCount {
        if count > maxCount {
            maxCount = count
        }
    }

    // Collect every Digest that reaches the maximum count
    var maxValues [][]byte
    for vStr, count := range valueCount {
        if count == maxCount {
            maxValues = append(maxValues, []byte(vStr))
        }
    }

    // Detect a tie: several distinct Digests sharing the highest vote count
    allEqual := len(maxValues) > 1

    // If there is a single maximum, look up its SenderID
    var maxVote Identify
    if len(maxValues) == 1 {
        for _, v := range b.QBAVoteBuffer {
            if v != nil && bytes.Equal(v.Result.Digest, maxValues[0]) {
                maxVote = v.Result.SenderID
                break
            }
        }
    }

    // Assemble the result
    result := &VoteResult{
        Digest:        maxValues[0], // take the first maximum
        SenderID:      maxVote,
        AllVoteEqual:  allEqual,
    }
    return result
}
