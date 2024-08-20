package gravityforms

import (
	"net/http"
	"time"
)

const TimeFormat = time.DateTime

type Service struct {
	key    string
	secret string
	url    string

	debug bool
	http  *http.Client
}

func New(url string, key string, secret string, timeout time.Duration, debug bool) *Service {
	return &Service{
		key:    key,
		secret: secret,
		url:    url,
		debug:  debug,
		http:   &http.Client{Timeout: timeout},
	}
}
