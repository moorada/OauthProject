#!/bin/bash
docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445/ \
--name MyAppOpenId \
--id myclientOpenId \
--secret mysecretOpenId \
--grant-types authorization_code,refresh_token \
--callbacks http://oidcclient.test/callbacksopenid,http://oidcclient.test,http://identityprovider.test/authentication/logout \
--response-types code,id_token \
--token-endpoint-auth-method client_secret_post \
--post-logout-callbacks http://oidcclient.test \
--scope openid,email,profile
