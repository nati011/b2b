package handler

import "net/http"

func successResponse(w http.ResponseWriter, status int, message interface{}) {
	env := Envelope{"body": message}

	err := WriteJSON(w, env, status)
	if err != nil {
		logError(err)
		w.WriteHeader(500)
	}
}

func OperationSuccessResponse(w http.ResponseWriter, message interface{}) {
	successResponse(w, http.StatusAccepted, message)
}
