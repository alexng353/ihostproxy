package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type SQLiteCredentialStore struct {
	db *sql.DB
}

type JsonCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (store *SQLiteCredentialStore) Load(cred_env string) (error error) {
	var parsed_env_creds []JsonCredentials
	err := json.Unmarshal([]byte(cred_env), &parsed_env_creds)
	if err != nil {
		return err
	}

	for _, cred := range parsed_env_creds {
		err := store.AddEntry(cred.Username, cred.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

// dataBaseFile defaults to "main.db"
func NewSQLiteCredentialStore(dataBaseFile ...string) *SQLiteCredentialStore {
	var dataSource string = "main.db"
	if len(dataBaseFile) >= 1 {
		dataSource = dataBaseFile[0]
	}

	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		slog.Error("Failed to open database", "database", dataSource, "error", err)
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (id TEXT, username TEXT, password TEXT);
						CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON users (username);
						CREATE UNIQUE INDEX IF NOT EXISTS idx_id ON users (id);`
	_, err = db.Exec(query)

	if err != nil {
		slog.Error("Failed to create table", "database", dataSource, "error", err)
		log.Fatal(err)
	}

	return &SQLiteCredentialStore{db: db}
}

func (store *SQLiteCredentialStore) Valid(user, password, _ string) bool {
	var hash []byte
	query := "SELECT id, password FROM users WHERE username=?"
	err := store.db.QueryRow(query, user).Scan(&user, &hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		slog.Error("Failed to compare password", "username", user, "error", err)
		return false
	}

	return true
}

func (store *SQLiteCredentialStore) AddEntry(user, password string) (err error) {
	hash, err := hashPw(password)

	uid := P.Gen("user")

	query := "INSERT INTO users(id, username, password) VALUES(?, ?, ?)"
	_, err = store.db.Exec(query, uid, user, hash)

	if err != nil {
		slog.Error("Failed to add entry", "username", user, "error", err)

		return err
	}

	return nil
}

func (store *SQLiteCredentialStore) RemoveEntry(user string) (err error) {
	query := "DELETE FROM users WHERE username=?"
	_, err = store.db.Exec(query, user)

	if err != nil {
		slog.Error("Failed to remove entry", "username", user, "error", err)
		return err
	}

	return nil
}

func (store *SQLiteCredentialStore) GetEntry(user, password string) (id string, err error) {
	var hash []byte
	query := "SELECT id, password FROM users WHERE username=?"
	err = store.db.QueryRow(query, user).Scan(&id, &hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Error("Failed to compare password", "username", user, "error", err)
		return "", err
	}

	if err != nil {
		slog.Error("Failed to get entry", "username", user, "error", err)
		return "", err
	}

	return id, nil
}

func hashPw(password string) (hash string, err error) {
	str_hash, hash_err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hash_err != nil {
		slog.Error("Failed to hash password", "error", err)
		return "", hash_err
	}

	return string(str_hash), nil
}
