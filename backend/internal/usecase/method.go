package usecase

import (
	"adeline/backend/internal/provider"
	"errors"
	"math/rand"
	"strconv"
)

func (u *Usecase) CreateUser(user provider.User) error {
	if user.Password != nil {
		*user.Password = hash(*user.Password)
	}
	return u.p.CreateUser(user)
}

func (u *Usecase) ChangeSettings(login *string, password *string, nickname *string, avatar *string, cookie string) error {
	user, err := u.p.GetUserInf(cookie)
	if err != nil {
		return err
	}
	if login != nil && user.Login != login {
		if err := u.p.ChangeUserLogin(*login, *user.Login); err != nil {
			return err
		}
	}
	if password != nil {
		if err := u.p.ChangeUserPassword(*user.Login, hash(*password)); err != nil {
			return err
		}
	}
	if nickname != nil {
		if err := u.p.ChangeUserNick(*user.Login, *nickname); err != nil {
			return err
		}
	}
	if avatar != nil {
		if err := u.p.ChangeUserAvatar(*user.Login, *avatar); err != nil {
			return err
		}
	}
	return nil

}

func (u *Usecase) LoginUser(login string, password string) (string, error) {
	flag, err := u.p.IsInBase(login, hash(password))
	if err != nil {
		return "", err
	}
	if flag == true {
		cookie := strconv.Itoa(rand.Int())
		u.p.CreateSession(cookie, login)
		return cookie, nil
	}

	return "", errors.New("No user in database")
}

func (u *Usecase) CheckSession(cookie string) (bool, error) {
	return u.p.GetSession(cookie)
}

func (u *Usecase) GetUser(cookie string) (*provider.User, error) {
	usinf, err := u.p.GetUserInf(cookie)
	if err != nil {
		return nil, err
	}
	user := &provider.User{UserInf: *usinf, Password: nil}
	return user, nil
}

func (u *Usecase) DelSession(login string, cookie string, all bool) error {
	if all {
		return u.p.DelAllSessions(login)
	}
	if err := u.p.DelSession(cookie); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) AddTestGroup(tg provider.TestGroup) (error, *int) {
	err, e, id := u.p.InsertTestGroup(tg)
	if err != nil {
		return err, nil
	}
	for key, val := range e {
		if val != nil {
			return errors.New(string(*key.Id) + val.Error()), nil
		}
	}
	return nil, &id
}

func (u *Usecase) AddTestGroupResult(tg provider.TestGroupResult) error {
	err, e := u.p.InsertTestGroupRezult(tg)
	if err != nil {
		return err
	}
	for key, val := range e {
		if val != nil {
			return errors.New(string(*key.Test_id) + val.Error())
		}
	}
	return nil
}

func (u *Usecase) GetTestGroup(id int) (*provider.TestGroup, error) {
	return u.p.GetTestGroupInfo(id)
}

func (u *Usecase) GetTestGroupwithLogin(login string) ([]provider.TestGroup, error) {
	return u.p.GetTestGroupInfoLOGIN(login)
}

func (u *Usecase) DeleteTestGroup(login string, id int) error {
	return u.p.DeleteTestGroup(id, login)
}

func (u *Usecase) GetTestGroupResult(login string) ([]provider.TestGroupResult, error) {
	return u.p.GetTestGroupResultInfo(login)
}
