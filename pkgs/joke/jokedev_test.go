package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var testSingleJokeDev = jokeApiSingleResponse{
	Error:    false,
	Category: "Any",
	Type:     "single",
	Flags:    jokeApiFlags{},
	Id:       0,
	Safe:     false,
	Lang:     "",
	Content:  "testContent",
}

var testJokeDevUrl = fmt.Sprintf(jokeDevAPIUrlTemplate, testJokeSearch.Category, testJokeSearch.Type)

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

	os.Setenv(humorAPIKey, "123")

	ctx := context.Background()
	service := NewJokeDevService(ctx)

	joke, err := service.Get(ctx, testJokeSearch)
	if err != nil {
		t.Fatal(err)
	}

	if joke.Answer != testSingleJokeDev.Content {
		t.Fatalf("Invalid joke response. Expected: '%s', Result: '%s'", testSingleJokeDev.Content, joke.Answer)
	}

	if joke.Category != Category(testSingleJokeDev.Category) {
		t.Fatalf("Invalid category. Expected: '%s', Result: '%s'", testJokeSearch, joke.Category)
	}
}

func TestDevService_GetButApiReturnInvalidStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	badResponses := []int{http.StatusTooManyRequests, http.StatusPaymentRequired, http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError}

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewStringResponder(http.StatusOK, ""))

	os.Setenv(humorAPIKey, "123")

	for _, status := range badResponses {
		t.Run("API responses status "+strconv.Itoa(status), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			service := NewJokeDevService(ctx)

			if _, err := service.Get(ctx, testJokeSearch); err == nil {
				t.Fatal("service didn't handle correct a bad/invalid http status")
			}
		})
	}
}

func TestDevService_GetButApiLimitWasExceeded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		testJokeDevUrl,
		httpmock.NewStringResponder(http.StatusTooManyRequests, "").HeaderAdd(http.Header{xAPIQuotaLeftHeaderName: []string{"0"}}))

	os.Setenv(humorAPIKey, "123")

	ctx := context.Background()
	service := NewJokeDevService(ctx)

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, DevServiceLimitExceededErr) {
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
	defer cancel()

	service := NewJokeDevService(ctx)

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if _, err := service.Get(ctx, testJokeSearch); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if httpmock.GetTotalCallCount() != 1 {
		t.Fatalf("service call DevService after got information about limitation exceeded. Call API %d times", httpmock.GetTotalCallCount())
	}
}
