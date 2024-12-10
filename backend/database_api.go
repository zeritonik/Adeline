package main

import (
	"database/sql"
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
	_, err := dp.db.Exec(`INSERT INTO user_inf(login,nickname,password,avatar) values($1,$2,$3,$4);`,
		*(user.Login), *(user.Nickname), *(user.Password), *(user.Avatar))
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
	row := dp.db.QueryRow(`SELECT login from user_inf where login = $1;`, login)
	if err := row.Scan(&log); err != nil || log != login {
		return false
	}
	row = dp.db.QueryRow(`SELECT password from user_inf where password = $1;`, password)

	if err := row.Scan(&pass); err != nil || pass != password {
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

func (dp *DatabaseProvider) Get_User(cookie string) (us *User, err error) {
	row := dp.db.QueryRow(`Select login, avatar, nickname from user_inf where login = (select user_login from sessions where Astiay_isos = $1);`, cookie)
	if row.Err() != nil {
		return nil, row.Err()
	}
	user := User{}
	var l, a, n string
	if err := row.Scan(&l, &a, &n); err != nil {
		return nil, err
	}
	user.Login = &l
	user.Avatar = &a
	user.Nickname = &n

	return &user, nil

}

func (dp *DatabaseProvider) Change_user_avatar(login string, avatar string) (err error) {
	_, err = dp.db.Exec("update user_inf set avatar = $1 where login = $2", avatar, login)
	if err != nil {
		return
	}
	return nil
}
