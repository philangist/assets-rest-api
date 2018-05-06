package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	Category  int64
	ProjectID int64
	ParentID  int64
	Offset    int64
	Limit     int64
}

func NewAssetsQuery(id, category, projectID, parentID, offset, limit string) (*AssetsQuery, error) {
	var id64, category64, projectID64, parentID64, offset64, limit64 int64

	id64, err := CoerceToInt64(id)
	if err != nil {
		return nil, err
	}

	category64, err = CoerceToInt64(category)
	if err != nil {
		return nil, err
	}

	parentID64, err = CoerceToInt64(parentID)
	if err != nil {
		return nil, err
	}

	projectID64, err = CoerceToInt64(projectID)
	if err != nil {
		return nil, err
	}

	offset64, err = CoerceToInt64(offset)
	if err != nil {
		return nil, err
	}

	limit64, err = CoerceToInt64(limit)
	if err != nil {
		return nil, err
	}

	return &AssetsQuery{
		ID:         id64,
		Category:  category64,
		ProjectID:  projectID64,
		ParentID:   parentID64,
		Offset:     offset64,
		Limit:      limit64,
	}, nil
}

func (aq *AssetsQuery) Validate() error {
	id := aq.ID
	category := aq.Category
	projectID := aq.ProjectID
	parentID := aq.ParentID
	limit := aq.Limit
	offset := aq.Offset

	if id < 0 {
		return fmt.Errorf("Error. Invalid ID value: %d", id)
	}

	if (id != 0) && (limit > 0 && offset > 0) {
		return fmt.Errorf(
			"Error. ID value cannot be present with Limit and Offset values")
	}

	if (category < 0) || (category >= 3) {
		return fmt.Errorf(
			"Error. Invalid category value %d", category)
	} 

	if projectID < 0 {
		return fmt.Errorf(
			"Error. Invalid projectID value %d", projectID)
	}

	if parentID < 0 {
		return fmt.Errorf(
			"Error. Invalid parentID value %d", parentID)
	}

	return nil
}

func (aq *AssetsQuery) Evaluate() string {
	query :=
`SELECT * FROM assets`
	if aq.ID > 0 {
		query += fmt.Sprintf(" WHERE id = %d;", aq.ID)
		return query
	}

	if aq.Category > 0 {
		query += fmt.Sprintf(" WHERE category=%d", aq.Category)
		if aq.ProjectID > 0 {
			query += fmt.Sprintf(" AND project_id=%d", aq.ProjectID)
		}
		if aq.ParentID > 0 {
			query += fmt.Sprintf(" AND parent_id=%d", aq.ParentID)
		}
		goto Pagination
	}

	if aq.ProjectID > 0 {
		query += fmt.Sprintf(" WHERE project_id=%d", aq.ProjectID)
		if aq.ParentID > 0 {
			query += fmt.Sprintf(" AND parent_id=%d", aq.ParentID)
		}
		goto Pagination
	}

	if aq.ParentID > 0 {
		query += fmt.Sprintf(" WHERE parent_id=%d", aq.ParentID)
	}

Pagination:
	if aq.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", aq.Limit)
	}

	if aq.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", aq.Offset)
	}

	query += ";"

	log.Println("query is ", query)
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

// JSON serializable alias of Asset
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
