package handler

import (
	"github.com/moorada/OauthProject/cmd/authc/repouser"

	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
)

type Handler struct {
	HydraAdmin hydraAdmin.ClientService
	UserRepo   repouser.Repository
}
