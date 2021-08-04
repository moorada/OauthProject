package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc"
	ll "github.com/evilsocket/islazy/log"
	echosession "github.com/moorada/OauthProject/pkg/session"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	casbin "github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/db"
	"github.com/moorada/OauthProject/pkg/render"
	"golang.org/x/oauth2"
)

// Endpoint is OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  hydraURL + "oauth2/auth",
	TokenURL: hydraURL + "oauth2/token",
}

const hydraURL = "http://hydra.test/"
const port = "2222"
const OIDCClientURLwithoutslash = "http://oidcclient.test"
const OIDCClientURL = OIDCClientURLwithoutslash + "/"
const RedirectURLOpenId = OIDCClientURL + "callbacksopenid"
const ErrNonSeiAutorizzato = "Non sei autorizzato ad accedere a questa risorsa"
const ErrCasbin = "Error enfoce casbin in get info corso: %s"

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var OpenIdConf = &oauth2.Config{
	RedirectURL:  RedirectURLOpenId,
	ClientID:     "myclientOpenId",
	ClientSecret: "mysecretOpenId",
	Scopes:       []string{"openid", "profile", "email"},
	Endpoint:     Endpoint,
}

var ApprovalForce oauth2.AuthCodeOption = oauth2.SetAuthURLParam("prompt", "none")
var ConsentedRequired oauth2.AuthCodeOption = oauth2.SetAuthURLParam("prompt", "consent")
var LoginRequired oauth2.AuthCodeOption = oauth2.SetAuthURLParam("prompt", "login")

var stateStore = map[string]bool{}

//authorization control with casbin
var casbin_enforcer casbin.Enforcer

const casbin_model_path = "./cmd/frontend_openid/auth_model.conf"
const casbin_policy_path = "./cmd/frontend_openid/policy.csv"
const openIdsession = "openId_session"

//https://www.dedgar.com/post/golang-echo-router-example

