package main

import (
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"net/http"
)

const IDstudente = "id_studente"

func InfoStudente(c echo.Context) error {
	idStudente := c.Param(IDstudente)

	sub := getUserFromSession(c)                          // the user that wants to access a resource.
	obj := "/studente/" + idStudente + "/info" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	//ll.Info("%s",casbin_enforcer.GetAllRoles())

	if ok == true {
		studente, err := db.GetAPIStudenteFromMatricola(idStudente)
		if err != nil{
			return c.String(http.StatusNotFound, err.Error())
		}
		if studente.Matricola == "" {
			return c.String(http.StatusNotFound, "Nessuno studente trovato")
		}
		return c.JSON(http.StatusOK, studente)
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}
}
func VotiStudente(c echo.Context) error {

	idStudente := c.Param(IDstudente)

	sub := getUserFromSession(c)                           // the user that wants to access a resource.
	obj := "/studente/" + idStudente + "/voti" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}

	if ok == true {
		studente, err := db.GetStudenteFromMatricola(idStudente)
		voti, err := db.GetAPIVotiFromStudente(studente)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, voti)
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}

