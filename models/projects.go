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

func DefaultProjectsManager() *ProjectsManager {
	return NewProjectsManager(ReadDBConfigFromEnv())
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
	queryString, values := query.Build()
	log.Printf("query is %s values are %v", queryString, values)
	if len(values) == 0 {
		rows, err = db.Query(queryString)
	} else {
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

	offset := query.Offset
	limit := query.Limit
	total := len(projects)

	if offset > int64(total) {
		offset = int64(total)
	}

	if (limit > 0) && (offset+limit <= int64(total)) {
		projects = projects[offset : offset+limit]
	} else {
		projects = projects[offset:]
	}

	return NewProjects(projects, total, query.Limit, query.Offset), nil
}

// EntityQuery
type ProjectsQuery struct {
	ID     int64
	Offset int64
	Limit  int64
}

func NewProjectsQuery(id, offset, limit string) (*ProjectsQuery, error) {
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

func (pq *ProjectsQuery) Build() (string, []interface{}) {
	query := `
SELECT DISTINCT p.id, p.name, a.id, p.created_at
FROM projects p
JOIN assets a ON a.project_id=p.id
WHERE a.category=1 AND a.parent_id is NULL`

	counter := 1
	parameters := make([]interface{}, 0)

	if pq.ID > 0 {
		query += fmt.Sprintf(" AND p.id=$%d", counter)
		parameters = append(parameters, pq.ID)
	}

	query += ";"
	return query, parameters
}

// SerializableEntity
type Projects struct {
	Projects []Project  `json:"data"`
	Page     Pagination `json:"page"`
}

func NewProjects(projects []Project, total int, limit, offset int64) *Projects {
	page := Pagination{total, limit, offset}
	return &Projects{projects, page}
}

func (p *Projects) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

type Project struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	RootFolderID int       `json:"root_folder_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p *Project) Serialize() ([]byte, error) {
	return json.Marshal(p)
}
