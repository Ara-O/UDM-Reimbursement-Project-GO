package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func validateStruct(userData *models.UserDataPreVerification) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(userData)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func storeUserDataInRedis(userData *models.UserDataPreVerification) (string, error) {
	db := database.GetRedisDatabaseConnection()
	userId := uuid.New().String()

	err := db.Set(context.Background(), userId, userData, 15*time.Minute).Err()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return userId, nil
}

func sanitizeUserData(userData *models.UserDataPreVerification) models.UserDataPreVerification {
	userData.FirstName = strings.TrimSpace(userData.FirstName)
	userData.LastName = strings.TrimSpace(userData.LastName)
	userData.WorkEmail = strings.ToLower(strings.TrimSpace(userData.WorkEmail))
	userData.WorkEmail = userData.WorkEmail + "@udmercy.edu"
	return *userData
}

func sendEmail(userData *models.UserDataPreVerification, id string) error {
	// Create and send a new message
	m := gomail.NewMessage()
	m.SetHeader("From", "UDM Reimbursement Team <ara@araoladipo.dev>")
	m.SetHeader("To", userData.WorkEmail)
	m.SetHeader("Subject", "Verify your UDM Email")

	url := fmt.Sprintf("http://localhost:5173/complete-verification/%s", id)

	mailTemplate := fmt.Sprintf(`
	<div style="background: white">
	<h3 style="font-weight: 500">Verify your Account</h3>
	<h4 style="font-weight: 300">Hello %s,</h4>
	<h4 style="font-weight: 300">Thanks for signing up for the University of Detroit Mercy Reimbursement System!</h4>
	<h4 style="font-weight: 300">You can verify your account with this link</h4>
	<a href="%s"><button style="font-weight: 300; cursor:pointer; color: white; text-decoration: none; background: #a5093e; padding: 7px 20px; border: none">Here</button></a>
	</div>`, userData.FirstName, url)

	// Set the email body
	m.SetBody("text/html", mailTemplate)

	// Create a new SMTP client
	d := gomail.NewDialer(os.Getenv("SENDGRID_URL"), 587, os.Getenv("SENDGRID_USERNAME"), os.Getenv("SENDGRID_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Main function
func SendConfirmationEmail(w http.ResponseWriter, r *http.Request) {
	// Read data from request
	var userData models.UserDataPreVerification

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//Validate struct format
	if err := validateStruct(&userData); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Sanitize user data
	userData = sanitizeUserData(&userData)

	// Store user data in redis to
	cacheID, err := storeUserDataInRedis(&userData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := sendEmail(&userData, cacheID); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email Sent!"))
}
