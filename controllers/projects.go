package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func ProjectsGet(id, offset, limit string) (*models.Projects, error) {
	dbConfig := models.ReadDBConfigFromEnv()

	pm := models.NewProjectsManager(dbConfig)
	query, err := models.NewProjectsQuery(id, offset, limit)
	if err != nil {
		return nil, err
	}

	return pm.Execute(query)
}

func ProjectsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ignore := ""

	projects, err := ProjectsGet(id, ignore, ignore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if projects.Total == 0 {
		http.NotFound(w, r)
		return
	}

	project := projects.Projects[0]
	serializedProject, err := project.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProject)
}

func ProjectsQueryController(w http.ResponseWriter, r *http.Request) {
	ignore := ""

	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	projects, err := ProjectsGet(ignore, offset, limit)
	if err != nil {
		log.Panic(err)
	}


	serializedProjects, err := projects.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
}
