package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func ListenMetrics(address string) error {
	//use separated ServeMux to prevent handling on the global Mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(address, mux)
}
