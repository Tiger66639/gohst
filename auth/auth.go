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

func isConnected(session *sessions.Session) bool {
	sessionID := session.Values["sessionID"].(string)
	_, existsID := sessionIDs[sessionID]
	return existsID
}

func Connect(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("error fetching session: ", err)
		return
	}

	if r.FormValue("state") != session.Values["state"].(string) {
		return
	}

	sessionID := session.Values["sessionID"].(string)
	if _, existsID := sessionIDs[sessionID]; !existsID {
		log.Println("Session already exists")
		return
	}

	// TODO: remove hard coding on this to actually be usable in the real world
	if r.FormValue("username") != "gohst" ||
		r.FormValue("password") != "thisisareallybadpassword" {
		log.Println("Wrong credentials used %s, %s",
			r.FormValue("username"), r.FormValue("password"))
		return
	}

	thisSession := RandomString(32)
	sessionIDs[thisSession] = 1

	session.Values["sessionID"] = thisSession
	session.Save(r, w)
	log.Println("login successful")
}

/**
 * This function remove the current session information and attempts
 * to delete the client's cookie
 */
func Disconnect(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("error fetching session: ", err)
		return
	}
	session.Values["sessionID"] = nil
	session.Values["state"] = nil
	return
}

// randomString returns a random string with the specified length
func RandomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
