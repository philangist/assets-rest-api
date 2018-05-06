package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	ID     sql.NullInt64
	Offset sql.NullInt64
	Limit  sql.NullInt64
}

func NewProductsQuery(id, offset, limit string) (*ProjectsQuery, error) {
	var id64, offset64, limit64 sql.NullInt64

	coerceToNullInt64 := func(value string) sql.NullInt64 {
		i, err := strconv.Atoi(value)
		return sql.NullInt64{Int64: int64(i), Valid: err == nil}
	}

	if id == "" {
		offset64 = coerceToNullInt64(offset)
		limit64 = coerceToNullInt64(limit)
		return &ProjectsQuery{
			Offset: offset64,
			Limit: limit64,
		}, nil
	}

	id64 = coerceToNullInt64(id)
	return &ProjectsQuery{ID: id64}, nil
}

func (pq *ProjectsQuery) Validate() error {
	id := pq.ID.Int64
	offset := pq.Offset.Int64
	limit := pq.Limit.Int64

	if pq.ID.Valid && (id < 0) {
		return fmt.Errorf("Error. Invalid ID value: %d", id)
	}

	if pq.ID.Valid && (limit > 0 && offset > 0) {
		return fmt.Errorf(
			"Error. ID value cannot be present with Limit and Offset values")
	}

	return nil
}

func (pq *ProjectsQuery) Evaluate() string {
	// generates query that left joins projects and assets tables to return
	// composite `Project` representation for the API layer
	query :=
`SELECT p.id, p.name, a.id, p.created_at FROM projects p JOIN assets a ON a.project_id=p.id WHERE a.category=1`
	if pq.ID.Valid && (pq.ID.Int64 > 0) {
		query += fmt.Sprintf(" AND p.id=%d", pq.ID.Int64)
	}
	if pq.Limit.Valid && (pq.Limit.Int64 > 0){
		query += fmt.Sprintf(" LIMIT %d", pq.Limit.Int64)
	}
	if pq.Offset.Valid && (pq.Offset.Int64 > 0){
		query += fmt.Sprintf(" OFFSET %d", pq.Offset.Int64)
	}

	query += ";"

	log.Println("query is ", query)

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
