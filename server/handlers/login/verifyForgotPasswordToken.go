package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
)

func checkIfTokenExistsInRedis(token string) error {
	db := database.GetRedisDatabaseConnection()

	_, err := db.Get(context.Background(), token).Result()
	if err != nil {
		return err
	}

	return nil
}

func VerifyForgotPasswordToken(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		UserToken string `json:"user_token"`
	}

	// Decoding user token
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := checkIfTokenExistsInRedis(reqData.UserToken); err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Token verified"))
}
