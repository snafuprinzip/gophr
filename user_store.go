package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

// UserStore is an abstraction interface each storage backend has to
// provide to store User structs
type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

// FileUserStore is a file storage implementation of the UserStore interface
type FileUserStore struct {
	filename string
	Users    map[string]User
}

var globalUserStore UserStore

// Save stores the user records on file
func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user

	// contents, err := json.MarshalIndent(store, "", "  ")
	contents, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(store.filename, contents, 0660)
	if err != nil {
		return err
	}
	return nil
}

// NewFileUserStore loads the user records from file or returns an empty one
// if the file does not exist
func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		Users:    map[string]User{},
		filename: filename,
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// if the file doesn't exist we return the fresh instance
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = yaml.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Find returns the user with the given id or nil if not found
func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

// FindByUsername returns the user with the given username or nil if not found
func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(username) == strings.ToLower(user.Username) {
			return &user, nil
		}
	}
	return nil, nil
}

// FindByEmail returns the user with the given email address or nil if not found
func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(email) == strings.ToLower(user.Email) {
			return &user, nil
		}
	}
	return nil, nil
}
