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
	db, err := pm.Connection()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = query.Validate()
	if err != nil {
		log.Panic(err)
	}

	var rows *sql.Rows
	queryString, values := query.Evaluate()
	log.Printf("query is %s values are %v", queryString, values)
	if len(values) == 0 {
		rows, err = db.Query(queryString)
	}else{
		rows, err = db.Query(queryString, values...)
	}
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var project Project
	var projects []Project

	for rows.Next() {
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
	ID     int64
	Offset int64
	Limit  int64
}

func NewProductsQuery(id, offset, limit string) (*ProjectsQuery, error) {
	var id64, offset64, limit64 int64
	var err error

	if id == "" {
		offset64, err = CoerceToInt64(offset)
		if err != nil {
			return nil, err
		}

		limit64, err = CoerceToInt64(limit)
		if err != nil {
			return nil, err
		}

		return &ProjectsQuery{
			Offset: offset64,
			Limit:  limit64,
		}, nil
	}

	id64, err = CoerceToInt64(id)
	if err != nil {
		return nil, err
	}

	return &ProjectsQuery{ID: id64}, nil
}

func (pq *ProjectsQuery) Validate() error {
	id := pq.ID
	offset := pq.Offset
	limit := pq.Limit

	if id < 0 {
		return fmt.Errorf("Error. Invalid ID value: %d", id)
	}

	if (id != 0) && (limit > 0 && offset > 0) {
		return fmt.Errorf(
			"Error. ID value cannot be present with Limit and Offset values")
	}

	return nil
}

func (pq *ProjectsQuery) Evaluate() (string, []interface{}) {
	query := `SELECT p.id, p.name, a.id, p.created_at
FROM projects p JOIN assets a ON a.project_id=p.id
WHERE a.category=1 AND a.parent_id is NULL`

	counter := 1
	parameters := make([]interface{}, 0)

	// Generates PostgreSQL prepared statements of the type
	// ("SELECT foo FROM bar WHERE foo.name = $1", "baz")
	addParameter := func(sqlFragment string, parameter interface{}){
		query += fmt.Sprintf(sqlFragment, counter)
		parameters = append(parameters, parameter)
		counter += 1
	}

	if pq.ID > 0 {
		addParameter(" AND p.id=$%d;", pq.ID)
		return query, parameters
	}
	if pq.Limit > 0 {
		addParameter(" LIMIT $%d", pq.Limit)
	}
	if pq.Offset > 0 {
		addParameter(" OFFSET $%d", pq.Offset)
	}

	query += ";"
	return query, parameters
}

// SerializableEntity
type Projects struct {
	Projects []Project `json:"data"`
	Total    int       `json:"total"`
}

type Project struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	RootFolderID int       `json:"root_folder_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p *Projects) Serialize() ([]byte, error) {
	return json.Marshal(p)
}
