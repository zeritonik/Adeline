package api

import (
	"adeline/backend/internal/provider"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (Srv *Server) AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !strings.Contains(c.Request().URL.String(), "api") {
			return next(c)
		}
		cc := &CustomCont{Context: c, User: nil, UserCookie: nil}
		cookie, err := c.Cookie("astiay_isos")
		if err != nil {
			return next(cc)
		}
		cc.UserCookie = &cookie.Value
		f, err := Srv.uc.CheckSession(cookie.Value)
		if err != nil {
			c.Error(echo.ErrInternalServerError)
			return nil
		}
		if !f {
			return next(cc)
		}
		u, _ := Srv.uc.GetUser(cookie.Value)
		cc.User = u
		return next(cc)
	}
}
func (srv *Server) PostCreateUser(c echo.Context) error {
	user := provider.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(500, Message{Err: err.Error()})
	}
	fmt.Println(user.Avatar, user.Login, user.Nickname, user.Password)
	if err := srv.uc.CreateUser(user); err != nil {
		return c.JSON(500, Message{Err: err.Error()})
	}
	return c.JSON(http.StatusOK, Message{Err: "Registered"})
}

func (srv *Server) PostLogin(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User != nil {
		return cc.JSON(200, cc.UserInf)
	}
	u := provider.User{}
	if err := cc.Bind(&u); err != nil {
		return cc.JSON(400, Message{err.Error()})
	}
	if u.Password == nil {
		return cc.JSON(401, Message{"No Data"})
	}
	cookieval, err := srv.uc.LoginUser(*u.Login, *u.Password)
	if err != nil {
		return cc.JSON(401, Message{err.Error()})
	}
	c.SetCookie(&http.Cookie{Value: cookieval, Name: "astiay_isos", MaxAge: 3600})
	return c.JSON(200, u.UserInf)
}

func (srv *Server) GetSettings(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{""})
	}
	return c.JSON(http.StatusOK, cc.User)
}

func (srv *Server) PostLogout(c echo.Context) error {
	s := struct {
		All *bool `json:"all"`
	}{}
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{""})
	}
	if err := c.Bind(&s); err != nil {
		return c.JSON(500, Message{Err: err.Error()})
	}
	if err := srv.uc.DelSession(*cc.Login, *cc.UserCookie, *s.All); err != nil {
		return c.JSON(500, Message{Err: err.Error()})
	}
	return c.JSON(http.StatusOK, Message{"deleted"})
}

func (srv *Server) PostSettings(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{""})
	}
	u := provider.User{}
	if err := c.Bind(&u); err != nil {
		return cc.JSON(400, Message{err.Error()})
	}
	if err := srv.uc.ChangeSettings(u.Login, u.Password, u.Nickname, u.Avatar, *cc.UserCookie); err != nil {
		return c.JSON(500, Message{err.Error()})
	}
	return c.JSON(http.StatusOK, u.UserInf)
}

func (srv *Server) PostTests(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{""})
	}
	tg := provider.TestGroup{}
	if err := c.Bind(&tg); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tg.Author = cc.Login
	tg.Kolvo = tg.CalcCol()
	_, tg.Id = srv.uc.AddTestGroup(tg)
	return c.JSON(http.StatusCreated, &tg)
}

func (srv *Server) GetTests(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{""})
	}
	id := cc.Param("id")
	if id != "" {
		var tg *provider.TestGroup
		Id, _ := strconv.Atoi(id)
		tg, err := srv.uc.GetTestGroup(Id, *cc.Login)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Message{err.Error()})
		}
		if tg == nil {
			return c.JSON(500, Message{"No data to this user"})
		}
		tg.Kolvo = tg.CalcCol()
		return c.JSON(http.StatusOK, tg)
	}
	tgs, err := srv.uc.GetTestGroupwithLogin(*cc.Login)

	if err != nil {
		return c.JSON(500, Message{err.Error()})
	}
	for i := range tgs {
		tgs[i].Kolvo = tgs[i].CalcCol()
	}
	return c.JSON(http.StatusOK, tgs)
}

func (srv *Server) GetTestGroupRez(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{"Unautorized"})
	}

	inf := struct {
		Code        *string `json:"source_code"`
		Language    *string `json:"language"`
		TestGroupId *int    `json:"test_group_id"`
	}{}
	if err := cc.Bind(&inf); err != nil {
		return cc.String(http.StatusInternalServerError, err.Error())
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
	ret := []provider.Rez{*tgr}

	return c.JSON(http.StatusOK, ret)

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
