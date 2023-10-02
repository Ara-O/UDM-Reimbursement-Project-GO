package router

import (
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/registration"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/middlewares"
	"github.com/go-chi/chi"
)

func DefineRoutes(r *chi.Mux) {

	r.Post("/api/send-confirmation-email", registration.SendConfirmationEmail)
	// Requires signup token authentication
	r.Group(func(r chi.Router) {
		r.Use(middlewares.UserRegistrationTokenAuth)
		r.Post("/api/verify-user-registration-token", registration.VerifyUserRegistrationToken)
		r.Post("/api/register", registration.Register)
	})
}
