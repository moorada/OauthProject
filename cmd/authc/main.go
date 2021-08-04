package main

import (
	"fmt"
	ll "github.com/evilsocket/islazy/log"
	"github.com/moorada/OauthProject/cmd/authc/handler"
	"github.com/moorada/OauthProject/cmd/authc/repouser"
	"github.com/moorada/OauthProject/pkg/template"
	"github.com/moorada/OauthProject/views"
	"io"
	"net/url"

	"github.com/labstack/echo/v4"
	hydra "github.com/ory/hydra-client-go/client"
)

var (
	adminURL, _ = url.Parse("http://localhost:4445")
	hydraClient = hydra.NewHTTPClientWithConfig(nil,
		&hydra.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
)

var userInfo = []repouser.UserInfo{
	{
		ID:       "S960483",
		Email:    "saverio@bergamaschi.com",
		Nome:     "Saverio",
		Cognome:  "Bergamaschi",
		Password: "password",
	},
	{
		ID:       "S960228",
		Email:    "boris@moretti.com",
		Nome:     "Boris",
		Cognome:  "Moretti",
		Password: "password",
	},
	{
		ID:       "D6002",
		Email:    "olindo@pirozzi.com",
		Nome:     "Olindo",
		Cognome:  "Pirozzi",
		Password: "password",
	},
	{
		ID:       "D6003",
		Email:    "valente@pinto.com",
		Nome:     "Valente",
		Cognome:  "Pinto",
		Password: "password",
	},
}

func main() {

	initLog()
	defer ll.Close()
	//fmt.Println("AUTHENTICATION")
	ll.Important("%s", "IDENTITY SERVER")

	controller := handler.Handler{
		HydraAdmin: hydraClient.Admin,
		UserRepo:   repouser.NewMemory(userInfo),
	}

	t := &Template{}
	var e = echo.New()
	e.Renderer = t


	e.GET("/authentication/login", controller.GetLogin)
	e.POST("/authentication/login", controller.PostLogin)
	e.GET("/authentication/consent", controller.GetConsent)
	e.POST("/authentication/consent", controller.PostConsent)
	e.GET("/authentication/logout", controller.GetLogout)
	e.POST("/authentication/logout", controller.PostLogout)

	if err := e.Start(":8000"); err != nil {
		ll.Fatal(err.Error())
	}

}

type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmpl := template.New("", &template.BinData{
		Asset:      views.Asset,
		AssetDir:   views.AssetDir,
		AssetNames: views.AssetNames,
	})

	tpl, err := tmpl.Parse(fmt.Sprintf("views/authc/%s", name))
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
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
