package agent

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(osqConfig *OSQueryConfig) error {
	metricMux := mux.NewRouter()
	aggregator := &Aggregator{Tables: map[string]interface{}{}}
	for query, _ := range osqConfig.Schedule {
		metricMux.HandleFunc(fmt.Sprintf("/%s", query), newHandler(query, aggregator))
	}
	return nil
}

func newHandler(metricsKey string, aggregator *Aggregator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if data, ok := aggregator.Tables[metricsKey]; ok {
			json.NewEncoder(w).Encode(data)
		}
	}
}
