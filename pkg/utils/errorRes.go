package utils

import (
	"encoding/json"
	"net/http"
	"tracker/pkg/model"
)

func ErrRes(w http.ResponseWriter, err error, errorStatus int) {
	encoder := json.NewEncoder(w)
	errMsg := err.Error()
	payload := model.GenericPayload{Msg: errMsg}
	w.WriteHeader(errorStatus)
	encoder.Encode(payload)
}
