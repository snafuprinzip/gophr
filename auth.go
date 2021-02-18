package main

import "net/http"

// AuthenticateRequest redirects the user to the login page if not
// properly authenticated
func AuthenticateRequest(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	if !authenticated {
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}
