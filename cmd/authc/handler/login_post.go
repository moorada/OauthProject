package handler

import (
	"fmt"
	ll "github.com/evilsocket/islazy/log"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func (h Handler) PostLogin(c echo.Context) error {
	ctx := c.Request().Context()

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta PostLogin a IDENTITY SERVER: \n%v",string(requestDump))

	formData := struct {
		LoginChallenge string `validate:"required"`
		Email          string `validate:"required"`
		Password       string `validate:"required"`
		RememberMe     string `validate:"required"`
	}{
		LoginChallenge: c.FormValue("login_challenge"),
		Email:          c.FormValue("email"),
		Password:       c.FormValue("password"),
		RememberMe:     c.FormValue("remember_me"),
	}

	// TODO validation

	var rememberMe = formData.RememberMe == "true"

	user, err := h.UserRepo.GetUserByEmail(c.Request().Context(), formData.Email)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	if user.Password != formData.Password {
		return c.String(http.StatusNotFound, "Wrong username and password")
	}

	//username e Psw corretti
	// Using Hydra Admin to accept login request!
	loginGetParam := admin.NewGetLoginRequestParams()
	loginGetParam.SetLoginChallenge(formData.LoginChallenge)

	_, err = h.HydraAdmin.GetLoginRequest(loginGetParam)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error GetLoginRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	subject := fmt.Sprint(user.ID)
	loginAcceptParam := admin.NewAcceptLoginRequestParams()
	loginAcceptParam.WithContext(ctx)
	loginAcceptParam.SetLoginChallenge(formData.LoginChallenge)
	loginAcceptParam.SetBody(&models.AcceptLoginRequest{
		Subject:  &subject,
		Remember: rememberMe,
	})

	respLoginAccept, err := h.HydraAdmin.AcceptLoginRequest(loginAcceptParam)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error AcceptLoginRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	// If success, it will redirect to consent page using handler GetConsent
	// It then show the consent form
	//ll.Important("RedirectTo", *respLoginAccept.GetPayload().RedirectTo )
	return c.Redirect(http.StatusFound, *respLoginAccept.GetPayload().RedirectTo)
}
