package registration

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
)

// Uses the userRegistrationAuth middleware so if this gets hit, then it was a success
// Returns the user data
func VerifyUserRegistrationToken(w http.ResponseWriter, r *http.Request) {
	var userData models.UserDataPreVerification

	parsedUserData, ok := r.Context().Value("userData").([]byte)
	if !ok {
		fmt.Println("Error converting to slice of bytes")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Takes the user data stored from the middleware, parses it to make sure no errors
	if err := json.Unmarshal(parsedUserData, &userData); err != nil {
		fmt.Println("Error parsing JSON")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Write(parsedUserData)
}
