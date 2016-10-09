package peerManager

import (
	"google.golang.org/grpc"
	"net"
	"hyperchain/consensus/testEnv/connectInit"
	"fmt"
	"log"
	"hyperchain/consensus/testEnv/eventManager"
	"hyperchain/consensus/testEnv/peerManager/serverService"
	"hyperchain/consensus/testEnv/statistics"
)


type ServerImpl struct{
	server *grpc.Server
	listener net.Listener
	localIP string
}



func  getServerWithService(service connectInit.TestEnvServer) *ServerImpl{
	s := grpc.NewServer()
	connectInit.RegisterTestEnvServer(s,service)
	serverInstance :=&ServerImpl{
		server:s,
	}
	return serverInstance
}
//ToDo new server 应该可以根据service创建不同机制的server
func NewServer(port string,manager eventManager.Manager,info *statistics.Jsoncat) *ServerImpl{
	sS:=&serverService.ConsensusService{
		Manager:manager,
		StatInfo:info,
	}

	serverInstance:=getServerWithService(sS)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serverInstance.listener=lis
	return serverInstance
}
func (s *ServerImpl) StopServer(){
	fmt.Println("stop")
	s.server.Stop()
}
func (s *ServerImpl) Start(){
	fmt.Println("server start")
	s.server.Serve(s.listener)
}
