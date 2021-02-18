package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleImageNew handles the new image GET requests
func HandleImageNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// display new image form
	RenderTemplate(w, r, "images/new", nil)

}

// HandleImageCreate is the new image POST handler and reads an image from url or file
func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
		return
	}

	HandleImageCreateFromFile(w, r)
}

// HandleImageCreateFromURL downloads an image from a given url
func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)

	image := NewImage(user)
	image.Description = r.FormValue("description")

	err := image.CreateFromURL(r.FormValue("url"))

	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error":    err,
				"ImageURL": r.FormValue("url"),
				"Image":    image,
			})
			return
		}
		panic(err)
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

// HandleImageCreateFromFile uploads an image from a given file
func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {

	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")

	file, headers, err := r.FormFile("file")

	// No file was uploaded
	if file == nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": image,
		})
		return
	}

	// A file was uploaded, but an error occurred
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = image.CreateFromFile(file, headers)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": err,
			"Image": image,
		})
		return
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}
