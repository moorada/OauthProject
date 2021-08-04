#!/bin/bash

docker-compose up &
gnome-terminal gnome-terminal -- /bin/sh -c "go run cmd/authc/main.go"
gnome-terminal gnome-terminal -- /bin/sh -c "go run cmd/frontend_oauth/main.go"
gnome-terminal gnome-terminal -- /bin/sh -c "go run cmd/RS_oauth/main.go"
gnome-terminal gnome-terminal -- /bin/sh -c "cd ergo && ergo run -domain .test"
