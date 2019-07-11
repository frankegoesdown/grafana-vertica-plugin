package main

import (
	"encoding/json"
	"net/http"

	gVS "github.com/frankegoesdown/grafana-vertica-plugin/internal/app/grafana-vertica-service-server"
	"github.com/go-chi/chi"
	_ "github.com/vertica/vertica-sql-go"
)

func main() {
	ConnectionString := ""
	r := chi.NewRouter()

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte("foo"))
	}))

	r.Get("/zip", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "applicaiton/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='grafana-vertica-datasource.zip'")

		http.ServeFile(w, r, "./bin/grafana-vertica-datasource.zip")
	}))

	r.Post("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		mapD := []string{"apple", "lettuce"}
		mapB, _ := json.Marshal(mapD)
		w.Write(mapB)
	}))

	r.Post("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		qr := gVS.QueryRequest{}

		err := json.NewDecoder(r.Body).Decode(&qr)
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

			resp, err := gVS.TableResponse(ConnectionString, qr)
			mapB, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(mapB)
		} else {
			resp, err := gVS.SeriesResponse(ConnectionString, qr)
			mapB, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(mapB)
		}
	}))

	http.ListenAndServe(":7000", r)
}
