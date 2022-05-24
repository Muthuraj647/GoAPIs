package vault

import (
	"RestSample/Codefiles/services/userservice/view"
	"RestSample/dgrijalva/jwt-go"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("security")

// encrypt and compare passwords, for simplicity added here

func EncryptingPassword(password []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		fmt.Println("Error when Encrypting")
		fmt.Println(err)
		return ""
	}
	return string(hashed)
}

//validating

func ValidatethePassword(original, given string) error {
	err := bcrypt.CompareHashAndPassword([]byte(original), []byte(given))
	if err != nil {
		fmt.Println("Wrong Password")
		return err
	}
	return nil
}

//token Validation

func TokenValidation(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	cookies, err := r.Cookie("JWTToken")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied! need to login before access this page"))
			return nil, err
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cookie Missing or Invalid"))
		return nil, err
	}

	tokenstr := cookies.Value
	claims := &view.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenstr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token or token expires,Invalid Signature"))
			return nil, err
		}
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte("Error Occured"))
		return nil, err
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token or token expires"))
		return nil, err
	}

	return claims, nil
}
