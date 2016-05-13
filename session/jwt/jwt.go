package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
	"fmt"
	"github.com/jaem/bounce"
)

const defaultTTL = 3600 * 24 * 7 // 1 week
const token_secret = "secret"

//ensure that AuthHandler implements http.Handler
var _ bounce.ISession = (*IdManager)(nil)

func NewSession() *IdManager {
	return &IdManager{ fp: fingerprint{
		key: []byte(token_secret),
		method: jwt.SigningMethodHS256,
		ttl: defaultTTL,
	}}
}

type fingerprint struct {
	key    []byte
	method jwt.SigningMethod
	ttl    int64
}

// IdManager is a JSON Web Token (JWT) Provider which create or retrieves tokens
// with a particular signing key and options.
type IdManager struct {
	fp fingerprint
}

func (m *IdManager)	GetIdentity(w http.ResponseWriter, r *http.Request) (*bounce.Identity, error) {
	//req, _ := httputil.DumpRequest(r, false)
	//fmt.Println(string(req))
	token, err := getToken(r, m.fp)
	if err != nil {
		// bad unverified jwt token -
		return nil, err
	}
	id := new(bounce.Identity)
	if val, ok := token.Claims["uid"].(string); val != "" && ok {
		id.Uid = val
	} else {
		return id, errors.New("missing uid in jwt token")
	}
	if val, ok := token.Claims["access"].(string); val != "" && ok {
		id.Access = val
	}
	return id, nil
}

func (m *IdManager) SaveIdentity(id *bounce.Identity, w http.ResponseWriter, r *http.Request) {
	jwtToken, err := newSignedString(id, m.fp)
	if err != nil {
		fmt.Println("Unable to generate new jwtToken in jwt.IdManager.SaveIdentity")
		return
	}
	w.Header().Set("Set-Cookie", "jwt_token=" + string(jwtToken) + "; httponly; path=/;")
	//w.Header().Set("Set-Cookie", "jwt_token=" + jwtToken + ";Secure;HttpOnly;") // https - need TLS (key, cert) during production
	//w.Header().Set("Authorization", "Bearer " + jwtToken)
}

func (m *IdManager) DeleteIdentity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", "jwt_token=deleted; HttpOnly; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT")
}

// Based on https://github.com/dghubble/jwts/jwts.go
// Get gets the signed JWT from the Authorization header. If the token is
// missing, expired, or the signature does not validate, returns an error.
func getToken(r *http.Request, fp fingerprint) (*jwt.Token, error) {
	//jwtString := r.Header.Get("Authorization")
	jwtToken, err := r.Cookie("jwt_token")
	if err != nil {
		return nil, jwt.ErrNoTokenInRequest
	}
	token, err := jwt.Parse(jwtToken.Value, keyFunc(fp))
	if err == nil && token.Valid {
		// token parsed, exp/nbf checks out, signature verified, Valid is true
		return token, nil
	}
	return nil, jwt.ErrNoTokenInRequest
}

// keyFunc accepts an unverified JWT and returns the signing/verification key.
// Also ensures tha the token's algorithm matches the signing method expected.
func keyFunc(fp fingerprint) func(
unverified *jwt.Token) (interface{}, error) {
	return func (unverified *jwt.Token) (interface{}, error) {
		// require token alg to match the set signing method, do not allow none
		if meth := unverified.Method; meth == nil || meth.Alg() != fp.method.Alg() {
			return nil, jwt.ErrHashUnavailable
		}
		return fp.key, nil
	}
}

// New returns a new *jwt.Token which has the prescribed signing method, issued
// at time, and expiration time set on it.
// Add claims to the Claims map and use the controller to sign digitally(token) to get
// the a JWT byte slice
func newSignedString(id *bounce.Identity, fp fingerprint) ([]byte, error) {
	token := jwt.New(fp.method)
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Duration(fp.ttl) * time.Second).Unix()
	token.Claims["uid"] = id.Uid
	token.Claims["access"] =  id.Access //"[{\"type\":\"repository\",\"action\":\"push\"}]"
	jwtString, err := token.SignedString(fp.key)
	return []byte(jwtString), err
}
