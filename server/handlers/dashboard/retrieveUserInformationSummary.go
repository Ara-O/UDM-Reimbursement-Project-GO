package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResData struct {
	FirstName        string `json:"first_name" bson:"first_name"`
	LastName         string `json:"last_name" bson:"last_name"`
	PhoneNumber      int64  `json:"phone_number" bson:"phone_number"`
	WorkEmail        string `json:"work_email" bson:"work_email"`
	EmploymentNumber int64  `json:"employment_number" bson:"employment_number"`
}

func getUserIdFromMiddleware(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		return "", err
	}

	userId := claims["user_id"].(string)

	return userId, nil
}

func RetrieveUserInformationSummary(w http.ResponseWriter, r *http.Request) {
	db := database.GetMongoDbConnection().Database("udm-go")
	coll := db.Collection("faculties")

	id, err := getUserIdFromMiddleware(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}

	var userData ResData

	if err = coll.FindOne(context.Background(), filter).Decode(&userData); err != nil {
		fmt.Println("User does not exist based on user id")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	resData, err := json.Marshal(userData)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(resData)
}
