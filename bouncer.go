package bouncer

import (
	"os"
	"log"
	"strconv"
	"net/http"
	"reflect"
	"github.com/jaem/nimble"
	"golang.org/x/net/context"
	//"fmt"
	"fmt"
)

var logger = log.New(os.Stdout, "[bounce.] ", 0)

type providerMap map[string]Provider

type Bouncer struct {
	idm  IdManager   // identity manager (eg. jwt, session etc)
	pmap providerMap // Hashmap of authoriy providers used by server [string, Policy]
}

func New(m IdManager) *Bouncer {
	return &Bouncer{ idm: m, pmap:providerMap{} }
}

// Register adds a policy to the hashmap.
func (b *Bouncer) Register(key string, p Provider) {
	if key == "" || p == nil {
		logger.Println("Failed to registered provider: " + key + " using " + reflect.TypeOf(p).String())
		return
	}
	b.pmap[key] = p
	logger.Println("Successfully registered provider: " + key + " using " + reflect.TypeOf(p).String())
}

// Deregister a policy from hashmap
func (b *Bouncer) Deregister(key string) {
	delete(b.pmap, key)
}

// IdentifyRequest gets the user identity for the request. Default method is jwt.
func (b*Bouncer) IdentifyRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := b.idm.GetIdentity(w, r)
	if err != nil {
		next(w, r)
		return
	}
	// TODO: validate userid with database
	c := nimble.GetContext(r)
	c = context.WithValue(c, "id", id)
	nimble.SetContext(r, c)
	next(w, r)
}

// Authenticate starts the authentication per request
func (b *Bouncer) Authenticate(vider string) func(w http.ResponseWriter, r *http.Request) {
	// sanity check to ensure that policies are registered
	//var providers []string
	//for _, vider := range viders {
	//	exist := b.pmap[vider]
	//	if exist != nil {
	//		providers = append(providers, vider)
	//	}
	//}

	provider := b.pmap[vider]
	if provider == nil {
		//panic(errors.New("No authentication provider specified"))
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if isLoggedIn(r) {
			// already logged in
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		id, err := provider.ResolveProvider(r)
		if err != nil || id == nil {
			fmt.Println("unauthorized")
			NotAuthorized(w, r) // unauthorized
			return
		}

		// successfully authenticated
		b.idm.SaveIdentity(id, w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

// Checks if the user is logged in, before proceeding with the next handler
func IsLoggedIn(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	loggedIn := isLoggedIn(r)
	logger.Println("isLoggedIn = " + strconv.FormatBool(loggedIn))
	if loggedIn {
		next(w, r)
	} else { // no identity
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	}
}

// Check if there is an id in the context.
func isLoggedIn(r *http.Request) bool {
	c := nimble.GetContext(r)
	if id, ok := c.Value("id").(*Identity); ok {
		logger.Println("identity.Uid = " + id.Uid)
		return true
	}
	return false
}
