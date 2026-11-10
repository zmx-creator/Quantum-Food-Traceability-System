package cmd

import (
	"flag"
	"log"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"os"
	"strconv"
	"strings"
	"time"
)

type SharedConfig struct {
	ClientServer   bool
	Port           int
	Id             message.Identify
	Table		   map[message.Identify]string
	FaultNum	   uint
	RoundTimeout time.Duration
	MaxRetry     int
	// ConsistencyRetry bounds how many times one node may be asked in a single
	// round to send an inconsistent message again
	ConsistencyRetry int
}

func ReadConfig() *SharedConfig {
	port, _ := GetConfigurePort()
	id, _   := GetConfigureID()
	table, _ := GetConfigureTable()

	t := make(map[message.Identify]string)
	for k, v := range table {
		t[message.Identify(k)] = v
	}

	flag.Parse()
	return &SharedConfig{
		Port: 		   port,
		Id:            message.Identify(id),
		Table:         t,
		// QBA broadcasts messages that are signed with QDS, so it sits in the
		// authenticated Byzantine agreement model: the bound is n >= 2f+1 rather
		// than the unsigned n >= 3f+1. Five nodes therefore tolerate f=2.
		FaultNum:      uint((len(t) - 1) / 2),
		RoundTimeout:  GetConfigureRoundTimeout(),
		MaxRetry:      GetConfigureMaxRetry(),
		ConsistencyRetry: GetConfigureConsistencyRetry(),
	}
}

// Read configuration
func GetConfigureID() (id int, err error){
	rawID  := os.Getenv("PBFT_NODE_ID")
	if id, err = strconv.Atoi(rawID); err != nil {
		return
	}
	return
}

func GetConfigureTable() (map[int]string, error){
	rawTable  := os.Getenv("PBFT_NODE_TABLE")
	nodeTable := make(map[int]string, 0)

	tables := strings.Split(rawTable, ";")
	for index, t := range tables {
		nodeTable[index] = t
	}
	return nodeTable, nil
}

func GetConfigurePort() (port int, err error){
	rawPort := os.Getenv("PBFT_LISTEN_PORT")
	if port, err = strconv.Atoi(rawPort); err != nil {
		return
	}
	return
}

// GetConfigureRoundTimeout reads QBA_ROUND_TIMEOUT_MS, defaulting to five seconds.
func GetConfigureRoundTimeout() time.Duration {
	raw := os.Getenv("QBA_ROUND_TIMEOUT_MS")
	if raw == "" {
		return 5 * time.Second
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		log.Printf("[Config] invalid QBA_ROUND_TIMEOUT_MS %q, using 5000", raw)
		return 5 * time.Second
	}
	return time.Duration(ms) * time.Millisecond
}

// GetConfigureMaxRetry reads QBA_MAX_RETRY, defaulting to five attempts.
//
// Five is the number the demonstration reports: a batch whose quorum never
// arrives is retransmitted five times and then discarded as untrustworthy.
func GetConfigureMaxRetry() int {
	raw := os.Getenv("QBA_MAX_RETRY")
	if raw == "" {
		return 5
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		log.Printf("[Config] invalid QBA_MAX_RETRY %q, using 5", raw)
		return 5
	}
	return n
}

// GetConfigureConsistencyRetry reads QBA_CONSISTENCY_RETRY, defaulting to one
// resend request per inconsistent signer per round. The bound is what keeps the
// protocol's "make it send again" step from becoming a message flood against an
// honest node.
func GetConfigureConsistencyRetry() int {
	raw := os.Getenv("QBA_CONSISTENCY_RETRY")
	if raw == "" {
		return 1
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		log.Printf("[Config] invalid QBA_CONSISTENCY_RETRY %q, using 1", raw)
		return 1
	}
	return n
}
