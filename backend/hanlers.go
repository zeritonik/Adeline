package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type BD_handlers struct {
	dp DatabaseProvider
}

type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func newResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}
func (p *pathResolver) Addpath(path string, fun http.HandlerFunc) {
	p.handlers[path] = fun
}
func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, check); ok && err == nil {
			handlerFunc(res, req)
			return
		} else if err != nil {
			fmt.Print(res, err)
		}

	}
	http.NotFound(res, req)
}

func homepage(w http.ResponseWriter, r *http.Request) {

	p := &Page{
		Title:   "Adeline",
		Content: "",
		err:     nil,
	}
	t := template.Must(template.ParseFiles("build/index.html"))

	var info CodeInfo
	switch r.Method {
	case "GET":
		info.Language = r.URL.Query().Get("code-language")
		info.CodeText = r.URL.Query().Get("code-text")
		info.CodeFile = r.URL.Query().Get("code-file")
		fmt.Println(info.Language)
		fmt.Println(info.CodeText)
		p = info.ExecuteProgram()

	}
	t.Execute(w, p)

}
func profilepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

type Page struct {
	Title   string
	Content string
	err     error
}

type CodeInfo struct {
	Language string
	CodeText string
	CodeFile string
}

func (code *CodeInfo) ExecuteProgram() (p *Page) {
	p = &Page{
		Title:   "Adeline",
		Content: "",
		err:     nil,
	}
	switch code.Language {
	case "Python":
		cmd := exec.Command("python3", "python/prog.py")
		var out bytes.Buffer
		cmd.Stdout = &out

		f, _ := os.Create("python/prog.py")
		f.WriteString(code.CodeText)
		f.Close()
		err := cmd.Run()
		if err == nil {

			var rez string
			for _, val := range out.Bytes() {
				rez += (string(val))
			}

			fmt.Println(err)

			p.Content = rez

			return p
		} else {
			p = &Page{
				Title:   "Adeline",
				Content: "Error",
				err:     err,
			}
			return p
		}
	case "C++":
		f, _ := os.Create(("c++/proj/main.cpp"))
		f.WriteString(code.CodeText)
		f.Close()
		f, _ = os.Create("c++/proj/CMakeLists.txt")
		f.WriteString("cmake_minimum_required(VERSION 3.22.1)\n")
		f.WriteString("project(main)\n")
		f.WriteString("add_executable(main main.cpp)")
		f.Close()
		var out bytes.Buffer
		cmd := exec.Command("cmake", "-S", "c++/proj", "-B", "c++/build-dir")
		err := cmd.Run()
		if err != nil {
			fmt.Println(1, err)
			log.Fatal(err)
		}
		cmd = exec.Command("cmake", "--build", "c++/build-dir")
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)

			cmd = exec.Command("c++/build-dir/main")
			cmd.Stdout = &out
			fmt.Println(cmd.Err, exec.ErrDot)
			if errors.Is(cmd.Err, exec.ErrDot) {
				cmd.Err = nil

			}

			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			fmt.Println(err)

			fmt.Println(out.String())
		}

	}
	return p
}

func (b *BD_handlers) Create_user(w http.ResponseWriter, r *http.Request) {
	user := User{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else if user.Avatar == nil || user.Login == nil || user.Nickname == nil || user.Password == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NO DATA"))
		} else if err = b.dp.CreateUser(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
	}

}

func (b *BD_handlers) Login_user(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Login    *string `json:"login"`
		Password *string `json:"password"`
	}{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else if user.Login == nil || user.Password == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No data"))
		} else if !b.dp.Is_In_Base(*(user.Login), *(user.Password)) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Login does not exist or data is incorrect"))
			return
		} else {
			cooki := (b.dp.Gen_coockie(*(user.Login)))
			s := Session{Astiay_isos: &cooki, User_login: user.Login}
			if err := b.dp.CreateSession(&s); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}
			cookie := http.Cookie{
				Name:     "token",
				Value:    *(s.Astiay_isos),
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)

		}

	}
}

func (b *BD_handlers) Delete_Session(w http.ResponseWriter, r *http.Request) {
	cookie := struct {
		Cookies *[]*string `json:"Astiay_isos"`
		All     *bool      `json:"all"`
	}{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&cookie)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else if cookie.Cookies == nil || cookie.All == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No data"))
		} else {
			if !(*(cookie.All)) {
				for i := 0; i < len(*cookie.Cookies); i++ {
					if err := b.dp.Del_session(*(*cookie.Cookies)[i]); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
						break
					}
				}
			} else {
				if err := b.dp.Del_All_Sessions(); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
				}
			}
		}
	}
}

func (b *BD_handlers) Settings(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Login  *string `json:"login"`
		Avatar *string `json:"avatar"`
	}{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&res); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if res.Login == nil || res.Avatar == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No data"))
			return
		}
		if err := b.dp.Change_user_avatar(*res.Login, *res.Avatar); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return

		}
		w.WriteHeader(http.StatusOK)
	case "GET":
		cookie, err := r.Cookie("Astiay_isos")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		fmt.Println(cookie.Value)
		us, err := b.dp.Get_User(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if us.Login == nil || us.Avatar == nil || us.Nickname == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Incorrect data in table"))
			return
		}

		u, err := json.Marshal(us)
		w.WriteHeader(http.StatusOK)
		w.Write(u)

	}
}
