// Here is my thinking behind session handling here:
// 1. User logs in
// 2. Perform authentication against what is in the DB to make sure it is the actual user
// 3. If auth fails, bounce to try again
// 4. If auth succeeds, generate a session ID using the username + a timestamp (TS), and concatenate those together like so:
//    test@test.com|2015-06-18 00:19:51.704738594 -0400 EDT
// 5. Generate a random key using gorilla/securecookie, and use it to generate an HMAC of the concatenated string
// 6. Store that hash, key, expiration TS, and the username in a sessions table
// 7. Send back a secure cookie containing the session id, username, and expiration TS
// 8. Hash those values and compare to what is in the DB
// 9. If valid and not expired, proceed; otherwise, take the appropriate action
package main

import (
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/gorilla/securecookie"
	"fmt"
)

// A useful function to generate a session id
func generateSessionId(message string, key []byte) string {
	key = securecookie.GenerateRandomKey(64)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	log.Println(string([]byte(message)))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func setSession(sd SessionData, r http.ResponseWriter) {
	value := map[string]string{
		"session_id": sd.SessionId,
		"username": sd.Username,
		"exp_ts": fmt.Sprint(sd.ExpTime),
	}
	if encoded, err := cookieHandler.Encode("session", value); err != nil {
		cookie := &http.Cookie{
			Name: "session",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(r, cookie)
	}
}