func main() {

	initLog()
	ll.Level = ll.INFO

	defer ll.Close()

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:2000")
	ll.Important("%s", "FRONTEND OPENID")

	db.InitDB("cmd/frontend_openid/database") //init db

	//init casbin enforcer
	init_casbin()

	t := &render.Template{}

	e := echo.New()
	e.Renderer = t
	//e.Use(middleware.Logger())

	e.Use(echosession.New())

	e.GET("/", Homepage)
	e.GET("/resetscopes", RemoveConsense)
	e.GET("/callbacksopenid", CallbackOpenId)
	e.GET("/logged", Logged)
	e.GET("/corso/:id_corso/info", InfoCorso)
	e.GET("/studente/:id_studente/info", InfoStudente)
	e.GET("/docente/:id_docente/info", InfoDocente)
	e.GET("/studente/:id_studente/voti", VotiStudente)
	e.GET("/corso/:id_studente/voti", VotiCorso)
	e.GET("/corso/:id_corso/voto/:id_studente", VotoStudente)
	e.PUT("/corso/:id_corso/voto/:id_studente", EditVotoStudente)
	e.POST("/corso/:id_corso/voto/:id_studente", AddVotoStudente)
	e.DELETE("/corso/:id_corso/voto/:id_studente", DeleteVotoStudente)
	e.GET("/corso/:id_corso/dispensa/", DispensaCorso)
	e.GET("/corso/:id_corso/dispensa/:n_Dispensa", DispensaCorsoNumero)
	e.PUT("/corso/:id_corso/dispensa/:n_Dispensa", EditDispensaCorsoNumero)
	e.DELETE("/corso/:id_corso/dispensa/:n_Dispensa", DeleteDispensaCorsoNumero)
	e.POST("/corso/:id_corso/dispensa/add", AddDispensaCorsoNumero)
	e.GET("/logout", Logout)

	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.String(http.StatusNotFound, "Pagina non trovata!")
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

func Homepage(c echo.Context) error {
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta homepage a FRONTEND_OPENID: \n%v", string(requestDump))

	stateOpenId := makeState()
	stateStore[stateOpenId] = true

	// Will return loginURLOauth,
	// for example: http://localhost:4444/oauth2/auth?client_id=myclient&prompt=consent&redirect_uri=http%3A%2F%2Fexample.com&response_type=code&scope=users.write+users.read+users.edit&state=XfFcFf7KL7ajzA2nBY%2F8%2FX3lVzZ6VZ0q7a8rM3kOfMM%3D
	loginURLOpeinId := OpenIdConf.AuthCodeURL(stateOpenId)

	return c.Render(http.StatusOK, "homepage_openid.html", map[string]interface{}{
		"LoginURLOpenId": loginURLOpeinId,
	})
}

func makeState() string {
	a := make([]byte, 32)
	_, err := rand.Read(a)
	if err != nil {
		ll.Fatal("%v", err)
	}
	return base64.StdEncoding.EncodeToString(a)
}

func validateOpenID(openId string) bool {

	ll.Info("Raw IDToken:%v\n", openId)
	j := hydraURL + "/.well-known/jwks.json"
	keys, err := http.Get(j)
	if err != nil {
		ll.Error("Err get: \n%v", err.Error())
		return false
	} else {
		requestDump, err := httputil.DumpRequest(keys.Request, true)
		if err != nil {
			ll.Error("Err dumpRequest: \n%v", err.Error())
		}
		ll.Debug("Richiesta get da frontend_oauth: \n%v", string(requestDump))
	}

	defer keys.Body.Close()
	bodykeys, err := io.ReadAll(keys.Body)
	if err != nil {
		ll.Error("Err readAll: \n%v", err.Error())
		return false
	}

	var jwksJSON json.RawMessage = bodykeys
	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.New(jwksJSON)
	if err != nil {
		ll.Error("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
		return false
	}

	// Parse the JWT.
	token, err := jwt.Parse(openId, jwks.KeyFuncLegacy)
	if err != nil {
		ll.Error("Failed to parse the JWT.\nError: %s", err.Error())
		return false
	}

	// Check if the token is valid.
	if !token.Valid {
		ll.Error("The token is not valid.")
		return false
	}

	claims := getClaims(openId)
	iss := fmt.Sprintf("%v", claims["iss"])
	aud := fmt.Sprintf("%v", claims["aud"])
	// TODO validation time
	if iss != hydraURL {
		ll.Error("Token invalid, iss isn't correct %v,%v", iss, hydraURL)
		return false
	}
	if aud != "["+OpenIdConf.ClientID+"]" {
		ll.Error("Token invalid, aud isn't correct %v,%v", aud, "["+OpenIdConf.ClientID+"]")
		return false
	}
	ll.Info("Token is valid")
	return true

}

func CallbackOpenId(c echo.Context) error {
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta CallbackOpenId a FRONTEND_OPENID: \n%v", string(requestDump))

	store := echosession.FromContext(c)
	_, ok := store.Get(openIdsession)
	if ok {
		return c.Redirect(http.StatusFound, OIDCClientURL+"logged")
	} else {
		ctx := c.Request().Context()

		code := c.QueryParam("code")
		state := c.QueryParam("state")

		if code == "" {
			return c.String(http.StatusOK, "authorization code is empty")
		}

		// If state is exist

		if _, exist := stateStore[state]; !exist {
			return c.String(http.StatusOK, "state is not generated by this Client")
		}

		delete(stateStore, state)

		// Exchange code for access token
		accessToken, err := OpenIdConf.Exchange(ctx, code)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		//converting jwt token

		tokenString := fmt.Sprintf("%v", accessToken.Extra("id_token"))
		if validateOpenID(tokenString) {
			claims := getClaims(tokenString)
			// ... error handling
			// do something with decoded claims
			for key, val := range claims {
				fmt.Printf("Key: %v, value: %v\n", key, val)
			}
			//	user_id := fmt.Sprintf("%v", claims["sub"])
			//	fmt.Println("---------------sub:" + user_id)

			store := echosession.FromContext(c)
			store.Set(openIdsession, tokenString)

			err = store.Save()
			if err != nil {
				return err
			}
			return c.Redirect(http.StatusFound, OIDCClientURL+"logged")
		} else {
			return c.String(http.StatusInternalServerError, "OpenId token non valido")
		}

	}
}

func getClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})
	return claims
}

