package peerManager

import (
	"testing"
	"time"
	pb "hyperchain/protos"
)

func TestSendViewChange(t *testing.T){
	exit := make(chan int)
	nodes:=startNodes()

	for i:=0;i<24;i++ {
		msg := &pb.Message{
			Type: pb.Message_TRANSACTION,
			Payload: []byte{0x00, 0x00, 0x05},
			Timestamp: time.Now().UnixNano(),
			Id: 1,
		}
		nodes[0].UniCastToSelf(msg)
	}

	<-exit
}

func TestViewChangeNewViewTimer(t *testing.T){
	exit := make(chan int)
	nodes:=startNodes()
	for {
		if nodes[0].NetState {
			break
		}
	}
	//close:=func(){
	//	temp:=time.After(3*time.Second)
	//	<-temp
	//	nodes[0].Close()
	//}
	//go close()

	time.Sleep(3*time.Second)

	for i:=0;i<1000;i++ {
		msg := &pb.Message{
			Type: pb.Message_TRANSACTION,
			Payload: []byte{0x00, 0x00, 0x05},
			Timestamp: time.Now().UnixNano(),
			Id: 1,
		}
		nodes[0].UniCastToSelf(msg)
	}
	//time.Sleep(5*time.Second)
	//for i:=0;i<1000;i++ {
	//	msg := &pb.Message{
	//		Type: pb.Message_TRANSACTION,
	//		Payload: []byte{0x00, 0x00, 0x05},
	//		Timestamp: time.Now().UnixNano(),
	//		Id: 1,
	//	}
	//	nodes[0].UniCastToSelf(msg)
	//}
	<-exit
}

