package peerManager

import (
	"github.com/Sirupsen/logrus"
	"sync"
	"io/ioutil"
	"os"
	"github.com/op/go-logging"
)

var loggerfile *logrus.Logger
var initOnce sync.Once
var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("hyperchain/consensus/testEnv/peerManager")
}


func init(){
	initLogger()
}

func initLogger() {
	initOnce.Do(func() {
		logFile := ioutil.Discard

		logFile, _ = os.Create("peer.log")

		loggerfile = &logrus.Logger{
			Out:       logFile,
			Formatter: new(logrus.TextFormatter),
			Level:     logrus.DebugLevel,
		}
	})
}
