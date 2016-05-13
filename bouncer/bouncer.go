package bouncer

import (
	"os"
	"log"
	"errors"
	"net/http"
	"reflect"
	"golang.org/x/net/context"

	"github.com/jaem/nimble"
	"github.com/jaem/bounce"
	"github.com/jaem/bounce/session/jwt"
)

var logger = log.New(os.Stdout, "[bounce.] ", 0)

// identity manager (eg. jwt, session etc)
var	session bounce.ISession

// Hashmap of authoriy providers used by server [string, Policy]
var providers bounce.ProviderMap

func init() {
	session = jwt.NewSession()
	providers = bounce.ProviderMap{}
}

// Register adds a policy to the hashmap.
func UseProvider(key string, p bounce.IProvider) {
	if key == "" || p == nil {
		logger.Println("Failed to registered provider: " + key + " using " + reflect.TypeOf(p).String())
		return
	}
	providers[key] = p
	logger.Println("Successfully registered provider: " + key + " using " + reflect.TypeOf(p).String())
}

// Deregister a policy from hashmap
func UnuseProvider(key string) {
	delete(providers, key)
}

// IdentifyRequest gets the user identity for the request. Default method is jwt.
func RestoreSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := session.GetIdentity(w, r)
	if err != nil || id == nil {
		next(w, r)
		return
	}

	// TODO: need to verify id??
	c := nimble.GetContext(r)
	c = context.WithValue(c, "id", id)
	nimble.SetContext(r, c)
	next(w, r)
}

// Disconnect an existing login
func InvalidateSession(w http.ResponseWriter, r *http.Request) {
	session.DeleteIdentity(w, r)
	onLogoutRedirect(w, r)
}

// Authenticate starts the authentication per request
func Authenticate(vider string) func(w http.ResponseWriter, r *http.Request) {
	// sanity check to ensure that the policy has been registered
	provider := providers[vider]
	if provider == nil {
		panic(errors.New("No authentication provider specified"))
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if isLoggedIn(r) {
			// already logged in
			onLoggedinRedirect(w, r)
			return
		}

		id, err := provider.ResolveProvider(r)
		if err != nil || id == nil {
			onUnauthorizedRedirect(w, r)
			return
		}

		// successfully authenticated
		session.SaveIdentity(id, w, r)
		onLoggedinRedirect(w, r)
		return
	}
}

// Checks if the user is logged in, before proceeding with the next handler
func IsLoggedIn(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if isLoggedIn(r) {
		next(w, r)
	} else { // no identity
		onLoginRedirect(w, r)
	}
}

// Check if there is an id in the context.
func isLoggedIn(r *http.Request) bool {
	c := nimble.GetContext(r)
	if id, ok := c.Value("id").(*bounce.Identity); ok {
		logger.Println("identity.Uid = " + id.Uid)
		return true
	}
	return false
}
