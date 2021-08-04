package handler

import (
	ll "github.com/evilsocket/islazy/log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func (h Handler) GetLogin(c echo.Context) error {
	ctx := c.Request().Context()
	initLog()

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta GetLogin a IDENTITY SERVER: \n%v",string(requestDump))

	loginChallenge := strings.TrimSpace(c.QueryParam("login_challenge"))
	if loginChallenge == "" {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"ErrorTitle":   "Login Challenge Is Not Exist!",
			"ErrorContent": "Login Challenge Is Not Exist!",
		})
	}

	// Using Hydra Admin to get the login challenge info
	loginGetParam := admin.NewGetLoginRequestParams()
	loginGetParam.WithContext(ctx)
	loginGetParam.SetLoginChallenge(loginChallenge)

	respLoginGet, err := h.HydraAdmin.GetLoginRequest(loginGetParam)
	//ll.Important("respLoginGet:\n %v", respLoginGet)
	if err != nil {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"ErrorTitle":   "Failed When Get Login Request Info",
			"ErrorContent": err.Error(),
		})
	}

	skip := false
	if respLoginGet.GetPayload().Skip != nil {
		skip = *respLoginGet.GetPayload().Skip
	}

	// If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate
	// the user.
	if skip {
		// Using Hydra Admin to accept login request!
		loginAcceptParam := admin.NewAcceptLoginRequestParams()
		loginAcceptParam.WithContext(ctx)
		loginAcceptParam.SetLoginChallenge(loginChallenge)
		loginAcceptParam.SetBody(&models.AcceptLoginRequest{
			Subject: respLoginGet.GetPayload().Subject,

		})

		respLoginAccept, err := h.HydraAdmin.AcceptLoginRequest(loginAcceptParam)
		if err != nil {
			return c.Render(http.StatusOK, "login.html", map[string]interface{}{
				"ErrorTitle":   "Cannot Accept Login Request",
				"ErrorContent": err.Error(),
			})
		}

		// If success, it will redirect to consent page using handler GetConsent
		// It then show the consent form
		return c.Redirect(http.StatusFound, *respLoginAccept.GetPayload().RedirectTo)
	}

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"LoginChallenge": loginChallenge,
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
