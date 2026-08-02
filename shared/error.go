package shared

import (
	"fmt"
	"net/http"
)

func WriteErrorBadRequest(w http.ResponseWriter, message string, err error) {
	http.Error(w, fmt.Sprintf("%s: %s", message, err), http.StatusBadRequest)
}

func WriteErrorServerError(w http.ResponseWriter, message string, err error) {
	http.Error(w, fmt.Sprintf("%s: %s", message, err), http.StatusInternalServerError)
}
