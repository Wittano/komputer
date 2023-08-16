package internal

import (
	"net/http"
	"time"
)

var Client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives:     true,
		IdleConnTimeout:       time.Second * 5,
		ResponseHeaderTimeout: time.Second * 5,
		ExpectContinueTimeout: time.Second * 5,
		ForceAttemptHTTP2:     true,
	},
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 2,
}
