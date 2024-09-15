package db

import (
	"context"
	"database/sql"
	"time"
)

type Layer struct {
	Digest    string
	CreatedAt time.Time
}

func (l *Layer) Insert(ctx context.Context, db *sql.DB) error {
	query := `
		INSERT INTO layers (digest, created_at)
		VALUES (?, ?)
	`
	_, err := db.ExecContext(ctx, query, l.Digest, l.CreatedAt)
	return err
}

func (l *Layer) Exists(ctx context.Context, db *sql.DB) (bool, error) {
	query := `
		SELECT COUNT(*) FROM layers WHERE digest = ?
	`
	var count int
	err := db.QueryRowContext(ctx, query, l.Digest).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
