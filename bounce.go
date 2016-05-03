package bounce

import (
	"net/http"
	"reflect"
	"fmt"
)

type Storage interface {
	Load(*http.Request)
	Save(http.ResponseWriter, *http.Request)
}

type Provider interface {
	ResolveProvider(*http.Request)
}

type providerMap map[string]Provider

type bounce struct {
	storage Storage
	pmap providerMap // Hashmap of providers used by server [string, Policy]
}

func New(storage Storage) *bounce {
 return &bounce{ storage: storage, pmap: providerMap{} }
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
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			http.NotFound(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		fmt.Println("...... running authenticate ")

		b.storage.Load(r)

		for _, p := range providers {
			b.pmap[p].ResolveProvider(r)
		}

		authenticated := false
		if authenticated {
			next(w, r)
		} else {
			http.NotFoundHandler()
		}
	}
}

func (b *bounce) Hoho(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("hohohohoh")
	next(w, r)
}