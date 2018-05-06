package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var DOMAIN_ROOT = "http://dev.frame.io"

// EntityManager
type AssetsManager struct {
	DBConfig *DBConfig
}

func NewAssetsManager(dbc *DBConfig) *AssetsManager {
	return &AssetsManager{dbc}
}

func DefaultAssetsManager() *AssetsManager {
	return NewAssetsManager(ReadDBConfigFromEnv())
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

	var rows *sql.Rows
	queryString, values := query.Build()

	if len(values) == 0 {
		rows, err = db.Query(queryString)
	} else {
		rows, err = db.Query(queryString, values...)
	}
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var asset Asset
	var assets []Asset

	for rows.Next() {
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

	serializableAssets := []*SerializableAsset{}
	for _, asset := range assets {
		serializableAssets = append(
			serializableAssets,
			NewSerializableAsset(&asset),
		)
	}

	offset := query.Offset
	limit := query.Limit
	total := len(serializableAssets)

	if offset > int64(total) {
		offset = int64(total)
	}

	if (limit > 0) && (offset+limit <= int64(total)) {
		serializableAssets = serializableAssets[offset : offset+limit]
	} else {
		serializableAssets = serializableAssets[offset:]
	}

	return NewAssets(
		serializableAssets,
		total,
		query.Limit,
		query.Offset,
	), nil
}

// EntityQuery
type AssetsQuery struct {
	ID          int64
	Category    int64
	ProjectID   int64
	ParentID    int64
	Descendants bool
	Limit       int64
	Offset      int64
}

func NewAssetsQuery(
	id, category, projectID, parentID, descendants, offset, limit string) (*AssetsQuery, error) {
	var id64, category64, projectID64, parentID64, offset64, limit64 int64
	var descendantsBool bool

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

	if descendants != IGNORE {
		descendantsBool, err = strconv.ParseBool(descendants)
		if err != nil {
			return nil, err
		}
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
		ID:          id64,
		Category:    category64,
		ProjectID:   projectID64,
		ParentID:    parentID64,
		Descendants: descendantsBool,
		Offset:      offset64,
		Limit:       limit64,
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
			"Error. Invalid Category value %d", category)
	}

	if projectID < 0 {
		return fmt.Errorf(
			"Error. Invalid ProjectID value %d", projectID)
	}

	if parentID < 0 {
		return fmt.Errorf(
			"Error. Invalid ParentID value %d", parentID)
	}

	return nil
}

func (aq *AssetsQuery) baseQuery() (string, []interface{}) {
	query := `SELECT * FROM assets`

	counter := 1
	parameters := make([]interface{}, 0)

	// Generates PostgreSQL prepared statements of the type
	// ("SELECT foo FROM bar WHERE foo.name = $1", "baz")
	addParameter := func(sqlFragment string, parameter interface{}) {
		query += fmt.Sprintf(sqlFragment, counter)
		parameters = append(parameters, parameter)
		counter += 1
	}

	if aq.ID > 0 {
		addParameter(" WHERE id = $%d", aq.ID)
		return query, parameters
	}

	if aq.Category > 0 {
		addParameter(" WHERE category=$%d", aq.Category)
		if aq.ProjectID > 0 {
			addParameter(" AND project_id=$%d", aq.ProjectID)
		}
		if aq.ParentID > 0 {
			addParameter(" AND parent_id=$%d", aq.ParentID)
		}
		return query, parameters
	}

	if aq.ProjectID > 0 {
		addParameter(" WHERE project_id=$%d", aq.ProjectID)
		if aq.ParentID > 0 {
			addParameter(" AND parent_id=$%d", aq.ParentID)
		}
		return query, parameters
	}

	if aq.ParentID > 0 {
		addParameter(" WHERE parent_id=$%d", aq.ParentID)
	}
	return query, parameters
}

func (aq *AssetsQuery) Build() (string, []interface{}) {
	query, parameters := aq.baseQuery()

	if aq.Descendants {
		descendantsQuery := `
WITH ancestor_nodes AS (
    {BASE_QUERY}
)
SELECT assets.id, assets.name, assets.parent_id, assets.media_url, assets.category,
       assets.project_id, assets.created_at
FROM assets
WHERE assets.id IN
    (SELECT ID FROM ancestor_nodes)
OR assets.parent_id IN
   (SELECT ID FROM ancestor_nodes)
`
		replacer := strings.NewReplacer("{BASE_QUERY}", query)
		query = replacer.Replace(descendantsQuery)
	}

	query += " ORDER BY assets.id ASC;"
	return query, parameters
}

// SerializableEntity
type Assets struct {
	Assets []*SerializableAsset `json:"data"`
	Page   Pagination           `json:"page"`
}

func NewAssets(assets []*SerializableAsset, total int, limit, offset int64) *Assets {
	page := Pagination{total, limit, offset}
	return &Assets{assets, page}
}

func (a *Assets) Serialize() ([]byte, error) {
	return json.Marshal(a)
}

type Asset struct {
	ID        int
	Name      string
	ParentID  sql.NullInt64
	MediaURL  sql.NullString
	Category  string
	ProjectID int
	CreatedAt time.Time
}

// JSON serializable alias of Asset
type SerializableAsset struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ParentID  int       `json:"parent_id"`
	MediaURL  string    `json:"media_url"`
	Category  string    `json:"type"`
	ProjectID int       `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (sa *SerializableAsset) Serialize() ([]byte, error) {
	return json.Marshal(sa)
}

func NewSerializableAsset(a *Asset) *SerializableAsset {
	var parentID int
	var mediaURL string

	if a.ParentID.Valid {
		parentID = (int)(a.ParentID.Int64)
	}

	if a.MediaURL.Valid {
		mediaURL = fmt.Sprintf(
			"%s/%s",  DOMAIN_ROOT, a.MediaURL.String)
	}

	return &SerializableAsset{
		ID:        a.ID,
		Name:      a.Name,
		ParentID:  parentID,
		MediaURL:  mediaURL,
		Category:  a.Category,
		ProjectID: a.ProjectID,
		CreatedAt: a.CreatedAt,
	}
}
