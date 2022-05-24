package controller

import (
	"RestSample/Codefiles/services/movieManagement/model"
	"RestSample/Codefiles/services/movieManagement/view"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//Handler Functions
func Routes() *http.ServeMux {
	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/addCategory", addCategory)
	multiplexer.HandleFunc("/addMovie", addMovie)
	multiplexer.HandleFunc("/movielist", movielist)
	return multiplexer
}

//Functions

func addCategory(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if r.Body != nil {
			category := view.MoviesCategory{}
			json.NewDecoder(r.Body).Decode(&category)
			//fmt.Printf("value of id %v, type %s\n", category.CategoryID, reflect.TypeOf(category.CategoryID))
			err := model.InsertMovieCategory(category)
			if err != nil {
				fmt.Print("Error Message ")
				fmt.Println(err)
				return
			}
			fmt.Println("Movie Category Created")
			w.WriteHeader(http.StatusCreated)

		} else {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Empty Body")
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Post Method Required")
	}

}

func addMovie(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if r.Body != nil {

			movie := view.Movies{}

			json.NewDecoder(r.Body).Decode(&movie)
			userID, ok := r.URL.Query()["UserID"]
			if !ok {
				fmt.Println("Query Param Missing")
				w.Write([]byte("Query Param Missing"))
				return
			}
			id, er := strconv.Atoi(userID[0])
			if er != nil {
				fmt.Println("UserID must be number")
				return
			}
			movie.UserID = id
			err := model.InsertMovies(movie)
			if err != nil {
				fmt.Print("Error Message ")
				fmt.Println(err)
				return
			}
			fmt.Println("Movie Category Created")
			w.WriteHeader(http.StatusCreated)

		} else {
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Empty Body")
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Post Method Required")
	}

}

func movielist(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		categoryID, ok := r.URL.Query()["categoryID"]
		userID, k := r.URL.Query()["userID"]

		if !ok || !k {
			fmt.Println("Query Param Missing")
			w.Write([]byte("Query Param Missing"))
			return
		}

		uID, a := strconv.Atoi(userID[0])
		catID, b := strconv.Atoi(categoryID[0])

		if a != nil || b != nil {
			fmt.Println("Wrong user id or category id")
			w.Write([]byte("Query Param Wrong"))
			return
		}
		fmt.Println("Debug")
		fmt.Println(uID)
		fmt.Println(catID)
		movies, err := model.FindAll(uID, catID)

		if err != nil {
			fmt.Println(err)
			return
		}

		//w.Write([]byte("This is the result from DB"))
		json.NewEncoder(w).Encode(&movies)

	}

}
