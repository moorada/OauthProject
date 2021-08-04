#!/bin/bash

#Avviare Hydra
docker-compose up &

gnome-terminal gnome-terminal -- /bin/sh -c "go run cmd/authc/main.go"
gnome-terminal gnome-terminal -- /bin/sh -c "go run cmd/frontend_openid/main.go cmd/frontend_openid/studente.go cmd/frontend_openid/votoStudente.go cmd/frontend_openid/corso.go cmd/frontend_openid/dispensa.go cmd/frontend_openid/docente.go"
gnome-terminal gnome-terminal -- /bin/sh -c "cd ergo && ergo run -domain .test"