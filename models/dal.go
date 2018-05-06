package models

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	User     string
	Password string
	Hostname string
	Name     string
	port     string
}

func ReadDBConfigFromEnv() *DBConfig {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	hostName := "db"
	port := "5432"

	return &DBConfig{user, password, hostName, name, port}
}

func (dbc *DBConfig) ConnectionString() string {
	replacer := strings.NewReplacer(
		"{user}", dbc.User,
		"{password}", dbc.Password,
		"{hostname}", dbc.Hostname,
		"{name}", dbc.Name,
		"{port}", dbc.port,
	)
	format := "postgresql://{user}:{password}@{hostname}:{port}/{name}?sslmode=disable"
	return replacer.Replace(format)
}

type EntityManager interface {
	Connection() (*sql.DB, error)
	Execute(query EntityQuery) (SerializableEntity, error)
}

type EntityQuery interface {
	Validate() error
	Evaluate() string // should return (preparedStatement string, values []string)
}

type SerializableEntity interface {
	Serialize() ([]byte, error)
}

func CoerceToInt64(value string) (int64, error) {
	if value != "" {
		return strconv.ParseInt(value, 10, 64)
	}
	return 0, nil
}

func isNullCheck() {}
