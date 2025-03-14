package handler

import (
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "api: ", log.LstdFlags)

func logError(r *http.Request, err error) {
	logger.Println(err)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := Envelope{"error": message}

	err := WriteJSON(w, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func RequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}
