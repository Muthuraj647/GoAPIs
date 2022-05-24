package model

import (
	"RestSample/Codefiles/services/userservice/vault"
	"RestSample/Codefiles/services/userservice/view"
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
		fmt.Println(err)
		return nil
	}
	fmt.Println("Connected with database " + DBname)
	conn = db

	return db

}

//handler functions

func InsertUsers(data view.Users) error {

	//to Encrypt the Password
	encryptedPassword := vault.EncryptingPassword([]byte(data.Password))

	InsertQuery, err := conn.Query("INSERT INTO Users (Name,Email,Gender,Password) VALUES(?,?,?,?)", data.Name, data.Email, data.Gender, encryptedPassword)

	if err != nil {
		fmt.Println("Error Occured when creating Users")
		fmt.Println(err)
		return err
	}
	defer InsertQuery.Close()
	return nil

}

//login User

func CheckUser(data view.Users) (view.Users, error) {
	loginQuery, err := conn.Query("SELECT Name,Password FROM Users WHERE Email=?", data.Email)
	if err != nil {
		fmt.Println("Error occured when Querying DB")
		fmt.Print("Error Message ")
		fmt.Println(err)
		return data, err
	}

	expectedData := view.Users{}
	for loginQuery.Next() {
		loginQuery.Scan(&expectedData.Name, &expectedData.Password)
	}

	fmt.Println("to debug")

	fmt.Println("FromDB -> " + expectedData.Password)
	fmt.Println("FromDB -> " + expectedData.Name)
	fmt.Println("FromUser ->" + data.Password)

	validate := vault.ValidatethePassword(expectedData.Password, data.Password)

	if validate != nil {
		return data, validate
	}
	fmt.Println("Authorized")
	defer loginQuery.Close()
	data.Name = expectedData.Name
	return data, nil

}

func ChangePassword(usr view.Users) error {

	encryptedPassword := vault.EncryptingPassword([]byte(usr.Password))

	passwordQuery, err := conn.Query("UPDATE Users SET Password=? WHERE Email=?", encryptedPassword, usr.Email)

	if err != nil {
		fmt.Println("Error occured when Querying DB")
		fmt.Print("Error Message ")
		fmt.Println(err)
		return err
	}
	defer passwordQuery.Close()
	fmt.Println("Password Changed...")
	return nil
}
