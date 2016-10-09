package systemTool

import (
	 "hyperchain/protos"
)

type SystemStack interface {
	BroadCast(msg *protos.Message)
	UniCastToID(id int,msg *protos.Message)
}