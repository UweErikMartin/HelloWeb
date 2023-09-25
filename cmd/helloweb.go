package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/pkg/flagutil"
	"k8s.io/klog"
)

var (
	applicationFlags = flag.NewFlagSet("helloweb", flag.ExitOnError)
	started          = time.Now()
)

func main() {
	var appConfig applicationConfig

	applicationFlags.IntVar(&appConfig.port, "health-port", 4000, "health endpoint port")
	applicationFlags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(applicationFlags, "helloweb")

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
	duration := time.Now().Sub(started)
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
