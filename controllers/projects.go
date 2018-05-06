package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func ProjectsGet(id, offset, limit string) (*models.Projects, error) {
	if limit == models.IGNORE {
		limit = "10" // default page size
	}

	pm := models.DefaultProjectsManager()
	query, err := models.NewProjectsQuery(id, offset, limit)
	if err != nil {
		return nil, err
	}

	return pm.Execute(query)
}

func ProjectsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	projects, err := ProjectsGet(id, models.IGNORE, models.IGNORE)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(projects.Projects) == 0 {
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
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	projects, err := ProjectsGet(models.IGNORE, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serializedProjects, err := projects.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
}
