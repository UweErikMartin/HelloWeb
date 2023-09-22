package main

import (
	"fmt"
	"net/http"

	"k8s.io/klog"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWeb)

	klog.Infoln("Start to listen on port 80")
	defer klog.Flush()
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		klog.Errorln(err)
	}
}

func helloWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Web!")
}
