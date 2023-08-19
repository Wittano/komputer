package internal

import (
	"net/http"
	"time"
)

var Client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
		ForceAttemptHTTP2: true,
	},
	Timeout: time.Second * 2,
}
