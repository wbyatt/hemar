package db

import (
	"context"
	"database/sql"
	"time"
)

type Layer struct {
	ImageDigest string
	Digest      string
	CreatedAt   time.Time
	Index       int
	Size        int
}

func (l *Layer) Insert(ctx context.Context, db *sql.DB) error {
	query := `
		INSERT INTO layers (image_digest, digest, created_at, index, size)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query, l.ImageDigest, l.Digest, l.CreatedAt, l.Index, l.Size)
	return err
}

func (l *Layer) Exists(ctx context.Context, db *sql.DB) (bool, error) {
	query := `
		SELECT COUNT(*) FROM layers WHERE image_digest = ? AND digest = ?
	`
	var count int
	err := db.QueryRowContext(ctx, query, l.ImageDigest, l.Digest).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
