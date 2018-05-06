package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func AssetsGet(id, category, projectID, parentID, offset, limit string) ([]byte, error) {
	dbConfig := models.ReadDBConfigFromEnv()

	pm := models.NewAssetsManager(dbConfig)
	query, err := models.NewAssetsQuery(
		id, category, projectID, parentID, offset, limit,
	)
	if err != nil {
		log.Panic(err)
	}

	assets, err := pm.Execute(query)
	if err != nil {
		log.Panic(err)
	}

	return assets.Serialize()
}

func AssetsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ignore := ""

	serializedAssets, err := AssetsGet(
		id, ignore, ignore, ignore, ignore, ignore,
	)

	/*
	        if len(assets) == 0 {
	            return 404
	        }
		asset := assets[0]
	        return asset.Serialize()
	*/

	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}

func AssetsQueryController(w http.ResponseWriter, r *http.Request) {
	ignore := ""

	category := r.FormValue("type")
	projectID := r.FormValue("project_id")
	parentID := r.FormValue("parent_id")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	log.Printf("limit is %s, offset is %s", limit, offset)

	serializedAssets, err := AssetsGet(
		ignore, category, projectID, parentID, offset, limit,
	)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}
