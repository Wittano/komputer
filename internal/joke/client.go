package joke

import (
	"net/http"
	"time"
)

var Client = http.Client{
	Timeout: time.Second * 2,
}
