package handler

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, data Envelope, headers http.Header) error {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// w.WriteHeader(http.StatusOK)
	w.Write(response)
	return nil
}
