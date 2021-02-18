package main

import "database/sql"

var globalImageStore ImageStore

const pageSize = 25

// DBImageStore is a database implementation of the ImageStore interface
type DBImageStore struct {
	db *sql.DB
}

// NewDBImageStore returns a newly created mysql DBImageStore
func NewDBImageStore() ImageStore {
	return &DBImageStore{
		db: globalMySQLDB,
	}
}

// Save image in mysql database
func (store *DBImageStore) Save(image *Image) error {
	_, err := store.db.Exec(`
	REPLACE INTO images
	  (id, user_id, name, location, description, size, created_at)
	VALUES
	   (?, ?, ?, ?, ?, ?, ?)
	`,
		image.ID,
		image.UserID,
		image.Name,
		image.Location,
		image.Description,
		image.Size,
		image.CreatedAt,
	)
	return err
}

// Find returns the image with the given id from the mysql database
func (store *DBImageStore) Find(id string) (*Image, error) {
	row := store.db.QueryRow(`
	SELECT id, user_id, name, location, description, size, created_at
	FROM images
	WHERE id = ?
	`,
		id,
	)

	image := Image{}
	err := row.Scan(
		&image.ID,
		&image.UserID,
		&image.Location,
		&image.Description,
		&image.Size,
		&image.CreatedAt,
	)
	return &image, err
}

// FindAll returns a list of images from the mysql database
func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	rows, err := store.db.Query(`
	SELECT id, user_id, name, location, description, size, created_at
	FROM images
	ORDER BY created_at DESC
	LIMIT ?
	OFFSET ?
	`,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

// FindAllByUser returns a list of images of that user from the mysql database
func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	rows, err := store.db.Query(`
		SELECT id, user_id, name, location, description, size, created_at
		FROM images
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
		`,
		user.ID,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}
