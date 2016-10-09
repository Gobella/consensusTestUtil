package consensusHelper

import (
	"hyperchain/consensus/testEnv/eventManager"
	"hyperchain/consensus"
	"github.com/golang/protobuf/proto"
	"fmt"
	"hyperchain/consensus/testEnv/statistics"

	pb "hyperchain/protos"
	"time"
)

type ConsensusMsgHelper struct {
	batch consensus.Consenter
	logstatic *statistics.LogStatic
}
func (cmh *ConsensusMsgHelper) ProcessEvent(e eventManager.Event) eventManager.Event{


	if v,b:=e.(*eventManager.ConsensusMsgEvent);b{
		if v.Msg.Type==pb.Message_TRANSACTION{
			cmh.logstatic.RecordCount("in ConsensusMsgHelper level","Message_TRANSACTION","ReceiveMsg")
		}

	}
	if cmh.batch.GetLocalID()==3{
		cmh.logstatic.RecordCount("in ConsensusMsgHelper level","ConsensusMsgHelper","ReceiveMsg")
		cmh.logstatic.RandomDelay(100,time.Microsecond,cmh.batch.GetLocalID(),3)
	}


	switch e.(type) {

	case *eventManager.ConsensusMsgEvent:
		cmh.logstatic.RecordCount("count","ConsensusMsgHelper","ConsensusMsgEvent")
		if v,b:=e.(*eventManager.ConsensusMsgEvent);b{
			c,_:=proto.Marshal(v.Msg)
			go cmh.batch.RecvMsg(c)
		}
	case *eventManager.StateUpdatedEvent:
		cmh.logstatic.RecordCount("count","ConsensusMsgHelper","StateUpdatedEvent")
		if v,b:=e.(*eventManager.StateUpdatedEvent);b{
			c,_:=proto.Marshal(v.Msg)
			cmh.batch.RecvMsg(c)
		}

	default:
		fmt.Println("there is no other operation!")
	}

	return nil
}

func NewMsgHelper(c consensus.Consenter,ls *statistics.LogStatic) *ConsensusMsgHelper{

	cmh:=&ConsensusMsgHelper{
		batch:c,
		logstatic:ls,
	}
	return  cmh
}