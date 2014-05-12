package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sender"
)

func main() {
	http.HandleFunc("/email", email)
	log.Println("Contact form server running...")
	http.ListenAndServe(":8889", nil)
}

func email(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder(r.Body)
	msg := sender.EmailMessage{}
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Unable to decode %v with error %v", r.Body, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = sender.SendEmail(msg, r.Host)
	if err != nil {
		log.Printf("Unable to send email, error %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Sending email!")
	w.WriteHeader(http.StatusOK)
}
