package shared

import "net/http"

func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
