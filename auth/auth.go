package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

/**
 * Initializes the store with a random value each time.
 * If the client doesn't have this, or if it has a different one,
 * a new cookie will then need to be generated for them.
 */
var store = sessions.NewCookieStore([]byte(RandomString(32)))

/**
 * Represent the IDs of logged in sessions.
 * Privileges are represented by an integer but currently
 * don't mean anything.
 */
var sessionIDs = make(map[string]int)

func IsConnected(session *sessions.Session) bool {
	if receipt, ok := session.Values["receipt"].(string); ok {
		if _, existsID := sessionIDs[receipt]; !existsID {
			return true
		}
	}
	return false
}

func Connect(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	if IsConnected(session) {
		http.Redirect(w, r, "backend/manage", http.StatusFound)
		return
	}

	// TODO: remove hard coding on this to actually be usable in the real world
	if r.FormValue("username") != "gohst" || r.FormValue("password") != "thisisareallybadpassword" {
		log.Printf("Wrong credentials used %s and %s", r.FormValue("username"), r.FormValue("password"))
		http.Redirect(w, r, "backend", http.StatusFound)
		return
	}

	log.Printf("Correct credentials used!")
	thisSession := RandomString(32)
	sessionIDs[thisSession] = 1
	session.Values["receipt"] = thisSession
	session.Save(r, w)
	http.Redirect(w, r, "backend/manage", http.StatusFound)
}

/**
 * This function remove the current session information and attempts
 * to delete the client's cookie
 */
func Disconnect(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	session.Values["receipt"] = nil
	return
}

// randomString returns a random string with the specified length
func RandomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
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
