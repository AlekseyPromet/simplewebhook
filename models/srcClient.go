package models

import (
	"net/url"
)

type Source struct {
	Url      string `json:"url"`
	Requests struct {
		Amount     uint64 `json:"amount"`
		PerSeconds uint64 `json:"per_second"`
	}
}

type ResponceWebhool struct {
	Iteration uint64 `json:"iteration"`
}

type Requests struct {
	Iteration  uint64 `json:"iteration"`
	Amount     uint64 `json:"amount"`
	PerSeconds uint64 `json:"per_second"`
}

type SourceStore struct {
	Url url.URL `json:"url"`
	Requests
}

type ApiKey struct {
	Key string
}
