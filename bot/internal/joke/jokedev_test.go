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

	ctx := context.Background()
	service := NewJokeDevService(ctx)

	joke, err := service.Joke(ctx, testParams)
	if err != nil {
		t.Fatal(err)
	}

	if joke.Answer != testSingleJokeDev.Content {
		t.Fatalf("Invalid joke response. Expected: '%s', Result: '%s'", testSingleJokeDev.Content, joke.Answer)
	}

	if joke.Category != Category(testSingleJokeDev.Category) {
		t.Fatalf("Invalid category. Expected: '%s', Result: '%s'", testParams, joke.Category)
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

			service := NewJokeDevService(ctx)

			if _, err := service.Joke(ctx, testParams); err == nil {
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

	ctx := context.Background()
	service := NewJokeDevService(ctx)

	if _, err := service.Joke(ctx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
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

	if _, err := service.Joke(ctx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if _, err := service.Joke(ctx, testParams); !errors.Is(err, DevServiceLimitExceededErr) {
		t.Fatal(err)
	}

	if httpmock.GetTotalCallCount() != 1 {
		t.Fatalf("service call DevService after got information about limitation exceeded. Call API %d times", httpmock.GetTotalCallCount())
	}
}
