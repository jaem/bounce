package bouncer

import (
	"net/http"
	"fmt"
)

var onUnauthorizedRedirect = func (w http.ResponseWriter, r *http.Request) {
	fmt.Println("unauthorized")
	http.Error(w, "404 page not found", http.StatusUnauthorized)
}

var onLogoutRedirect = func (w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

// redirect to main page after logged in
var onLoggedinRedirect = func (w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

// redirect to login page
var onLoginRedirect = func (w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", http.StatusFound)
}