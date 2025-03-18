package handler

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, data Envelope, status int) error {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Write(response)
	return nil
}
