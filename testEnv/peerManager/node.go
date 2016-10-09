package peerManager

import (
	pb "hyperchain/protos"
	"hyperchain/consensus/testEnv/connectInit"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"strings"
	"strconv"
	"hyperchain/consensus/controller"
	"hyperchain/consensus/testEnv/eventManager"
	"hyperchain/consensus/consensusHelper"
	"time"
	"hyperchain/consensus/testEnv/statistics"
	"hyperchain/consensus/helper"
	"hyperchain/consensus/testEnv/externalBus"
)


type Address struct{
	IP string
	port int
	address string
}
type Node struct {
	server  *ServerImpl
	connectionPool map[string]*grpc.ClientConn //和所有节点的连接池
	replicaIDMap   map[int] string             //ip地址和id的一个映射
	addr []string                              //全网ip地址
	nodeSize int                               //维护一个全网节点数量,包括自己
	nodeNumberInNet int

	LocalAddr *Address
	LocalID int
	NetState bool
	Faulty int                                 //为测试预留，比如2代表网络堵塞，1代表发送错误消息

	Statics *statistics.LogStatic

	manager eventManager.Manager               //一个多向事件流
	EB externalBus.External
}

func (n *Node)Close(){
	n.server.StopServer()
}


func (n *Node) NodeNetInit() error{

	n.nodeNumberInNet=0
	n.connectionPool=make(map[string]*grpc.ClientConn)
	addr:=n.addr

	fc:=0
	for c:=len(n.replicaIDMap);c>0;c--{
		if strings.EqualFold(addr[c-1],n.LocalAddr.address){
			continue
		}
		err:=n.getConnectionWithAddr(addr[c-1])
		if err!=nil &&fc<3{
			fc++
			c++
			logger.Info("get connection failed! retry ",fc+1," times;err:",err,"\n this will try again!")
		}else{
			if fc>=3{
				logger.Infof("get connection %v failed! err: %v",addr[c-1],err)
				return err
			}
			fc=0
		}
	}
	n.NetState=true

	return nil
}

func (n *Node) getConnectionWithAddr(addr string) error{
	var err error
	n.connectionPool[addr],err=GetConnectionWithAddr(addr)
	n.nodeNumberInNet++
	if err!=nil{
		n.nodeNumberInNet--
		delete(n.connectionPool,addr)
		return err
	}

	return nil
}

func (n *Node) assignID(){

	n.replicaIDMap=make(map[int] string)


	if n.addr==nil||len(n.addr)==0{
		logger.Info("please input addr array")
		return
	}

	n.nodeSize=len(n.addr)

	for c:=len(n.addr);c>0;c--{
		tmp:=n.nodeSize-c
		addrTemp:=n.addr[tmp]

		n.replicaIDMap[tmp+1]=addrTemp
		if strings.EqualFold(n.addr[tmp],n.LocalAddr.address){
			n.LocalID=tmp+1
		}
	}
}

func  NewNode(localIP string,localPort int,addr []string) *Node{

	n:=&Node{
		NetState:false,
		addr:addr,
	}

	addrString:=localIP+":"+strconv.Itoa(localPort)

	n.LocalAddr=&Address{
		IP:localIP,
		port:localPort,
		address:addrString,
	}

	n.assignID()

	n.manager=eventManager.NewManagerImpl()

	n.RegistLogger()

	n.EB=externalBus.NewExternal(n.manager,n)

	h:=helper.NewHelper(n,n.EB,n.Statics)

	n.EB.Setlog(n.Statics)

	n.RegistConsenter(h)

	//go n.Statics.Start()

	n.manager.Start()


	p:=":"+strconv.Itoa(localPort)

	formate:=&statistics.Jsoncat{
		StaticInfo:n.Statics,
	}
	n.server=NewServer(p,n.manager,formate)
	return n
}

func (n *Node) RegistLogger(){
	info:=&statistics.NodeBasicInfo{
		Ip:n.LocalAddr.address,
		Id:strconv.Itoa(n.LocalID),
	}
	n.Statics=statistics.NewLogStatic(info,"general.log_prefix",4*time.Second)
	go n.Statics.Start()
}
func (n *Node) RegistConsenter( h helper.Stack){
	id:=uint64(n.LocalID)
	consenter:=controller.NewConsenter(id,h)
	ch:=consensusHelper.NewMsgHelper(consenter,n.Statics)
	n.manager.RegistReceiver(eventManager.ConsensusMsgEventConst,ch)
}

func (n *Node) retryNodeNetInit() {

	if n.nodeNumberInNet==n.nodeSize-1{
		logger.Infof("retry initial the replica ID is %d, port:%v failed ,all node connection has been initialed",n.LocalID,n.LocalAddr.port)
		return
	}
	n.nodeNumberInNet=0
	for _,addr :=range n.addr{
		if strings.EqualFold(addr,n.LocalAddr.address){
			continue
		}
		ex:=time.After(10*time.Second)
		for range time.Tick(3 * time.Second){
			if n.getConnectionWithAddr(addr)==nil{
				break
			}
			select {
			case <-ex:
				break
			}
		}
	}
	n.NetState=true
}

func (n *Node) BroadCast(msg *pb.Message){
	//logger.Info("broad cast n.nodeNumberInNet",n.nodeNumberInNet)
	if n.nodeNumberInNet<n.nodeSize-1 || !n.NetState{
		n.retryNodeNetInit()
	}
	for _,con:=range n.connectionPool{
		client:=GetClientWithConn(con)
		go n.UniCast(msg,client)
	}
}
func (n *Node) UniCast(msg *pb.Message,client connectInit.TestEnvClient){
	msgResponse,err:=client.Chat(context.Background(),msg)
	if err!=nil{
		//todo
		n.NetState=false
		logger.Info("msg:",msgResponse,"err:",err)
		n.retryNodeNetInit()
	}
}
func (n *Node) UniCastToSelf(msg *pb.Message){

	conn,err:=GetConnectionWithAddr(n.LocalAddr.address)
	if err!=nil{
		logger.Error("get connection failed !")
		return
	}
	client:=GetClientWithConn(conn)
	_,err=client.Chat(context.Background(),msg)

	if err!=nil{
		logger.Error("get client failed!!")
	}
}

func (n *Node) UniCastToID(id int,msg *pb.Message){
	target,ok:=n.replicaIDMap[id]
	if !ok{
		logger.Error("it doesn't exit for replica which id is",id)
		return
	}
	conn:=n.connectionPool[target]

	client:=GetClientWithConn(conn)
	_,err:=client.Chat(context.Background(),msg)

	if err!=nil{
		logger.Error("get client failed!!")
	}
}

func (n *Node) StartNode(){
	go n.server.Start()

	if n.NodeNetInit()==nil {
		logger.Info("//////////---------replica ID",n.LocalID,"网络初始化成功!---------//////////")
	}else{
		logger.Info("//////////---------replica ID",n.LocalID,"网络初始化失败!---------//////////")
		n.retryNodeNetInit()
	}
}


func (n *Node) GetLocalID() int{
	return n.LocalID
}
