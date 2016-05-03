package jwt

import (
	"net/http"
	"github.com/jaem/bounce"
	"fmt"
)

//ensure that AuthHandler implements http.Handler
var _ bounce.Storage = (*Storage)(nil)

func NewStorage() *Storage {
	return &Storage{}
}

type Storage struct {}

func (s *Storage)	Load(*http.Request) {
		fmt.Println("loading jwt ")
}

func (s *Storage) Save(http.ResponseWriter, *http.Request) {
}