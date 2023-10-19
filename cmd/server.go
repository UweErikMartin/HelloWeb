package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// create the http server if insecure connections are allowed
	httpSrv := &http.Server{
		Addr:              app.GetInsecureAddrAsString(),
		Handler:           app.Routes(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	httpsSrv := &http.Server{
		Addr:              app.GetAddrAsSring(),
		Handler:           app.Routes(),
		TLSConfig:         app.GetTLSConfig(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		if app.AllowInsecureConnections() {
			// create the http Endpoint
			app.Logger.Printf("Start listening on http://%s\n", app.GetInsecureAddrAsString())
			if err := httpSrv.ListenAndServe(); err != nil {
				app.Logger.Printf("KistenAndServe: %v", err)
			}
		} else {
			app.Logger.Println("http is disabled")
		}
	}()

	go func() {
		app.Logger.Printf("Start listening on https://%s\n", app.GetAddrAsSring())
		if err := httpsSrv.ListenAndServeTLS("", ""); err != nil {
			app.Logger.Printf("ListenAndServeTLS: %v", err)
		}
	}()

	defer func() {
		if err := httpSrv.Shutdown(ctx); err != nil {
			app.Logger.Println("error when shutting down the http server: ", err)
		}
		if err := httpsSrv.Shutdown(ctx); err != nil {
			app.Logger.Println("error when shutting down the https server: ", err)
		}
	}()

	sig := <-sigs
	app.Logger.Printf("sig %v", sig)
	cancel()
}
