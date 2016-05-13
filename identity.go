package bounce

type Identity struct {
	Uid			string
	Access 	string
	Store 	map[string]string // to store app specific data

	//https://docs.docker.com/registry/spec/auth/jwt/
	//iss (Issuer) The issuer of the token, typically the fqdn of the authorization server.
	//sub (Subject) The subject of the token; the name or id of the client which requested it. This should be empty (`""`) if the client did not authenticate.
	//aud (Audience) The intended audience of the token; the name or id of the service which will verify the token to authorize the client/subject.
	//exp (Expiration) The token should only be considered valid up to this specified date and time.
	//nbf (Not Before) The token should not be considered valid before this specified date and time.
	//iat (Issued At) Specifies the date and time which the Authorization server generated this token.
	//jti (JWT ID) A unique identifier for this token. Can be used by the intended audience to prevent replays of the token.
	//The Claim Set will also contain a private claim name unique to this authorization server specification:
	//{
	//	"iss": "auth.docker.com",
	//	"sub": "jlhawn",
	//	"aud": "registry.docker.com",
	//	"exp": 1415387315,
	//	"nbf": 1415387015,
	//	"iat": 1415387015,
	//	"jti": "tYJCO1c6cnyy7kAn0c7rKPgbV1H1bFws",
	//	"access": [
	//		{
	//			"type": "repository",
	//			"name": "samalba/my-app",
	//			"actions": [
	//			"pull",
	//			"push"
	//			]
	//		}
	//	]
	//}

}
