package main

import (
	"context"
	"net/http"
	"strings"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	jwtSecretKey = []byte("your_secret_key")
	oauthConfig  = &oauth2.Config{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func (m *Middleware) Authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next(w, r.WithContext(ctx), ps)
	}
}

func (m *Middleware) Authorize(roles ...string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claims, ok := r.Context().Value("claims").(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				next(w, r, ps)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func (m *Middleware) OAuthLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (m *Middleware) OAuthCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Store user information in the database
	if err := m.StoreUserInfo(userInfo); err != nil {
		http.Error(w, "Failed to store user info", http.StatusInternalServerError)
		return
	}

	// Process user info and generate JWT token
	// Placeholder for user info processing and JWT token generation
}

func (m *Middleware) GetUserInfo(username string) (map[string]interface{}, error) {
	var userInfo map[string]interface{}
	err := m.dbPool.QueryRow(context.Background(), "SELECT info FROM users WHERE username=$1", username).Scan(&userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (m *Middleware) StoreUserInfo(userInfo map[string]interface{}) error {
	_, err := m.dbPool.Exec(context.Background(), "INSERT INTO users (username, info) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET info = $2", userInfo["email"], userInfo)
	return err
}
