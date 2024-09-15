package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"

	duckdb "github.com/marcboeker/go-duckdb"
)

func init() {
	connector, err := duckdb.NewConnector("/home/byatt/hemar/.hemar/.duckdb", func(execer driver.ExecerContext) error {
		bootQueries := []string{
			"INSTALL json",
			"LOAD 'json'",
		}

		for _, query := range bootQueries {
			_, err := execer.ExecContext(context.Background(), query, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create connector: %v", err)
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	ctx := context.Background()

	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS images (repository TEXT, tag TEXT, digest TEXT, manifest JSON, created_at TIMESTAMP, size INTEGER)")
	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS layers (digest TEXT, created_at TIMESTAMP)")
	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS containers (image_digest TEXT, digest TEXT, created_at TIMESTAMP)")
}

func DB() *sql.DB {
	db, err := sql.Open("duckdb", "/home/byatt/hemar/.hemar/.duckdb")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
}
