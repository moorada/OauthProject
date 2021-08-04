#!/bin/bash

docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445/ \
--name MyAppOauth \
--id myclientOauth \
--secret mysecretOauth \
--grant-types authorization_code,refresh_token \
--response-types code,id_token \
--callbacks http://oaclient.test/callbacksOauth,http://oaclient.test \
--token-endpoint-auth-method client_secret_post \
--post-logout-callbacks http://oaclient.test \
--scope offline,articles.write,articles.read
