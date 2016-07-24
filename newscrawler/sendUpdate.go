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
	_, err = http.Post("http://localhost:8443/wl08ncvrqisnv1wu8unwl08k05vo81j9", "application/json", bytes.NewReader(data))
	if err != nil {
		return
	}

	return
}
