package repository

import (
	"context"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

type Repository[T primitives.VerifiableAndRecordable] interface {
	CreateVersion(ctx context.Context, record T) error
	GetById(ctx context.Context, record T, id string) error
	GetLatestByPrefix(ctx context.Context, record T, prefix string) error
	ListByPrefix(ctx context.Context, records *[]T, prefix string) error

	Get(
		ctx context.Context,
		record T,
		condition data.ClauseOrExpression,
		order data.Ordering,
	) error

	Select(
		ctx context.Context,
		records *[]T,
		condition data.ClauseOrExpression,
		order data.Ordering,
		limit *uint,
	) error

	ListLatestByPrefix(
		ctx context.Context,
		records *[]T,
		preFilter data.ClauseOrExpression,
		condition data.ClauseOrExpression,
		order data.Ordering,
		limit *uint,
	) error
}
