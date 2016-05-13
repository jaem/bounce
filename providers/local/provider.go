package local

import (
	"net/http"

	"github.com/jaem/bounce"
	"fmt"
	"errors"
)

//ensure that AuthHandler implements http.Handler
var _ bounce.IProvider = (*Provider)(nil)

type VerifyFunc func(username string, password string) (*bounce.Identity, error)

func NewProvider(fn VerifyFunc) *Provider {
	p := &Provider{
		usernameField: "username",
		passwordField: "password",
		verify: fn,
	}

	return p
}

type Provider struct {
	*bounce.Provider
	usernameField string
	passwordField string
	verify VerifyFunc
}

func (p *Provider) ResolveProvider(r *http.Request) (*bounce.Identity, error) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	fmt.Println("username = " + username)
	fmt.Println("password = " + password)

	if username == "" || password == "" {
		return nil, errors.New("Missing credentials")
	}

	identity, err := p.verify(username, password)

	if err != nil {
		return nil, err
	}

	return identity, nil
}

