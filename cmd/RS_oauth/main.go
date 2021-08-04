package main

import (
	"bytes"
	"encoding/json"
	ll "github.com/evilsocket/islazy/log"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/render"
	hydra "github.com/ory/hydra-client-go/client"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const hydraURL = "http://hydra.test/"

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Author  string `json:"Author"`
	Subject string `json:"Subject"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Introspection struct {
	Active    bool          `json:"active"`
	Scope     string        `json:"scope"`
	ClientID  string        `json:"client_id"`
	Sub       string        `json:"sub"`
	Exp       int           `json:"exp"`
	Iat       int           `json:"iat"`
	Nbf       int           `json:"nbf"`
	Aud       []interface{} `json:"aud"`
	Iss       string        `json:"iss"`
	TokenType string        `json:"token_type"`
	TokenUse  string        `json:"token_use"`
}

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

const port = "10000"
const PATH = "./cmd/RS_oauth/articles.json"

var autore_matricola = map[string]string{
	"Boris":   "S960228",
	"Saverio": "S960483",
}

func introspection(token string) Introspection {
	headers := map[string][]string{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Accept":       []string{"application/json"},
	}

	var body []byte
	body = []byte("token=" + token)

	req, err := http.NewRequest("POST", "http://localhost:4445/oauth2/introspect", bytes.NewBuffer(body))
	req.Header = headers
	if err != nil {
		ll.Error("%v", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("INTROSPECTION Richiesta da RS a Hydra: \n%v", string(requestDump))

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ll.Error("%v", err.Error())
	}
	var m = Introspection{}
	err = json.Unmarshal(resp_body, &m)
	if err != nil {
		ll.Error("%v", err.Error())
	}
	return m
}

var ArticlesSlice Articles

type Articles struct {
	Articles []Article `json:"articles"`
}

var counter int

func HomePage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the HomePage!")
}

func ReturnAllArticles(c echo.Context) error {
	token := c.QueryParam("token")
	m := introspection(token)
	autore := c.Param("author")

	if autore_matricola[autore] == m.Sub {
		var ArticlesMapAuthor = map[int]Article{}
		if strings.Contains(m.Scope, "articles.read") {
			for i, article := range ArticlesSlice.Articles {
				if article.Author == autore && article.Subject == m.Sub {
					if article.Subject == m.Sub {
						ArticlesMapAuthor[i] = article
					}
				}
			}
			if len(ArticlesMapAuthor) > 0 {
				return c.JSON(http.StatusOK, ArticlesMapAuthor)
			} else {
				return c.String(http.StatusNotFound, "Nessun articolo per questo autore")
			}
		} else {
			return c.String(http.StatusForbidden, "Questo access token non ha lo scope articles.read")
		}
	} else {
		return c.String(http.StatusForbidden, "Questo access token non corrisponde a questo autore")
	}

}

func ResetScopes(c echo.Context) error {
	token := c.QueryParam("token")
	m := introspection(token)
	err := removeConsense(m.Sub, m.ClientID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Errore: "+ err.Error())
	} else {
		return c.String(http.StatusOK, "Scopes resettati")
	}
}
func removeConsense(subject string, clientId string) error {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", "http://localhost:4445/oauth2/auth/sessions/consent?subject="+subject+"&client="+clientId, nil)
	if err != nil {
		ll.Error("%v", err)
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		ll.Error("%v", err)
		return err
	}
	defer resp.Body.Close()

	//Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ll.Error("%v", err)
		return err
	}
	ll.Important("Risposta delete consenso: %v, %v", resp.Status, string(respBody))
	return nil
}

func CreateNewArticle(c echo.Context) error {

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta POST da frontend_oauth a RS: \n%v", string(requestDump))

	//	var myJsonString string
	var article Article
	err = json.NewDecoder(c.Request().Body).Decode(&article)
	if err != nil {
		return err
	}

	token := c.QueryParam("token")
	m := introspection(token)
	autore := c.Param("author")
	if strings.Contains(m.Scope, "articles.write") {
		if article.Author == autore && article.Subject == m.Sub {
			addArticle(article)
			ll.Debug("Il resource server ha aggiunto l'articolo:\n%v", article)
			return c.String(http.StatusCreated, "Articolo aggiunto")
		} else {
			ll.Debug("Il resource server non ha aggiunto l'articolo perchè l'access token non corrisponde al proprietario della risorsa")
			return c.String(http.StatusForbidden, "Non puoi perchè l'access token appartiene a un altro")
		}
	} else {
		ll.Debug("Il resource server non ha aggiunto l'articolo perchè non è presente lo scope articles.write per questo access token")
		return c.String(http.StatusForbidden, "Non puoi per lo scope")
	}
}

func init() {
	counter = 0
	jsonFile, err := os.Open(PATH)
	// if we os.Open returns an error then handle it
	if err != nil {
		ll.Error("Error to open articles.json, %v", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var articles Articles
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &articles)
	for i, s := range articles.Articles {
		s.Id = strconv.Itoa(i)
		addArticle(s)
	}

}

func addArticle(a Article) {
	a.Id = strconv.Itoa(counter)
	ArticlesSlice.Articles = append(ArticlesSlice.Articles, a)
	counter++

	//jsonFile, err := os.Open(PATH)

	b, err := json.Marshal(ArticlesSlice)
	if err != nil {
		ll.Error("Error to open articles.json, %v", err)
	}
	_ = ioutil.WriteFile(PATH, b, 0644)

}

func main() {
	initLog()
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:2000")
	defer ll.Close()
	ll.Important("%s", "RESOURCE SERVER")
	t := &render.Template{}

	e := echo.New()
	e.Renderer = t
	//e.Use(middleware.Logger())

	//e.GET(("/", HomePage)
	e.GET("/:author/all", ReturnAllArticles)
	e.GET("/resetscopes", ResetScopes)
	//e.GET("/{author}/article/{id}", ReturnSingleArticle)

	e.POST("/:author/article/add", CreateNewArticle)
	//e.DELETE("/{author}/article/{id}", DeleteArticle)

	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.String(http.StatusNotFound, "Pagina non trovata!")
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}

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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
