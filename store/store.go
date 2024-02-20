package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
)

type Store interface {
	Create(ctx context.Context, src models.Source) (models.ApiKey, error)
	Increment(ctx context.Context, key models.ApiKey, valueInc uint64) error
	Get(ctx context.Context, key models.ApiKey) (models.SourceStore, error)
}
