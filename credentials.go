package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	// "github.com/things-go/go-socks5"
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

	query := "CREATE TABLE IF NOT EXISTS users (username TEXT, password TEXT); CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON users (username)"
	_, err = db.Exec(query)

	if err != nil {
		slog.Error("Failed to create table", "database", dataSource, "error", err)
		log.Fatal(err)
	}

	return &SQLiteCredentialStore{db: db}
}

func (store *SQLiteCredentialStore) Valid(user, password, _ string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username=? AND password=?)"
	err := store.db.QueryRow(query, user, password).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}

func (store *SQLiteCredentialStore) AddEntry(user, password string) (err error) {
	query := "INSERT INTO users(username, password) VALUES(?, ?)"
	_, err = store.db.Exec(query, user, password)

	if err != nil {
		slog.Error("Failed to add entry", "username", user, "error", err)

		return err
	}

	return nil
}

func (store *SQLiteCredentialStore) RemoveEntry(user, password string) (err error) {
	query := "DELETE FROM users WHERE username=? AND password=?"
	_, err = store.db.Exec(query, user, password)

	if err != nil {
		slog.Error("Failed to remove entry", "username", user, "error", err)
		return err
	}

	return nil
}

// func main() {
// 	store := NewSQLiteCredentialStore("your_database_file.db")
//
// 	// Use the Valid method
// 	if store.Valid("username", "password", "userAddress") {
// 		fmt.Println("Credentials are valid")
// 	} else {
// 		fmt.Println("Invalid credentials")
// 	}
// }
