package handler

import (
	ll "github.com/evilsocket/islazy/log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ory/hydra-client-go/client/admin"
)

func (h Handler) GetLogout(c echo.Context) error {
	ctx := c.Request().Context()

	ll.Important("Hydra Logout")
	ll.Debug("Hydra Logout")
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta GetLogout a IDENTITY SERVER: \n%v", string(requestDump))

	logoutChallenge := strings.TrimSpace(c.QueryParam("logout_challenge"))
	if logoutChallenge == "" {
		return c.Render(http.StatusOK, "logout.html", map[string]interface{}{
			"ErrorTitle":   "Logout Challenge Is Not Exist!",
			"ErrorContent": "Logout Challenge Is Not Exist!",
		})
	}

	// Using Hydra Admin to get the logout challenge info
	logoutGetParam := admin.NewGetLogoutRequestParams()
	logoutGetParam.WithContext(ctx)
	logoutGetParam.SetLogoutChallenge(logoutChallenge)

	_, err = h.HydraAdmin.GetLogoutRequest(logoutGetParam)
	if err != nil {
		return c.Render(http.StatusOK, "logout.html", map[string]interface{}{
			"ErrorTitle":   "Failed When Get Logout Request Info",
			"ErrorContent": err.Error(),
		})
	}
	logoutAcceptParam := admin.NewAcceptLogoutRequestParams()
	logoutAcceptParam.WithContext(ctx)
	logoutAcceptParam.SetLogoutChallenge(logoutChallenge)

	respLogoutAccept, err := h.HydraAdmin.AcceptLogoutRequest(logoutAcceptParam)
	if err != nil {
		return c.Render(http.StatusOK, "logout.html", map[string]interface{}{
			"ErrorTitle":   "Cannot Accept Logout Request",
			"ErrorContent": err.Error(),
		})
	}

	// If success, it will redirect to postlogout page using
	return c.Redirect(http.StatusFound, *respLogoutAccept.GetPayload().RedirectTo)
}
