#!/bin/bash

#Avviare Hydra
docker-compose up &

#eseguire identity provider
go run cmd/authc/main.go &

#eseguire client
go run cmd/frontend_oauth/main.go &
go run cmd/frontend_openid/main.go cmd/frontend_openid/studente.go cmd/frontend_openid/votoStudente.go cmd/frontend_openid/corso.go cmd/frontend_openid/dispensa.go cmd/frontend_openid/docente.go &
#eseguire RS_oauth
go run cmd/RS_oauth/main.go &
#eseguire proxy
cd ergo && ergo run -domain .test &

