package serverService

import (
	"hyperchain/consensus/testEnv/eventManager"
	"hyperchain/protos"
	"golang.org/x/net/context"
	"github.com/op/go-logging"
	"hyperchain/consensus/testEnv/connectInit"
	"hyperchain/consensus/testEnv/statistics"
)
var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("hyperchain/consensus/testEnv/peerManager/serverService")
}

type ConsensusService struct {
	Manager eventManager.Manager
	StatInfo *statistics.Jsoncat

}

func (s *ConsensusService ) Chat(ctx context.Context, c *protos.Message) (*protos.Message,error)  {
	s.Manager.Queue() <- &eventManager.ConsensusMsgEvent{
		Msg:c,
	}
	return &protos.Message{},nil
}

func (s *ConsensusService ) Statistic(ctx context.Context, in *connectInit.StatMessage) (*connectInit.StatMessage, error){
	var info []byte

	switch in.Type {

	case "info":
		info=s.StatInfo.ToJson()

	}

	return &connectInit.StatMessage{
		Type:"info",
		JSONInfo:info,
	},nil
}