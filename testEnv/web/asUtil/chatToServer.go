package asUtil

import (
	"hyperchain/consensus/testEnv/peerManager"
	"github.com/op/go-logging"
	"context"
	"hyperchain/consensus/testEnv/connectInit"
)

var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("hyperchain/consensus/testEnv/web")
}

//id 从1开始写起
func ReadServerConfig(id int) string{
	addr:=[]string{
		"127.0.0.1:2333",
		"127.0.0.1:2334",
		"127.0.0.1:2335",
		"127.0.0.1:2336",
	}
	return addr[id-1]
}

func ChatWithServer(addr string) ([]byte,error){

	conn,err:=peerManager.GetConnectionWithAddr(addr)

	if err!=nil {
		logger.Error("obtain a connetion failed! connection:",addr)
		return nil,err
	}

	client:=peerManager.GetClientWithConn(conn)

	msg1:=&connectInit.StatMessage{
		Type:"info",
	}
	response,err2:=client.Statistic(context.Background(),msg1)
	if err2!=nil{
		logger.Error("obtain infomation from connection:",addr,"failed!")
		return nil,err2
	}
	return response.JSONInfo,nil
}