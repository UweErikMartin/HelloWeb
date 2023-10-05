package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "plain/text")
	hn, _ := os.Hostname()
	app.Logger.Printf("Host: %s(%s) - health endpoint called from %s\n", hn, runtime.GOOS, r.RemoteAddr)
	fmt.Fprintf(w, "Host: %s(%s) - health endpoint called from %s\n", hn, runtime.GOOS, r.RemoteAddr)
}

type envelope map[string]any

func (app *Application) MemStats(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	js, err := json.MarshalIndent(envelope{"MemStats": memStats}, "", "\t")

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}
