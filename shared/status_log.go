package shared

import (
	"log"
	"net/http"
	"time"
)

func LogOK(method, url string, duration time.Duration) {
	log.Printf("%s %s %d %s", method, url, http.StatusOK, duration)
}

func LogBadRequest(method, url string, duration time.Duration) {
	log.Printf("%s %s %d %s", method, url, http.StatusBadRequest, duration)
}

func LogInternalServerError(method, url string, duration time.Duration) {
	log.Printf("%s %s %d %s", method, url, http.StatusInternalServerError, duration)
}

func LogNotFound(method, url string, duration time.Duration) {
	log.Printf("%s %s %d %s", method, url, http.StatusNotFound, duration)
}

func LogUnauthorized(method, url string, duration time.Duration) {
	log.Printf("%s %s %d %s", method, url, http.StatusUnauthorized, duration)
}
