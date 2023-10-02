package registration

import (
	"net/http"
)

// Uses the userRegistrationAuth middleware so if this gets hit, then it was a success
func VerifyUserRegistrationToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Registration token verified"))
}
