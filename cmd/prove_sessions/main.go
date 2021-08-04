package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	echosession "github.com/moorada/OauthProject/pkg/session"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("Secret"))))

	// Add the name "Steve" to the session
	e.GET("/login", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   0,
			HttpOnly: false,
			Secure:   true,
		}
		sess.Values["name"] = "Steve"
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusOK)
	})

	// Reply with the name saved in the session
	e.GET("/whoami", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, sess.Values["name"])
	})



	//e.GET("/", func(ctx echo.Context) error {
	//	if sessioned(ctx) {
	//		return ctx.String(http.StatusOK, "sessionato")
	//	} else {
	//		store := echosession.FromContext(ctx)
	//		store.Set("foo", "bar")
	//		err := store.Save()
	//		if err != nil {
	//			return err
	//		}
	//		return ctx.Redirect(302, "/foo")
	//	}
	//})
	//
	//e.GET("/foo", func(ctx echo.Context) error {
	//	store := echosession.FromContext(ctx)
	//	foo, ok := store.Get("foo")
	//	if !ok {
	//		return ctx.String(http.StatusNotFound, "not found")
	//	}
	//	return ctx.String(http.StatusOK, fmt.Sprintf("foo:%s", foo))
	//})
	//e.GET("/foo2", func(ctx echo.Context) error {
	//	return echosession.Destroy(ctx)
	//})

	e.Logger.Fatal(e.Start(":9999"))
}

func sessioned(ctx echo.Context) bool {
	store := echosession.FromContext(ctx)
	_, ok := store.Get("foo")
	return ok
}
