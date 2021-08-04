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

func (h Handler) PostConsent(c echo.Context) error {
	ctx := c.Request().Context()

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta PostConsent a IDENTITY SERVER: \n%v", string(requestDump))

	formData := struct {
		ConsentChallenge string   `validate:"required"`
		GrantScope       []string `validate:"required"`
		RememberMe       string   `validate:"required"`
	}{
		ConsentChallenge: c.FormValue("consent_challenge"),
		GrantScope:       c.Request().Form["grant_scope"],
		RememberMe:       c.FormValue("remember_me"),
	}

	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.WithContext(ctx)
	consentGetParams.SetConsentChallenge(formData.ConsentChallenge)

	consentGetResp, err := h.HydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		// if error, redirects to ...
		str := fmt.Sprint("error GetConsentRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}
	var rememberMe = formData.RememberMe == "true"

	id := consentGetResp.GetPayload().Subject
	user, err := h.UserRepo.GetUserById(c.Request().Context(), id)

	idtoken := map[string]interface{}{}

	scopes := formData.GrantScope
	if contains(scopes, "email") {
		idtoken["email"] = user.Email
	}
	if contains(scopes, "profile") {
		idtoken["name"] = user.Nome
		idtoken["family_name"] = user.Cognome
	}

	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	consentAcceptBody := &models.AcceptConsentRequest{
		GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		GrantScope:               formData.GrantScope,
		Remember:                 rememberMe,
		Session: &models.ConsentRequestSession{
			// Sets session data for the OpenID Connect ID token.
			IDToken: idtoken,
		},
	}

	consentAcceptParams := admin.NewAcceptConsentRequestParams()
	consentAcceptParams.WithContext(ctx)
	consentAcceptParams.SetConsentChallenge(formData.ConsentChallenge)
	consentAcceptParams.WithBody(consentAcceptBody)

	consentAcceptResp, err := h.HydraAdmin.AcceptConsentRequest(consentAcceptParams)
	if err != nil {
		str := fmt.Sprint("error AcceptConsentRequest", err.Error())
		return c.String(http.StatusUnprocessableEntity, str)
	}

	return c.Redirect(http.StatusFound, *consentAcceptResp.GetPayload().RedirectTo)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
