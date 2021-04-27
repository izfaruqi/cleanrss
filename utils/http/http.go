package http

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	errJson, err := json.Marshal(err)
	if err != nil {
		errJson, _ = json.Marshal(err)
	}
	w.Write(errJson)
}

func WriteJson(w http.ResponseWriter, d interface{}) {
	out, err := json.Marshal(d)
	if err != nil {
		out, _ = json.Marshal(err)
	}
	w.Write(out)
}
