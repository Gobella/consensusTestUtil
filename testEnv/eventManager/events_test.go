package eventManager

import (
	"testing"
	"fmt"
	"reflect"
	"time"
	"sync"
	"hyperchain/consensus/testEnv/connectInit"
)

var mutex sync.Mutex
var mutex1 sync.Mutex


type testEvent0 struct {
}

type testReceiver0 struct{
	count int
}

func (tr *testReceiver0) ProcessEvent(e Event) Event{
	mutex1.Lock()

	defer mutex1.Unlock()
	time.Sleep(1*time.Millisecond)
	tr.count++
	fmt.Println("testEvent0",reflect.TypeOf(e),tr.count)
	return e
}



type testEvent1 struct {
}

type testReceiver1 struct{
	count int
}

//func (tr *testReceiver1) ProcessEvent(e Event) Event{
//	mutex.Lock()
//
//	defer mutex.Unlock()
//	tr.count++
//	fmt.Println("testEvent1",reflect.TypeOf(e),tr.count)
//	return e
//}


////test the mutiEvent receiver
//func TestEvent(t *testing.T){
//	m := NewManagerImpl()
//	reciver:=&testReceiver0{}
//	m.RegistReceiver(&testEvent0{},reciver)
//	m.RegistReceiver(&testEvent1{},&testReceiver1{})
//
//	m.Start()
//
//	go concurrent(m)
//	go concurrent(m)
//	go concurrent(m)
//	go concurrent(m)
//
//	time.Sleep(100*time.Second)
//}

func concurrent(m Manager){
	fmt.Println("time start:",time.Now().UnixNano())
	for i:=0;i<10000;i++ {
		m.Queue()<-&testEvent0{}
		m.Queue()<-&testEvent1{}
	}
	fmt.Println("time end:",time.Now().UnixNano())
}

func TestEventC(t *testing.T){
	m := NewManagerImpl()
	reciver:=&testReceiver0{}
	m.RegistReceiver(ConsensusMsgEventConst,reciver)

	m.Start()
	c:=&ConsensusMsgEvent{
		Msg:&connectInit.Message{
			Timestamp:time.Now().Unix(),
		},
	}
	m.Queue()<-c

	time.Sleep(10*time.Second)
}