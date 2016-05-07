package local

import (
	"net/http"

	"github.com/jaem/bouncer"
	"fmt"
	"github.com/gorilla/context"
)

//ensure that AuthHandler implements http.Handler
var _ bouncer.Provider = (*Provider)(nil)

func NewProvider() *Provider {
	p := &Provider{
		usernameField: "username",
		passwordField: "password",
	}

	return p
}

type Provider struct {
	usernameField string
	passwordField string
}

func (p *Provider) ResolveProvider(r *http.Request) (*bouncer.Identity, error) {
	fmt.Println("local provider")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from LocalProvider, value is " + value.(string))
	} else {
		fmt.Println("from LocalProvider, value is nil")
	}

	return &bouncer.Identity{ Uid:"jaem", Access:"some access" }, nil
}


//func (self *Provider) Verified(err, user, info) {
//	if (err) { return self.error(err); }
//	if (!user) { return self.fail(info); }
//	self.success(user, info);
//
//}
//
//func (self *Provider) Verify() {
//
//}


