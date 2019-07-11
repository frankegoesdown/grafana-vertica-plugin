package main

import (
	"context"
	"encoding/json"
	"net/http"
	"log"
	
	gVS "github.com/frankegoesdown/grafana-vertica-plugin/internal/app/grafana-vertica-service-server"
	
	_ "github.com/vertica/vertica-sql-go"
)

func main() {
	ConnectionString := ""

	a.PublicServer().Method(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte("foo"))
	}))

	a.PublicServer().Method(http.MethodPost, "/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		mapD := []string{"apple", "lettuce"}
		mapB, _ := json.Marshal(mapD)
		w.Write(mapB)
	}))

	a.PublicServer().Method(http.MethodGet, "/zip", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "applicaiton/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='ozon-grafana-simple-json-datasource.zip'")

		http.ServeFile(w, r, "./bin/ozon-grafana-simple-json-datasource.zip")
	}))
	a.PublicServer().Method(http.MethodPost, "/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		qr := gVS.QueryRequest{}

		err = json.NewDecoder(r.Body).Decode(&qr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isTableQuery := false
		for _, q := range qr.Targets {
			if q.Type == "table" {
				isTableQuery = true
			}
		}

		if isTableQuery {

			resp, err := gVS.TableResponse(config.ConnectionString, qr)
			mapB, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(mapB)
		} else {
			resp, err := gVS.SeriesResponse(config.ConnectionString, qr)
			mapB, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(mapB)
		}
	}))

	if err := a.Run(gVS.NewGrafanaVerticaServiceServer()); err != nil {
		log.Fatalf(context.Background(), "can't run app: %s", err)
	}
}
