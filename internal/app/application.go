package application

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	args   cmdlineArgs
	Logger *log.Logger
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// add the health endpoint ROOT/health
	mux.HandleFunc(fmt.Sprintf("%shealth", app.args.RootPath), app.Health)
	// add the memstats endpoint
	mux.HandleFunc(fmt.Sprintf("%smemstats", app.args.RootPath), app.MemStats)

	return mux
}

func (app *Application) GetTLSConfig() *tls.Config {
	CertFilePath := fmt.Sprintf("%s/%s", app.args.CertDir, app.args.TLSCertFile)
	KeyFilePath := fmt.Sprintf("%s/%s", app.args.CertDir, app.args.TLSKeyFile)
	CAFilePath := fmt.Sprintf("%s/%s", app.args.CertDir, app.args.MTLSCACertFile)

	app.Logger.Printf("Loading Certificate %s, %s", CertFilePath, KeyFilePath)

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)

	if err != nil {
		app.Logger.Printf("cannot load TLS Certificate files %s and %s\n", CertFilePath, KeyFilePath)
		return nil
	}

	caCert, err := os.ReadFile(CAFilePath)

	if err != nil {
		return &tls.Config{Certificates: []tls.Certificate{serverTLSCert}}
	} else {
		app.Logger.Printf("loaded mTLS Certificate Authority file %s - using mTLS\n", CAFilePath)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		return &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
	}
}
