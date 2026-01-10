package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	cb "github.com/hyperledger/fabric/protos/common"
	"log"
)

var test_reqeust_num uint64 = 0

func (n *Node) executeAndReplyThread() {
	for {
		select {
		case <-n.executeNotify:
			// execute batch
			batchs:= n.buffer.BatchExecute()
			if len(batchs) == 0 {
				log.Printf("[Reply] lost sequence now")
				continue
			}
                       requestBatchs := make([]*message.Request, 0)
			for _, b := range batchs {
				requestBatchs = append(requestBatchs, b.Message.Requests...)
			}
			// map the digest to request
			test_reqeust_num = test_reqeust_num + uint64(len(requestBatchs))
			log.Printf("[Reply]execute request(%d)", test_reqeust_num)
			// pending state
			pending := make(map[string]bool)
			for _, r := range requestBatchs {
				msg		  := r.Op.Envelope
				channel   := r.Op.ChannelID
				configSeq := r.Op.ConfigSeq
				switch r.Op.Type {
				case message.TYPECONFIG:
					var err error
					seq := n.supports[channel].Sequence()
					if configSeq < seq {
						if msg, _, err = n.supports[r.Op.ChannelID].ProcessConfigMsg(r.Op.Envelope); err != nil {
							log.Println(err)
						}
					} 
					batch := n.supports[channel].BlockCutter().Cut()
					if batch != nil {
						block := n.supports[channel].CreateNextBlock(batch)
						n.supports[channel].WriteBlock(block, nil)
					}
					pending[channel] = false
					// write config block
					block := n.supports[channel].CreateNextBlock([]*cb.Envelope{msg})
					log.Printf("Create Block Number:%v",block.Header.Number)
					n.supports[channel].WriteConfigBlock(block, nil)
				case message.TYPENORMAL:
					seq := n.supports[channel].Sequence()
					if configSeq < seq {
						if _, err := n.supports[channel].ProcessNormalMsg(msg); err != nil {
						}
					}
					batches, p := n.supports[channel].BlockCutter().Ordered(msg)
					for _, batch := range batches {
						block := n.supports[channel].CreateNextBlock(batch)
						log.Printf("Create Block Number:%v",block.Header.Number)
						n.supports[channel].WriteBlock(block, nil)
					}
					pending[channel] = p
				}
			}
			for k, v := range pending {
				if v {
					batch := n.supports[k].BlockCutter().Cut()
					if batch != nil {
						block := n.supports[k].CreateNextBlock(batch)
						log.Printf("Create Block Number:%v",block.Header.Number)
						n.supports[k].WriteBlock(block, nil)
					}
				}
			}
			log.Printf("Finish Write Block")
			n.ConsensusCount++
		}
	}
}
