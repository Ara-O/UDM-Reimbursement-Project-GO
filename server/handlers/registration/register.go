package registration

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var userData models.UserDataPostVerification
	json.NewDecoder(r.Body).Decode(&userData)

	fmt.Printf("%+v", userData)
}
