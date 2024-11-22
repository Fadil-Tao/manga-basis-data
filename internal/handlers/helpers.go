package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func statusCode(err error) int {
	var statusCode int
	if strings.Contains(fmt.Sprint(err), "Unauthorized"){
		statusCode = http.StatusUnauthorized
	}else{
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}

