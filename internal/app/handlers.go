package application

import (
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog"
)

func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "plain/text")
	hn, _ := os.Hostname()
	klog.Infof("Host: %s - health endpoint called from %s\n", hn, r.RemoteAddr)
	fmt.Fprintf(w, "Host: %s - health endpoint called from %s\n", hn, r.RemoteAddr)
}
