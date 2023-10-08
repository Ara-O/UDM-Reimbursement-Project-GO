package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func storeEmailInRedis(email string) (string, error) {
	db := database.GetRedisDatabaseConnection()

	id := uuid.NewString()

	err := db.Set(context.Background(), id, email, 15*time.Minute).Err()
	if err != nil {
		return "", err
	}

	return id, nil
}

func sendEmail(email string, id string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "UDM Reimbursement Team <ara@araoladipo.dev>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Forgot Password - UDM Reimbursement System")

	url := fmt.Sprintf("http://localhost:5173/forgot-password/%s", id)

	mailTemplate := fmt.Sprintf(`
	<div style="background: white">
	<h3 style="font-weight: 500">Forgot password</h3>
	<h4 style="font-weight: 300">Hello!</h4>
	<h4 style="font-weight: 300">Here is the recovery link for your password</h4>
	<a href="%s"><button style="font-weight: 300; cursor:pointer; color: white; text-decoration: none; background: #a5093e; padding: 7px 20px; border: none">Here</button></a>
	</div>`, url)

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

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		WorkEmail string `json:"work_email"`
	}

	// Decoding body
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Sanitizing work email
	reqData.WorkEmail = strings.ToLower(strings.TrimSpace(reqData.WorkEmail))
	formattedEmail := fmt.Sprintf("%s@udmercy.edu", reqData.WorkEmail)

	// Storing email in redis
	id, err := storeEmailInRedis(formattedEmail)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Send id to email
	if err := sendEmail(formattedEmail, id); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println(reqData)
}
