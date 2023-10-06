package registration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestData struct {
	WorkEmail        string `json:"work_email"`
	EmploymentNumber int64  `json:"employment_number"`
}

func validateEmail(collection *mongo.Collection, email string) error {
	formattedEmail := fmt.Sprintf("%s@udmercy.edu", email)
	filter := bson.M{"work_email": formattedEmail}

	var existingFaculty models.UserDataPostVerification
	if err := collection.FindOne(context.Background(), filter).Decode(&existingFaculty); err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing faculty was found
			return nil
		}

		//Another error that needs to be handled occured
		return err
	} else {
		return errors.New("Faculty with the same email already exists")
	}
}

func validateEmploymentNumber(collection *mongo.Collection, employmentNumber int64) error {
	filter := bson.M{"employment_number": employmentNumber}

	var existingFaculty models.UserDataPostVerification
	if err := collection.FindOne(context.Background(), filter).Decode(&existingFaculty); err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing faculty was found
			return nil
		}

		//Another error that needs to be handled occured
		return err
	} else {
		return errors.New("Faculty with the same employment number already exists")
	}
}

func VerifyUserInformation(w http.ResponseWriter, r *http.Request) {
	var userData RequestData
	client := database.GetMongoDbConnection()
	db := client.Database("udm-go").Collection("faculties")

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err := validateEmail(db, userData.WorkEmail)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	err = validateEmploymentNumber(db, userData.EmploymentNumber)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

}
