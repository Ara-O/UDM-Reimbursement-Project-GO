package foapa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type retrievedFoapaData struct {
	ID           primitive.ObjectID    `bson:"_id"`
	FoapaDetails []models.FoapaDetails `bson:"foapa_details" json:"foapa_details"`
}

func getUserIdFromMiddleware(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		return "", err
	}

	userId := claims["user_id"].(string)

	return userId, nil
}

func convertStringToId(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func RetrieveFoapaDetails(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIdFromMiddleware(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("id:", id)

	db := database.GetMongoDbConnection()
	coll := db.Database("udm-go").Collection("faculties")

	// Convert user id from
	objectId, err := convertStringToId(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	filter := bson.M{"_id": objectId}
	projection := bson.M{"foapa_details": 1}

	var foapaInformation retrievedFoapaData
	coll.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&foapaInformation)

	fmt.Printf("%+v", foapaInformation)

	formattedJson, err := json.Marshal(foapaInformation.FoapaDetails)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write(formattedJson)
}
