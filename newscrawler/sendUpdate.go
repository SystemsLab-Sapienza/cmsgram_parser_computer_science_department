package newscrawler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func sendUpdate(id string) (err error) {
	var (
		payload = struct {
			Key   string
			Value string
		}{"crawler", id}
	)

	// Encode the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Send the payload
	_, err = http.Post(conf.BotURI, "application/json", bytes.NewReader(data))
	if err != nil {
		return
	}

	return
}
