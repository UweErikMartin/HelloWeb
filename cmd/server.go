package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/coreos/pkg/flagutil"
	"k8s.io/klog"
)

func main() {
	var port int
	var rootUrl string

	// Read the commandline and environment variables into the application config
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.IntVar(&port, "port", 80, "port the server is listening on")
	flags.StringVar(&rootUrl, "root-url", "/", "root url the server is serving")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "SERVER")

	defer klog.Flush()

	// Start health endpoint
	addr := fmt.Sprintf(":%d", port)
	healthEndPoint := http.NewServeMux()

	healthEndPoint.HandleFunc(rootUrl+"/health", healthProbeHandler)
	klog.Infof("Listening on :%d%s\n", port, rootUrl)

	klog.Fatal(http.ListenAndServe(addr, healthEndPoint))
}

func healthProbeHandler(w http.ResponseWriter, r *http.Request) {
	hn, err := os.Hostname()

	if err != nil {
		klog.Errorln(err)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	klog.Infof("serving health for %s - MemoryAllocated: %d\n", r.RemoteAddr, m.TotalAlloc)
	fmt.Fprint(w, hn+" feels well!")
}
