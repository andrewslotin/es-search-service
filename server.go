package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type searcher interface {
	Search(ctx context.Context, query string) ([]json.RawMessage, error)
}

// SearchHandler returns an http.Handler that server search requests and responds
// with a list of results
func SearchHandler(s searcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		if q == "" {
			writeError(w, http.StatusBadRequest, "missing query parameter")
			return
		}

		results, err := s.Search(req.Context(), q)
		if err != nil {
			log.Printf("failed to perform search: %s", err)
			writeError(w, http.StatusInternalServerError, "")
			return
		}

		json.NewEncoder(w).Encode(struct {
			Status  string            `json:"status"`
			Results []json.RawMessage `json:"results"`
		}{
			Status:  "success",
			Results: append([]json.RawMessage{}, results...), // make sure "results" is always an array
		})
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	if message == "" {
		message = http.StatusText(code)
	}

	http.Error(
		w,
		fmt.Sprintf(`{"status": "error", "code": %d, "error": %s}`, code, strconv.Quote(message)),
		code,
	)
}
