package reimbursement

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

type Activity struct {
	ActivityName    string `json:"activityName"`
	Cost            string `json:"cost"`
	ActivityDate    string `json:"activityDate"`
	ActivityID      string `json:"activityId"`
	ActivityReceipt string `json:"activityReceipt"`
	FoapaNumber     string `json:"foapaNumber"`
}

type Ticket struct {
	UDMPUVoucher           bool       `json:"udmpuVoucher"`
	Activities             []Activity `json:"activities"`
	Destination            string     `json:"destination"`
	PaymentRetrievalMethod string     `json:"paymentRetrievalMethod"`
	ReimbursementDate      string     `json:"reimbursementDate"`
	ReimbursementName      string     `json:"reimbursementName"`
	ReimbursementReason    string     `json:"reimbursementReason"`
	ReimbursementReceipts  []string   `json:"reimbursementReceipts"`
	ReimbursementStatus    string     `json:"reimbursementStatus"`
	TotalCost              int64      `json:"totalCost"`
}

func getUserIdFromMiddleware(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		return "", err
	}

	userId := claims["user_id"].(string)

	return userId, nil
}

func AddReimbursement(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	json.NewDecoder(r.Body).Decode(&ticket)

	db := database.GetMongoDbConnection()
	coll := db.Database("udm-go").Collection("reimbursement-claims")
	faculty_coll := db.Database("udm-go").Collection("faculties")

	res, err := coll.InsertOne(context.Background(), ticket)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId, _ := getUserIdFromMiddleware(r)

	id, err := primitive.ObjectIDFromHex(userId)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	faculty_coll.FindOneAndUpdate(context.Background(), filter, bson.M{"$push": bson.M{"reimbursement_tickets": res.InsertedID}})
	fmt.Println(res)
	// fmt.Printf("%+v", ticket)
}
