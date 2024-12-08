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
	t := template.Must(template.ParseFiles("static/testing_request.html"))

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
		}
		err = b.dp.CreateUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
	}

}

func (b *BD_handlers) Login_user(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&user)
		sr := b.dp.Is_In_Base(user.Login, user.Password)
		fmt.Println(sr)
		if !b.dp.Is_In_Base(user.Login, user.Password) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Login does not exist or data is incorrect"))
			return
		} else {
			s := Session{Astiay_isos: b.dp.Gen_coockie(user.Login), User_login: user.Login}
			b.dp.CreateSession(&s)
			cookie := http.Cookie{
				Name:     "token",
				Value:    s.Astiay_isos,
				Path:     "/login",
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
		Cookie string `json:"Astiay_isos"`
	}{}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&cookie)
		if err := b.dp.Del_session(cookie.Cookie); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Session does not exist"))
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
}
