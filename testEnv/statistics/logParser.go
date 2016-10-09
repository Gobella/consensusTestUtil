package statistics

import (
	"sync"
	"github.com/Sirupsen/logrus"
	"time"
	"container/list"
	"fmt"
	"github.com/spf13/viper"
)

type LogStatic struct {
	Nodeinfo *NodeBasicInfo
	LogRecord map[string] *InfoFormate
	LogRecordList map[string] *list.List
	LogRecordInfo map[string] string
	mutex *sync.Mutex
	loggerfile *logrus.Logger
	config *viper.Viper
	model int
	timeDuration time.Duration
	Blockbar chan time.Time
	Blockbool bool
}


type InfoFormate struct {
	Des string
	Count uint64
}


//record the basic info for node
type NodeBasicInfo struct {
	Id string
	Ip string
}

func NewLogStatic(info *NodeBasicInfo,level string,time time.Duration) *LogStatic{
	ls:=&LogStatic{
		Nodeinfo:info,
		LogRecord:make(map[string] *InfoFormate),
		LogRecordList:make(map[string] *list.List),
		LogRecordInfo:make(map[string] string),
		mutex:&sync.Mutex{},
		config:loadConfig(),
		timeDuration:time,
		Blockbool:false,
	}
	ls.getLoggerWithID(ls.Nodeinfo.Id,level)
	return ls
}



//it will record the history of parameter ve in a list
func (ls *LogStatic) RecordList(ve interface{}, tags ...interface{}){
	key:="c"
	for _,nametag:=range tags{
		str:=fmt.Sprint(nametag)
		key=key+"_"+str
	}

	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	value:=fmt.Sprint(ve)

	if v,ok:=ls.LogRecordList[key];ok{
		v.PushBack(value)
		return
	}

	p:=list.New()
	p.PushBack(value)
	ls.LogRecordList[key]=p
}

//it will record the info once if the tags is the same
func (ls *LogStatic) RecordInfo(a interface{}, tags ...interface{}){
	key:="c"
	for _,nametag:=range tags{
		str:=fmt.Sprint(nametag)
		key=key+"_"+str
	}

	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	value:=fmt.Sprint(a)
	ls.LogRecordInfo[key]=value
}

//it will record the count with unique tags
func (ls *LogStatic) RecordCount(des1 string, tags ...interface{}){
	key:="c"
	for _,nametag:=range tags{
		name:=fmt.Sprint(nametag)
		key=key+"_"+name
	}

	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if v,ok:=ls.LogRecord[key];ok{
		v.Des=des1
		v.Count++
		return
	}

	c:=&InfoFormate{
		Des:des1,
		Count:1,
	}
	ls.LogRecord[key]=c
}


func (ls *LogStatic) PrintRecord(i int){
	ls.loggerfile.Info("//////////////////////////////////////////////////////////////////////////////////////////////////")
	ls.loggerfile.Info("                                         ",i)
	ls.loggerfile.Info("//////////////////////////////////////////////////////////////////////////////////////////////////")

	ls.loggerfile.Info("<<-----------------node basic infomation----------------->>")
	ls.loggerfile.Info("replica id:",ls.Nodeinfo.Id)
	for key,value:=range ls.LogRecordInfo{
		ls.loggerfile.Info(key,":",value)
	}
	ls.loggerfile.Info("<<----------------------record count--------------------->>")
	for key,value:=range ls.LogRecord{
		ls.loggerfile.Info(key,":",value)
	}
	ls.loggerfile.Info("<<----------------------record list---------------------->>")
	for key,value:=range ls.LogRecordList{
		st:=parseList(value)
		ls.loggerfile.Info(key,st)
	}
	//after printing ,the record should be cleared to prevent growing memory
	//ls.clear()
}


func (ls *LogStatic) clear(){
	ls.LogRecordInfo=make(map[string] string)
	ls.LogRecordList=make(map[string] *list.List)
}

//start a timer to write the log into file
func (ls *LogStatic) Start(){
	i:=0
	for tick := range time.Tick(ls.timeDuration){
		logger.Infof("<<<<<<<<<<<<per %v second to print the log!>>>>>>>>>>%v",ls.timeDuration,tick)
		i++
		go ls.PrintRecord(i)
	}
}

func parseList(ls *list.List) string{
	var st string
	for iter := ls.Front();iter != nil ;iter = iter.Next() {
		v,ok:=iter.Value.(string)
		if ok{
			st=st+"-"+v
		}
	}
	return st
}
