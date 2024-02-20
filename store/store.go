package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
)

type Store interface {
	Create(ctx context.Context, src models.Source) (models.SourceStore, error)
	Increment(ctx context.Context, key string, valueInc uint64) error
	Get(ctx context.Context, key string) (models.SourceStore, error)
	GetAll(ctx context.Context) chan *models.SourceStore
	Delete(ctx context.Context, key string) (models.SourceStore, error)
}
