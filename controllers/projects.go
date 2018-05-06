package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func ProjectsGet(id, offset, limit string) ([]byte, error) {
	dbConfig := models.ReadDBConfigFromEnv()

	pm := models.NewProjectsManager(dbConfig)
	query, err := models.NewProductsQuery(id, offset, limit)
	if err != nil {
		log.Panic(err)
	}

	projects, err := pm.Execute(query)
	if err != nil {
		log.Panic(err)
	}
	return projects.Serialize()
}

func ProjectsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	serializedProjects, err := ProjectsGet(id, "", "")

	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
}

func ProjectsQueryController(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	log.Printf("limit is %s, offset is %s", limit, offset)

	serializedProjects, err := ProjectsGet("", offset, limit)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
}
