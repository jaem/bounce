package local

import (
	"net/http"

	"github.com/jaem/bouncer/models"
	"fmt"
)

//ensure that AuthHandler implements http.Handler
var _ models.IProvider = (*Provider)(nil)

func NewProvider() *Provider {
	p := &Provider{
		usernameField: "username",
		passwordField: "password",
	}

	return p
}

type Provider struct {
	*models.Provider
	usernameField string
	passwordField string
}

func (p *Provider) ResolveProvider(r *http.Request) (*models.Identity, error) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	fmt.Println("username = " + username)
	fmt.Println("password = " + password)

	if username == "" || password == "" {
		//return this.fail({ message: options.badRequestMessage || 'Missing credentials' }, 400);
	}

	return &models.Identity{ Uid: username, Access:"some access" }, nil
}

