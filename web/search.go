package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andrewslotin/es-search-service/storage"
)

type searcher interface {
	Search(ctx context.Context, query string, opts storage.SearchOptions) ([]json.RawMessage, error)
}

// SearchHandler returns an http.Handler that server search requests and responds
// with a list of results
func SearchHandler(s searcher) SecureHandler {
	return func(w http.ResponseWriter, req AuthenticatedRequest) {
		q := req.URL.Query().Get("q")
		if q == "" {
			writeError(w, http.StatusBadRequest, "missing query parameter")
			return
		}

		var from int
		if s := req.URL.Query().Get("from"); s != "" {
			v, err := strconv.Atoi(s)
			if err != nil || v < 0 {
				writeError(w, http.StatusBadRequest, "malformed from parameter")
				return
			}
			from = v
		}

		var size int
		if s := req.URL.Query().Get("size"); s != "" {
			v, err := strconv.Atoi(s)
			if err != nil || v < 0 {
				writeError(w, http.StatusBadRequest, "malformed size parameter")
				return
			}
			size = v
		}

		results, err := s.Search(req.Context(), q, storage.SearchOptions{
			From:   from,
			Size:   size,
			Sort:   req.URL.Query()["sort"], // allow multiple "sort" parameters
			Filter: req.URL.Query().Get("filter"),
		})
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
	}
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
