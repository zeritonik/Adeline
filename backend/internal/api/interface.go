package api

import "adeline/backend/internal/provider"

type Usecase interface {
	CreateUser(provider.User) error
	LoginUser(string, string) (string, error)
	DelSession([]string, bool) error
	GetUser(string) (*provider.UserInf, error)
	ChangeSettings(*string, *string, *string, *string, string) error
	CheckSession(string) (bool, error)
	AddTestGroup(tg provider.TestGroup) error
	GetTestGroup(int, string) (*provider.TestGroup, error)
	GetTestGroupwithLogin(login string) ([]provider.TestGroup, error)
	DeleteTestGroup(login string, id int) error
}
