package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func checkIfUserExistsInMongo(email string) error {
	db := database.GetMongoDbConnection()
	coll := db.Database("udm-go").Collection("faculties")

	projection := bson.M{"_id": 1}

	var user struct {
		ID interface{} `bson:"_id"`
	}

	if err := coll.FindOne(context.Background(), bson.M{"work_email": email}, options.FindOne().SetProjection(projection)).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("No user found")
		}

		return err
	}

	return nil
}

func hashNewPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func updatePassword(email string, password string) error {
	db := database.GetMongoDbConnection()
	coll := db.Database("udm-go").Collection("faculties")

	_, err := coll.UpdateOne(context.Background(), bson.M{"work_email": email}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "password", Value: password}}}})

	return err
}

func getEmailFromRedis(token string) (string, error) {
	db := database.GetRedisDatabaseConnection()

	email, err := db.Get(context.Background(), token).Result()
	if err != nil {
		return "", err
	}

	return email, nil
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	// Decoding user token
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if token exists in redis
	userEmail, err := getEmailFromRedis(reqData.Token)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Check if user exists in mongo
	if err := checkIfUserExistsInMongo(userEmail); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Hash new password
	hash, err := hashNewPassword(reqData.NewPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Updating the user's password
	if err := updatePassword(userEmail, hash); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println(reqData)
}
