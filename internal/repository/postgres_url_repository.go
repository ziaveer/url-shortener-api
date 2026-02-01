package repository

import (
	"context"
	"database/sql"
)

type PostgresURLRepository struct {
	db *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{db: db}
}

func (r *PostgresURLRepository) Save(ctx context.Context, code string, longURL string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO urls (code, long_url) VALUES ($1, $2)`,
		code,
		longURL,
	)
	return err
}

func (r *PostgresURLRepository) FindByCode(ctx context.Context, code string) (string, bool, error) {

	var longURL string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT long_url FROM urls WHERE code = $1`,
		code,
	).Scan(&longURL)

	if err == sql.ErrNoRows {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return longURL, true, nil
}
