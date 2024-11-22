package main

import (
	"fmt"
	"net/http"
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
	}
	t := template.Must(template.ParseFiles("static/testing_request.html"))
	t.Execute(w, p)
}
func profilepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

type Page struct {
	Title   string
	Content string
}
