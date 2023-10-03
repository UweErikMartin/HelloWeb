package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	app "github.com/UweErikMartin/HelloWeb/internal/app"

	"k8s.io/klog"
)

func main() {

	// Initialize the application and parse the commandline
	// and environment variables
	app := &app.Application{}
	app.ParseCommandlineAndEnvironment(os.Args[1:])

	finish := make(chan struct{})

	// start the http server
	go func() {
		klog.Infoln("Start listening on http port")
		// create the http Endpoint
		httpSrv := &http.Server{
			Addr:              app.GetInsecureAddrAsString(),
			Handler:           app.Routes(),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		}

		klog.Fatal(httpSrv.ListenAndServe())
	}()

	go func() {
		if tlsConfig := app.GetTLSConfig(); tlsConfig != nil {
			httpsSrv := &http.Server{
				Addr:              app.GetAddrAsSring(),
				Handler:           app.Routes(),
				TLSConfig:         app.GetTLSConfig(),
				ReadTimeout:       5 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       5 * time.Second,
			}
			klog.Infoln("Start listening on https port")
			klog.Fatal(httpsSrv.ListenAndServeTLS("", ""))
		}
	}()

	go func() {
		klog.Infoln("Start profiling endpoint")
		klog.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// wait forever
	<-finish
}
