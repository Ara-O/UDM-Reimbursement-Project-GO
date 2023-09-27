package registration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/gomail.v2"
)

type UserData struct {
	FirstName        string `json:"first_name" validate:"required,min=2,max=50" mapstructure:"first_name"`
	LastName         string `json:"last_name" validate:"required,min=2,max=50" mapstructure:"last_name"`
	PhoneNumber      int64  `json:"phone_number" validate:"required,number" mapstructure:"phone_number"`
	WorkEmail        string `json:"work_email" validate:"required,alphanum" mapstructure:"work_email"`
	EmploymentNumber int64  `json:"employment_number" validate:"required,number" mapstructure:"employment_number"`
	Department       string `json:"department" validate:"required" mapstructure:"department"`
}

func loadEnvironmentVariables() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error reading environment variables")
		return err
	}
	return nil
}

func validateStruct(userData *UserData) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(userData)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func createJWT(userData *UserData) (string, error) {
	var tokenAuth *jwtauth.JWTAuth

	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_TOKEN_KEY")), nil)

	formattedUserData := make(map[string]interface{})

	//Convert the userData struct to a map[string]interface{}
	if err := mapstructure.Decode(userData, &formattedUserData); err != nil {
		fmt.Println(err)
		return "", err
	}

	_, tokenString, err := tokenAuth.Encode(formattedUserData)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func sanitizeUserData(userData *UserData) UserData {
	userData.FirstName = strings.TrimSpace(userData.FirstName)
	userData.LastName = strings.TrimSpace(userData.LastName)
	userData.WorkEmail = strings.ToLower(strings.TrimSpace(userData.WorkEmail))
	userData.WorkEmail = userData.WorkEmail + "@udmercy.edu"
	return *userData
}

func sendEmail(userData *UserData, token string) error {
	// Create and send a new message
	m := gomail.NewMessage()
	m.SetHeader("From", "UDM Reimbursement Team <ara@araoladipo.dev>")
	m.SetHeader("To", userData.WorkEmail)
	m.SetHeader("Subject", "Verify your UDM Email")

	// Alter token so that it is embeddable in a url
	token = strings.ReplaceAll(token, ".", "$")

	url := fmt.Sprintf("http://localhost:5173/complete-verification/%s", token)

	mailTemplate := fmt.Sprintf(`
	<div style="background: white">
	<h3 style="font-weight: 500">Verify your Account</h3>
	<h4 style="font-weight: 300">Hello %s,</h4>
	<h4 style="font-weight: 300">Thanks for signing up for the University of Detroit Mercy Reimbursement System!</h4>
	<h4 style="font-weight: 300">You can verify your account with this link</h4>
	<a href="%s"><button style="font-weight: 300; color: white; text-decoration: none; background: #a5093e; padding: 7px 20px">Here</button></a>
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
	// Load environment variables
	if err := loadEnvironmentVariables(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Read data from request
	var userData UserData

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

	//TODO: Replace with redis
	//Create JWTs
	token, err := createJWT(&userData)

	fmt.Println(token)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := sendEmail(&userData, token); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email Sent!"))

}
