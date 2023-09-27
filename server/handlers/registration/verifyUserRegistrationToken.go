package registration

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenData struct {
	Token string `json:"token"`
}

func ParseToken(token TokenData) (TokenData, error) {
	// token := strings.ReplaceAll(token.Token, "$", ".")
	return token, nil
}

func VerifyUserRegistrationToken(w http.ResponseWriter, r *http.Request) {
	var token TokenData

	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	fmt.Printf("%+v", token)
}
