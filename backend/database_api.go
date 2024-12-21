package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func hash(value string) string {
	h := sha256.New()
	h.Write([]byte(value))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type Conn_inf struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}
type User struct {
	Login    *string `json:"login"`
	Nickname *string `json:"nickname"`
	Password *string `json:"password"`
	Avatar   *string `json:"avatar"`
}

type Session struct {
	Astiay_isos *string `json:"Astiay_isos"`
	User_login  *string `json:"User_login"`
	Creation    *string `json:"Creation"`
	Expire      *string `json:"Expire"`
}
type test struct {
	id             *int
	input          *string
	correct_output *string
	author         *string
}

type test_result struct {
	test_id        *int
	output         *string
	verdict        *string
	execution_time *int
	max_memory     *int
}
type test_group struct {
	id           *int
	name         *string
	author       *string
	tests        *[]test
	time_limit   *int
	memory_limit *int
}

type test_group_result struct {
	group_id           *int
	source_code        *string
	language           *string
	test_results       *[]test_result
	verdict            *string
	max_execution_time *int
	max_memory         *int
}
type DatabaseProvider struct {
	db *sql.DB
}

func (dp *DatabaseProvider) CreateUser(user *User) (er error) {
	pass := hash(*(user.Password))

	_, err := dp.db.Exec(`INSERT INTO user_inf(login,nickname,password,avatar) values($1,$2,$3,$4);`,
		*(user.Login), *(user.Nickname), pass, *(user.Avatar))
	return err
}

func (dp *DatabaseProvider) CreateSession(session *Session) (er error) {
	_, err := dp.db.Exec(`INSERT INTO sessions(astiay_isos,user_login) values($1,$2);`,
		*(session.Astiay_isos), *(session.User_login))
	return err
}
func (dp *DatabaseProvider) Is_In_Base(login string, password string) (status bool) {
	var log string
	var pass string
	row := dp.db.QueryRow(`SELECT login, password from user_inf where login = $1 and password = $2;`, login, hash(password))
	if err := row.Scan(&log, &pass); err != nil || log != login || pass != hash(password) {
		fmt.Println(hash(password))
		fmt.Println(pass)
		return false
	}

	return true
}
func (dp *DatabaseProvider) Gen_coockie(login string) (coockie string) {
	row := dp.db.QueryRow("select count(*) from sessions;")
	var k string
	err := row.Scan(&k)
	if err != nil {
		log.Fatal(err)
	}
	return login + k
}

func (dp *DatabaseProvider) Del_session(cookie string) (er error) {
	var cooki string
	row := dp.db.QueryRow(`SELECT Astiay_isos from sessions where Astiay_isos = $1;`, cookie)
	if err := row.Scan(&cooki); err != nil || cookie != cookie {
		return err
	}
	_, err := dp.db.Exec(`DELETE FROM sessions WHERE Astiay_isos = $1;`, cookie)

	return err
}

func (dp *DatabaseProvider) Del_All_Sessions() (err error) {
	_, err = dp.db.Exec("truncate sessions;")
	return
}

func (dp *DatabaseProvider) Get_User(cookie string) (user *User, err error) {
	row := dp.db.QueryRow(`Select login, avatar, nickname from user_inf where login = (select user_login from sessions where Astiay_isos = $1);`, cookie)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var l, a, n string
	if err := row.Scan(&l, &a, &n); err != nil {
		return nil, err
	}
	u := User{Login: new(string), Nickname: new(string), Password: new(string), Avatar: new(string)}

	*(u.Login) = l
	*(u.Avatar) = a
	*(u.Nickname) = n
	return &u, nil

}

func (dp *DatabaseProvider) Change_user_avatar(login string, avatar string) (err error) {
	_, err = dp.db.Exec(`update user_inf set avatar = $1 where login = $2`, avatar, login)
	if err != nil {
		return
	}
	return nil
}
func (dp *DatabaseProvider) IsSessionActive(cookie string) bool {
	row := dp.db.QueryRow(`select astiay_isos from sessions where astiay_isos = $1`, cookie)
	var r string
	if err := row.Scan(&r); err != nil {
		return false
	}
	if r != cookie {
		return false
	}
	return true
}

func (dp *DatabaseProvider) ChangeUserPassword(login string, password string) error {
	pass := hash(password)
	_, err := dp.db.Exec(`update user_inf set password = $1 where login = $2`, pass, login)
	if err != nil {
		return err
	}
	return nil
}
func (dp *DatabaseProvider) ChangeUserNick(login string, nickname string) error {
	_, err := dp.db.Exec(`update user_inf set nickname = $1 where login = $2`, nickname, login)
	if err != nil {
		return err
	}
	return nil
}
