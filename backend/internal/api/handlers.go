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

	return c.String(http.StatusOK, "Registered")
}

func (srv *Server) PostLogin(c echo.Context) error {
	user := provider.User{Login: new(string), Password: new(string)}
	cook, err := c.Cookie("astiay_isos")
	if err == nil {
		check, err := srv.uc.CheckSession(cook.Value)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if check == true {
			return c.String(http.StatusOK, "Autorized")
		}
	}
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
	return c.String(http.StatusOK, "Autorized")
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
		All     *bool
		Cookies []string
	}{}
	if err := c.Bind(&s); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err := srv.uc.DelSession(s.Cookies, *s.All); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "")
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
	tg := provider.TestGroup{}
	if err := c.Bind(&tg); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	*(tg.Author) = *(user.Login)
	if err := srv.uc.AddTestGroup(tg); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, tg)
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
	id := c.QueryParam("group-id")
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
		u := struct {
			Id           *int    `json:"int"`
			Memory_limit *int    `json:"memory_limit"`
			Time_limit   *int    `json:"time_limit"`
			Kolvo        *int    `json:"quantity_tests"`
			Author       *string `json:"author"`
			Name         *string `json:"name"`
		}{Memory_limit: new(int), Time_limit: new(int), Kolvo: new(int), Author: new(string), Name: new(string), Id: new(int)}
		*(u.Memory_limit) = *(tg.Memory_limit)
		*(u.Time_limit) = *(tg.Time_limit)
		*(u.Author) = *(tg.Author)
		*(u.Kolvo) = len(tg.Tests)
		*(u.Name) = *(tg.Name)
		*(u.Id) = *(tg.Id)

		return c.JSON(http.StatusOK, u)
	}
	tgs, err := srv.uc.GetTestGroupwithLogin(*user.Login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var a []struct {
		Id           *int    `json:"int"`
		Memory_limit *int    `json:"memory_limit"`
		Time_limit   *int    `json:"time_limit"`
		Kolvo        *int    `json:"quantity_tests"`
		Author       *string `json:"author"`
		Name         *string `json:"name"`
	}
	for _, val := range tgs {
		u := struct {
			Id           *int    `json:"int"`
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
