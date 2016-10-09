package handler

import (
	"net/http"
	"strings"
	"reflect"
	"hyperchain/consensus/testEnv/web/controller"
	"fmt"
)

func AjaxHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")

	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	fmt.Println("action:",action)
	ajax := &controller.AjaxController{}
	controller := reflect.ValueOf(ajax)

	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)

	method.Call([]reflect.Value{responseValue, requestValue})
}