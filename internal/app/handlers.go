package application

import (
	"net/http"

	"k8s.io/klog"
)

func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
	klog.Infoln("Serving health endpoint")
	w.WriteHeader(http.StatusOK)
}
