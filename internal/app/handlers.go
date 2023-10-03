package application

import (
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog"
)

func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
	klog.Infoln("Serving health endpoint")
	w.WriteHeader(http.StatusOK)
	hn, _ := os.Hostname()
	fmt.Fprintf(w, "Host: %s - health endpoint called from %s\n", hn, r.RemoteAddr)
}
