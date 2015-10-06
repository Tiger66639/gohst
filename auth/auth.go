package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

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

func GetConnectedUser(r *http.Request) *data.User {
	session, _ := store.Get(r, "session-name")
	if receipt, ok := session.Values["receipt"].(string); ok {
		user, _ := data.GetUserFromId(sessionIDs[receipt])
		return user
	}
	return nil
}

// Connect creates a valid session if the correct authentication params
// are provided
func Connect(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("username")) < 6 || len(r.FormValue("password")) < 6 {
		log.Printf("username or password too small!")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	username := strings.Trim(r.FormValue("username"), " ")
	password := strings.Trim(r.FormValue("password"), " ")
	session, _ := store.Get(r, "session-name")
	title := r.FormValue("redirect")
	salt, err := data.GetSalt(username)
	if err != nil {
		http.Redirect(w, r, title, http.StatusFound)
		return
	}

	hash := Hash(salt, password)

	if data.DoHashesMatch(username, hash) {
		receipt := randomString(32)
		userid := data.GetUserId(username)
		sessionIDs[receipt] = userid
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
	salt := randomString(64)

	if len(r.FormValue("password")) < 6 || len(r.FormValue("username")) < 6 || len(r.FormValue("email")) < 6 {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	user := &data.User{Salt: salt, Hash: Hash(salt, r.FormValue("password")), Username: r.FormValue("username"), Email: r.FormValue("email"), Role: "user"}

	rows, _ := data.AddNewUser(user).RowsAffected()

	if rows != 1 {
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
