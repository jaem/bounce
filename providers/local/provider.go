package local

import (
	"net/http"

	"github.com/jaem/bounce"
	"fmt"
	"github.com/gorilla/context"
)

//ensure that AuthHandler implements http.Handler
var _ bounce.Provider = (*Provider)(nil)

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

func (self *Provider) ResolveProvider(r *http.Request) {
	fmt.Println("local provider")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from LocalProvider, value is " + value.(string))
	} else {
		fmt.Println("from LocalProvider, value is nil")
	}
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


