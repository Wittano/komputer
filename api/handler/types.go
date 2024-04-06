package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type apiError struct {
	Msg    string
	Status int
	Err    error
}

func (a apiError) Error() string {
	errMsg := a.Msg
	if a.Err != nil {
		errMsg = a.Err.Error()
	}

	return errMsg
}

func (a apiError) MarshalJSON() ([]byte, error) {
	msg := a.Msg
	if msg == "" {
		msg = "INTERNAL SERVER ERROR"
	}

	res := map[string]string{
		"message": msg,
	}

	return json.Marshal(res)
}

func newInternalApiError(err error) apiError {
	return apiError{Status: http.StatusInternalServerError, Err: err}
}

type ApiHandler func(w http.ResponseWriter, r *http.Request) error

func MakeHttpHandler(handler ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiErr apiError

		err := handler(w, r)

		if errors.As(err, &apiErr) {
			w.WriteHeader(apiErr.Status)

			resBody, err := apiErr.MarshalJSON()
			if err != nil {
				slog.ErrorContext(r.Context(), "failed send error message: %s", err)
			} else {
				_, err := w.Write(resBody)
				if err != nil {
					slog.ErrorContext(r.Context(), "failed send error message: %s", err)
				}
			}
		}

		if err != nil {
			slog.ErrorContext(r.Context(), err.Error())
		}
	}
}
