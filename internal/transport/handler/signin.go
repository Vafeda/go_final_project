package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Vafeda/go_final_project/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func signIn(w http.ResponseWriter, r *http.Request) {
	password := models.Password{}

	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		encodeResponse(w, &models.ErrorResponse{Error: err.Error()})
		return
	}

	todoPassword, exist := os.LookupEnv("TODO_PASSWORD")
	if !exist {
		return
	}
	fmt.Println(password)
	if password.Password != todoPassword {
		return
	}

	key := "anybody"
	t := jwt.New(jwt.SigningMethodHS256)
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return
	}
	fmt.Println(s)
	encodeResponse(w, models.JSONToken{Token: s})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var Jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				Jwt = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			// ...
			s := []byte("anybody")

			jwtToken, err := jwt.Parse(Jwt, func(t *jwt.Token) (interface{}, error) {
				// секретный ключ для всех токенов одинаковый, поэтому просто возвращаем его
				return s, nil
			})

			valid = jwtToken.Valid
			fmt.Println("Аунтификация епта")
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
