package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philangist/frameio-assets/models"
)

func AssetsGet(id, category, projectID, parentID, descendants, offset, limit string) (*models.Assets, error) {
	pm := models.DefaultAssetsManager()

	query, err := models.NewAssetsQuery(
		id, category, projectID, parentID, descendants, offset, limit,
	)
	if err != nil {
		return nil, err
	}

	return pm.Execute(query)
}

func AssetsGetController(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	assets, err := AssetsGet(
		id,
		models.IGNORE,
		models.IGNORE,
		models.IGNORE,
		models.IGNORE,
		models.IGNORE,
		models.IGNORE,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(assets.Assets) == 0 {
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
	category := r.FormValue("type")
	projectID := r.FormValue("project_id")
	parentID := r.FormValue("parent_id")
	descendants := r.FormValue("descendants")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	assets, err := AssetsGet(
		models.IGNORE,
		category,
		projectID,
		parentID,
		descendants,
		offset,
		limit,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serializedAssets, err := assets.Serialize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedAssets)
}
