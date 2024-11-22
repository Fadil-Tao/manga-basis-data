package handlers

import (
	"fmt"
	"net/http"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Development")
	fmt.Fprint(w, "server is ready")
}
