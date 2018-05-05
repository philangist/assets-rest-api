package main

import (
	"fmt"
	"log"

	"github.com/philangist/frameio-assets/models"
)


func main (){
	fmt.Println("Hej, Världen!")

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
	log.Printf("projects are %s", serializedProjects)
}
