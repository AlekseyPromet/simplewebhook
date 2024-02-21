package store

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestTestStore_Get(t *testing.T) {
	type fields struct {
		smap   *sync.Map
		logger *zap.Logger
	}
	type args struct {
		ctx context.Context
		key string
	}

	fieldsTest := fields{
		logger: zap.NewNop(),
		smap:   &sync.Map{},
	}

	key := "25c5f9e1-8661-48b5-8eaa-d01bba03a3df"
	resultStoreModel := models.SourceStore{
		Key: key,
		Url: "http://example.com",
		Requests: models.Requests{
			Iteration:  0,
			Amount:     100,
			PerSeconds: 1,
		},
	}
	fieldsTest.smap.Store(key, resultStoreModel)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantM   models.SourceStore
		wantErr bool
	}{
		0: {
			name:   "test store 1",
			fields: fieldsTest,
			args: args{
				ctx: context.Background(),
				key: key,
			},
			wantM: resultStoreModel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TestStore{
				smap:   tt.fields.smap,
				logger: *tt.fields.logger,
			}
			gotM, err := s.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("TestStore.Get() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}
