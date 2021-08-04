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

func (h Handler) PostLogout(c echo.Context) error {
	ctx := c.Request().Context()

	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		ll.Error("Err dumpRequest: \n%v", err.Error())
	}
	ll.Debug("Richiesta PostLogout a IDENTITY SERVER: \n%v", string(requestDump))


	formData := struct {
		ConsentChallenge string   `validate:"required"`
		GrantScope       []string `validate:"required"`
	}{
		ConsentChallenge: c.FormValue("consent_challenge"),
		GrantScope:       c.Request().Form["grant_scope"],
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

	// If a user has granted this application the requested scope, hydra will tell us to not show the UI.

	// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong

	id := consentGetResp.GetPayload().Subject
	user, err := h.UserRepo.GetUserById(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	consentAcceptBody := &models.AcceptConsentRequest{
		GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		GrantScope:               formData.GrantScope,
		Session: &models.ConsentRequestSession{
			// Sets session data for the OpenID Connect ID token.
			IDToken: map[string]interface{}{
				"extra_vars": map[string]interface{}{
					"id":      user.ID,
					"email":   user.Email,
					"nome":    user.Nome,
					"cognome": user.Cognome,
				},
			},
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
