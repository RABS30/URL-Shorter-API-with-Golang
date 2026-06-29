package middleware

import (
	"log"
	"net/http"
	"shorter-url/internal/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func GuestOnly(secretKey string) func(httprouter.Handle) httprouter.Handle {
	if secretKey == "" {
		log.Fatal("secret key not found")
	}
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			token, _ := r.Cookie("token")
			if token != nil {
				valid := TokenIsValid(secretKey, token.Value)
				if valid {
					helper.BadResponse(w, http.StatusBadRequest, "already authenticated")

					if wrapper, ok := w.(*ResponseWriterWrapper); ok {
						wrapper.WriteError("user already authenticated")
					}

					return
				}
			}

			next(w, r, p)
		}
	}
}

func TokenIsValid(secretKey string, tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return false
	} else {
		return token != nil && token.Valid
	}
}
