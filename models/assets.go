package models

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"log"
	// "strconv"
	"time"
)

// EntityManager
type AssetsManager struct {
	DBConfig *DBConfig
}

func NewAssetsManager(dbc *DBConfig) *AssetsManager {
	return &AssetsManager{dbc}
}

func (am *AssetsManager) Connection() (*sql.DB, error) {
	return sql.Open("postgres", am.DBConfig.ConnectionString())
}

func (am *AssetsManager) Execute(query *AssetsQuery) (*Assets, error) {
	db, err := am.Connection()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = query.Validate()
	if err != nil {
		log.Panic(err)
	}
	rows, err := db.Query(query.Evaluate())
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var asset Asset
	var assets []Asset
	var serializableAssets []*SerializableAsset

	for rows.Next(){
		err = rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.ParentID,
			&asset.MediaURL,
			&asset.Category,
			&asset.ProjectID,
			&asset.CreatedAt,
		)

		if err != nil {
			log.Panic(err)
		}

		assets = append(assets, asset)
	}

	for _, asset := range assets {
		serializableAssets = append(
			serializableAssets,
			NewSerializableAsset(&asset),
		)
	}

	return &Assets{serializableAssets, len(serializableAssets)}, nil
}

// EntityQuery
type AssetsQuery struct {
	ID        sql.NullInt64
	Type      sql.NullString
	ProjectID sql.NullInt64
	ParentID  sql.NullInt64
	Offset    sql.NullInt64
	Limit     sql.NullInt64
}

func NewAssetsQuery() *AssetsQuery {
	return &AssetsQuery{}
}

func (aq *AssetsQuery) Validate() error {
	return nil
}

func (aq *AssetsQuery) Evaluate() string {
	query :=
`SELECT * FROM assets;`
	return query
}

// SerializableEntity
type Assets struct {
	Assets []*SerializableAsset `json:"data"`
	Total                   int `json:"total"`
}

type Asset struct {
	ID                    int
	Name               string
	ParentID    sql.NullInt64
	MediaURL   sql.NullString
	Category           string
	ProjectID             int
	CreatedAt       time.Time
}

type SerializableAsset struct {
	ID               int `json:"id"`
	Name          string `json:"name"`
	ParentID         int `json:"parent_id"`
	MediaURL      string `json:"media_url"`
	Category      string `json:"type"`
	ProjectID        int `json:"project_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewSerializableAsset(a *Asset) *SerializableAsset {
	var parentID int
	var mediaURL string

	if a.ParentID.Valid {
		parentID = (int)(a.ParentID.Int64)
	}

	if a.MediaURL.Valid {
		mediaURL = a.MediaURL.String
	}

	return &SerializableAsset{
		ID:         a.ID,
		Name:       a.Name,
		ParentID:   parentID,
		MediaURL:   mediaURL,
		Category:   a.Category,
		ProjectID:  a.ProjectID,
		CreatedAt:  a.CreatedAt,
	}
}

func (a *Assets) Serialize() ([]byte, error) {
	return json.Marshal(a)
}
