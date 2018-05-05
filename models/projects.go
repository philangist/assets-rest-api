package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

// EntityManager
type ProjectsManager struct {
	DBConfig *DBConfig
}

func NewProjectsManager(dbc *DBConfig) *ProjectsManager {
	return &ProjectsManager{dbc}
}

func (pm *ProjectsManager) Connection() (*sql.DB, error) {
	return sql.Open("postgres", pm.DBConfig.ConnectionString())
}

func (pm *ProjectsManager) Execute(query *ProjectsQuery) (*Projects, error) {
	db, err := pm.Connection()
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

	var project Project
	var projects []Project

	for rows.Next(){
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.RootFolderID,
			&project.CreatedAt,
		)

		if err != nil {
			log.Panic(err)
		}

		projects = append(projects, project)
	}

	return &Projects{projects, len(projects)}, nil
}

// EntityQuery
type ProjectsQuery struct {
	ID sql.NullInt64
}

func NewProductsQuery(id string) *ProjectsQuery {
	if id == ""{
		return &ProjectsQuery{}
	}
	i, err := strconv.Atoi(id)
	id64 := sql.NullInt64{Int64: int64(i), Valid: err == nil}
	return &ProjectsQuery{id64}
}

func (pq *ProjectsQuery) Validate() error {
	value := pq.ID.Int64
	if value < 0 {
		return fmt.Errorf("Error. Invalid ID value: %d", value)
	}
	return nil
}

func (pq *ProjectsQuery) Evaluate() string {
	// generates query that left joins projects and assets tables to return
	// composite `Project` representation for the API layer
	query :=
`SELECT p.id, p.name, a.id, p.created_at FROM projects p JOIN assets a ON a.project_id=p.id WHERE a.category=1`
	if (reflect.TypeOf(pq.ID) != nil) && (pq.ID.Int64 > 0) {
		query += fmt.Sprintf(" AND p.id=%d", pq.ID.Int64)
	}

	query += ";"
	return query
}

// SerializableEntity
type Projects struct {
	Projects []Project `json:"data"`
	Total         int `json:"total"`
}

type Project struct {
	ID                  int `json:"id"`
	Name             string `json:"name"`
	RootFolderID        int `json:"root_folder_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (p *Projects) Serialize() ([]byte, error) {
	return json.Marshal(p)
}
