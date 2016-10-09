package peerManager
import (
	"hyperchain/consensus/testEnv/connectInit"
	"google.golang.org/grpc"
)


type Client struct {
	Addr *Address
	Connection *grpc.ClientConn
	peer connectInit.TestEnvClient
	id int

}

//address formate "localhost:50051"
func GetConnectionWithAddr(address string) (*grpc.ClientConn,error){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err !=nil{
		logger.Warning("abtain an connection failed:",err)
		return nil,err
	}
	return conn,err
}

func GetClientWithConn(c *grpc.ClientConn) connectInit.TestEnvClient{
	return connectInit.NewTestEnvClient(c)
}

func (c *Client) clientInit() bool{
	conn,err:=GetConnectionWithAddr(c.Addr.address)
	if err==nil && conn!=nil{
		GetClientWithConn(conn)
		return true
	}
	return false
}
func NewClient(addr *Address,id int) *Client{

	client:=&Client{
		Addr:addr,
		id:id,
	}

	if client.clientInit(){
		return client
	}
	logger.Warning("new cliet failed!")
	return nil

}