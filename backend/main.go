package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"text/template"
)

func main() {
	pr := newResolver()
	pr.Addpath("*/home", homepage)
	pr.Addpath("*/profile", profilepage)

	var port string = ":8000"

	fmt.Println("Server is listening...")

	fmt.Println("http://localhost" + port)
	err := http.ListenAndServe(port, pr)
	if err != nil {
		fmt.Println("Ошибка запуска сервера")
	}

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
