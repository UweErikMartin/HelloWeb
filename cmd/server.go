package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	app "github.com/UweErikMartin/HelloWeb/internal/app"
)

func main() {

	logger := log.New(os.Stdout, os.Args[0]+": ", log.Ldate|log.Ltime)

	// Initialize the application and parse the commandline
	// and environment variables
	app := &app.Application{
		Logger: logger,
	}
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

			app.Logger.Printf("Start listening on http://%s\n", app.GetInsecureAddrAsString())
			if err := httpSrv.ListenAndServe(); err != nil {
				app.Logger.Fatalln(err)
			}
			wg.Done()
		}(wg)
	} else {
		app.Logger.Println("HTTP is disabled")
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
			app.Logger.Printf("Start listening on https://%s\n", app.GetAddrAsSring())
			httpsSrv.ListenAndServeTLS("", "")
		}(wg)
		wg.Done()
	} else {
		app.Logger.Println("HTTPS is disabled")
	}

	// wait until the goroutines complete
	wg.Wait()
}
