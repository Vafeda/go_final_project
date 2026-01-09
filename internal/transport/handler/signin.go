package handler

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"

	"github.com/Vafeda/TODO-List/internal/models"
)

const (
	EnvTodoPassword = "TODO_PASSWORD"

	jwtKey = "golang"
)

func signIn(w http.ResponseWriter, r *http.Request) {
	pass, err := decode[models.Password](r)
	if err != nil {
		encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	todoPassword, exist := os.LookupEnv(EnvTodoPassword)
	if !exist {
		encode(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Server configuration error: TODO_PASSWORD not set",
		})
		return
	}

	if pass.Password != todoPassword {
		encode(w, http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid password",
		})
		return
	}

	t := jwt.New(jwt.SigningMethodHS256)
	s, err := t.SignedString([]byte(jwtKey))
	if err != nil {
		encode(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	encode(w, http.StatusOK, models.JSONToken{Token: s})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login.html" {
			fmt.Println(r.URL.Path)
			next(w, r)
			return
		}

		todoPassword, exist := os.LookupEnv(EnvTodoPassword)
		if !exist || strings.TrimSpace(todoPassword) == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("Redirect 1")
			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
				strings.Contains(r.Header.Get("Accept"), "application/json") ||
				strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnauthorized)
				encode(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
				return
			}

			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		jwtToken, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if !jwtToken.Valid {
			fmt.Println("Redirect 2")
			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
				strings.Contains(r.Header.Get("Accept"), "application/json") ||
				strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnauthorized)
				encode(w, http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
				return
			}
			
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		next(w, r)
	})
}
