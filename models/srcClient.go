package models

type Source struct {
	Url      string `json:"url,required"`
	Requests struct {
		Amount     uint64 `json:"amount,required"`
		PerSeconds int64  `json:"per_seconds,required"`
	}
}

type ResponseWebhook struct {
	Iteration uint64 `json:"iteration"`
}

type Requests struct {
	Iteration  uint64 `json:"iteration"`
	Amount     uint64 `json:"amount"`
	PerSeconds int64  `json:"per_seconds"`
}

type SourceStore struct {
	Key string
	Url string
	Requests
}

type ApiKey struct {
	Key string
}
