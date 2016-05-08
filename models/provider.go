package models

import "net/http"

// Provider needs to be implemented for every authentication provider
type IProvider interface {
	ResolveProvider(*http.Request) (*Identity, error)
}

type ProviderMap map[string]IProvider

type Provider struct {}

//func (ap *Provider) Complete(id *Identity) {
//	if (err) { return self.error(err); }
//	if (!user) { return self.fail(info); }
//	self.success(user, info);
//}
