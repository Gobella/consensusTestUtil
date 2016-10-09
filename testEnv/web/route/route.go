package route

import (
	"net/http"
	"hyperchain/consensus/testEnv/web/handler"
	"log"
	"html/template"
	"path/filepath"
	"os"
)


func render(w http.ResponseWriter, tmplName string, context map[string]interface{}) {
	tmpl, err := template.ParseFiles(tmplName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, context)
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var fpath string
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		fpath = filepath.Join(p, "src/hyperchain/consensus/testEnv/web/template")
	}
	f:=fpath+"/index.html"
	locals := make(map[string]interface{})
	render(w, f, locals)

	return
}

func ServerStart(){

	var fpath string
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		fpath = filepath.Join(p, "src/hyperchain/consensus/testEnv/web/template")
	}
	http.Handle("/index/css/", http.FileServer(http.Dir(fpath)))
	http.Handle("/index/js/", http.FileServer(http.Dir(fpath)))
	http.HandleFunc("/statistic/",handler.AjaxHandler)
	http.HandleFunc("/index/",indexHandler)
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
