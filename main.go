package main

import (
	"net/http"
	"reCaptcha/recaptcha"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Got a user"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/getuser", recaptcha.RecaptchaMiddleware(GetUserHandler))

	http.ListenAndServe(":8000", mux)
}
