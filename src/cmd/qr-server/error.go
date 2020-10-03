package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendError(w http.ResponseWriter, code int, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, string(genErrorJSON(code, err)))
}

func genErrorJSON(code int, err error) []byte {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	resp := struct {
		Code int    `json:"code"`
		Err  string `json:"error"`
	}{
		Code: code,
		Err:  errorMsg,
	}

	b, _ := json.Marshal(resp)
	return b
}
