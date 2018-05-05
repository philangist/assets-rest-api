package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	/*
        In this part of pm, we'll invoke pm.Connection(), defer connection.Close(), use query to build out a sql statement, run the sql statement with connection, and return it as a Projects struct.
        */
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

	return &Projects{projects}, nil
}

// EntityQuery
type ProjectsQuery struct {
	ID int  // TODO: replace ID int with ID string.
	// Typecast to assert it's a valid int when value is not nil
}

func NewProductsQuery(id int) *ProjectsQuery {
	return &ProjectsQuery{id}
}

func (pq *ProjectsQuery) Validate() error {
	if pq.ID <= 0 {
		return fmt.Errorf("Error. Invalid ID value: %d", pq.ID)
	}
	return nil
}

func (pq *ProjectsQuery) Evaluate() string {
	// more representative select statement:
	query :=
`SELECT p.id, p.name, a.id root_folder_id, p.created_at
FROM projects p JOIN assets a ON a.project_id=p.id
WHERE p.id=1 AND a.category=1;
`
	/*query := "SELECT * FROM projects"

	if pq.ID != 0 {
		conditional := fmt.Sprintf(" WHERE ID=%d", pq.ID)
		query += conditional
	}*/
	return query
}

// SerializableEntity
type Projects struct {
	Projects []Project `json:"data"`
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
