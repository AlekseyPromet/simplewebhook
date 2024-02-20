package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TestStore struct {
	smap   *sync.Map
	logger zap.Logger
}

func NewTestStore(logger zap.Logger) *TestStore {
	return &TestStore{
		smap:   new(sync.Map),
		logger: logger,
	}
}

func (s *TestStore) Create(ctx context.Context, src models.Source) (models.SourceStore, error) {

	key := uuid.New().String()

	result := models.SourceStore{
		Key: key,
		Url: src.Url,
		Requests: models.Requests{
			Amount:     src.Requests.Amount,
			PerSeconds: src.Requests.PerSeconds,
		},
	}
	s.smap.Store(key, result)

	return result, nil
}

func (s *TestStore) Increment(ctx context.Context, key string, valueInc uint64) error {

	_, err := uuid.Parse(key)
	if err != nil {
		return err
	}

	value, ok := s.smap.Load(key)
	if !ok {
		return errors.New("key not found")
	}

	m, ok := value.(models.SourceStore)
	if !ok {
		return errors.New("key not found")
	}

	m.Requests.Iteration += valueInc

	s.smap.Store(key, m)

	return nil
}

func (s *TestStore) Get(ctx context.Context, key string) (m models.SourceStore, err error) {

	_, err = uuid.Parse(key)
	if err != nil {
		return m, err
	}

	value, ok := s.smap.Load(key)
	if !ok {
		return m, errors.New("key not found")
	}

	m, ok = value.(models.SourceStore)
	if !ok {
		return m, errors.New("model is not parced")
	}

	return m, nil
}

func (s *TestStore) GetAll(ctx context.Context) chan *models.SourceStore {

	out := make(chan *models.SourceStore, 1)

	go func() {
		defer close(out)

		s.smap.Range(func(key, value any) bool {

			m, ok := value.(models.SourceStore)
			if !ok {
				s.logger.Error("model is not parced")
				return false
			}

			out <- &m

			return true
		})
	}()

	return out
}

func (s *TestStore) Delete(ctx context.Context, key string) (m models.SourceStore, err error) {
	_, err = uuid.Parse(key)
	if err != nil {
		return m, err
	}

	value, ok := s.smap.LoadAndDelete(key)
	if !ok {
		return m, errors.New("key not found")
	}

	m, ok = value.(models.SourceStore)
	if !ok {
		return m, errors.New("model is not parced")
	}

	return m, nil
}
