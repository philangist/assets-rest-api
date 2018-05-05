package models

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	// "log"
	// "reflect"
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

// SerializableEntity
type Assets struct {
	Assets []Asset `json:"data"`
	Total         int `json:"total"`
}

type Asset struct {
	ID                  int `json:"id"`
	Name             string `json:"name"`
	RootFolderID        int `json:"root_folder_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (a *Assets) Serialize() ([]byte, error) {
	return json.Marshal(a)
}
