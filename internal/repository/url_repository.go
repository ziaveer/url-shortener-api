package repository

import "context"

type URLRepository interface {
	Save(ctx context.Context, code string, longURL string) error
	FindByCode(ctx context.Context, code string) (string, bool, error)
}
