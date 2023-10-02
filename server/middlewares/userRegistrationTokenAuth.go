package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
)

type TokenData struct {
	Token string `json:"token"`
}

func getUserDataFromRedis(token TokenData) (models.UserDataPreVerification, error) {
	var userData models.UserDataPreVerification
	db := database.GetRedisDatabaseConnection()

	//Get the stored user data from redis
	val, err := db.Get(context.Background(), token.Token).Result()
	if err != nil {
		return models.UserDataPreVerification{}, err
	}

	err = json.Unmarshal([]byte(val), &userData)
	if err != nil {
		return models.UserDataPreVerification{}, err
	}

	return userData, nil
}

func UserRegistrationTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token TokenData

		if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		userData, err := getUserDataFromRedis(token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid token, please restart the registration process"))
			return
		}

		formattedUserData, err := json.Marshal(&userData)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// Storing the user data in the context to be used
		ctx := context.WithValue(r.Context(), "userData", formattedUserData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
