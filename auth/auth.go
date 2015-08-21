package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/cosban/gohst/data"
	"github.com/gorilla/sessions"
)

/**
 * Initializes the store with a random value each time.
 * If the client doesn't have this, or if it has a different one,
 * a new cookie will then need to be generated for them.
 */
var store = sessions.NewCookieStore([]byte(randomString(32)))

/**
 * Represent the IDs of logged in sessions.
 * Privileges are represented by an integer but currently
 * don't mean anything.
 */
var sessionIDs = make(map[string]int)

// Hash takes a salt and provided string and returns their corresponding
// combined sha256 string
func Hash(salt, provided string) string {
	hasher := sha256.New()
	hasher.Write([]byte(provided))
	first := hex.EncodeToString(hasher.Sum(nil))

	hasher = sha256.New()
	hasher.Write(append([]byte(first), salt...))

	return hex.EncodeToString(hasher.Sum(nil))
}

// IsConnected returns true if a client is currently logged in
func IsConnected(r *http.Request) bool {
	session, _ := store.Get(r, "session-name")
	if receipt, ok := session.Values["receipt"].(string); ok {
		if _, existsID := sessionIDs[receipt]; existsID {
			return true
		}
	}
	return false
}

// Connect creates a valid session if the correct authentication params
// are provided
func Connect(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("username")) < 6 || len(r.FormValue("password")) < 6 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	session, _ := store.Get(r, "session-name")
	title := r.FormValue("redirect")

	salt := data.GetSalt(r.FormValue("username"))
	hash := Hash(salt, r.FormValue("password"))
	if !data.DoHashesMatch(r.FormValue("username"), hash) {
		log.Printf("invalid login attempt")
	} else {
		receipt := randomString(32)
		sessionIDs[receipt] = 1
		session.Values["receipt"] = receipt
		session.Save(r, w)
		// on success we don't want to to the login page again...
		if title == "login" {
			title = "backend/manage"
		}
	}
	http.Redirect(w, r, title, http.StatusFound)
}

// Disconnect removes the current session information and attempts
// to delete the client's cookie
func Disconnect(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if receipt, ok := session.Values["receipt"].(string); ok {
		delete(sessionIDs, receipt)
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "", http.StatusFound)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	salt := randomString(32)
	if len(r.FormValue("password")) < 6 || len(r.FormValue("username")) < 6 || len(r.FormValue("email")) < 6 {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	user := &data.User{Salt: salt,
		Hash:     Hash(salt, r.FormValue("password")),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Role:     "user"}
	result := data.AddNewUser(user)
	rows, _ := result.RowsAffected()
	if rows != 1 {
		// TODO: check if username is taken etc
		log.Printf("OH NO THE ADDING THE USER DIDN'T AFFECT THE CORRECT AMOUNT OF ROWS!")
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// randomString returns a random string with the specified length
func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:length]
}

func decodeBase64(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
