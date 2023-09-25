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
	applicationFlags = flag.NewFlagSet("server", flag.ExitOnError)
)

func main() {
	var appConfig applicationConfig

	// Read the commandline and environment variables into the application config
	applicationFlags.IntVar(&appConfig.healthEndpointPort, "health-port", 4000, "health endpoint port")
	applicationFlags.StringVar(&appConfig.healthEndpointPath, "health-path", "/health", "health endpoint path")
	applicationFlags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(applicationFlags, "server")

	defer klog.Flush()

	// Start the listeners asynchronously
	finish := make(chan bool)

	// Start health endpoint
	go func() {
		addr := fmt.Sprintf(":%d", appConfig.healthEndpointPort)
		healthEndPoint := http.NewServeMux()
		healthEndPoint.HandleFunc(appConfig.healthEndpointPath, healthProbeHandler)
		klog.Infof("Start serving health endpoint :%d%s\n", appConfig.healthEndpointPort, appConfig.healthEndpointPath)
		klog.Fatal(http.ListenAndServe(addr, healthEndPoint))
	}()

	<-finish
}

func healthProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Web!")
}
