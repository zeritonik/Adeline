package provider

import (
	"database/sql"
	"fmt"
	"log"
)

type UserInf struct {
	Login    *string `json:"login"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}
type ConnInf struct {
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
type Test struct {
	Id             *int    `json:"id"`
	Input          *string `json:"input"`
	Correct_output *string `json:"correct_output"`
}

type TestResult struct {
	test_id        *int
	output         *string
	verdict        *string
	execution_time *int
	max_memory     *int
}
type TestGroup struct {
	Id           *int    `json:"id"`
	Name         *string `json:"name"`
	Author       *string `json:"author"`
	Tests        []Test  `json:"tests"`
	Time_limit   *int    `json:"time_limit"`
	Memory_limit *int    `json:"memory_limit"`
}

type TestGroupResult struct {
	group_id           *int
	source_code        *string
	language           *string
	test_results       *[]TestResult
	verdict            *string
	max_execution_time *int
	max_memory         *int
}
type DatabaseProvider struct {
	db *sql.DB
}

func NewDatabaseProvider(host string, port int, user, password, dbName string) *DatabaseProvider {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Создание соединения с сервером postgres
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return &DatabaseProvider{db: conn}
}
