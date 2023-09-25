package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/pkg/flagutil"
	"k8s.io/klog"
)

var (
	applicationFlags = flag.NewFlagSet("helloweb", flag.ExitOnError)
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
	fmt.Fprint(w, "Hello Web!")
}
