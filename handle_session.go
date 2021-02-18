package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleSessionNew is the /login GET handler and displays the login form
func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Next": next,
	})
}

// HandleSessionCreate is the /login POST handler and checks for the correct password
// before opening a new session
func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// extract form values
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")

	// find user and check for validation errors and password credentials
	user, err := FindUser(username, password)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"Error": err,
				"User":  user,
				"Next":  next,
			})
			return
		}
		panic(err)
	}

	// find an existing session for the user or generate a new one
	session := FindOrCreateSession(w, r)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	if next == "" {
		next = "/"
	}

	// redirect the user to the intended page
	http.Redirect(w, r, next+"?flash=Signed+in", http.StatusFound)
}

// HandleSessionDestroy is the /signout POST handler and deletes the session from the
// global store
func HandleSessionDestroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(w, r, "sessions/destroy", nil)
}
