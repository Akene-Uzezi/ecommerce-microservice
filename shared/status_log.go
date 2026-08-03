package shared

import (
	"log"
	"net/http"
)

func LogOK(method, url string) {
	log.Printf("%s %s %d", method, url, http.StatusOK)
}

func LogBadRequest(method, url string) {
	log.Printf("%s %s %d", method, url, http.StatusBadRequest)
}

func LogInternalServerError(method, url string) {
	log.Printf("%s %s %d", method, url, http.StatusInternalServerError)
}
