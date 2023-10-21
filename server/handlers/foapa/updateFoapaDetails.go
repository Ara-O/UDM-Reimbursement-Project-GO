package foapa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateFoapaDetails(w http.ResponseWriter, r *http.Request) {
	reqData := []models.FoapaDetails{}
	db := database.GetMongoDbConnection()
	coll := db.Database("udm-go").Collection("faculties")

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := getUserIdFromMiddleware(r)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId, err := primitive.ObjectIDFromHex(id)

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}

	_, err = coll.UpdateOne(context.Background(), filter, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "foapa_details", Value: reqData}}}})

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
