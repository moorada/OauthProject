package main

import (
	"fmt"

	casbin "github.com/casbin/casbin/v2"
)

func main() {
	e, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(e.GetAllRoles())

	sub := "Saverio"       // the user that wants to access a resource.
	obj := "/Saverio/voti" // the resource that is going to be accessed.
	act := "GET"          // the operation that the user performs on the resource.

	/*GET /id_studente/* -> può accedere solo studente #ed eventuale segreteria

	  GET /id_corso/informazioni -> studenti e prof
	  GET /id__corso/materiale -> solo prof tit e studenti iscritti
	  GET /id_corso/voti -> solo prof tit
	  POST /id_corso/informazioni -> solo prof tit
	  PUT|POST|GET|DELETE /id_corso/materiale/* -> solo prof tit
	  POST|PUT|DELETE|GET /id_corso/voto/* -> solo prof tit

	  GET /informazioni/docenti/* -> solo professori
	  GET /informazioni/studenti/* -> solo professori e studenti
	  GET /informazioni/esterni/* -> tutti


	  #POST|GET|PUT|DELETE /informazioni/* -> segreteria
	*/

}
