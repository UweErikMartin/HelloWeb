package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	app "github.com/UweErikMartin/HelloWeb/internal/app"

	"k8s.io/klog"
)

func main() {

	// Initialize the application and parse the commandline
	// and environment variables
	app := &app.Application{}
	app.ParseCommandlineAndEnvironment(os.Args[1:])

	wg := &sync.WaitGroup{}

	// start the http server
	if app.AllowInsecureConnections() {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			// create the http Endpoint
			httpSrv := &http.Server{
				Addr:              app.GetInsecureAddrAsString(),
				Handler:           app.Routes(),
				ReadTimeout:       5 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       5 * time.Second,
			}

			klog.Infoln("Start listening on http port")
			klog.Fatal(httpSrv.ListenAndServe())
			wg.Done()
		}(wg)
	} else {
		klog.Infoln("HTTP is disabled")
	}

	if tlsConfig, err := app.GetTLSConfig(); err == nil {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			httpsSrv := &http.Server{
				Addr:              app.GetAddrAsSring(),
				Handler:           app.Routes(),
				TLSConfig:         tlsConfig,
				ReadTimeout:       5 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       5 * time.Second,
			}
			klog.Infoln("Start listening on https port")
			klog.Fatal(httpsSrv.ListenAndServeTLS("", ""))
			wg.Done()
		}(wg)
	} else {
		klog.Infoln("HTTPS is disabled")
	}

	// wait until the goroutines complete
	wg.Wait()
}
