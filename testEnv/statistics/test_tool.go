package statistics

import (
	"github.com/spf13/viper"
	"strings"
	"fmt"
	"time"
	"math/rand"
	"os"
	"path/filepath"
)

//it will use the default value if the param is 0,it will block,if con1 and con2 are equal or one of them are nil
func (ls *LogStatic) RandomDelay(i int32,du time.Duration,con1 interface{},con2 interface{}){

	s1:=fmt.Sprint(con1)
	s2:=fmt.Sprint(con2)
	if !(con1==nil||con2==nil)&&!strings.EqualFold(s1, s2) {
		return
	}

	var err error
	if du==0{
		du,err=time.ParseDuration( ls.config.GetString("delay.defaultTime"))
	}
	if err!=nil {
		logger.Errorf("delay.defaultTime read error :%v",err)
		return
	}
	if i==0 {
		i=100
	}
	c:=rand.Int31n(i)
	//time.Sleep(time.Duration(c)*du)
	tick:=time.Tick(time.Duration(c)*du)
	<-tick
}

func (ls *LogStatic) Block(timed time.Duration){

	timevar:=time.After(timed)
	ti:=<-timevar
	ls.Blockbar<-ti

}

func (ls *LogStatic) BlockLock(){
	select {
	case <-ls.Blockbar:
		ls.RecordCount("","开始堵塞")
		ls.Blockbool=true
	default:
		ls.RecordCount("no","BlockLock")
	}

	if ls.Blockbool {
		ls.RecordCount("","堵塞")
		return
	}
}
// loadConfig load the config in the tool.yaml
func loadConfig() (config *viper.Viper) {

	config = viper.New()

	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)
	config.SetConfigName("tool")
	config.AddConfigPath("./")
	config.AddConfigPath("../consensus/testEnv/statistics")
	config.AddConfigPath("../../consensus/testEnv/statistics")
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		pbftpath := filepath.Join(p, "src/hyperchain/consensus/testEnv/statistics")
		config.AddConfigPath(pbftpath)
	}

	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Error reading plugin config: %s", err))
	}

	return config
}
func (ls *LogStatic) ExeInFixedTime(){

}

