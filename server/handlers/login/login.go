package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type RequestData struct {
	WorkEmail string `json:"work_email" bson:"work_email"`
	Password  string `json:"password" bson:"password"`
}

func validatePassword(textPassword string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(textPassword))
}

func createJwtToken(id interface{}) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_TOKEN_KEY")), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": id})
	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := database.GetMongoDbConnection()
	collection := db.Database("udm-go").Collection("faculties")

	//Decode user request
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check for existing faculty
	userEmail := fmt.Sprintf("%s@udmercy.edu", reqData.WorkEmail)
	var userData models.UserDataPostVerification

	start := time.Now()
	err := collection.FindOne(context.Background(), bson.M{"work_email": userEmail}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No faculty found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	end := time.Now()
	fmt.Println("Getting full data: ", end.Sub(start))

	start = time.Now()

	//Validate password
	if err = validatePassword(reqData.Password, userData.Password); err != nil {
		fmt.Println("Incorrect password")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	end = time.Now()

	fmt.Println("Validating password: ", end.Sub(start))

	//Create jwt with the user id
	var userId struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	projection := bson.M{"_id": 1}

	start = time.Now()
	if err = collection.FindOne(context.Background(), bson.M{"work_email": userEmail}, options.FindOne().SetProjection(projection)).Decode(&userId); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	end = time.Now()

	fmt.Println("Finding only id: ", end.Sub(start))

	start = time.Now()
	token, err := createJwtToken(userId.ID)
	end = time.Now()

	fmt.Println("Creating token: ", end.Sub(start))

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Return token
	w.WriteHeader(200)
	w.Write([]byte(token))
}
