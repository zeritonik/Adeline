package main

import (
	"database/sql"

	"fmt"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "maksim"
	password = "123"
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

	pr.Addpath("POST */api/register", h.Create_user)
	pr.Addpath("POST */api/login", h.Login_user)
	pr.Addpath("POST */api/logout", h.Delete_Session_Post)
	pr.Addpath("POST */api/profile/settings", h.Post_Settings)
	pr.Addpath("GET */api/profile/settings", h.Get_Settings)
	pr.Addpath("GET */static/*/*", get_Static)
	pr.Addpath("GET */static/*", get_Static)
	pr.Addpath("GET */profile/settings", get_HTML)
	pr.Addpath("GET */", get_HTML)
	pr.Addpath("GET */register", get_HTML)
	pr.Addpath("GET */login", get_HTML)
	pr.Addpath("GET */profile", get_HTML)
	pr.Addpath("GET */profile/settings", get_HTML)

	var port string = ":8080"

	fmt.Println("Server is listening...")

	fmt.Println("http://localhost" + port)
	err = http.ListenAndServe(port, pr)
	if err != nil {
		fmt.Println("Ошибка запуска сервера")
	}

}
