package joke

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var (
	testParams = internal.SearchParams{
		Type:     joke.Single,
		Category: joke.Any,
	}
)

var testHumorAPIResponse = humorAPIResponse{
	Content: "testJokeRes",
	ID:      213,
}

func TestHumorAPIService_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response, err := json.Marshal(testHumorAPIResponse)
	if err != nil {
		t.Fatal(err)
	}

	httpmock.RegisterResponder("GET",
		humorApiURL+humorAPICategory(testParams.Category),
		httpmock.NewBytesResponder(http.StatusOK, response))

	if err = os.Setenv(humorAPIKey, "123"); err != nil {
		t.Fatal(err)
	}

	ctx := log.NewContext(context.Background(), "")
	service := NewHumorAPIService(ctx)

	j, err := service.RandomJoke(ctx, testParams)
	if err != nil {
		t.Fatal(err)
	}

	if j.Answer != testHumorAPIResponse.Content {
		t.Fatalf("Invalid j response. Expected: '%s', Result: '%s'", testHumorAPIResponse.Content, j.Answer)
	}

	if j.Category != testParams.Category {
		t.Fatalf("Invalid category. Expected: '%s', Result: '%s'", testParams, j.Category)
	}
}

func TestHumorAPIService_GetWithMissingApiKey(t *testing.T) {
	ctx := log.NewContext(context.Background(), "")
	service := NewHumorAPIService(ctx)
	if _, err := service.RandomJoke(ctx, testParams); err == nil {
		t.Fatal("service found API key, but it didn't set")
	}
}

func TestHumorAPIService_GetWithApiReturnInvalidStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responses := []int{http.StatusTooManyRequests, http.StatusPaymentRequired, http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError}

	httpmock.RegisterResponder("GET",
		humorApiURL+humorAPICategory(testParams.Category),
		httpmock.NewStringResponder(http.StatusOK, ""))

	if err := os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	for _, status := range responses {
		t.Run("API responses status "+strconv.Itoa(status), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logCtx := log.NewContext(ctx, "")
			service := NewHumorAPIService(ctx)

			if _, err := service.RandomJoke(logCtx, testParams); err == nil {
				t.Fatal("service didn't handle correct a bad/invalid http status")
			}
		})
	}
}

func TestHumorAPIService_GetWithApiLimitExceeded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		humorApiURL+humorAPICategory(testParams.Category),
		httpmock.NewStringResponder(http.StatusPaymentRequired, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	if err := os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	ctx := log.NewContext(context.Background(), "")
	service := NewHumorAPIService(ctx)

	if _, err := service.RandomJoke(ctx, testParams); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}
}

func TestHumorAPIService_Active(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		humorApiURL+humorAPICategory(testParams.Category),
		httpmock.NewStringResponder(http.StatusPaymentRequired, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	if err := os.Setenv(humorAPIKey, "123"); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logCtx := log.NewContext(ctx, "")
	service := NewHumorAPIService(logCtx)

	if _, err := service.RandomJoke(logCtx, testParams); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}

	if _, err := service.RandomJoke(logCtx, testParams); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}

	if httpmock.GetTotalCallCount() != 1 {
		t.Fatalf("service call HumorAPI after got information about limitation exceeded. Call API %d times", httpmock.GetTotalCallCount())
	}
}
