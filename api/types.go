package api

import (
	"fmt"
	"time"
)

type AudioFileInfo struct {
	Filename string `json:"filename"`
	ID       string `json:"id"`
}

func (a AudioFileInfo) String() string {
	return fmt.Sprintf("%s - %s", a.Filename, a.ID)
}

type GetAudioIdsResponse struct {
	Files []AudioFileInfo `json:"files"`
}

type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
