package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func init() {
	// Initialize the Firebase app with service account
	opt := option.WithCredentialsFile("path to the json file service account.json")
	var err error
	firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}

func validateToken(idToken string) (*auth.Token, error) {
	client, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return nil, err
	}

	log.Printf("Verified ID token: %v\n", token)
	return token, nil
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		fmt.Println("Yaha A gaya")
		w.WriteHeader(http.StatusOK)
		return
	}

	idToken := r.Header.Get("Authorization")
	if idToken == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	fmt.Println(idToken)

	if len(idToken) > 7 && idToken[:7] == "Bearer " {
		idToken = idToken[7:]
	}

	fmt.Println("Yaha A gaya 3")

	token, err := validateToken(idToken)
	if err != nil {
		fmt.Println("token verification failed")
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"uid": token.UID,
		// Add more user data as needed
	}
	fmt.Println(response)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/verify", myHandler)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
