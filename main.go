package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/philangist/frameio-assets/controllers"
)


func main (){
	fmt.Println("Hej, Världen!")

	http.HandleFunc("/projects", controllers.ProjectsControllerGET)
	log.Println("Nu lyssna på :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
