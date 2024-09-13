package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

func init() {
	db, err := sql.Open("duckdb", "/home/byatt/hemar/.hemar/.duckdb")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS images (repository TEXT, tag TEXT, digest TEXT, created_at TIMESTAMP, size INTEGER)")
	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS layers (image_digest TEXT, digest TEXT, created_at TIMESTAMP, index INTEGER, size INTEGER)")
	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS containers (image_digest TEXT, digest TEXT, created_at TIMESTAMP)")
}

func DB() *sql.DB {
	db, err := sql.Open("duckdb", "/home/byatt/hemar/.hemar/.duckdb")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
}
