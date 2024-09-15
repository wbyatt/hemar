package db

import (
	"context"
	"database/sql"
	"time"
)

type Image struct {
	Repository string
	Tag        string
	Digest     string
	CreatedAt  time.Time
	Size       int
	Manifest   []byte
}

func (i *Image) Insert(ctx context.Context, db *sql.DB) error {
	query := `
		INSERT INTO images (repository, tag, digest, created_at, size)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query, i.Repository, i.Tag, i.Digest, i.CreatedAt, i.Size)
	return err
}

func (i *Image) Exists(ctx context.Context, db *sql.DB) (bool, error) {
	query := `
		SELECT COUNT(*) FROM images WHERE repository = ? AND tag = ? AND digest = ?
	`
	var count int
	err := db.QueryRowContext(ctx, query, i.Repository, i.Tag, i.Digest).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (i *Image) ExistsByRepositoryAndTag(ctx context.Context, db *sql.DB) (bool, error) {
	query := `
		SELECT COUNT(*) FROM images WHERE repository = ? AND tag = ?
	`
	var count int
	err := db.QueryRowContext(ctx, query, i.Repository, i.Tag).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (i *Image) HydrateByRepositoryAndTag(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT repository, tag, digest, created_at, size FROM images WHERE repository = ? AND tag = ? LIMIT 1
	`

	row := db.QueryRowContext(ctx, query, i.Repository, i.Tag)
	err := row.Scan(&i.Repository, &i.Tag, &i.Digest, &i.CreatedAt, &i.Size)
	if err != nil {
		return err
	}

	return nil
}

func ListImages(ctx context.Context, db *sql.DB) ([]Image, error) {
	query := `
		SELECT repository, tag, digest, created_at, size FROM images
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.Repository, &image.Tag, &image.Digest, &image.CreatedAt, &image.Size)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func GetImage(ctx context.Context, db *sql.DB, repository string, tag string) (Image, error) {
	query := `
		SELECT repository, tag, digest, created_at, size FROM images WHERE repository = ? AND tag = ? LIMIT 1
	`

	row := db.QueryRowContext(ctx, query, repository, tag)
	var image Image
	err := row.Scan(&image.Repository, &image.Tag, &image.Digest, &image.CreatedAt, &image.Size)
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

func (i *Image) ExtractLayers(ctx context.Context) error {

}
