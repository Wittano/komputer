package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/internal/joke"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var testSingleJokeDev = singleResponse{
	Error:    false,
	Category: "Any",
	Type:     "single",
	Flags:    flags{},
	Id:       0,
	Safe:     false,
	Lang:     "",
	Content:  "testContent",
}

var testJokeDevUrl = fmt.Sprintf(jokeDevAPIUrlTemplate, testParams.Category, testParams.Type)

func TestDevService_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response, err := json.Marshal(testSingleJokeDev)
	if err != nil {
		t.Fatal(err)
	}

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewBytesResponder(http.StatusOK, response))

	if err = os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	ctx := log.NewContext(context.Background(), "")
	service := NewJokeDevService(ctx)

	res, err := service.RandomJoke(ctx, testParams)
	if err != nil {
		t.Fatal(err)
	}

	if res.Answer != testSingleJokeDev.Content {
		t.Fatalf("Invalid res response. Expected: '%s', Result: '%s'", testSingleJokeDev.Content, res.Answer)
	}

	if res.Category != joke.Category(testSingleJokeDev.Category) {
		t.Fatalf("Invalid category. Expected: '%s', Result: '%s'", testParams, res.Category)
	}
}

func TestDevService_GetAndApiReturnInvalidStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	badResponses := []int{http.StatusTooManyRequests, http.StatusPaymentRequired, http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError}

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewStringResponder(http.StatusOK, ""))

	if err := os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	for _, status := range badResponses {
		t.Run("API responses status "+strconv.Itoa(status), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logCtx := log.NewContext(ctx, "")
			service := NewJokeDevService(logCtx)

			if _, err := service.RandomJoke(logCtx, testParams); err == nil {
				t.Fatal("service didn't handle correct a bad/invalid http status")
			}
		})
	}
}

func TestDevService_GetWithApiLimitExceeded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewStringResponder(http.StatusTooManyRequests, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	if err := os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	ctx := log.NewContext(context.Background(), "")
	service := NewJokeDevService(ctx)

	if _, err := service.RandomJoke(ctx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}
}

func TestDevService_Active(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewStringResponder(http.StatusTooManyRequests, "").HeaderAdd(http.Header{rateLimitRemainingHeaderName: []string{"0"}}))

	os.Setenv(humorAPIKey, "123")

	ctx, cancel := context.WithCancel(context.Background())
	logCtx := log.NewContext(ctx, "")
	defer cancel()

	service := NewJokeDevService(ctx)

	if _, err := service.RandomJoke(logCtx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if _, err := service.RandomJoke(logCtx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if httpmock.GetTotalCallCount() != 1 {
		t.Fatalf("service call DevService after got information about limitation exceeded. Call API %d times", httpmock.GetTotalCallCount())
	}
}
