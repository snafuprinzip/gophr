package main

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// SessionStore is an abstraction interface to easily switch between storage
// backends e.g. files, database etc.
type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

// global list of server sessions
var globalSessionStore SessionStore

// FileSessionStore is a file based implementation of the SessionStore interface
type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

// NewFileSessionStore loads the SessionStore from file or returns a new one
// if the file doesn't exist
func NewFileSessionStore(name string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		Sessions: map[string]Session{},
		filename: name,
	}

	contents, err := ioutil.ReadFile(name)

	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = yaml.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, err
}

// Find returns the Session with the given id or nil if not found
func (store *FileSessionStore) Find(id string) (*Session, error) {
	session, exists := store.Sessions[id]
	if !exists {
		return nil, nil
	}

	return &session, nil
}

// Save stores the Session in a yaml file
func (store *FileSessionStore) Save(session *Session) error {
	store.Sessions[session.ID] = *session
	//	contents, err := json.MarshalIndent(store, "", "  ")
	contents, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

// Delete removes a Session from the store
func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)
	contents, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}
