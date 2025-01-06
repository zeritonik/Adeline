package api

import (
	"adeline/backend/internal/provider"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
		// fmt.Println(err.Error())

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
	file, err := cc.FormFile("avatar")
	if err != nil {
		return cc.JSON(500, Message{err.Error()})
	}
	if err := SaveImg(file); err != nil {
		return cc.JSON(500, Message{err.Error()})
	}
	path := "/media/avatars/" + file.Filename
	u.Avatar = &path
	fmt.Println(*u.Avatar)
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
		return c.JSON(500, Message{err.Error()})
	}
	if tg.Name == nil {
		return c.JSON(400, Message{"No data"})
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
		tg, err := srv.uc.GetTestGroup(Id)
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

func (srv *Server) SendCode(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{"Unautorized"})
	}
	inf := CodeInf{}
	if err := cc.Bind(&inf); err != nil {
		return cc.JSON(500, Message{err.Error()})
	}
	id, err := strconv.Atoi(cc.Param("id"))
	if err != nil {
		return cc.JSON(500, Message{err.Error()})
	}
	tg, err := srv.uc.GetTestGroup(id)
	if err != nil {
		return cc.JSON(400, Message{err.Error()})
	}
	// fmt.Println(*tg.Author)
	// fmt.Println(*inf.Language)
	// fmt.Println(*inf.Source)
	tr := provider.TestGroupResult{}
	tr.Group_id = tg.Id
	tr.Language = inf.Language
	tr.Source_code = inf.Source
	tr.Test_results = make([]provider.TestResult, 0)
	if err := ExecutePython(tg, &tr); err != nil {
		return cc.JSON(500, Message{err.Error()})
	}

	return c.JSON(200, Message{"ok"})
}

func (srv *Server) DeleteGroup(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return cc.JSON(401, Message{"Not autorized"})
	}
	id := cc.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "No_id")
	}
	i, _ := strconv.Atoi(id)
	if err := srv.uc.DeleteTestGroup(*cc.Login, i); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (srv *Server) GetResults(c echo.Context) error {
	cc := c.(*CustomCont)
	if cc.User == nil {
		return c.JSON(401, Message{"Not autorized"})
	}
	tgr, err := srv.uc.GetTestGroupResult(*cc.Login)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, err.Error())
	}
	return cc.JSON(200, tgr)

}
func ExecutePython(tg *provider.TestGroup, tr *provider.TestGroupResult) error {
	f, _ := os.Create("backend/tests/prog.py")
	f.WriteString(*tr.Source_code)
	f.Close()
	for i, val := range tg.Tests {
		rez := provider.TestResult{Verdict: new(string), Output: new(string)}
		id := i + 1
		rez.Test_id = &id
		cmd := exec.Command("timeout", strconv.FormatFloat(float64(100000000)/1000000, 'g', 1, 64), "python3", "backend/tests/prog.py")
		var out bytes.Buffer
		cmd.Stdout = &out

		start := time.Now()
		cmd.Stdin = bytes.NewBufferString(*val.Input)
		err := cmd.Run()
		duration := time.Since(start).Abs().Milliseconds()
		if err.Error() == "exit status 124" {
			*rez.Verdict = "TL"
			*rez.Output = ""
			rez.Execution_time = &duration
			fmt.Println(err)
			return nil
		} else if err != nil {

		}

		o := out.String()
		rez.Output = &o

		rez.Execution_time = &duration
		// fmt.Println(*rez.Execution_time)

		tr.Test_results = append(tr.Test_results, rez)

	}

	return nil
}

func SaveImg(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	path := "media/avatars/" + file.Filename
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
