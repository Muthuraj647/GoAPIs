package model

import (
	"RestSample/Codefiles/services/movieManagement/view"
	"database/sql"
	"fmt"
)

//global
var conn *sql.DB

//Connect with DB

func RegisterWithDB() *sql.DB {
	Driver := "mysql"
	User := "root"
	Password := "Admin@123"
	DBname := "Microservices"
	db, err := sql.Open(Driver, User+":"+Password+"@/"+DBname)
	if err != nil {
		fmt.Println("Error Occured when try to connect with DB")
		return nil
	}
	fmt.Println("Connected with database " + DBname)
	conn = db

	return db

}

//handler functions

func InsertMovieCategory(category view.MoviesCategory) error {

	insertQuery, err := conn.Query("INSERT INTO MovieCategory(CategoryID,CategoryName) VALUES(?,?)", category.CategoryID, category.CategoryName)

	if err != nil {
		fmt.Println("Error Occured when creating Users")
		fmt.Println(err)
		return err
	}
	defer insertQuery.Close()
	return nil
}

func InsertMovies(movie view.Movies) error {

	insertQuery, err := conn.Query("INSERT INTO Movies VALUES(?,?,?,?,?)", movie.MoviesID, movie.MovieName, movie.CategoryID, movie.URL, movie.UserID)

	if err != nil {
		fmt.Println("Error Occured when creating Users")
		fmt.Println(err)
		return err
	}
	defer insertQuery.Close()
	return nil
}

func FindAll(userID, categoryID int) ([]view.Movies, error) {
	findQuery, err := conn.Query("SELECT * FROM Movies WHERE CategoryID=? && UserID=?", categoryID, userID)

	var movies = []view.Movies{}
	if err != nil {
		fmt.Println("Error Occured when creating Users")
		fmt.Println(err)
		return movies, err
	}

	for findQuery.Next() {
		movie := view.Movies{}
		findQuery.Scan(&movie.MoviesID, &movie.MovieName, &movie.CategoryID, &movie.URL, &movie.UserID)
		fmt.Println(movie.MovieName)
		movies = append(movies, movie)
	}
	defer findQuery.Close()
	return movies, nil
}
