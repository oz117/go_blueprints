package auth

import "net/http"

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, req)
	}
}

// MusthAuth simply creates an authHandler structure
func MusthAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
