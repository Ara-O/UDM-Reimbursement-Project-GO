package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/application"
	"github.com/Ara-Oladipo/UDM-Reimbursement-Project-Go/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error reading environment variables")
		log.Fatal(err)
	}

	if err := database.StartDatabase(); err != nil {
		log.Fatal(err)
	}

	app := application.New(":8080")

	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe(app.Addr, app.Routes))
}
