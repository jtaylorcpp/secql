package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func StartServer(osqConfig *OSQueryConfig, aggregator *Aggregator, address, port string) error {
	router := mux.NewRouter()

	routeMap := map[string]bool{}
	for query, _ := range osqConfig.Schedule {
		routeMap[query] = true
	}

	for _, pack := range osqConfig.Packs {
		for query, _ := range pack.Queries {
			routeMap[query] = true
		}
	}

	for query, _ := range routeMap {
		logrus.Infof("adding in route for osquery query %s", query)
		router.HandleFunc(fmt.Sprintf("/%s", query), newHandler(query, aggregator))
	}

	srv := &http.Server{
		Handler: router,
		Addr:    address + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func newHandler(metricsKey string, aggregator *Aggregator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("API request for metrics path: %v", *r.URL)
		if data, ok := aggregator.Tables[metricsKey]; ok {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "\t")
			encoder.Encode(data.Rows)
		} else {
			logrus.Infof("handler path %s doesnt have any records", metricsKey)
		}
	}
}
