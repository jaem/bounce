package bounce

import (
	"net/http"
	"reflect"
	"fmt"
	"github.com/jaem/nimble"
	"golang.org/x/net/context"
	"gopkg.in/check.v1"
)

type providerMap map[string]Provider

type bounce struct {
	idm  IdManager   // identity manager (eg. jwt, session etc)
	pmap providerMap // Hashmap of authoriy providers used by server [string, Policy]
}

func New(m IdManager) *bounce {
	return &bounce{ idm: m, pmap:providerMap{} }
}

// Register adds a policy to the hashmap.
func (b *bounce) Register(key string, p Provider) {
	if key == "" || p == nil {
		fmt.Println("Failed to registered provider: " + key + " using " + reflect.TypeOf(p).String())
		return
	}
	b.pmap[key] = p
	fmt.Println("Successfully registered provider: " + key + " using " + reflect.TypeOf(p).String())
}

func (b *bounce) Deregister(key string) {
	delete(b.pmap, key)
}

func (b *bounce) getIdentity(w http.ResponseWriter, r *http.Request) bool {
	id, err := b.idm.GetIdentity(w, r)
	if err != nil {
		// go back to login page
		return false
	}

	c := nimble.GetContext(r)
	c = context.WithValue(c, "id", id)
	nimble.SetContext(r, c)

	return true
}

func (b* bounce) Reauthenticate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if b.getIdentity(w, r) {
		next(w, r)
	} else {
		http.Redirect(w, r, "/auth/login", 301)
	}
}

// Authenticate starts the authentication per request
func (b *bounce) Authenticate(viders ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// sanity check to ensure that policies are registered
	var providers []string
	for _, provi := range viders {
		exist := b.pmap[provi]
		if exist != nil {
			providers = append(providers, provi)
		}
	}

	// return a 404 error if a policy is not specified
	if len(providers) == 0 {
		Fatal(nothing to authenticate against)
		return
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if b.getIdentity(w, r) {
			// already logged in
			http.Redirect(w, r, "/", 200)
		}

		var id Identity
		// means id is probably not valid.
		// lets get identity from various providers.
		for _, p := range providers {
			id, _ := b.pmap[p].ResolveProvider(r)
			if id != nil {
				break
			}
		}

		if id == nil {
			http.NotFound(w, r)
			return
		}

		// successfully authenticated
		b.idm.SaveIdentity(id, w, r)

		next(w, r)
	}
}

func (b *bounce) Hoho(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("hohohohoh")
	next(w, r)
}