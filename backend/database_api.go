package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Conn_inf struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}
type User struct {
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type Session struct {
	Astiay_isos string
	User_login  string
	Creation    string
	Expire      string
}

type DatabaseProvider struct {
	db *sql.DB
}

func (dp *DatabaseProvider) CreateUser(user *User) (er error) {
	_, err := dp.db.Exec(fmt.Sprintf("INSERT INTO user_inf(login,nickname,password,avatar) values('%s','%s','%s','%s');",
		user.Login, user.Nickname, user.Password, user.Avatar))
	return err
}

func (dp *DatabaseProvider) CreateSession(session *Session) {
	_, err := dp.db.Exec(fmt.Sprintf("INSERT INTO sessions(astiay_isos,user_login) values('%s','%s');",
		session.Astiay_isos, session.User_login))
	if err != nil {
		log.Fatal(err)
	}
}

func Gen_coockie(login string, dp *DatabaseProvider) (coockie string) {
	row, err := dp.db.Exec("select count(*) from sessions;")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(row)
	return login
}