func Logged(c echo.Context) error {
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta Logged a FRONTEND_OPENID: \n%v", string(requestDump))

	store := echosession.FromContext(c)
	s, ok := store.Get(openIdsession)
	if ok {
		claims := getClaims(fmt.Sprintf("%v", s))
		user_id := fmt.Sprintf("%v", claims["sub"])
		email := fmt.Sprintf("%v", claims["email"])
		name := fmt.Sprintf("%v", claims["name"])
		familyName := fmt.Sprintf("%v", claims["family_name"])
		ll.Important("email: %v,name: %v,familyname: %v", email, name, familyName)

		nome := name + " "+ familyName

		//logout_url := "http://localhost:2222/logout"

		return c.Render(http.StatusOK, "after_login_idToken.html", map[string]interface{}{
			"user_id":     user_id,
			"nome_utente": nome,
			"email": email,
		})
	} else {
		stateOpenId := makeState()
		stateStore[stateOpenId] = true
		loginURLOpeinId := OpenIdConf.AuthCodeURL(stateOpenId)
		return c.Redirect(http.StatusFound, loginURLOpeinId)
		//return c.Redirect(http.StatusFound, "http://hydra.test/oauth2/auth?id_token_hint="+fmt.Sprintf("%v", s))
	}
}

func Logout(c echo.Context) error {
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Important("Richiesta Logout a FRONTEND_OPENID: \n%v", string(requestDump))

	store := echosession.FromContext(c)
	s, _ := store.Get(openIdsession)
	ll.Important("%v", s)
	err = echosession.Destroy(c)
	if err != nil {
		ll.Error("Err to destroy echosession: \n%v", err.Error())
	}
	oauthLogout := hydraURL + "oauth2/sessions/logout?id_token_hint=" + fmt.Sprintf("%v", s) + "&post_logout_redirect_uri=" + OIDCClientURLwithoutslash
	//oauthLogout := "http://localhost:4444/oauth2/sessions/logout"
	err = c.Redirect(http.StatusFound, oauthLogout)
	if err != nil {
		ll.Error("Errore del redirect: %v", err)
	}
	return err
}

func RemoveConsense(c echo.Context) error {
	store := echosession.FromContext(c)
	s, _ := store.Get(openIdsession)
	claims := getClaims(fmt.Sprintf("%v", s))
	userId := fmt.Sprintf("%v", claims["sub"])
	err := removeConsense(userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Errore: "+ err.Error())
	} else {
		return c.String(http.StatusOK, "Scopes resettati")
	}
}

func removeConsense(subject string) error {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", "http://localhost:4445/oauth2/auth/sessions/consent?subject="+subject+"&client="+OpenIdConf.ClientID, nil)
	if err != nil {
		ll.Error("%v",err)
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		ll.Error("%v",err)
		return err
	}
	defer resp.Body.Close()

	//Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ll.Error("%v",err)
		return err
	}
	ll.Important("Risposta delete consenso: %v, %v",resp.Status, string(respBody))
	return nil

}

func init_casbin() {
	e, err_casbin1 := casbin.NewEnforcer(casbin_model_path, casbin_policy_path)
	if err_casbin1 != nil {
		ll.Error("Errore di casbin: %v", err_casbin1)
	}
	casbin_enforcer = *e
}

//creazione delle funzioni che simulano le chiamate rest

func initLog() {
	ll.Output = "/dev/stdout"
	ll.Level = ll.DEBUG
	ll.OnFatal = ll.NoneOnFatal
	ll.DateFormat = "06-Jan-02"
	ll.TimeFormat = "15:04:05"
	ll.DateTimeFormat = "2006-01-02 15:04:05"
	ll.Format = "{datetime} {level:color}{level:name}{reset} {message}"

	if err := ll.Open(); err != nil {
		panic(err)
	}
}

func getUserFromSession(c echo.Context) string {
	store := echosession.FromContext(c)
	s, _ := store.Get(openIdsession)
	claims := getClaims(fmt.Sprintf("%v", s))
	userId := fmt.Sprintf("%v", claims["sub"])
	return userId
}
