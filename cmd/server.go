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
	var intHttpPort int
	var intHttpsPort int
	var strPath string
	var strHost string
	var strTLSCert string
	var strTLSKey string

	// Read the commandline and environment variables into the application config
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.IntVar(&intHttpPort, "http-port", 80, "http-port the server is listening on")
	flags.IntVar(&intHttpsPort, "https-port", 443, "https-port the server is listening on")
	flags.StringVar(&strPath, "path", "", "path the server is serving")
	flags.StringVar(&strHost, "host", "", "host the server is serving")
	flags.StringVar(&strTLSCert, "tls-cert", "", "certificate file location for the tls communication")
	flags.StringVar(&strTLSKey, "tls-key", "", "public key for TLS communication")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "SERVER")

	strHttpAddr := fmt.Sprintf("%s:%d", strHost, intHttpPort)
	strHttpsAddr := fmt.Sprintf("%s:%d", strHost, intHttpsPort)
	strHandlerPath := fmt.Sprintf("%s/health", strPath)
	runtime.ReadMemStats(&memStats)

	http.HandleFunc(strHandlerPath, handleHealth)

	// create the http and https server
	httpServer := &http.Server{
		Addr: strHttpAddr,
	}

	httpsServer := &http.Server{
		Addr: strHttpsAddr,
	}

	finish := make(chan bool)

	go func() {
		klog.Infof("Listening on %s%s - MemoryAllocated: %d\n", strHttpAddr, strHandlerPath, memStats.TotalAlloc)
		klog.Fatal(httpServer.ListenAndServe())
	}()

	go func() {
		if strTLSCert != "" && strTLSKey != "" {
			klog.Infof("Listening on %s%s - MemoryAllocated: %d\n", strHttpsAddr, strHandlerPath, memStats.TotalAlloc)
			klog.Fatal(httpsServer.ListenAndServeTLS(strTLSCert, strTLSKey))
		}
	}()

	<-finish
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	runtime.ReadMemStats(&memStats)
	msg := fmt.Sprintf("MemoryAllocated: %d\n", memStats.TotalAlloc)
	hn, _ := os.Hostname()
	klog.Infof(msg)
	fmt.Fprintf(w, "%s feels well %s", hn, msg)
}
