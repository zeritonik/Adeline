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
			// fmt.Println(pattern)
			// fmt.Println(check)
			// ok, _ := path.Match(pattern, check)
			// fmt.Println(ok)
			handlerFunc(res, req)
			return
		}

	}
	http.ServeFile(res, req, "./build/index.html")

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
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if user.Avatar == nil {
		user.Avatar = new(string)
		*user.Avatar = "null"
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else if user.Login == nil || user.Nickname == nil || user.Password == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("NO DATA"))
	} else if err = b.dp.CreateUser(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)

		cookie := http.Cookie{
			Name:     "astiay_isos",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Value:    b.dp.Gen_coockie(*user.Login),
		}
		user_inf := struct {
			Nickname *string
			Login    *string
			Avatar   *string
		}{
			Nickname: user.Nickname,
			Login:    user.Login,
			Avatar:   user.Avatar,
		}
		http.SetCookie(w, &cookie)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(user_inf); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
		}
	}

}

func (b *BD_handlers) Login_user(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Login    *string `json:"login"`
		Password *string `json:"password"`
	}{}

	rcook, err := r.Cookie("astiay_isos")

	if err != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if user.Login == nil || user.Password == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No data"))
			return
		}
		if !b.dp.Is_In_Base(*(user.Login), *(user.Password)) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Login does not exist or data is incorrect"))
			return
		}
		rcook = &http.Cookie{
			Name:     "astiay_isos",
			Value:    b.dp.Gen_coockie(*user.Login),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
		}
		s := Session{Astiay_isos: &rcook.Value, User_login: user.Login}
		if err := b.dp.CreateSession(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

	}

	us, err := b.dp.Get_User(rcook.Value)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			w.Write([]byte("Not autorized"))
			return
		}
		fmt.Print(err.Error())
	}

	u, err := json.Marshal(&us)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, rcook)
	w.Write(u)

}

func (b *BD_handlers) Delete_Session_Post(w http.ResponseWriter, r *http.Request) {
	cookie := struct {
		Cookies *[]*string `json:"astiay_isos"`
		All     *bool      `json:"all"`
	}{}
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

func (b *BD_handlers) Get_Settings(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	cookie, err := r.Cookie("astiay_isos")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not autorized"))
		return
	}
	user, err = b.dp.Get_User(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not autorized"))
		return
	}
	us, err := json.Marshal(user)
	w.Write([]byte(us))
}

func (b *BD_handlers) Post_Settings(w http.ResponseWriter, r *http.Request) {
	user := User{}
	cookie, err := r.Cookie("astiay_isos")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not autorized"))
		return
	}
	if b.dp.IsSessionActive(cookie.Value) == false {
		fmt.Println(cookie.Value)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not autorized"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user.Avatar == nil && user.Nickname == nil && user.Password == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No data sent"))
		return
	}
	old, _ := b.dp.Get_User(cookie.Value)

	if user.Avatar != nil {
		b.dp.Change_user_avatar(*old.Login, *user.Avatar)
	}
	if user.Password != nil {
		b.dp.ChangeUserPassword(*old.Login, *user.Password)
	}
	if user.Nickname != nil {
		b.dp.ChangeUserNick(*old.Login, *user.Nickname)
	}

}

func get_HTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./build/index.html")
}
func get_Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./build"+r.URL.Path)
}
