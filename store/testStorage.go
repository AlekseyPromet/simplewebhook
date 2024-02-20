package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/google/uuid"
)

type TestStore struct {
	smap *sync.Map
}

func NewTestStore() *TestStore {
	return &TestStore{smap: new(sync.Map)}
}

func (s *TestStore) Create(ctx context.Context, src models.Source) (models.ApiKey, error) {

	key := uuid.New().String()

	path, err := url.Parse(src.Url)
	if err != nil {
		return models.ApiKey{}, fmt.Errorf("parse url address %w", err)
	}

	s.smap.Store(key, models.SourceStore{
		Url: *path,
		Requests: models.Requests{
			Amount:     src.Requests.Amount,
			PerSeconds: src.Requests.PerSeconds,
		},
	})

	return models.ApiKey{
		Key: key,
	}, nil
}

func (s *TestStore) Increment(ctx context.Context, key models.ApiKey, valueInc uint64) error {

	keyUuid, err := uuid.Parse(key.Key)
	if err != nil {
		return err
	}

	value, ok := s.smap.Load(keyUuid.String())
	if !ok {
		return errors.New("key not found")
	}

	m, ok := value.(models.SourceStore)
	if !ok {
		return errors.New("key not found")
	}

	m.Requests.Iteration = valueInc

	s.smap.Store(keyUuid.String(), m)

	return nil
}

func (s *TestStore) Get(ctx context.Context, key models.ApiKey) (m models.SourceStore, err error) {

	keyUuid, err := uuid.Parse(key.Key)
	if err != nil {
		return m, err
	}

	value, ok := s.smap.Load(keyUuid.String())
	if !ok {
		return m, errors.New("key not found")
	}

	m, ok = value.(models.SourceStore)
	if !ok {
		return m, errors.New("model is not parced")
	}

	return m, nil
}
