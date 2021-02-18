package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleUserNew handles the new user requests
func HandleUserNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	RenderTemplate(w, r, "users/new", nil)
}

// HandleUserCreate handles the new user requests
func HandleUserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := NewUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}

	// Create a new session
	session := NewSession(w)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

// HandleUserEdit is the /account GET handler that show the user's account page
// for him to edit his personal data
func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

// HandleUserUpdate is the /account POST handler and takes the form's data
// to update the account data
func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)

	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}
