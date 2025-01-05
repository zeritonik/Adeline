package provider

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func (dp *DatabaseProvider) CreateUser(user User) error {
	_, err := dp.db.Exec(`INSERT INTO user_inf(login,nickname,password,avatar) values($1,$2,$3,null);`,
		*(user.Login), *(user.Nickname), *(user.Password))
	return err
}

func (dp *DatabaseProvider) CreateSession(cookie string, login string) error {
	_, err := dp.db.Exec(`INSERT INTO sessions(astiay_isos,user_login) values($1,$2);`, cookie, login)
	return err
}
func (dp *DatabaseProvider) IsInBase(login string, password string) (bool, error) {
	row := dp.db.QueryRow(`select exists (select login from user_inf where login = $1 and password = $2)`, login, password)
	var k string
	if err := row.Scan(&k); err != nil {
		return false, err
	}
	if k == "true" {
		return true, nil
	}
	return false, nil
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

func (dp *DatabaseProvider) DelSession(cookie string) (er error) {
	// var cooki string
	// row := dp.db.QueryRow(`SELECT astiay_isos from sessions where astiay_isos = $1;`, cookie)
	// if err := row.Scan(&cooki); err != nil {
	// 	return err
	// }
	_, err := dp.db.Exec(`DELETE FROM sessions WHERE astiay_isos = $1;`, cookie)
	return err
}

func (dp *DatabaseProvider) DelAllSessions(login string) (err error) {
	_, err = dp.db.Exec("DELETE FROM sessions WHERE user_login = $1;", login)
	return err
}

func (dp *DatabaseProvider) GetUserInf(cookie string) (user *UserInf, err error) {
	row := dp.db.QueryRow(`Select login, avatar, nickname from user_inf where login = (select user_login from sessions where astiay_isos = $1);`, cookie)
	u := UserInf{
		Login:    new(string),
		Nickname: new(string),
		Avatar:   new(string)}
	if err := row.Scan(u.Login, u.Avatar, u.Nickname); err != nil {
		return nil, err
	}
	return &u, nil

}

func (dp *DatabaseProvider) ChangeUserAvatar(login string, avatar string) error {
	_, err := dp.db.Exec(`update user_inf set avatar = $1 where login = $2`, avatar, login)
	return err
}
func (dp *DatabaseProvider) GetSession(cookie string) (bool, error) {
	row := dp.db.QueryRow(`select exists( select astiay_isos from sessions where astiay_isos = $1)`, cookie)
	var r string
	if err := row.Scan(&r); err != nil {
		return false, err
	}
	if r == "true" {
		return true, nil
	}
	return false, nil
}

func (dp *DatabaseProvider) ChangeUserPassword(login string, password string) error {
	_, err := dp.db.Exec(`update user_inf set password = $1 where login = $2`, password, login)
	return err
}
func (dp *DatabaseProvider) ChangeUserNick(login string, nickname string) error {
	_, err := dp.db.Exec(`update user_inf set nickname = $1 where login = $2`, nickname, login)
	return err
}
func (dp *DatabaseProvider) ChangeUserLogin(new string, old string) error {
	_, err := dp.db.Exec(`update user_inf set login = $1 where login = $2`, new, old)
	return err
}

func (dp *DatabaseProvider) InsertTestGroup(tg TestGroup) (error, map[Test]error, int) {
	var id int
	row := dp.db.QueryRow(`insert into test_group (name,author,time_limit,memory_limit) values ($1,$2,$3,$4) returning id;`, tg.Name, tg.Author, tg.Time_limit, tg.Memory_limit)
	if err := row.Scan(&id); err != nil {
		return err, nil, 0
	}

	e := make(map[Test]error)
	for _, val := range tg.Tests {
		if _, err := dp.db.Exec(`update test_group set tests = array_append(tests,($1,$2,$3)::test) where id = $4 `, val.Id, val.Input, val.Correct_output, id); err != nil {
			e[val] = err
		} else {
			e[val] = nil
		}
	}
	return nil, e, id
}

func (dp *DatabaseProvider) GetTestGroupInfo(id int, login string) (*TestGroup, error) {
	tg := TestGroup{Id: new(int), Name: new(string), Time_limit: new(int), Memory_limit: new(int), Author: new(string), Tests: *new([]Test)}
	var r []string
	row := dp.db.QueryRow(`select id,name,author,time_limit,memory_limit,tests from test_group where id = $1 and author = $2;`, id, login)
	if err := row.Scan(tg.Id, tg.Name, tg.Author, tg.Time_limit, tg.Memory_limit, pq.Array(&r)); err != nil {
		return nil, err
	}
	for i, val := range r {
		rez := strings.Split(val, ",")
		tg.Tests = append(tg.Tests, Test{Id: new(int), Input: new(string), Correct_output: new(string)})
		id, _ = (strconv.Atoi(strings.Trim(rez[0], "(")))
		tg.Tests[i].Id = &id
		input := strings.Trim(rez[1], "\"")
		tg.Tests[i].Input = &input
		output := strings.Trim(strings.Trim(rez[2], ")"), "\"")
		tg.Tests[i].Correct_output = &output
	}
	return &tg, nil
}

func (dp *DatabaseProvider) GetTestGroupInfoLOGIN(login string) ([]TestGroup, error) {
	var arrtg []TestGroup = make([]TestGroup, 0)
	var r []string
	row, err := dp.db.Query(`select id,name,author,time_limit,memory_limit,tests from test_group where author = $1;`, login)
	if err != nil {
		return arrtg, err
	}
	for row.Next() {
		tg := TestGroup{Id: new(int), Name: new(string), Time_limit: new(int), Memory_limit: new(int), Author: new(string), Tests: *new([]Test)}
		if err := row.Scan(tg.Id, tg.Name, tg.Author, tg.Time_limit, tg.Memory_limit, pq.Array(&r)); err != nil {
			return arrtg, err
		}
		for i, val := range r {
			rez := strings.Split(val, ",")
			tg.Tests = append(tg.Tests, Test{Id: new(int), Input: new(string), Correct_output: new(string)})
			id, _ := (strconv.Atoi(strings.Trim(rez[0], "(")))
			tg.Tests[i].Id = &id
			input := strings.Trim(rez[1], "\"")
			tg.Tests[i].Input = &input
			output := strings.Trim(strings.Trim(rez[2], ")"), "\"")
			tg.Tests[i].Correct_output = &output
		}
		arrtg = append(arrtg, tg)
	}

	return arrtg, nil
}

// func (dp *DatabaseProvider) InsertTestGroupRezult(tg TestGroupResult) (error, map[TestResult]error) {
// 	var id int
// 	row := dp.db.QueryRow(`insert into test_group_result (group_id,source_code,language,verdict,max_execution_time,max_memory) values ($1,$2,$3,$4,$5,$6) returning id;`, tg.Group_id, tg.Source_code, tg.Language, tg.Verdict, tg.Max_execution_time, tg.Max_memory)
// 	if err := row.Scan(&id); err != nil {
// 		return err, nil
// 	}

// 	e := make(map[TestResult]error)
// 	for _, val := range *tg.Test_results {
// 		if _, err := dp.db.Exec(`update test_group set tests = array_append(tests,($1,$2,$3,$4,$5)::test_result) where id = $6 `, val.Test_id, val.Output, val.Execution_time, val.Verdict, val.Max_memory, id); err != nil {
// 			e[val] = err
// 		} else {
// 			e[val] = nil
// 		}
// 	}
// 	return nil, e
// }

func (dp *DatabaseProvider) DeleteTestGroup(id int, login string) error {
	if _, err := dp.db.Exec(`delete from test_group where author = $1 and id = $2`, login, id); err != nil {
		return err
	}
	return nil
}

func (dp *DatabaseProvider) GetTestGroupResultInfo(login string) ([]TestGroupResult, error) {
	tgr := make([]TestGroupResult, 0)
	rows, err := dp.db.Query(`select test_group_result.id,group_id,language,max_execution_time, max_memory,source_code,results, verdict from test_group_result inner join test_group on group_id = test_group.id where author = $1;`, login)
	if err != nil {
		return tgr, err
	}
	for rows.Next() {
		tg := TestGroupResult{Id: new(int), Group_id: new(int), Source_code: new(string), Language: new(string), Verdict: new(string), Max_execution_time: new(int), Max_memory: new(int), Test_results: make([]TestResult, 0), String_results: make([]string, 0)}
		rows.Scan(tg.Id, tg.Group_id, tg.Language, tg.Max_execution_time, tg.Max_memory, tg.Source_code, tg.String_results, tg.Verdict)
		fmt.Println(tg.String_results)

	}
	return tgr, nil
}
