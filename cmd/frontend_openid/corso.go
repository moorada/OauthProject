package main

import (
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"net/http"
	"net/http/httputil"
)

const IDcorso = "id_corso"

func InfoCorso(c echo.Context) error {
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta homepage a FRONTEND_OPENID: \n%v",string(requestDump))


	idCorso := c.Param(IDcorso)

	sub := getUserFromSession(c)         // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/info" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error("Error enfoce casbin in get info corso: %s", err.Error())
	}

	if ok == true {
		corso, err := db.GetAPICorsoFromId(idCorso)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, corso)
	} else {
		return c.String(http.StatusForbidden, ErrNonSeiAutorizzato)
	}

}
func VotiCorso(c echo.Context) error {
	idCorso := c.Param(IDcorso)

	sub := getUserFromSession(c)         // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/voti" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error("Error enfoce casbin in get info corso: %s", err.Error())
	}

	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		voti, err := db.GetAPIVotiFromCorso(corso)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, voti)
	} else {
		return c.String(http.StatusForbidden, ErrNonSeiAutorizzato)
	}

}
