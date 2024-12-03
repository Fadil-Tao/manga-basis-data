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
	switch {
	case strings.Contains(fmt.Sprint(err), "Unauthorized"):
		return http.StatusUnauthorized
	case strings.Contains(fmt.Sprint(err), "not found"):
		return http.StatusNotFound
	case strings.Contains(fmt.Sprint(err),"conflict"):
		return http.StatusConflict
	case strings.Contains(fmt.Sprint(err),"already"):
		return http.StatusConflict
	case strings.Contains(fmt.Sprint(err), "malformed"):
		return http.StatusBadRequest 
	case strings.Contains(fmt.Sprint(err) ,"invalid"):
		return http.StatusBadRequest
	case strings.Contains(fmt.Sprint(err), "empty"):
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}

type responseWrapper struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func JSONMarshaller(message string,object interface{})([]byte, error){
	if object == nil {
		object = []interface{}{}
	}
	
	wrapped := responseWrapper{
		Message: message,
		Data: object,
	}
	return json.Marshal(wrapped)
}
