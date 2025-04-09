package handler

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRequestBody = errors.New("oopsy, invalid request data")
)

func errorResponse(w http.ResponseWriter, status int, message interface{}) {
	env := Envelope{"message": message}

	err := WriteJSON(w, env, status)
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, err error) {
	logError(err)
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, http.StatusInternalServerError, message)
}

func RequestErrorResponse(w http.ResponseWriter, err error) {
	logError(err)
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter) {
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	errorResponse(w, http.StatusUnauthorized, nil)
}
