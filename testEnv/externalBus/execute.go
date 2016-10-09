package externalBus

import (
	"hyperchain/consensus/testEnv/eventManager"
	pb "hyperchain/protos"

	"time"
	"hyperchain/core/types"
	"github.com/op/go-logging"
	"github.com/golang/protobuf/proto"
	"sync"
	"hyperchain/consensus/testEnv/statistics"
)
var logger *logging.Logger // package-level logger
//var currentChain *types.Chain
func init() {
	logger = logging.MustGetLogger("hyperchain/consensus/testEnv/externalBus")
	//currentChain=&types.Chain{}
}
type External interface{
	ExecuteTX(b []byte,seq uint64)
	UpdateState(updateState *pb.UpdateStateMessage)
	GetBlockchainInfo() *BlockchainInfo
	Setlog(logstatic *statistics.LogStatic)
}

type externalBus struct {
	manager eventManager.Manager
	nodeInfo ToExternalInfo
	currentChain *types.Chain
	logstatic	*statistics.LogStatic
	mutex        *sync.Mutex
}
type BlockchainInfo struct {
	Height            uint64 `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	CurrentBlockHash  []byte `protobuf:"bytes,2,opt,name=currentBlockHash,proto3" json:"currentBlockHash,omitempty"`
	PreviousBlockHash []byte `protobuf:"bytes,3,opt,name=previousBlockHash,proto3" json:"previousBlockHash,omitempty"`
}
func (eb *externalBus) GetBlockchainInfo() *BlockchainInfo{
	tmp:=&BlockchainInfo{
		Height:eb.currentChain.Height,
		CurrentBlockHash:eb.currentChain.LatestBlockHash,
		PreviousBlockHash:eb.currentChain.ParentBlockHash,
	}
	return tmp
}
func (eb *externalBus) ExecuteTX(b []byte,seq uint64){

	block:=&types.Block{
		Number:seq,
		BlockHash:b,
	}
	eb.WriteBlock(block,time.Now().UnixNano())
}
func (eb *externalBus) UpdateState(updateState *pb.UpdateStateMessage){
	time.Sleep(3*time.Second)
	logger.Info("*****************UpDate in the excute************************",updateState.SeqNo)
	updatedMsg:=&pb.StateUpdatedMessage{
		SeqNo:updateState.SeqNo,
	}

	tmpMsg, err := proto.Marshal(updatedMsg)
	temp:=eb.GetBlockchainInfo()
	start:=temp.Height
	for i:= start+1;i<updateState.SeqNo+1;i++{
		eb.logstatic.RecordList(i,"UpdateState","reWrite")
		t:=&types.Block{
			BlockHash:[]byte{0x00, 0x00, 0x05},
			Number:i,
		}
		eb.WriteBlock(t,time.Now().Unix())
	}
	if err != nil {
		eb.logstatic.RecordInfo("mashall faliled-cl")
		return
	}
	if updateState.TargetId!=nil{
		eb.currentChain.ParentBlockHash=eb.currentChain.LatestBlockHash
		eb.currentChain.LatestBlockHash=updateState.TargetId
		eb.currentChain.Height=updateState.SeqNo
	}
	msg:=&pb.Message{
		Type:pb.Message_STATE_UPDATED,
		Payload:tmpMsg,
		Timestamp:time.Now().Unix(),
		Id:uint64(eb.nodeInfo.GetLocalID()),
	}

	time.Sleep(100*time.Millisecond)
	eb.logstatic.RecordInfo("send StateUpdatedEvent")
	eb.manager.Queue()<-&eventManager.StateUpdatedEvent{
		Msg:msg,
	}
}

func (eb *externalBus) WriteBlock(block *types.Block,commitTime int64)  {

	block.ParentHash = eb.currentChain.LatestBlockHash
	block.WriteTime = time.Now().UnixNano()
	block.CommitTime = commitTime
	eb.currentChain.LatestBlockHash=block.BlockHash
	eb.currentChain.Height=block.Number
	eb.currentChain.ParentBlockHash=block.ParentHash
	eb.logstatic.RecordList(block.Number,"WriteBlock")

}
func  NewExternal(m eventManager.Manager,p ToExternalInfo) External{
	return &externalBus{
		manager:m,
		nodeInfo:p,
		currentChain:&types.Chain{},
		mutex:&sync.Mutex{},
	}
}

func (eb *externalBus) Setlog(logstatic	*statistics.LogStatic){
	eb.logstatic=logstatic
}

type ToExternalInfo interface {
	GetLocalID() int
}