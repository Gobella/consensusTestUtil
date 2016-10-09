package statistics

import (
	//"github.com/Sirupsen/logrus"
	//"sync"
	//"io/ioutil"
	//"os"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"github.com/Sirupsen/logrus"
)


var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("hyperchain/statistics")
}

func (ls *LogStatic) getLoggerWithID(id string,level string) {
	logFile := ioutil.Discard
	levelstr:=ls.config.GetString(level)
	fileName:=levelstr+"_"+id+"_peer.log"
	logFile, _ = os.Create(fileName)

	loggerfile := &logrus.Logger{
		Out:       logFile,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}
	ls.loggerfile=loggerfile
}



