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

var memStats runtime.MemStats

func main() {
	var intPort int
	var strRootUrl string

	// Read the commandline and environment variables into the application config
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.IntVar(&intPort, "port", 80, "port the server is listening on")
	flags.StringVar(&strRootUrl, "root-url", "", "root url the server is serving")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "SERVER")

	strPort := fmt.Sprintf(":%d", intPort)
	strPath := fmt.Sprintf("%s/health", strRootUrl)
	runtime.ReadMemStats(&memStats)

	mux := http.NewServeMux()

	mux.HandleFunc(strPath, handleHealth)
	klog.Infof("Listening on %s%s - MemoryAllocated: %d\n", strPort, strPath, memStats.TotalAlloc)
	klog.Fatal(http.ListenAndServe(strPort, mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	runtime.ReadMemStats(&memStats)
	msg := fmt.Sprintf("MemoryAllocated: %d\n", memStats.TotalAlloc)
	hn, _ := os.Hostname()
	klog.Infof(msg)
	fmt.Fprintf(w, "%s feels well %s", hn, msg)
}
