package application

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"k8s.io/klog"
)

type Application struct {
	args cmdlineArgs
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	path := fmt.Sprintf("%shealth", app.args.argRootPath)
	mux.HandleFunc(path, app.Health)

	return mux
}

func (app *Application) GetTLSConfig() *tls.Config {
	CertFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSCertFile)
	KeyFilePath := fmt.Sprintf("%s/%s", app.args.argCertDir, app.args.argTLSKeyFile)

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		klog.Errorf("cannot load TLS Certificate files %s and %s\n", CertFilePath, KeyFilePath)
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	return tlsConfig
}