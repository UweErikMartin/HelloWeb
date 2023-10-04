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

	path := fmt.Sprintf("%shealth", app.args.argRootPath)
	mux.HandleFunc(path, app.Health)

	if app.args.argEnableProfiling {
		// add the profile endpoints
		path := "/debug/pprof/"
		klog.Infof("Adding profiling Endpoint %s\n", path)
		mux.HandleFunc(path, prof.Index)
	}

	return mux
}

func (app *Application) GetTLSConfig() *tls.Config {
	CertFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSCertFile)
	KeyFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSKeyFile)
	CAFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argMTLSCACertFile)

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		klog.Errorf("cannot load TLS Certificate files %s and %s\n", CertFilePath, KeyFilePath)
		return nil
	}

	caCert, err := os.ReadFile(CAFilePath)

	if err != nil {
		klog.Errorf("cannot load mTLS Certificate Authority file %s\n", CAFilePath)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
		}
		return tlsConfig
	} else {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		return tlsConfig
	}
}
