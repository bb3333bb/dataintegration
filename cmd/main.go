package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)

type UserInfo struct {
	Picture       string
	EmailVerified int
	Sub           string
	Email         string
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/auth/google/login", handleGoogleLogin)
	http.HandleFunc("/auth/google/callback", handleGoogleCallback)
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/auth/google/login">Google Log In</a></body></html>`
	fmt.Fprintf(w, html)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Printf("client token: '%s'\n:", token)

	client := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	//resp, err := client.Get("https://www.googleapis.com/auth/userinfo.email")

	if err != nil {
		log.Printf("client.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var data UserInfo
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &data)
	defer resp.Body.Close()
	log.Printf("Response: %v\n", data)

	fmt.Fprintf(w, "Response: %v\n", resp.Status)
}
