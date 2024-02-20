package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
)

type Store interface {
	Create(ctx context.Context, src models.Source) (models.ApiKey, error)
	Update(ctx context.Context, key models.ApiKey, src models.Source) error
	Get(ctx context.Context, key models.ApiKey) (models.Source, error)
}
