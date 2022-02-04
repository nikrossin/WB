package main

import (
	"html/template"
	"log"
	"net/http"
)

func HttpServerStart(cash Cashe) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["id"]
		if !ok || len(id[0]) < 1 {
			http.Error(w, "empty ID or not correct request", http.StatusBadRequest)
		} else {
			if val, ok := cash[id[0]]; ok {
				tmpl, err := template.ParseFiles("static/data.html")
				if err != nil {
					http.Error(w, "server error:"+err.Error(), http.StatusInternalServerError)
				}
				err = tmpl.Execute(w, val)
				if err != nil {
					http.Error(w, "server error:"+err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, "Error ID", http.StatusNotFound)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
