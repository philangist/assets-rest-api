package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func AssetsGet(id, category, projectID, parentID, offset, limit string) (*models.Assets, error) {
	dbConfig := models.ReadDBConfigFromEnv()

	pm := models.NewAssetsManager(dbConfig)
	query, err := models.NewAssetsQuery(
		id, category, projectID, parentID, offset, limit,
	)
	if err != nil {
		return nil, err
	}

	return pm.Execute(query)
}

func AssetsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ignore := ""

	assets, err := AssetsGet(
		id, ignore, ignore, ignore, ignore, ignore,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if assets.Total == 0 {
		http.NotFound(w, r)
		return
	}

	asset := assets.Assets[0]
	serializedAsset, err := asset.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAsset)
}

func AssetsQueryController(w http.ResponseWriter, r *http.Request) {
	ignore := ""

	category := r.FormValue("type")
	projectID := r.FormValue("project_id")
	parentID := r.FormValue("parent_id")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	assets, err := AssetsGet(
		ignore, category, projectID, parentID, offset, limit,
	)
	if err != nil {
		log.Panic(err)
	}

	serializedAssets, err := assets.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}
