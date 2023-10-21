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

func DeleteFoapaDetail(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		FoapaDetail models.FoapaDetails `json:"foapa_detail"`
	}

	db := database.GetMongoDbConnection()
	id, err := getUserIdFromMiddleware(r)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	coll := db.Database("udm-go").Collection("faculties")

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}

	fmt.Println("FIlter", filter)
	_, err = coll.UpdateOne(context.Background(), filter)

	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(200)
}
