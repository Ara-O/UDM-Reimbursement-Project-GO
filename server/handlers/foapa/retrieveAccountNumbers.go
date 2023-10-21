package foapa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"go.mongodb.org/mongo-driver/bson"
)

func RetrieveAccountNumbers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("baaaaak")
	db := database.GetMongoDbConnection()

	coll := db.Database("udm-go").Collection("foapa-account-numbers")

	_, err := getUserIdFromMiddleware(r)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var accountData struct {
		AccountNumbers []struct {
			Number      string `json:"number"`
			Description string `json:"description"`
		} `json:"accountNumbers"`
	}

	if err := coll.FindOne(context.Background(), bson.M{}).Decode(&accountData); err != nil {
		fmt.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resData, err := json.Marshal(accountData.AccountNumbers)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(resData)
}
