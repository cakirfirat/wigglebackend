package main

import (
	"log"
	"net/http"
	. "wigglebackend/handlers"

	. "wigglebackend/helpers"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("Proje başlatıldı...")

	r := mux.NewRouter()

	/* AUTH ROUTES FOR GIHUB */
	r.HandleFunc("/api/v1/genesis", GenesisHandler).Methods("POST")
	r.HandleFunc("/api/v1/verify-code", VerifyCodeHandler).Methods("POST")
	r.Handle("/api/v1/register", ValidateJwt(RegisterHandler)).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/forgot-password", ForgotPasswordHandler).Methods("POST")

	/* PARTNER ROUTES */
	r.Handle("/api/v1/add-partner", ValidateJwt(AddPartnerHandler)).Methods("POST")

	server := &http.Server{

		Addr:    ":8081",
		Handler: r,
	}

	server.ListenAndServe()

}
