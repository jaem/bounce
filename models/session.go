package models

import "net/http"

type ISession interface {
	GetIdentity(http.ResponseWriter, *http.Request) (*Identity, error)
	SaveIdentity(*Identity, http.ResponseWriter, *http.Request)
	DeleteIdentity(http.ResponseWriter, *http.Request)
}
