package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrPathVariableNotFound = errors.New("oopsy, id not provided")
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

func GetPathParam(r *http.Request, param_position int) (int, error) {

	pathSegments := strings.Split(r.URL.Path, "/")
	var typedParamId int
	if len(pathSegments) >= 4 {
		paramIdValue := pathSegments[4]
		var err error
		typedParamId, err = strconv.Atoi(paramIdValue)

		if err != nil {
			return 0, ErrPathVariableNotFound
		}
	} else {
		return 0, ErrPathVariableNotFound
	}
	return typedParamId, nil
}
