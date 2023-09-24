package main

import (
	"flag"
	"fmt"
	"net/http"

	"k8s.io/klog"
)

func main() {
	var appConfig applicationConfig

	flag.IntVar(&appConfig.port, "port", 4000, "health endpoint port")

	flag.Parse()

	finish := make(chan bool)

	healthEndPoint := http.NewServeMux()
	healthEndPoint.HandleFunc("/", helloWeb)

	defer klog.Flush()

	addr := fmt.Sprintf(":%d", appConfig.port)

	go func() {
		klog.Infof("Start to listen on port %d\n", appConfig.port)
		klog.Fatal(http.ListenAndServe(addr, healthEndPoint))
	}()

	<-finish
}

func helloWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Web!")
}
