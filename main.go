package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/controllers"
)


func main (){
	fmt.Println("Hej, Världen!")

	router := mux.NewRouter()
	router.HandleFunc(
		"/projects/{id:[0-9]+}",
		controllers.ProjectsControllerGET,
	).Methods("GET")

	http.Handle("/", router)

	log.Println("Nu lyssna på :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
