package web

import "net/http"

// AuthenticatedRequest is an http.Request with accompanied by the name of the user
// on whose behalf it had been made
type AuthenticatedRequest struct {
	*http.Request
	Username string
}

// SecureHandler is an http.Handler that requires requests to be authenticated first
type SecureHandler func(w http.ResponseWriter, req AuthenticatedRequest)

// AuthMiddleware performs authentication before passing the request to
// the underlying handler. It responds with HTTP 401 if there was no
// Authorization header provided and stops request handling
func AuthMiddleware(next SecureHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, _, ok := req.BasicAuth()
		if !ok {
			writeError(w, http.StatusUnauthorized, "")
			return
		}

		next(w, AuthenticatedRequest{
			Request:  req,
			Username: user,
		})
	})
}
