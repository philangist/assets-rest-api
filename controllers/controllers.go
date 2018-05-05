package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func ProjectsGet(id string) ([]byte, error) {
	dbConfig := models.ReadDBConfigFromEnv()
	log.Println("connection string is ", dbConfig.ConnectionString())

	pm := models.NewProjectsManager(dbConfig)
	query := models.NewProductsQuery(id)
	projects, err := pm.Execute(query)
	if err != nil {
		log.Panic(err)
	}
	return projects.Serialize()
}

func ProjectsControllerGET(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	log.Printf("requestParams are %v", requestParams)

	id := requestParams["id"]
	serializedProjects, err := ProjectsGet(id)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
}
