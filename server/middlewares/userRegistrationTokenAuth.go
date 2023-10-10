package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
)

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	tokenType, token := strings.Split(authHeader, " ")[0], strings.Split(authHeader, " ")[1]

	if strings.ToUpper(strings.TrimSpace(tokenType)) != "BEARER" {
		fmt.Println("Invalid authorization type")
		return "", errors.New("Invalid authorization type")
	}

	return token, nil
}

func getUserDataFromRedis(token string) (models.UserDataPreVerification, error) {
	var userData models.UserDataPreVerification
	db := database.GetRedisDatabaseConnection()

	//Get the stored user data from redis
	val, err := db.Get(context.Background(), token).Result()
	if err != nil {
		return models.UserDataPreVerification{}, err
	}

	err = json.Unmarshal([]byte(val), &userData)
	if err != nil {
		return models.UserDataPreVerification{}, err
	}

	return userData, nil
}

// Middleware
func UserRegistrationTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		defer r.Body.Close()

		userData, err := getUserDataFromRedis(token)
		if err != nil {
			fmt.Println("err", err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid token, please restart the registration process"))
			return
		}

		formattedUserData, err := json.Marshal(&userData)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Storing the user data in the context to be used
		ctx := context.WithValue(r.Context(), "userData", formattedUserData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
