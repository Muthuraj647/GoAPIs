package controller

import (
	"RestSample/Codefiles/services/userservice/model"
	"RestSample/Codefiles/services/userservice/vault"
	"RestSample/Codefiles/services/userservice/view"
	"RestSample/dgrijalva/jwt-go"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//Handler Functions
func Routes() *http.ServeMux {
	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/home", home)
	multiplexer.HandleFunc("/signin", createUser)
	multiplexer.HandleFunc("/login", login)
	multiplexer.HandleFunc("/addCategory", addCategory)
	multiplexer.HandleFunc("/addMovie", addMovie)
	multiplexer.HandleFunc("/movielist", movielist)
	multiplexer.HandleFunc("/changePassword", changePwd)
	multiplexer.HandleFunc("/logout", logout)

	return multiplexer
}

//Functions

var users = view.Users{}

var jwtKey = []byte("security")

func createUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if r.Body != nil {

			data := users
			json.NewDecoder(r.Body).Decode(&data)

			//talk to db

			err := model.InsertUsers(data)
			if err != nil {
				fmt.Print("Error Message ")
				fmt.Println(err)
				return
			}
			fmt.Println("User Created\nPlease Login to Continue...")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("User Created\nPlease Login to Continue...\n"))

		} else {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Empty Body")
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Post Method Required")
	}

}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		if r.Body != nil {
			data := users
			json.NewDecoder(r.Body).Decode(&data)
			user, err := model.CheckUser(data)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Wrong Password or Username"))
				fmt.Println(err)
				return
			}

			//create jwt token

			expirationTime := time.Now().Add(time.Minute * 5)
			claim := view.Claims{
				Name:  user.Name,
				Email: user.Email,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

			tokenStr, err := token.SignedString(jwtKey)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Internal Server Error")
				fmt.Println(err)
				return
			}

			//create cookies with token

			http.SetCookie(w, &http.Cookie{
				Name:    "JWTToken",
				Value:   tokenStr,
				Expires: expirationTime,
			})

			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Login Successfull!!!"))
			w.Write([]byte(fmt.Sprintf("Welcome %s\n", user.Name)))
			fmt.Println("Login success")

		} else {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Empty Body")
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Post Method Required")
	}

}

func home(w http.ResponseWriter, r *http.Request) {

	claims, err := vault.TokenValidation(w, r)
	if err == nil {
		cl := claims.(*view.Claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Welcome %s! You're authoraized Person..", cl.Name)))

		return
	}
	w.Write([]byte("You're not authoraized Person.."))
	fmt.Println(err)

}

func addCategory(w http.ResponseWriter, r *http.Request) {
	claims, err := vault.TokenValidation(w, r)
	if err == nil {
		cl := claims.(*view.Claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Welcome %s! You're authoraized Person..", cl.Name)))
		http.Post("http://localhost:8081/addCategory", "JSON", r.Body)
		return
	}
	w.Write([]byte("You're not authoraized Person.."))
	fmt.Println(err)

}

func addMovie(w http.ResponseWriter, r *http.Request) {

	claims, err := vault.TokenValidation(w, r)
	if err == nil {
		cl := claims.(*view.Claims)
		w.WriteHeader(http.StatusOK)
		userID, ok := r.URL.Query()["UserID"]
		if !ok {
			fmt.Println("Query Param Missing")
			w.Write([]byte("Query Param Missing"))
			return
		}

		w.Write([]byte(fmt.Sprintf("Welcome %s! You're authoraized Person..", cl.Name)))
		http.Post("http://localhost:8081/addMovie?UserID="+userID[0], "JSON", r.Body)
		return
	}
	w.Write([]byte("You're not authoraized Person.."))
	fmt.Println(err)

}
func movielist(w http.ResponseWriter, r *http.Request) {

	claims, err := vault.TokenValidation(w, r)
	if err == nil {
		categoryID, ok := r.URL.Query()["categoryID"]
		userID, k := r.URL.Query()["userID"]

		if !ok || !k {
			fmt.Println("Query Param Missing")
			w.Write([]byte("Query Param Missing"))
			return
		}

		cl := claims.(*view.Claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Welcome %s! You're authoraized Person..", cl.Name)))
		http.Get("http://localhost:8081/movielist?userID=" + userID[0] + "&categoryID=" + categoryID[0])
		return
	}
	w.Write([]byte("You're not authoraized Person.."))
	fmt.Println(err)

}

func changePwd(w http.ResponseWriter, r *http.Request) {

	claims, err := vault.TokenValidation(w, r)
	if err == nil {
		cl := claims.(*view.Claims)
		user := users
		json.NewDecoder(r.Body).Decode(&user)
		user.Email = cl.Email
		er := model.ChangePassword(user)
		if er != nil {
			fmt.Println("error occured")
			fmt.Println(er)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(fmt.Sprintf("%s! Your Password changed..", cl.Name)))
		//logout(w, r)
		return
	}
	w.Write([]byte("You're not authoraized Person.."))
	fmt.Println(err)

}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("JWTToken")
	if cookie != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "JWTToken",
			Value:  "",
			MaxAge: -1,
		})
		fmt.Println("cookie " + cookie.Value)
		w.Write([]byte("Logged out"))
	}
	if err != nil {
		w.Write([]byte("No cookie found so you already logged out!"))
	}
}
