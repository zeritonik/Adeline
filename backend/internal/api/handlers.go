package api

import (
	"adeline/backend/internal/provider"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (srv *Server) PostCreateUser(c echo.Context) error {
	user := provider.User{Login: new(string), Password: new(string), Avatar: new(string), Nickname: new(string)}
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(user.Avatar, user.Login, user.Nickname, user.Password)
	if err := srv.uc.CreateUser(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rez := struct {
		Ans string `json:"ans"`
	}{Ans: "Registered"}
	return c.JSON(http.StatusOK, rez)
}

func (srv *Server) PostLogin(c echo.Context) error {
	user := provider.User{Login: new(string), Password: new(string)}
	cook, err := c.Cookie("astiay_isos")
	rez := struct {
		Ans string `json:"ans"`
	}{}
	if err == nil {
		check, err := srv.uc.CheckSession(cook.Value)
		if err != nil {
			rez.Ans = err.Error()
			return c.JSON(http.StatusInternalServerError, rez)
		}
		if check {
			return c.JSON(http.StatusOK, user)
		}
	}

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err != nil && user.Login == nil && user.Password == nil {
		return c.JSON(401, "No_user")
	}
	cooc, err := srv.uc.LoginUser(*user.Login, *user.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	cookie := http.Cookie{
		MaxAge:   3600,
		Value:    cooc,
		Name:     "astiay_isos",
		HttpOnly: false,
	}
	c.SetCookie(&cookie)
	return c.JSON(http.StatusOK, user)
}

func (srv *Server) GetSettings(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
func (srv *Server) PostLogout(c echo.Context) error {
	s := struct {
		All *bool `json:"all"`
	}{}
	rez := struct {
		Ans string `json:"ans"`
	}{}
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		rez.Ans = err.Error()
		return c.JSON(http.StatusInternalServerError, rez)
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		rez.Ans = err.Error()
		return c.JSON(http.StatusInternalServerError, rez)
	}
	if err := c.Bind(&s); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err := srv.uc.DelSession(*user.Login, cookie.Value, *s.All); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rez.Ans = "DELETED"
	return c.JSON(http.StatusOK, rez)
}

func (srv *Server) PostSettings(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user := provider.User{}
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err := srv.uc.ChangeSettings(user.Login, user.Password, user.Nickname, user.Avatar, cookie.Value); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "")
}

func (srv *Server) PostTests(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tg := provider.TestGroup{Author: new(string)}
	if err := c.Bind(&tg); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	*(tg.Author) = *(user.Login)
	if err := srv.uc.AddTestGroup(tg); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rez := struct {
		Id           *int    `json:"id"`
		Name         *string `json:"name"`
		Author       *string `json:"author"`
		Kolvo        *int    `json:"quantityOfTests"`
		Time_limit   *int    `json:"time_limit"`
		Memory_limit *int    `json:"memory_limit"`
	}{Id: new(int), Name: new(string), Author: new(string), Kolvo: new(int), Time_limit: new(int), Memory_limit: new(int)}
	*(rez.Author) = *(tg.Author)
	*(rez.Memory_limit) = *(tg.Memory_limit)
	*(rez.Kolvo) = (len(tg.Tests))
	*(rez.Time_limit) = *(tg.Time_limit)
	*(rez.Name) = *(tg.Name)

	return c.JSON(http.StatusCreated, rez)
}

func (srv *Server) GetTests(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	id := c.Param("id")
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if id != "" {
		var tg *provider.TestGroup
		Id, _ := strconv.Atoi(id)
		tg, err = srv.uc.GetTestGroup(Id, *user.Login)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if tg == nil {
			return c.String(http.StatusBadRequest, "No data to this user")
		}

		return c.JSON(http.StatusOK, tg)
	}
	tgs, err := srv.uc.GetTestGroupwithLogin(*user.Login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var a []struct {
		Id           *int    `json:"id"`
		Memory_limit *int    `json:"memory_limit"`
		Time_limit   *int    `json:"time_limit"`
		Kolvo        *int    `json:"quantity_tests"`
		Author       *string `json:"author"`
		Name         *string `json:"name"`
	} = make([]struct {
		Id           *int    `json:"id"`
		Memory_limit *int    `json:"memory_limit"`
		Time_limit   *int    `json:"time_limit"`
		Kolvo        *int    `json:"quantity_tests"`
		Author       *string `json:"author"`
		Name         *string `json:"name"`
	}, 0)
	for _, val := range tgs {
		u := struct {
			Id           *int    `json:"id"`
			Memory_limit *int    `json:"memory_limit"`
			Time_limit   *int    `json:"time_limit"`
			Kolvo        *int    `json:"quantity_tests"`
			Author       *string `json:"author"`
			Name         *string `json:"name"`
		}{Memory_limit: new(int), Time_limit: new(int), Kolvo: new(int), Author: new(string), Name: new(string), Id: new(int)}
		*(u.Memory_limit) = *(val.Memory_limit)
		*(u.Time_limit) = *(val.Time_limit)
		*(u.Author) = *(val.Author)
		*(u.Kolvo) = len(val.Tests)
		*(u.Name) = *(val.Name)
		*(u.Id) = *(val.Id)
		a = append(a, u)
	}
	return c.JSON(http.StatusOK, a)

}

func (srv *Server) GetTestGroupRez(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	inf := struct {
		Code        *string `json:"source_code"`
		Language    *string `json:"language"`
		TestGroupId *int    `json:"test_group_id"`
	}{}
	if err := c.Bind(&inf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(*inf.Code, *inf.Language, *inf.TestGroupId)
	test, err := srv.uc.GetTestGroup(*inf.TestGroupId, *user.Login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// var wg sync.WaitGroup
	// wg.Add(len(test.Tests))
	for _, val := range test.Tests {
		// go func() error {
		// defer wg.Done()
		output, _ := ExecutePython(*val.Id, *inf.Code, *val.Input)
		fmt.Println(output)
		// }()

	}
	// wg.Wait()
	return c.JSON(http.StatusOK, test)
}

func (srv *Server) DeleteGroup(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "No_id")
	}
	i, _ := strconv.Atoi(id)
	if err := srv.uc.DeleteTestGroup(*user.Login, i); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (srv *Server) GetResults(c echo.Context) error {
	cookie, err := c.Cookie("astiay_isos")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f, err := srv.uc.CheckSession(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !f {
		return c.String(401, "")
	}
	user, err := srv.uc.GetUser(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tgr, err := srv.uc.GetTestGroupResult(*user.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tgr)

}
func ExecutePython(id int, code string, input string) (string, error) {
	cmd := exec.Command("python3", "backend/tests/prog"+strconv.Itoa(id)+".py")
	var out bytes.Buffer
	cmd.Stdin = bytes.NewBufferString(input)
	cmd.Stdout = &out

	f, _ := os.Create("backend/tests/prog" + strconv.Itoa(id) + ".py")
	f.WriteString(code)
	f.Close()
	if err := cmd.Run(); err != nil {
		return "", err
	}

	output := out.String()
	return output, nil
}
