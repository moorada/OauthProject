# aggiungere client Openid
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
--scope offline,articles.write,articles.read,username.read

docker-compose exec hydra hydra clients update myclientOpenId \
--frontchannel-logout-callback http://localhost:2222/ \
--frontchannel-logout-session-required \
--endpoint http://127.0.0.1:4445/

docker-compose exec hydra hydra clients get myclientOpenId --endpoint http://127.0.0.1:4445/

docker-compose exec hydra hydra clients delete myclientOpenId --endpoint http://127.0.0.1:4445/
docker-compose exec hydra hydra clients delete myclientOauth --endpoint http://127.0.0.1:4445/
