package usecase

import (
	"crypto/sha256"
	"fmt"
)

type Usecase struct {
	p DatabaseProvider
}

func NewUsecase(p DatabaseProvider) *Usecase {
	return &Usecase{
		p: p,
	}
}

func hash(value string) string {
	h := sha256.New()
	h.Write([]byte(value))

	return fmt.Sprintf("%x", h.Sum(nil))
}
