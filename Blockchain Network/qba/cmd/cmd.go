package cmd

import (
	"flag"
	"github.com/hyperledger/fabric/orderer/consensus/qba/message"
	"os"
	"strconv"
	"strings"
)

type SharedConfig struct {
	ClientServer   bool
	Port           int
	Id             message.Identify
	Table		   map[message.Identify]string
	FaultNum	   uint
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
		FaultNum:      uint(len(t)/3),
	}
}


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

