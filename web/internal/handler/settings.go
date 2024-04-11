package handler

import (
	"encoding/json"
	"github.com/wittano/komputer/web/internal/settings"
	"io"
	"net/http"
)

func UpdateSettings(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var newSetting settings.Settings
	if err = json.Unmarshal(rawBody, &newSetting); err != nil {
		return err
	}

	return settings.Config.Update(newSetting)
}

func GetSettings(w http.ResponseWriter, _ *http.Request) error {
	data, err := json.Marshal(settings.Config)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	return err
}
