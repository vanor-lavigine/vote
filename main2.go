package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"net/http"
)

type ApiResponse struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"errorCode"`
	ErrorMsg  string      `json:"ErrorMsg"`
}

var db *sql.DB
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
	// Setup MySQL connection
	var err error
	db, err = sql.Open("mysql", "<user>:<password>@/<dbname>")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   14 * 24 * 60 * 60, // 14 days
		HttpOnly: true,
	}

	http.HandleFunc("/register", errorHandler(registerHandler))
	http.HandleFunc("/login", errorHandler(loginHandler))
	http.HandleFunc("/logout", errorHandler(logoutHandler))
	http.HandleFunc("/createCandidate", withAuth(adminOnly(errorHandler(createCandidateHandler))))
	http.HandleFunc("/deleteCandidate", withAuth(adminOnly(errorHandler(deleteCandidateHandler))))
	http.HandleFunc("/listCandidates", errorHandler(listCandidatesHandler))

	http.ListenAndServe(":8080", nil)
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				json.NewEncoder(w).Encode(ApiResponse{
					Code:      500,
					ErrorCode: "internal_error",
					ErrorMsg:  err.Error(),
				})
			}
		}()
		fn(w, r)
	}
}

func withAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := checkLogin(r); err != nil {
			json.NewEncoder(w).Encode(ApiResponse{
				Code:      401,
				ErrorCode: "not_logged_in",
				ErrorMsg:  "You need to log in to perform this action.",
			})
			return
		}
		fn(w, r)
	}
}

func adminOnly(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, _ := checkLogin(r)
		if username != "admin" {
			json.NewEncoder(w).Encode(ApiResponse{
				Code:      403,
				ErrorCode: "forbidden",
				ErrorMsg:  "Only admin can perform this action.",
			})
			return
		}
		fn(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if exists {
		json.NewEncoder(w).Encode(ApiResponse{
			Code:      400,
			ErrorCode: "user_exists",
			ErrorMsg:  "Username is already taken.",
		})
		return
	}

	hash := md5.New()
	hash.Write([]byte(username))
	hashValue := hex.EncodeToString(hash.Sum(nil))

	_, err = db.Exec("INSERT INTO users (username, password, hash) VALUES (?, ?, ?)", username, hashedPassword, hashValue)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "User registered successfully.",
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Hash the password
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=? AND password=?)", username, hashedPassword).Scan(&exists)
	if err != nil {
		panic(err)
	}
	if !exists {
		json.NewEncoder(w).Encode(ApiResponse{
			Code:      400,
			ErrorCode: "invalid_credentials",
			ErrorMsg:  "Invalid username or password.",
		})
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["username"] = username
	session.Save(r, w)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-name",
		Value:    session.ID,
		Path:     "/",
		MaxAge:   14 * 24 * 60 * 60,
		HttpOnly: true,
	})


	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Login successful.",
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.SetCookie(w, &http.Cookie{
		Name:   "session-name",
		Path:   "/",
		MaxAge: -1,
	})

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Logout successful.",
	})
}

func createCandidateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	_, err := db.Exec("INSERT INTO candidates (username, valid) VALUES (?, true)", username)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Candidate created successfully.",
	})
}

func deleteCandidateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	_, err := db.Exec("UPDATE candidates SET valid = false WHERE username = ?", username)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: "Candidate deleted successfully.",
	})
}

func listCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username FROM candidates WHERE valid=true")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var candidates []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		err := rows.Scan(&id, &username)
		if err != nil {
			panic(err)
		}
		candidates = append(candidates, map[string]interface{}{
			"id":       id,
			"username": username,
		})
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Code: 200,
		Data: candidates,
	})
}

func checkLogin(r *http.Request) (string, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return "", err
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return "", errors.New("not logged in")
	}

	return username, nil
}
