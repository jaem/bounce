package bouncer

import (
	"os"
	"log"
	"strconv"
	"net/http"
	"reflect"
	"github.com/jaem/nimble"
	"golang.org/x/net/context"
	"fmt"
	"github.com/jaem/bouncer/session/jwt"
	"github.com/jaem/bouncer/models"
)

var logger = log.New(os.Stdout, "[bounce.] ", 0)

// identity manager (eg. jwt, session etc)
var	session models.ISession

// Hashmap of authoriy providers used by server [string, Policy]
var providers models.ProviderMap

func init() {
	session = jwt.NewSession()
	providers = models.ProviderMap{}
}

// Register adds a policy to the hashmap.
func UseProvider(key string, p models.IProvider) {
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

// Disconnect an existing login
func InvalidateSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println("....... logging out")
	session.DeleteIdentity(w, r)
	onLogoutRedirect(w, r)
}

// Authenticate starts the authentication per request
func Authenticate(vider string) func(w http.ResponseWriter, r *http.Request) {
	// sanity check to ensure that policies are registered
	//var providers []string
	//for _, vider := range viders {
	//	exist := b.pmap[vider]
	//	if exist != nil {
	//		providers = append(providers, vider)
	//	}
	//}

	provider := providers[vider]
	if provider == nil {
		//panic(errors.New("No authentication provider specified"))
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
	loggedIn := isLoggedIn(r)
	logger.Println("isLoggedIn = " + strconv.FormatBool(loggedIn))
	if loggedIn {
		next(w, r)
	} else { // no identity
		onLoginRedirect(w, r)
	}
}

// Check if there is an id in the context.
func isLoggedIn(r *http.Request) bool {
	c := nimble.GetContext(r)
	if id, ok := c.Value("id").(*models.Identity); ok {
		logger.Println("identity.Uid = " + id.Uid)
		return true
	}
	return false
}
