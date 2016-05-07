package bounce

import (
	"os"
	"log"
	"strconv"
	"net/http"
	"reflect"
	"github.com/jaem/nimble"
	"golang.org/x/net/context"
)

var logger = log.New(os.Stdout, "[bounce.] ", 0)

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
		logger.Println("Failed to registered provider: " + key + " using " + reflect.TypeOf(p).String())
		return
	}
	b.pmap[key] = p
	logger.Println("Successfully registered provider: " + key + " using " + reflect.TypeOf(p).String())
}

func (b *bounce) Deregister(key string) {
	delete(b.pmap, key)
}

func (b* bounce) IdentifyRequest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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


// Check if there is an id in the context.
func (b* bounce) isLoggedIn(r *http.Request) bool {
	c := nimble.GetContext(r)
	if id, ok := c.Value("id").(*Identity); ok {
		logger.Println("identity.Uid = " + id.Uid)
		return true
	}
	return false
}

func (b* bounce) IsLoggedIn(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	loggedIn := b.isLoggedIn(r)
	logger.Println("isLoggedIn = " + strconv.FormatBool(loggedIn))
	if loggedIn {
		next(w, r)
	} else { // no identity
		//http.Redirect(w, r, "/auth/login", 301)
	}
}

// Authenticate starts the authentication per request
func (b *bounce) Authenticate(vider string) func(w http.ResponseWriter, r *http.Request) {
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
		if b.isLoggedIn(r) {
			// already logged in
			http.Redirect(w, r, "/", 200)
			return
		}

		id, err := provider.ResolveProvider(r)
		if err != nil || id == nil {
			NotAuthorized(w, r) // unauthorized
			return
		}

		// successfully authenticated
		b.idm.SaveIdentity(id, w, r)
	}
}