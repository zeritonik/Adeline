package api

import (
	"adeline/backend/internal/provider"
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
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
	us, err := srv.uc.GetUser(cookieval)
	if err != nil {
		return c.JSON(500, Message{err.Error()})
	}
	return c.JSON(200, us)
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
	if err != nil && err.Error() != "http: no such file" {
		return cc.JSON(500, Message{err.Error()})
	} else if err == nil {
		if *cc.Avatar != "" {
			err := os.Remove("." + *cc.Avatar)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := SaveImg(file, *cc.Login); err != nil {
			return cc.JSON(400, Message{err.Error()})
		}
		path := "/media/avatars/" + *cc.Login + "_" + time.Now().Format("2006-01-02-15-04-05") + ".png"
		u.Avatar = &path
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
	tr := provider.TestGroupResult{}
	tr.Group_id = tg.Id
	tr.Language = inf.Language
	tr.Source_code = inf.Source
	tr.Test_results = make([]provider.TestResult, 0)
	switch *inf.Language {
	case "python":
		if err := ExecutePython(tg, &tr); err != nil {
			return cc.JSON(500, Message{err.Error()})
		}
	case "go":
		if err := ExecuteGO(tg, &tr); err != nil {
			return cc.JSON(500, Message{err.Error()})
		}
	default:
		return cc.JSON(400, Message{"Not supported language"})
	}

	if err := srv.uc.AddTestGroupResult(tr); err != nil {
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
	path := "backend/tests/" + *tg.Author + "_prog.py"
	f, _ := os.Create(path)
	f.WriteString(*tr.Source_code)
	f.Close()
	var maxTime int64 = 0
	maxMemory := 0
	flag := "OK"
	for i, val := range tg.Tests {
		rez := provider.TestResult{}
		id := i + 1
		rez.Test_id = &id
		cmd := exec.Command("python3", "backend/scripts/time_mem_run.py", "-t", strconv.Itoa(*tg.Time_limit), "-m", strconv.Itoa(*tg.Memory_limit), "python3", path)
		cmd.Stdin = bytes.NewBufferString(*val.Input + "\n$$\n" + *val.Correct_output + "\n$$")
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		out := strings.Split(string(output[:]), "\n")

		if out[0] != "OK" {
			flag = out[0]
		}
		rez.Verdict = &out[0]
		rez.Output = &out[9]
		t, _ := strconv.ParseInt(out[2], 10, 64)
		rez.Execution_time = &t
		m, _ := strconv.Atoi(out[4])
		rez.Max_memory = &m
		maxMemory = max(maxMemory, *rez.Max_memory)
		maxTime = max(maxTime, *rez.Execution_time)

		tr.Test_results = append(tr.Test_results, rez)
	}
	tr.Max_memory = &maxMemory
	tr.Max_execution_time = &maxTime
	tr.Verdict = &flag
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
func ExecuteGO(tg *provider.TestGroup, tr *provider.TestGroupResult) error {
	path := "backend/tests/" + *tg.Author + "_prog.go"
	f, _ := os.Create(path)
	f.WriteString(*tr.Source_code)
	f.Close()
	var maxTime int64 = 0
	maxMemory := 0
	flag := "OK"
	for i, val := range tg.Tests {
		rez := provider.TestResult{}
		id := i + 1
		rez.Test_id = &id
		cmd := exec.Command("go", "run", "backend/scripts/time_mem_run.go")
		cmd.Stdin = bytes.NewBufferString(path + " " + strconv.Itoa(*tg.Memory_limit) + " " + strconv.Itoa(*tg.Time_limit) + "\n" + *val.Input + "\n#\n" + *val.Correct_output + "\n#")
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		out := strings.Split(string(output[:]), "\n")
		fmt.Println(out)

		if out[0] != "OK" {
			flag = out[0]
		}
		rez.Verdict = &out[0]
		rez.Output = &out[2]
		t, _ := strconv.ParseInt(out[8], 10, 64)
		rez.Execution_time = &t
		m, _ := strconv.Atoi(out[6])
		fmt.Println(m)
		rez.Max_memory = &m
		maxMemory = max(maxMemory, *rez.Max_memory)
		maxTime = max(maxTime, *rez.Execution_time)

		tr.Test_results = append(tr.Test_results, rez)
	}
	tr.Max_memory = &maxMemory
	tr.Max_execution_time = &maxTime
	tr.Verdict = &flag
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
func SaveImg(file *multipart.FileHeader, login string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}
	src.Seek(0, 0)
	if http.DetectContentType(buffer) == "image/png" {
		img, err := png.Decode(src)
		if err != nil {
			return err
		}
		if img.Bounds().Dx() != 128 && img.Bounds().Dy() != 128 {
			return echo.NewHTTPError(400, Message{"This size is not supported"})
		}
	} else {
		return echo.NewHTTPError(400, Message{"This image format is not supported."})

	}

	src.Seek(0, 0)
	path := "media/avatars/" + login + "_" + time.Now().Format("2006-01-02-15-04-05") + ".png"
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
