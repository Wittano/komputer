package joke

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jarcoal/httpmock"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var testHumorAPIResponse = humorAPIResponse{
	JokeRes: "testJokeRes",
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
		humorApiURL+toHumorAPICategory(testJokeSearch.Category),
		httpmock.NewBytesResponder(http.StatusOK, response))

	os.Setenv(humorAPIKey, "123")

	ctx := context.Background()
	service := NewHumorAPIService(ctx)

	joke, err := service.Get(ctx, testJokeSearch)
	if err != nil {
		t.Fatal(err)
	}

	if joke.Answer != testHumorAPIResponse.JokeRes {
		t.Fatalf("Invalid joke response. Expected: '%s', Result: '%s'", testHumorAPIResponse.JokeRes, joke.Answer)
	}

	if joke.Category != testJokeSearch.Category {
		t.Fatalf("Invalid category. Expected: '%s', Result: '%s'", testJokeSearch, joke.Category)
	}
}

func TestHumorAPIService_GetButMissingApiKey(t *testing.T) {
	ctx := context.Background()
	service := NewHumorAPIService(ctx)
	if _, err := service.Get(ctx, testJokeSearch); err == nil {
		t.Fatal("service found API key, but it didn't set")
	}
}

func TestHumorAPIService_GetButApiReturnInvalidStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	badResponses := []int{http.StatusTooManyRequests, http.StatusPaymentRequired, http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError}

	httpmock.RegisterResponder("GET",
		humorApiURL+toHumorAPICategory(testJokeSearch.Category),
		httpmock.NewStringResponder(http.StatusOK, ""))

	os.Setenv(humorAPIKey, "123")

	for _, status := range badResponses {
		t.Run("API responses status "+strconv.Itoa(status), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			service := NewHumorAPIService(ctx)

			if _, err := service.Get(ctx, testJokeSearch); err == nil {
				t.Fatal("service didn't handle correct a bad/invalid http status")
			}
		})
	}
}

func TestHumorAPIService_GetButApiLimitWasExceeded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		humorApiURL+toHumorAPICategory(testJokeSearch.Category),
		httpmock.NewStringResponder(http.StatusPaymentRequired, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	os.Setenv(humorAPIKey, "123")

	ctx := context.Background()
	service := NewHumorAPIService(ctx)

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}
}

func TestHumorAPIService_Active(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		humorApiURL+toHumorAPICategory(testJokeSearch.Category),
		httpmock.NewStringResponder(http.StatusPaymentRequired, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	os.Setenv(humorAPIKey, "123")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := NewHumorAPIService(ctx)

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, HumorAPILimitExceededErr) {
		t.Fatal(err)
	}

	if httpmock.GetTotalCallCount() != 1 {
		t.Fatalf("service call HumorAPI after got information about limitation exceeded. Call API %d times", httpmock.GetTotalCallCount())
	}
}
