package handler

import (
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "api: ", log.LstdFlags)

func logError(err error) {
	logger.Println(err)
}

func errorResponse(w http.ResponseWriter, status int, message interface{}) {
	env := Envelope{"message": message}

	err := WriteJSON(w, env, status)
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, http.StatusInternalServerError, message)
}

func RequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, http.StatusUnauthorized, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}
