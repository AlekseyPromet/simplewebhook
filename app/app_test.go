package app

import (
	"AlekseyPromet/examples/simplewebhook/store"
	"context"
	"net/http"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type serviceInterface interface {
	GetServeMux() *http.ServeMux
	Run(fx.Lifecycle) *http.Server
	WebhookCycle(context.Context, chan error)
}

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

	var want serviceInterface

	tests := []struct {
		name   string
		fields fields
		want   serviceInterface
	}{
		0: {
			name:   "test server 1",
			fields: *fieldsTest,
			want:   want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Service{
				port:   tt.fields.port,
				logger: tt.fields.logger,
				store:  tt.fields.store,
			}
			if tt.want = got; tt.want.GetServeMux() == nil {
				t.Errorf("Service.GetServeMux() is nil, want %v", tt.want)
			}
		})
	}
}
