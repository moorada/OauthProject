# Project on OAUTH2 and Open ID Connect using Hydra

This repository is an educational project and shows how to implement:

* An Oauth 2.0 Server using hydra and an Identity Provider
* An Oauth 2.0 Client, which requires permissions to access resources of a resource owner and a Resource Server
* An Open ID Connect Client and REST server which requires LogIn via OIDC and implements a ORM database and an authorization server by the Casbin framework


## User guide

This guide is for Linux users.

### Prerequisites

1. [Install docker](https://docs.docker.com/engine/install/ubuntu/)
2. [Install Ergo Proxy](https://github.com/cristianoliveira/ergo)
   ```shell script
    $ curl -s https://raw.githubusercontent.com/cristianoliveira/ergo/master/install.sh | sh
    ```
3. Clone this repository and download dependencies `go mod download`

### Start

Before starting the examples, set up automatic proxy configuration with: `http://localhost:2000/proxy.pac`

To simplify, there are two scripts:

1. `./upOA.sh` for the example of Oauth 2.0. To access the client app [http://oaclient.test](http://oaclient.test)
2. `./upOIDC.sh`for the Open ID Connect example. To access the client app [http://oidcclient.test](http://oidcclient.test)

If you are running Hydra for the first time, you need to add the clients

1. `addClientOA.sh`
2. `addClientOIDC.sh`

Test Credentials:

For example 1:

* saverio@bergamaschi.com
* boris@moretti.com

For example 2:

* Students:
    * saverio@bergamaschi.com
    * boris@moretti.com
* Teachers:
    * olindo@pirozzi.com
    * valente@pinto.com

The password for all accounts is "password"

To remove all running processes you can run `./down`
To update HTML code `./bind`
