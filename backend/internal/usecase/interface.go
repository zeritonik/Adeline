package usecase

import (
	"adeline/backend/internal/provider"
)

type DatabaseProvider interface {
	CreateUser(provider.User) error
	CreateSession(string, string) error
	DelSession(string) error
	DelAllSessions(string) error
	GetUserInf(string) (*provider.UserInf, error)
	GetSession(string) (bool, error)
	ChangeUserAvatar(string, string) error
	ChangeUserNick(string, string) error
	ChangeUserPassword(string, string) error
	ChangeUserLogin(string, string) error
	IsInBase(string, string) (bool, error)
	InsertTestGroup(tg provider.TestGroup) (error, map[provider.Test]error)
	GetTestGroupInfo(id int, login string) (*provider.TestGroup, error)
	GetTestGroupInfoLOGIN(login string) ([]provider.TestGroup, error)
	DeleteTestGroup(id int, login string) error
	GetTestGroupResultInfo(login string) (*provider.Rez, error)
}
