package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	ll "github.com/evilsocket/islazy/log"
	"github.com/moorada/OauthProject/pkg/render"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// Endpoint is OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  hydraURL + "oauth2/auth",
	TokenURL: hydraURL + "oauth2/token",
}

const hydraURL = "http://hydra.test/"

const port = "1111"
const RedirectURLOauth = "http://oaclient.test/callbacksOauth"
const ResourceServerURL = "http://resourceserver.test/"
const ResourceServerURLBack = "http://localhost:10000/"

var OAuthConf = &oauth2.Config{
	RedirectURL:  RedirectURLOauth,
	ClientID:     "myclientOauth",
	ClientSecret: "mysecretOauth",
	Scopes:       []string{"articles.write", "articles.read", "offline"},
	Endpoint:     Endpoint,
}

var stateStore = map[string]bool{}

var accessToken *oauth2.Token

func main() {
	initLog()
	defer ll.Close()
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:2000")

	ll.Important("%s", "FRONTEND OAUTH")

	t := &render.Template{}

	e := echo.New()
	e.Renderer = t
	//e.Use(middleware.Logger())

	e.GET("/", Homepage)
	e.GET("/callbacksOauth", CallbackOauth)
	e.GET("/cosavedo", CosaVedo)
	e.POST("/cosavedo", CosaVedoPost)

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Author  string `json:"Author"`
	Subject string `json:"Subject"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func CosaVedoPost(c echo.Context) error {

	myJsonString := c.FormValue("obj")
	var article Article
	err := json.Unmarshal([]byte(myJsonString), &article)
	if err != nil {
		ll.Error("errore nel unmarshallare, %v", err)
	}
	jsonValue, err := json.Marshal(article)
	if err != nil {
		ll.Error("errore nel marshallare, %v", err)
	}

	link := ResourceServerURL + article.Author + "/article/add?token=" + accessToken.AccessToken
	req, err := http.NewRequest("POST", link, bytes.NewBuffer(jsonValue))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		ll.Error("Err post da Oauth a RS: \n%v", err.Error())
		return c.String(http.StatusOK, "Errore"+err.Error())
	} else {
		requestDump, err := httputil.DumpRequest(resp.Request, true)
		if err != nil {
			ll.Error("Err dumpRequest: \n%v", err.Error())
		}
		ll.Debug("Richiesta post da frontend_oauth a RS: \n%v", string(requestDump))
		return c.String(resp.StatusCode, resp.Status)
	}
}

func Homepage(c echo.Context) error {

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta Homepage a frontend_oauth: \n%v", string(requestDump))

	// Generate random stateOauth
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	stateOauth := base64.StdEncoding.EncodeToString(b)

	stateStore[stateOauth] = true
	// Will return loginURLOauth,
	// for example: http://localhost:4444/oauth2/auth?client_id=myclient&prompt=consent&redirect_uri=http%3A%2F%2Fexample.com&response_type=code&scope=users.write+users.read+users.edit&state=XfFcFf7KL7ajzA2nBY%2F8%2FX3lVzZ6VZ0q7a8rM3kOfMM%3D
	loginURLOauth := OAuthConf.AuthCodeURL(stateOauth)

	return c.Render(http.StatusOK, "homepage_oauth.html", map[string]interface{}{
		"LoginURLOauth": loginURLOauth,
	})
}

func CosaVedo(c echo.Context) error {
	linkBoris := ResourceServerURL + "Boris/all?token=" + accessToken.AccessToken
	respBoris, err := http.Get(linkBoris)

	if err != nil {
		ll.Error("Err get: \n%v", err.Error())
	} else {
		requestDump, err := httputil.DumpRequest(respBoris.Request, true)
		if err != nil {
			ll.Error("Err dumpRequest: \n%v", err.Error())
		}
		ll.Debug("Richiesta get da frontend_oauth: \n%v", string(requestDump))
	}

	linkSaverio := ResourceServerURL + "Saverio/all?token=" + accessToken.AccessToken
	respSaverio, err := http.Get(linkSaverio)
	if err != nil {
		ll.Error("Err get: \n%v", err.Error())
	} else {
		requestDump, err := httputil.DumpRequest(respSaverio.Request, true)
		if err != nil {
			ll.Error("Err dumpRequest: \n%v", err.Error())
		}
		ll.Debug("Richiesta get da frontend_oauth: \n%v", string(requestDump))
	}
	headerBoris := fmt.Sprintf("%v", respBoris.StatusCode) + " " + respBoris.Status
	bodyBoris := []byte("{}")
	if respBoris.StatusCode == 200 {
		bodyBoris, err = io.ReadAll(respBoris.Body)
		defer respBoris.Body.Close()
		if err != nil {
			ll.Error("Err ReadAll: \n%v", err.Error())
		}
	}

	headerSaverio := fmt.Sprintf("%v", respSaverio.StatusCode) + " " + respSaverio.Status
	bodySaverio := []byte("{}")
	if respSaverio.StatusCode == 200 {
		bodySaverio, err = io.ReadAll(respSaverio.Body)
		defer respSaverio.Body.Close()
		if err != nil {
			ll.Error("Err ReadAll: \n%v", err.Error())
		}
	}
	defer respSaverio.Body.Close()
	if err != nil {
		ll.Error("Err readAll: \n%v", err.Error())
	}

	return c.Render(http.StatusOK, "cosavedo.html", map[string]interface{}{

		"LinkBoris":       linkBoris,
		"HeaderBoris":     headerBoris,
		"ArticlesBoris":   string(bodyBoris),
		"LinkSaverio":     linkSaverio,
		"HeaderSaverio":   headerSaverio,
		"ArticlesSaverio": string(bodySaverio),
		"ResetScopes":     ResourceServerURL + "resetscopes?token=" + accessToken.AccessToken,
	})

}

func CallbackOauth(c echo.Context) error {
	ctx := c.Request().Context()

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta CallbackOauth a frontend_oauth: \n%v", string(requestDump))

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	//scopes := c.QueryParam("scope")

	if code == "" {
		return c.String(http.StatusOK, "authorization code is empty")
	}

	// If state is exist

	if _, exist := stateStore[state]; !exist {
		return c.String(http.StatusOK, "state is generated by this Client")
	}
	delete(stateStore, state)

	// Exchange code for access token
	accessToken, err = OAuthConf.Exchange(ctx, code)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	//ll.Info("Access Token \n%v", *accessToken)

	return c.Render(http.StatusOK, "after_login.html", map[string]interface{}{
		"AccessToken": accessToken,
	})

}

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
