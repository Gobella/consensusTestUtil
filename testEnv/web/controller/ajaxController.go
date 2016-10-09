package controller
import (
	"net/http"
	"encoding/json"
	"fmt"
	"hyperchain/consensus/testEnv/web/asUtil"
	"strconv"
	"github.com/op/go-logging"
)
var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("hyperchain/consensus/testEnv/web")
}

type Result struct{
	Ret int
	Reason string
	Data interface{}
}

type AjaxController struct {}

func (this *AjaxController) CountAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 0, "参数错误", nil)
		return
	}
	//msgType:=r.FormValue("type")

	node:=r.FormValue("id")

	id,err:=strconv.Atoi(node)
	if err!=nil{
		logger.Error("fail to transfer string to int.",err)
		return
	}

	result,err1:=asUtil.ChatWithServer(asUtil.ReadServerConfig(id))
	if err1!=nil{
		logger.Error("failed to get the info from server.",err1)
		return
	}
	w.Write(result)
	fmt.Println("enter a Action,the value is",string(result))
	return
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}