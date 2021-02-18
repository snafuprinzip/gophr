package main

import "golang.org/x/crypto/bcrypt"

// User contains all user account information
type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength
)

// NewUser creates a new user account record, validates the user input
// and hashes the password
func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}
	if username == "" {
		return user, errNoUsername
	}
	if email == "" {
		return user, errNoEmail
	}
	if password == "" {
		return user, errNoPassword
	}
	if len(password) < passwordLength {
		return user, errPasswordTooShort
	}

	// Check if the username exists
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errUsernameExists
	}

	// Check if the email exists
	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)
	return user, err
}

// FindUser looks for a user with a specific username/password combination
// If both are matched the user and no error (nil) will be returned, if the user isn't found
// or the password doesn't match a newly created user will be returned with the
// form's values filled in and an errCredentialsIncorrect will that indicate a mismatch
func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	// find the user
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return out, errCredentialsIncorrect
	}

	// check for a password match
	if bcrypt.CompareHashAndPassword([]byte(existingUser.HashedPassword), []byte(password)) != nil {
		return out, errCredentialsIncorrect

	}

	// return a full match
	return existingUser, nil
}

// UpdateUser updates the User record with a new email and password
func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	// create shallow copy from user
	out := *user
	out.Email = email

	// Check if the email exists
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return out, errEmailExists
	}

	// current password empty
	if currentPassword == "" {
		return out, errNoPassword
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword),
	) != nil {
		return out, errPasswordIncorrect
	}

	// update the real user's email address
	user.Email = email

	if newPassword == "" {
		return out, nil
	}

	if len(newPassword) < passwordLength {
		return out, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}
