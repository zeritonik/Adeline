package main

import (
	"database/sql"

	"fmt"
	"log"
	"net/http"
)

const (
	host     = "46.19.67.80"
	port     = 5432
	user     = "maks"
	password = "YaEmo123!"
	dbname   = "adeline"
)

func main() {

	connstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	dp := DatabaseProvider{db: db}
	defer db.Close()

	h := BD_handlers{dp: dp}
	pr := newResolver()
	pr.Addpath("*/home", homepage)
	pr.Addpath("*/profile", profilepage)
	pr.Addpath("POST */register", h.Create_user)
	pr.Addpath("POST */login", h.Login_user)
	pr.Addpath("*/logout", h.Delete_Session)
	pr.Addpath("*/profile/setting", h.Settings)
	pr.Addpath("GET */*/*/*", HTML1)
	pr.Addpath("GET */*/*", HTML1)
	pr.Addpath("GET */", HTML)
	pr.Addpath("GET */register", HTML)
	pr.Addpath("GET */login", HTML)
	var port string = ":8080"

	fmt.Println("Server is listening...")

	fmt.Println("http://localhost" + port)
	err = http.ListenAndServe(port, pr)
	if err != nil {
		fmt.Println("Ошибка запуска сервера")
	}

}
