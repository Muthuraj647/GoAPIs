package main

import (
	"RestSample/Codefiles/services/userservice/controller"
	"RestSample/Codefiles/services/userservice/model"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//connect to db

	db := model.RegisterWithDB()
	defer db.Close()

	multiplexer := controller.Routes()

	fmt.Println("User service ready to Serve...")

	http.ListenAndServe(":8080", multiplexer)
}
