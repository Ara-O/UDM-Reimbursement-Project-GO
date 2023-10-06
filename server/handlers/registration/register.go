package registration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func checkIfFacultyAlreadyExists(w http.ResponseWriter, collection *mongo.Collection, employmentNumber int64) error {
	filter := bson.M{"employment_number": employmentNumber}

	var faculty models.UserDataPostVerification

	//Returns an error when it finds nothing
	err := collection.FindOne(context.Background(), filter).Decode(&faculty)

	// A document was found
	if err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return errors.New("Faculty already exists")
	}

	// A document was not found
	if err == mongo.ErrNoDocuments {
		return nil
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func createJwtToken(id interface{}) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_TOKEN_KEY")), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": id})
	return tokenString, nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	mongoDb := database.GetMongoDbConnection()
	db := mongoDb.Database("udm-go").Collection("faculties")

	// Decode request body
	var userData models.UserDataPostVerification
	json.NewDecoder(r.Body).Decode(&userData)

	// Validate request body
	if err := userData.ValidateStruct(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check for error or if faculry already exists
	if err := checkIfFacultyAlreadyExists(w, db, userData.EmploymentNumber); err != nil {
		return
	}

	// Hash password
	hash, err := hashPassword(userData.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Store hashed password
	userData.Password = hash

	// Insert data into database
	res, err := db.InsertOne(context.Background(), userData)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Create jwt
	token, err := createJwtToken(res.InsertedID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(token))

	fmt.Print("TOKEN: ", token)
}
