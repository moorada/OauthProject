package main

import (
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"net/http"
	"strconv"
)
func DispensaCorso(c echo.Context) error {
	idCorso := c.Param("id_corso")

	sub := getUserFromSession(c)                        // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/dispensa" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())
	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		dispense, err := db.GetAPIDispense(corso)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, dispense)
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}
func DispensaCorsoNumero(c echo.Context) error {

	idCorso := c.Param("id_corso")
	nDispensa, err := strconv.Atoi(c.Param("n_Dispensa"))

	if err != nil {
		ll.Error("Error to convert n. dispensa: %s", err.Error())
	}
	sub := getUserFromSession(c)                        // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/dispensa" // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)
	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())
	if ok == true {
		idCorsoInt, err := strconv.Atoi(idCorso)
		if err != nil {
			ll.Error("Error to convert n. corso: %s", err.Error())
		}
		dispensa, err := db.GetAPIDispensa(idCorsoInt, nDispensa)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, dispensa)
		//return c.String(http.StatusOK, "Non esiste questo dispensa")

	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}
func AddDispensaCorsoNumero(c echo.Context) error {

	idCorso := c.Param("id_corso")

	titolo := c.FormValue("titolo")
	testoDispensa := c.FormValue("testoDispensa")

	sub := getUserFromSession(c)                       // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/dispensa" // the resource that is going to be accessed.
	act := POST
	ok, err := casbin_enforcer.Enforce(sub, obj, act)
	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())
	if ok == true {
		db.AddDispensa(idCorso, titolo, testoDispensa)
		return c.String(http.StatusCreated, "Aggiunta dispensa")
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}
}

func EditDispensaCorsoNumero(c echo.Context) error {

	idCorso := c.Param("id_corso")
	idCorsoInt, err := strconv.Atoi(idCorso)
	nDispensa, err := strconv.Atoi(c.Param("n_Dispensa"))
	titolo := c.FormValue("titolo")
	testoDispensa := c.FormValue("testoDispensa")

	if err != nil {
		ll.Error("Error to convert n. dispensa in int: %s", err.Error())
	}
	sub := getUserFromSession(c)                        // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/dispensa" // the resource that is going to be accessed.
	act := PUT
	ok, err := casbin_enforcer.Enforce(sub, obj, act)
	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())
	if ok == true {

		dispensa, err := db.GetDispensa(idCorsoInt, nDispensa)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			dispensa.Titolo = titolo
			dispensa.Contenuto = testoDispensa
			db.UpdateDispensa(dispensa)
			return c.String(http.StatusOK, "Dispensa aggiornata")
		}
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}
func DeleteDispensaCorsoNumero(c echo.Context) error {

	idCorso := c.Param("id_corso")

	nDispensa, err := strconv.Atoi(c.Param("n_Dispensa"))

	if err != nil {
		ll.Error("Error to convert n. dispensa in int: %s", err.Error())
	}
	sub := getUserFromSession(c)                           // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/dispensa" // the resource that is going to be accessed.
	act := DELETE
	ok, err := casbin_enforcer.Enforce(sub, obj, act)
	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())
	if ok == true {
		idCorsoInt, err := strconv.Atoi(idCorso)
		if err != nil {
			ll.Error("Error to convert corso in int: %s", err.Error())
		}
		dispensa, err := db.GetDispensa(idCorsoInt, nDispensa)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			db.RemoveDispensa(dispensa)
			return c.String(http.StatusOK, "Dispensa rimossa")
		}
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}

