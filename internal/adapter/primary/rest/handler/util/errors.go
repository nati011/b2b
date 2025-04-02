package handler

import (
<<<<<<< HEAD
	"errors"
	"net/http"
)

var (
	ErrInvalidRequestBody = errors.New("oopsy, invalid request data")
)

=======
	"net/http"
)

>>>>>>> 9deb76f5 (+ add success response writer util)
func errorResponse(w http.ResponseWriter, status int, message interface{}) {
	env := Envelope{"message": message}

	err := WriteJSON(w, env, status)
	if err != nil {
		w.WriteHeader(500)
	}
}

<<<<<<< HEAD
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
=======
func ServerErrorResponse(w http.ResponseWriter, err error) {
	logError(err)
>>>>>>> 2461a09b (+ fix handler errors)
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, http.StatusInternalServerError, message)
}

<<<<<<< HEAD
func RequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, http.StatusUnauthorized, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
=======
func RequestErrorResponse(w http.ResponseWriter, err error) {
	logError(err)
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter) {
>>>>>>> 2461a09b (+ fix handler errors)
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	errorResponse(w, http.StatusUnauthorized, nil)
}
