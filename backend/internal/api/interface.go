package api

import (
	"adeline/backend/internal/provider"

	"github.com/labstack/echo/v4"
)

type Usecase interface {
	CreateUser(provider.User) error
	LoginUser(string, string) (string, error)
	DelSession(string, string, bool) error
	GetUser(string) (*provider.User, error)
	ChangeSettings(*string, *string, *string, *string, string) error
	CheckSession(string) (bool, error)
	AddTestGroup(tg provider.TestGroup) (error, *int)
	GetTestGroup(int, string) (*provider.TestGroup, error)
	GetTestGroupwithLogin(login string) ([]provider.TestGroup, error)
	DeleteTestGroup(login string, id int) error
	GetTestGroupResult(login string) (*provider.Rez, error)
}

type Message struct {
	Err string `json:"response"`
}
type CustomCont struct {
	echo.Context
	*provider.User
	UserCookie *string
}
