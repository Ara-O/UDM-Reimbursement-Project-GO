package router

import (
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/handlers/registration"
	"github.com/go-chi/chi"
)

func DefineRoutes(r *chi.Mux) {
	r.Post("/api/send-confirmation-email", registration.SendConfirmationEmail)
	r.Post("/api/verify-user-registration-token", registration.VerifyUserRegistrationToken)
}
