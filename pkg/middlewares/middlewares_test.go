package middlewares

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHideDotFiles(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler = HideDotFiles(http.StatusForbidden, handler)

	t.Run("some path element starts with a dot", func(t *testing.T) {
		assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "/test/.dir/secret", nil, http.StatusForbidden)
	})

	t.Run("final path element starts with a dot", func(t *testing.T) {
		assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "/test/.secret", nil, http.StatusForbidden)
	})

	t.Run("no element starting with a dot", func(t *testing.T) {
		assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "/test/dir/file", nil, http.StatusTeapot)
	})
}
