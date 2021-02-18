package main

import "errors"

// ValidationError is a subtype for user input error messages
type ValidationError error

var (
	errNoUsername           = ValidationError(errors.New("You must supply a user name"))
	errNoEmail              = ValidationError(errors.New("You must supply an email address"))
	errNoPassword           = ValidationError(errors.New("You must supply a password"))
	errPasswordTooShort     = ValidationError(errors.New("Your password is too short"))
	errUsernameExists       = ValidationError(errors.New("That username is already taken"))
	errEmailExists          = ValidationError(errors.New("That email address has already registered an account"))
	errCredentialsIncorrect = ValidationError(errors.New("We couldn't find a user with the supplied username and password combination"))
	errPasswordIncorrect    = ValidationError(errors.New("Password did not match"))

	// Image Manipulation Errors
	errInvalidImageType = ValidationError(errors.New("Please upload only jpeg, gif or png images"))
	errNoImage          = ValidationError(errors.New("Please select an image to upload"))
	errImageURLInvalid  = ValidationError(errors.New("Couldn't download image from the URL you provided"))
)

// IsValidationError returns true if the given error is a user input validation error
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}
