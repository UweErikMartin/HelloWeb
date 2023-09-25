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

	applicationFlags.IntVar(&appConfig.healthPort, "health-port", 4000, "health endpoint port")
	applicationFlags.StringVar(&appConfig.healthPath, "health-path", "/health", "health endpoint path")
	applicationFlags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(applicationFlags, "helloweb")

	defer klog.Flush()

	finish := make(chan bool)

	go func() {
		addr := fmt.Sprintf(":%d", appConfig.healthPort)
		healthEndPoint := http.NewServeMux()
		healthEndPoint.HandleFunc(appConfig.healthPath, healthProbeHandler)
		klog.Infof("Start to listen on port %d\n", appConfig.healthPort)
		klog.Fatal(http.ListenAndServe(addr, healthEndPoint))
	}()

	<-finish
}

func healthProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Web!")
}
