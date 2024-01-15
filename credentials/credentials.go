package credentials

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexng353/ihostproxy/pika"
)

type StaticCredentials map[string]struct{}

func (s StaticCredentials) Valid(user string) bool {
	_, ok := s[user]
	if ok {
		return true
	}

	return false
}

type SQLiteCredentialStore struct {
	db *sql.DB
	// cache map[string]bool
	cache StaticCredentials
}

type JsonCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var P = pika.Get()

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

var credentials = NewSQLiteCredentialStore()

func Get() *SQLiteCredentialStore {
	return credentials
}

// dataBaseFile defaults to "main.db"
func NewSQLiteCredentialStore(dataBaseFile ...string) *SQLiteCredentialStore {
	var dataSource string = "main.db"

	var dbpath = os.Getenv("DB_PATH")
	if dbpath != "" {
		dataSource = dbpath
	}

	if len(dataBaseFile) >= 1 {
		dataSource = dataBaseFile[0]
	}

	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		slog.Error("Failed to open database", "database", dataSource, "error", err)
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (id TEXT, username TEXT, password TEXT, admin BOOL);
						CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON users (username);
						CREATE UNIQUE INDEX IF NOT EXISTS idx_id ON users (id);`
	_, err = db.Exec(query)

	if err != nil {
		slog.Error("Failed to create table", "database", dataSource, "error", err)
		log.Fatal(err)
	}

	var usernames []string
	query = `SELECT username FROM users`

	rows, err := db.Query(query)
	if err != nil {
		slog.Error("Failed to get ids", "database", dataSource, "error", err)
	}

	for rows.Next() {
		var usrnm string
		err = rows.Scan(&usrnm)
		if err != nil {
			slog.Error("Failed to scan id", "database", dataSource, "error", err)
		}

		usernames = append(usernames, usrnm)
	}

	slog.Info("Loaded usernames", "usernames", usernames)

	cache := make(StaticCredentials)

	for _, id := range usernames {
		cache[id] = struct{}{}
	}

	return &SQLiteCredentialStore{db: db, cache: cache}
}

// var inmemorycache = make(map[string]bool)

func (store *SQLiteCredentialStore) Valid(user, password, _ string) bool {
	slog.Info("validating user", "username", user)
	// fmt.Println("cache", store.cache)
	slog.String("cache", fmt.Sprintf("%v", store.cache))

	if store.cache.Valid(user) {
		slog.Info("user already validated", "username", user)
		return true
	}
	slog.Warn("uh oh, we are going into database logic land :(")
	// func (store *SQLiteCredentialStore) Valid(user, password string) bool {

	var hash []byte
	query := "SELECT id, password FROM users WHERE username=?"
	err := store.db.QueryRow(query, user).Scan(&user, &hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		slog.Error("Failed to compare password", "username", user, "error", err)
		return false
	}

	store.cache[user] = struct{}{}

	return true
}

func (store *SQLiteCredentialStore) AddEntry(user, password string) (err error) {
	hash, err := hashPw(password)

	uid := P.Gen("user")

	query := "INSERT INTO users(id, username, password, admin) VALUES(?, ?, ?, ?)"
	_, err = store.db.Exec(query, uid, user, hash, false)

	if err != nil {
		slog.Error("Failed to add entry", "username", user, "error", err)

		return err
	}

	return nil
}

func (store *SQLiteCredentialStore) AddAdmin(user, password string) (err error) {
	hash, err := hashPw(password)

	uid := P.Gen("user")

	query := "INSERT INTO users(id, username, password, admin) VALUES(?, ?, ?, ?)"
	_, err = store.db.Exec(query, uid, user, hash, true)

	if err != nil {
		slog.Error("Failed to add entry", "username", user, "error", err)

		return err
	}

	return nil
}

func (store *SQLiteCredentialStore) RemoveEntry(id string) (err error) {
	query := "DELETE FROM users WHERE id=?"
	_, err = store.db.Exec(query, id)

	if err != nil {
		slog.Error("Failed to remove entry", "id", id, "error", err)
		return err
	}

	return nil
}

type User struct {
	Id             string
	Username       string
	HashedPassword string
	Admin          bool
}

func (store *SQLiteCredentialStore) GetEntry(username, password string) (*User, error) {
	var hash []byte
	query := "SELECT * FROM users WHERE username=?"

	var user = &User{}
	err := store.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &hash, &user.Admin)
	user.HashedPassword = string(hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Error("Failed to compare password", "username", username, "error", err)
		return nil, err
	}

	if err != nil {
		slog.Error("Failed to get entry", "username", username, "error", err)
		return nil, err
	}

	return user, nil
}

func (store *SQLiteCredentialStore) GetUser(id string) (*User, error) {
	var user = &User{}

	query := "SELECT * FROM users WHERE id=?"
	err := store.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.HashedPassword, &user.Admin)

	if err != nil {
		slog.Error("Failed to get entry", "id", id, "error", err)
		return nil, err
	}

	return user, nil
}

func (store *SQLiteCredentialStore) GetUsers() ([]*User, error) {
	var users []*User

	query := "SELECT * FROM users"
	rows, err := store.db.Query(query)

	if err != nil {
		slog.Error("Failed to get users", "error", err)
		return nil, err
	}

	for rows.Next() {
		var user = &User{}
		err = rows.Scan(&user.Id, &user.Username, &user.HashedPassword, &user.Admin)
		if err != nil {
			slog.Error("Failed to scan user", "error", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func hashPw(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return "", err
	}

	return string(hash), nil
}
