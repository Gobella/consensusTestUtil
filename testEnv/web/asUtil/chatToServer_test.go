package asUtil

import (
	"hyperchain/consensus/testEnv/peerManager"
	"hyperchain/protos"
	"strings"
	"testing"
	"time"
)

func startNodes() map[int] *peerManager.Node{
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
	}
	nodes:=make(map[int] *peerManager.Node)
	port:=2333
	for key,value:=range addr{
		a:=strings.Split(value,":")
		nodes[key]=peerManager.NewNode(a[0],port,addr)
		port++
	}
	for _,v:=range nodes{
		go v.StartNode()
	}
	return nodes
}

func TestChat(t *testing.T){
	exit := make(chan int)
	nodes:=startNodes()
	time.Sleep(2*time.Second)
	msg := &protos.Message{
		Type: protos.Message_TRANSACTION,
		Payload: []byte{0x00, 0x00, 0x03},
		Timestamp: time.Now().UnixNano(),
		Id: 1,
	}

	nodes[2].BroadCast(msg)
	time.Sleep(2*time.Second)
	info,_:=ChatWithServer(ReadServerConfig(2))
	logger.Info("the statistic information is ",string(info))
	<-exit
}