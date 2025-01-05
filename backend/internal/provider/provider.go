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
	UserInf
	Password *string `json:"password"`
}
type Test struct {
	Id             *int    `json:"id"`
	Input          *string `json:"input"`
	Correct_output *string `json:"correct_output"`
}

type TestResult struct {
	Test_id        *int    `json:"test_id"`
	Output         *string `json:"output"`
	Verdict        *string `json:"verdict"`
	Execution_time *int    `json:"execution_time"`
	Max_memory     *int    `json:"max_memory"`
}
type TestGroup struct {
	Id           *int    `json:"id"`
	Name         *string `json:"name"`
	Author       *string `json:"author"`
	Tests        []Test  `json:"tests"`
	Time_limit   *int    `json:"time_limit"`
	Memory_limit *int    `json:"memory_limit"`
	Kolvo        int     `json:"quantity_tests"`
}
type TestGroupResult struct {
	Id                 *int         `json:"id"`
	Group_id           *int         `json:"group_id"`
	Source_code        *string      `json:"code"`
	Language           *string      `json:"language"`
	Test_results       []TestResult `json:"test_results"`
	Verdict            *string      `json:"verdict"`
	Max_execution_time *int         `json:"max_exec_time"`
	Max_memory         *int         `json:"max_memory"`
	String_results     []string     `json:"results"`
}
type DatabaseProvider struct {
	db *sql.DB
}

func NewDatabaseProvider(host string, port int, user, password, dbName string) *DatabaseProvider {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return &DatabaseProvider{db: conn}
}

func (tg *TestGroup) CalcCol() int {
	k := len(tg.Tests)
	return k
}
