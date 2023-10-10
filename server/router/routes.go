package router

import (
	"os"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/foapa"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/login"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/registration"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func DefineRoutes(r *chi.Mux) {
	// Registration
	r.Group(func(r chi.Router) {
		r.Post("/api/send-confirmation-email", registration.SendConfirmationEmail)
		r.Post("/api/verify-user-information", registration.VerifyUserInformation)

		// Requires authenticating token for pre-auth data
		r.Group(func(r chi.Router) {
			r.Use(middlewares.UserRegistrationTokenAuth)
			r.Get("/api/verify-user-registration-token", registration.VerifyUserRegistrationToken)
			r.Post("/api/register", registration.Register)
		})
	})

	// Login
	r.Group(func(r chi.Router) {
		r.Post("/api/login", login.Login)
		r.Post("/api/forgot-password", login.ForgotPassword)
		r.Post("/api/verify-forgot-password-token", login.VerifyForgotPasswordToken)
		r.Post("/api/reset-password", login.ResetPassword)
	})

	// Requires main token authentication
	r.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_TOKEN_KEY")), nil)

		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/retrieve-foapa-details", foapa.RetrieveFoapaDetails)
	})
}
