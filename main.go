package main

import (
	"fmt"

	"github.com/philangist/frameio-assets/api"
)

func main (){
	dbName, dbUser, dbPassword := api.ReadDBCredentials()

	fmt.Println("Hej, Världen!")
	fmt.Println(
		"Postgres access credentials are: ", dbName, dbUser, dbPassword)
}
