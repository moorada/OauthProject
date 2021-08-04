package main

import (
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"net/http"
)
const idDocente = "id_docente"
func InfoDocente(c echo.Context) error {

	idDocente := c.Param(idDocente)

	sub := getUserFromSession(c)                        // the user that wants to access a resource.
	obj := "/docente/" + idDocente + "/info" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)
	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	//ll.Info("%s",casbin_enforcer.GetAllRoles())

	if ok == true {
		docente, err := db.GetAPIDocenteFromMatricola(idDocente)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		if docente.Matricola == "" {
			return c.String(http.StatusNotFound, "Nessun docente trovato")
		}
		return c.JSON(http.StatusOK, docente)
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}
}
