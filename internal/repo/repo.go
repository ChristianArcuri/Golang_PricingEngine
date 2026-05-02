package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// Repo executes PricingEngine stored procedures against SQL Server.
type Repo struct {
	db *sql.DB
}

func Open(connectionString string) (*Repo, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("SQLSERVER_CONNECTION_STRING is empty")
	}
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

func (r *Repo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) DB() *sql.DB {
	return r.db
}
