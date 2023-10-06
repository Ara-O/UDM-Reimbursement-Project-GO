package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func checkIfFacultyAlreadyExists(w http.ResponseWriter, collection *mongo.Collection, employmentNumber string) {
	filter := bson.M{"employment_number": employmentNumber}

	var faculty models.UserDataPostVerification

	err := collection.FindOne(context.Background(), filter).Decode(&faculty)

	if err == nil {
		// There were no documents
		if err != mongo.ErrNoDocuments {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	mongoDb := database.GetMongoDbConnection()
	db := mongoDb.Database("udm-go").Collection("faculties")

	// Decode request body
	var userData models.UserDataPostVerification
	json.NewDecoder(r.Body).Decode(&userData)

	formattedEmploymentNumber := fmt.Sprintf("T%d", userData.EmploymentNumber)
	checkIfFacultyAlreadyExists(w, db, formattedEmploymentNumber)

	fmt.Printf("%+v", userData)
}
