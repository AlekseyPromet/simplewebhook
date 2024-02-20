package app

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

const IterationDuration = time.Second

func (s *Service) WebhookCicle(ctx context.Context, errChan chan error) {
	s.logger.Sugar().Infof("star webhook cicle at %v", time.Now().Format(time.RFC3339))

	for {
		eg, ctx := errgroup.WithContext(ctx)

		for source := range s.store.GetAll(ctx) {
			s.postWebhook(eg, source, errChan)
		}

		if err := eg.Wait(); err != nil {
			s.logger.Sugar().Errorln(err)
			errChan <- err
		}
	}
}

func (s *Service) postWebhook(eg *errgroup.Group, source *models.SourceStore, errChan chan error) {
	eg.Go(func() error {

		client := resty.New()
		freqency := time.Duration(1000/source.PerSeconds) * time.Millisecond
		endTime := time.Now().Add(IterationDuration)
		body := &models.ResponseWebhook{Iteration: source.Amount}

		for i := time.Now(); i.Before(endTime); i = time.Now() {

			s.logger.Sugar().Debugf("POST to %v body %v\n", source.Url, body)

			// POST JSON
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(body).
				Post(source.Url)
			if err != nil {
				errChan <- err
				break
			}
			if resp != nil && resp.IsError() {
				errChan <- err
				break
			}
			if resp != nil {
				s.logger.Sugar().Debugf("post %v status %v", source.Url, resp.Status())
			}
			time.Sleep(freqency)
		}
		return nil
	})
}
