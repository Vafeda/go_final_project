package handler

import (
	"net/http"
	"os"

	"github.com/Vafeda/go_final_project/internal/models"
	"github.com/golang-jwt/jwt/v5"
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
		todoPassword, exist := os.LookupEnv(EnvTodoPassword)
		if !exist {
			encode(w, http.StatusInternalServerError, models.ErrorResponse{
				Error: "Server configuration error",
			})
			return
		}

		if len(todoPassword) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login.html", http.StatusFound)
				return
			}

			jwtToken, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})

			if !jwtToken.Valid {
				http.Redirect(w, r, "/login.html", http.StatusFound)
				return
			}
		}
		next(w, r)
	})
}
