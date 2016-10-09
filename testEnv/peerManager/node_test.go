package peerManager

import (
	"testing"
	//"time"
	//"fmt"
	//"fabric/peer/node"
	//"hyperchain/consensus/testEnv/connectInit"
	//"hyperchain/consensus/testEnv/eventManager"
	"time"
	"hyperchain/protos"
	"hyperchain/consensus/pbft"
	"github.com/golang/protobuf/proto"
	"strings"
	"context"
	"hyperchain/consensus/testEnv/connectInit"
	"fmt"
)

//func TestNode(t *testing.T){
//	addr:=[]string{
//		"127.0.0.1:2333",
//		"127.0.0.1:2334",
//		"127.0.0.1:2335",
//		"127.0.0.1:2336",
//	}
//	//mapc:=make(map[string] int)
//	//for c:=len(addr);c>0;c--{
//	//	p:=string(addr[c-1])
//	//	mapc[p]=c-1
//	//}
//	//fmt.Println("*****>start")
//	//mac1:=make(map[string] int)
//	n1:=NewNode("127.0.0.1","2333",addr)
//	n2:=NewNode("127.0.0.1","2334",addr)
//	n3:=NewNode("127.0.0.1","2335",addr)
//	n4:=NewNode("127.0.0.1","2336",addr)
//
//	//go n1.Start()
//	//go n2.Start()
//	//go n3.Start()
//	//go n4.Start()
//	fmt.Println("*****>netInit")
//	n1.NodeNetInit()
//	n2.NodeNetInit()
//	n3.NodeNetInit()
//	n4.NodeNetInit()
//	fmt.Println("*****>broadCast")
//	time.Sleep(3*time.Second)
//	n1.BroadCast(&connectInit.Message{})
//	n2.BroadCast(&connectInit.Message{})
//
//	time.Sleep(100*time.Second)
//
//
//
//}

func TestCheckPoint(t *testing.T){
	exit := make(chan int)
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
		}

	node0:= NewNode("127.0.0.1",2333,addr)
	node1:= NewNode("127.0.0.1",2334,addr)
	node2:= NewNode("127.0.0.1",2335,addr)
	node3:= NewNode("127.0.0.1",2336,addr)

	go node0.StartNode()
	go node1.StartNode()
	go node2.StartNode()
	go node3.StartNode()

	//txMsg := &protos.Message{
	//	Type: protos.Message_TRANSACTION,
	//	Payload: []byte{0x00, 0x00, 0x03},
	//	Timestamp: time.Now().UnixNano(),
	//	Id: 1,
	//}
	txMsg1 := &protos.Message{
	Type: protos.Message_TRANSACTION,
	Payload: []byte{0x00, 0x00, 0x04},
	Timestamp: time.Now().UnixNano(),
	Id: 1,
	}
	//for i:=0;i<2;i++{

	//go node2.UniCastToSelf(txMsg)
	//time.Sleep(2*time.Second)

	go node2.UniCastToSelf(txMsg1)

//}
	for i:=3;i>0;i--{
		chkpt := &pbft.Checkpoint{
		SequenceNumber: 6,
		ReplicaId:      uint64(i),
		Id:             "hsdkhskhfaere",
		}
		msg:=&pbft.Message{Payload: &pbft.Message_Checkpoint{Checkpoint: chkpt}}

		consensusMsg := &pbft.ConsensusMessage{Payload: &pbft.ConsensusMessage_PbftMessage{PbftMessage: msg}}

		msgPayload, err := proto.Marshal(consensusMsg)

		if err != nil {
		logger.Errorf("ConsensusMessage Marshal Error", err)
		return
		}
		pbMsg := &protos.Message{
		Type:		protos.Message_CONSENSUS,
		Payload:	msgPayload,
		Timestamp:	time.Now().UnixNano(),
		Id:		uint64(i),
		}
		node2.UniCastToID(4,pbMsg)
	}

	//to ensure new request batch
	time.Sleep(5*time.Second)
	for i:=0;i<18;i++{
		txMsg2 := &protos.Message{
			Type: protos.Message_TRANSACTION,
			Payload: []byte{0x00, 0x00, 0x05},
			Timestamp: time.Now().UnixNano(),
			Id: 1,
		}

		node3.UniCastToSelf(txMsg2)
	}


	time.Sleep(5*time.Second)

	txMsg8 := &protos.Message{
		Type: protos.Message_TRANSACTION,
		Payload: []byte{0x00, 0x00, 0x08},
		Timestamp: time.Now().UnixNano(),
		Id: 1,
	}

	node3.UniCastToSelf(txMsg8)

	<-exit
}


