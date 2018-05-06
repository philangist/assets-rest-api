package models

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"log"
	"strconv"
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
	ID        int64
	Type      int64
	ProjectID int64
	ParentID  int64
	Offset    int64
	Limit     int64
}

func NewAssetsQuery(id, category, projectID, parentID, offset, limit string) (*AssetsQuery, error) {
	var id64, category64, projectID64, parentID64, offset64, limit64 int64
	coerceToInt64 := func(value string) (int64, error) {
		if value != "" {
			return strconv.ParseInt(value, 10, 64)
		}
		return 0, nil
	}



	id64, err := coerceToInt64(id)
	if err != nil {
		return nil, err
	}

	category64, err = coerceToInt64(category)
	if err != nil {
		return nil, err
	}

	parentID64, err = coerceToInt64(parentID)
	if err != nil {
		return nil, err
	}

	projectID64, err = coerceToInt64(projectID)
	if err != nil {
		return nil, err
	}

	offset64, err = coerceToInt64(offset)
	if err != nil {
		return nil, err
	}

	limit64, err = coerceToInt64(limit)
	if err != nil {
		return nil, err
	}

	return &AssetsQuery{
		ID:         id64,
		Type:       category64,
		ProjectID:  projectID64,
		ParentID:   parentID64,
		Offset:     offset64,
		Limit:      limit64,
	}, nil
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
