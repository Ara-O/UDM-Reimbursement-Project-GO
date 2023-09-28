package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
)

type TokenData struct {
	Token string `json:"token"`
}

func GetUserDataFromRedis(token TokenData) (UserData, error) {
	db := database.GetRedisDatabaseConnection()

	//Get the stored user data from redis
	val, err := db.Get(context.Background(), token.Token).Result()
	if err != nil {
		return UserData{}, err
	}

	fmt.Println(val)

	return UserData{}, nil
}

func VerifyUserRegistrationToken(w http.ResponseWriter, r *http.Request) {
	var token TokenData

	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	userData, err := GetUserDataFromRedis(token)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token, please try to sign up again"))
		return
	}

	fmt.Println(userData)
}
