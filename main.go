package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

var templates = template.Must(template.New("t").ParseGlob("templates/**/*.html"))

func init() {
	// Assign a user store
	store, err := NewFileUserStore("./data/users.yaml")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store

	// Assign a session store
	sessionStore, err := NewFileSessionStore("./data/sessions.yaml")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	// Assign a sql database
	db, err := NewMySQLDB("gophr:SOTWIZiniw@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	// Assign an image store
	globalImageStore = NewDBImageStore()
}

func main() {
	router := NewRouter()
	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/signout", HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)

	middleware := Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)

	log.Fatal(http.ListenAndServe(":3000", middleware))
}

// NewRouter creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return router
}
