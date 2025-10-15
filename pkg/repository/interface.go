package repository

import (
	"context"
)

type Repository[T any] interface {
	CreateVersion(ctx context.Context, record T) error
	GetById(ctx context.Context, record T, id string) error
	GetLatestByPrefix(ctx context.Context, record T, prefix string) error
	ListByPrefix(ctx context.Context, records *[]T, prefix string) error
}
