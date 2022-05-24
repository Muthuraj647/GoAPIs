package main

import (
	"RestSample/Codefiles/services/movieManagement/controller"
	"RestSample/Codefiles/services/movieManagement/model"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//connect to db

	db := model.RegisterWithDB()
	defer db.Close()

	multiplexer := controller.Routes()

	fmt.Println("Movie service ready to Serve...")

	http.ListenAndServe(":8081", multiplexer)
}
