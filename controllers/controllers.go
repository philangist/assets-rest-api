package controllers

import (
	"log"
	"net/http"

	"github.com/philangist/frameio-assets/models"
)

func ProjectsControllerGET(w http.ResponseWriter, r *http.Request) {
	log.Println("Wrote json content-type header")

	dbConfig := models.ReadDBConfigFromEnv()
	log.Println("connection string is ", dbConfig.ConnectionString())

	pm := models.NewProjectsManager(dbConfig)
	query := models.NewProductsQuery("")
	projects, err := pm.Execute(query)
	if err != nil {
		log.Panic(err)
	}
	serializedProjects, err := projects.Serialize()
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedProjects)
	log.Println("Wrote json serialized response: ", projects)
}
