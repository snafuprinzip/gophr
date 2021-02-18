package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleHome handles the app's homepage
func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// display home page
	RenderTemplate(w, r, "index/home", nil)
}
