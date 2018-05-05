package api

import (
	_ "github.com/lib/pq"
	"os"
)

func ReadDBCredentials() (dbName, dbUser, dbPassword string){
	dbName = os.Getenv("POSTGRES_DB")
	dbUser = os.Getenv("POSTGRES_USER")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")

	return dbName, dbUser, dbPassword
}