func TestAssign(t *testing.T){

	exit := make(chan int)

	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
		}

	node0:= NewNode("127.0.0.1",2333,addr)
	node1:= NewNode("127.0.0.1",2334,addr)
	node2:= NewNode("127.0.0.1",2335,addr)
	node3:= NewNode("127.0.0.1",2336,addr)

	go node0.StartNode()
	go node1.StartNode()
	go node2.StartNode()
	go node3.StartNode()

	time.Sleep(2*time.Second)

	for i:=0;i<1000;i++ {
		msg := &protos.Message{
			Type: protos.Message_TRANSACTION,
			Payload: []byte{0x00, 0x00, 0x05},
			Timestamp: time.Now().UnixNano(),
			Id: 1,
		}
		node1.UniCastToSelf(msg)
	}
	<-exit
	//fmt.Println("addr:",node.LocalAddr.address,"cc",node.LocalID)
}
func TestUpState(t *testing.T){
	exit := make(chan int)
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
	}

	node0:= NewNode("127.0.0.1",2333,addr)
	node1:= NewNode("127.0.0.1",2334,addr)
	node2:= NewNode("127.0.0.1",2335,addr)
	node3:= NewNode("127.0.0.1",2336,addr)

	go node0.StartNode()
	go node1.StartNode()
	go node2.StartNode()
	go node3.StartNode()

	msg := &protos.Message{
		Type: protos.Message_TRANSACTION,
		Payload: []byte{0x00, 0x00, 0x03},
		Timestamp: time.Now().UnixNano(),
		Id: 1,
	}

	node1.BroadCast(msg)

	time.Sleep(2*time.Second)
	upMsg :=&protos.UpdateStateMessage{
		SeqNo: 3,
		TargetId: []byte{0x00, 0x00, 0x03},
		Replicas: []uint64{0,2},
	}
	//tmpMsg, err := proto.Marshal(upMsg)

	node2.EB.UpdateState(upMsg)

	<-exit
}

func startNodes() map[int] *Node{
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
	}
	nodes:=make(map[int] *Node)
	port:=2333
	for key,value:=range addr{
		a:=strings.Split(value,":")
		nodes[key]=NewNode(a[0],port,addr)
		port++
	}
	for _,v:=range nodes{
		go v.StartNode()
	}
	return nodes
}
func TestUnicast(t *testing.T){
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
	}
	nodes:=make(map[int] *Node)
	port:=2333
	for key,value:=range addr{
		a:=strings.Split(value,":")
		nodes[key]=NewNode(a[0],port,addr)
		port++
	}
	for _,v:=range nodes{
		go v.StartNode()
	}
	msg := &protos.Message{
		Type: protos.Message_TRANSACTION,
		Payload: []byte{0x00, 0x00, 0x03},
		Timestamp: time.Now().UnixNano(),
		Id: 1,
	}
	for {
		if nodes[0].NetState {
			break
		}
	}
	for i:=1000;i>0;i++{
		nodes[0].UniCastToID(2,msg)
	}


}
func TestStatistic(t *testing.T){
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

	conn,err:=GetConnectionWithAddr("127.0.0.1:2336")
	if err!=nil{
		return
	}
	client:=GetClientWithConn(conn)

	msg1:=&connectInit.StatMessage{
		Type:"info",
	}

	response,err0:=client.Statistic(context.Background(),msg1)

	if err0!=nil{
		return
	}
	fmt.Println("",string(response.JSONInfo))

	<-exit

}
