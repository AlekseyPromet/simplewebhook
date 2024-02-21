package app

import (
	"AlekseyPromet/examples/simplewebhook/store"
	"net/http"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestService_GetServeMux(t *testing.T) {
	type fields struct {
		port   string
		logger *zap.Logger
		store  store.Store
	}

	logger := zap.NewNop()

	fieldsTest := &fields{
		port:   "8089",
		logger: logger,
		store:  store.NewTestStore(*logger),
	}

	tests := []struct {
		name   string
		fields fields
		want   *http.ServeMux
	}{
		0: {
			name:   "test 1",
			fields: *fieldsTest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				port:   tt.fields.port,
				logger: tt.fields.logger,
				store:  tt.fields.store,
			}
			if got := s.GetServeMux(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetServeMux() = %v, want %v", got, tt.want)
			}
		})
	}
}
