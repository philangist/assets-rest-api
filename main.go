package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/controllers"
)

func main() {
	fmt.Println("Hej, Världen!")

	router := mux.NewRouter()

	// routers for /projects resource
	router.Path("/projects/").
		HandlerFunc(controllers.ProjectsQueryController).
		Methods("GET")

	router.Path("/projects/{id:[0-9]+}").
		HandlerFunc(controllers.ProjectsGetController).
		Methods("GET")

	// routers for /assets resource
	router.Path("/assets/").
		HandlerFunc(controllers.AssetsQueryController).
		Methods("GET")

	router.Path("/assets/{id:[0-9]+}").
		HandlerFunc(controllers.AssetsGetController).
		Methods("GET")

	http.Handle("/", router)

	log.Println("Nu lyssna på :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
