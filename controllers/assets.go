package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func AssetsGet(id, offset,  limit string) ([]byte, error) {
	dbConfig := models.ReadDBConfigFromEnv()

	pm := models.NewAssetsManager(dbConfig)
	query := models.NewAssetsQuery()
	assets, err := pm.Execute(query)
	if err != nil {
		log.Panic(err)
	}
	return assets.Serialize()
}

func AssetsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	serializedAssets, err := AssetsGet(id, "", "")

	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}

func AssetsQueryController(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	log.Printf("limit is %s, offset is %s", limit, offset)

	serializedAssets, err := AssetsGet("", offset, limit)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}
