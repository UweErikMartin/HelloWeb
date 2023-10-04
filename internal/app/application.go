package application

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	prof "net/http/pprof"

	"k8s.io/klog"
)

type Application struct {
	args cmdlineArgs
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// add the health endpoint ROOT/health
	mux.HandleFunc(fmt.Sprintf("%shealth", app.args.argRootPath), app.Health)
	// add the memstats endpoint
	mux.HandleFunc(fmt.Sprintf("%smemstats", app.args.argRootPath), app.MemStats)
	// add the profiling endpoint /debug/pprof
	if app.args.argEnableProfiling {
		klog.Infoln("Adding profiling Endpoint /debug/pprof/")
		mux.HandleFunc("/debug/pprof/", prof.Index)
	}

	return mux
}

func (app *Application) GetTLSConfig() (*tls.Config, error) {
	CertFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSCertFile)
	KeyFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSKeyFile)
	CAFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argMTLSCACertFile)

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		klog.Errorf("cannot load TLS Certificate files %s and %s\n", CertFilePath, KeyFilePath)
		return nil, err
	}

	caCert, err := os.ReadFile(CAFilePath)

	if err != nil {
		klog.Errorf("cannot load mTLS Certificate Authority file %s - using normal tls\n", CAFilePath)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
		}
		return tlsConfig, nil
	} else {
		klog.Infof("loaded mTLS Certificate Authority file %s - using mTLS\n", CAFilePath)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		return tlsConfig, nil
	}
}
