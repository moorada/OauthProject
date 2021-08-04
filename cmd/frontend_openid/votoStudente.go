package main

import (
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"net/http"
	"strconv"
)


func VotoStudente(c echo.Context) error {


	idCorso := c.Param(IDcorso)
	idStudente := c.Param(IDstudente)

	sub := getUserFromSession(c)                                     // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/voto/" + idStudente // the resource that is going to be accessed.
	act := GET
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())

	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		studente, err := db.GetStudenteFromMatricola(idStudente)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		voti, err := db.GetAPIVotoFromCorso(corso, studente)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, voti)
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}

func EditVotoStudente(c echo.Context) error {

	idCorso := c.Param(IDcorso)
	idStudente := c.Param(IDstudente)
	voto, err := strconv.Atoi(c.FormValue("voto"))

	if err != nil {
		ll.Error("Error to get voto: %s", err.Error())
	}

	sub := getUserFromSession(c)                                     // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/voto/" + idStudente // the resource that is going to be accessed.
	act := PUT
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())

	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		studente, err := db.GetStudenteFromMatricola(idStudente)
		if err != nil {
			return c.String(http.StatusNoContent, err.Error())
		}
		v, err := db.GetVoto(studente, corso)
		if err != nil {
			return c.String(http.StatusNoContent, err.Error())
		}
		v.Voto = voto
		db.UpdateVoto(v)
		return c.String(http.StatusOK, "Voto Modificato correttamente")
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}

}

func AddVotoStudente(c echo.Context) error {

	idCorso := c.Param(IDcorso)
	idStudente := c.Param(IDstudente)
	voto, err := strconv.Atoi(c.FormValue("voto"))
	if err != nil {
		ll.Error("Error to get voto: %s", err.Error())
	}

	sub := getUserFromSession(c)                                    // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/voto/" + idStudente // the resource that is going to be accessed.
	act := POST
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())

	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		studente, err := db.GetStudenteFromMatricola(idStudente)
		if err != nil {
			return c.String(http.StatusNoContent, err.Error())
		}
		voti, err := db.GetAPIVotoFromCorso(corso, studente)
		if voti.Voto == 0 {
			v := db.Voto{Voto: voto, Studente: studente.ID, Corso: corso.ID}
			db.AddVoto(v)
			return c.String(http.StatusCreated, "Voto inserito correttamente")
		}
		return c.String(http.StatusOK, "Voto già esistente, modificarlo o eliminarlo")
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}
}

func DeleteVotoStudente(c echo.Context) error {

	idCorso := c.Param(IDcorso)
	idStudente := c.Param(IDstudente)

	sub := getUserFromSession(c)                                  // the user that wants to access a resource.
	obj := "/corso/" + idCorso + "/voto/" + idStudente // the resource that is going to be accessed.
	act := DELETE
	ok, err := casbin_enforcer.Enforce(sub, obj, act)

	if err != nil {
		ll.Error(ErrCasbin, err.Error())
	}
	ll.Info("%s", casbin_enforcer.GetAllRoles())

	if ok == true {
		corso := db.GetCorsoFromId(idCorso)
		studente, err := db.GetStudenteFromMatricola(idStudente)
		if err != nil {
			return c.String(http.StatusNoContent, err.Error())
		}
		v, err := db.GetVoto(studente, corso)
		if err != nil {
			return c.String(http.StatusNoContent, err.Error())
		}
		err = db.RemoveVoto(v)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "Voto Eliminato correttamente")
		}
	} else {
		return c.String(http.StatusForbidden, "Non sei autorizzato")
	}
}
