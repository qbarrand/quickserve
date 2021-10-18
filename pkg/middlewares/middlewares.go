package middlewares

import (
	"net/http"
	"strings"
)

// HideDotFiles goes over elements in the URL's path and returns code as a response if one of them starts
// with a dot ('.').
// If code is inferior to 0, it defaults to 403.
func HideDotFiles(code int, next http.Handler) http.HandlerFunc {
	if code <= 0 {
		code = http.StatusForbidden
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, elem := range strings.Split(r.URL.Path, "/") {
			if strings.HasPrefix(elem, ".") {
				w.WriteHeader(code)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
