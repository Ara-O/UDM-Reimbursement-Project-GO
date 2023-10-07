package router

import (
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/login"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/registration"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/middlewares"
	"github.com/go-chi/chi"
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
	})

	// Requires signup token authentication
}
