package main

import (
	"fmt"
	ll "github.com/evilsocket/islazy/log"
	"github.com/moorada/OauthProject/pkg/db"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const lorem string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

var popoDocCorsi = [][]string{
	[]string{"Olindo", "Pirozzi", "D6002", "Crittografia"},
	[]string{"Valente", "Pinto", "D6003", "Architetture orientate ai servizi"},
	[]string{"Emanuele", "Baresi", "D6004", "Architetture orientate ai servizi"},
	[]string{"Giovanni", "Moretti", "D6005", "Architetture orientate ai servizi"},
	[]string{"Prisco", "Monaldo", "D6006", "Sistemi biometrici"},
	[]string{"Prisco", "Monaldo", "D6006", "Intelligenza Artificiale"},
	[]string{"Arcangela", "Calabrese", "D6007", "Privatezza e protezione dei dati"},
}
var popoStuVoti = [][]string{
	[]string{"Saverio", "Bergamaschi", "S960483", "Crittografia,Architetture orientate ai servizi,Sistemi biometrici,Privatezza e protezione dei dati"},
	[]string{"Boris", "Moretti", "S960228", "Crittografia,Architetture orientate ai servizi"},
	[]string{"Rossi", "Francesco", "S960312", "Crittografia,Sistemi biometrici,Privatezza e protezione dei dati"},
	[]string{"Verdi", "Antonio", "S942222", "Crittografia,Architetture orientate ai servizi,Privatezza e protezione dei dati"},
	[]string{"Bianchi", "Maria", "S911111", "Privatezza e protezione dei dati,Sistemi biometrici"},
}

func populeDB() {
	for _, v := range popoDocCorsi {
		docente := db.Docente{Cognome: v[0], Nome: v[1], Matricola: v[2]}
		_, err := db.GetDocenteFromNome(v[0], v[1])
		if err != nil {
			db.AddDocente(docente)
		}
		docente, err = db.GetDocenteFromNome(v[0], v[1])
		if err != nil {
			panic(err)
		}
		ll.Important("Docente id: %v", docente.ID)
		corso := db.Corso{NomeCorso: v[3], DocenteID: docente.ID, Informazioni: "Corso tenuto da " + docente.Nome + " " + docente.Cognome + ". Questo corso si chiama " + v[3] + lorem}
		db.AddCorso(corso)
		corso, err = db.GetCorsoFromName(v[3])
		if err != nil {
			ll.Error("%s", err)
		}
		c := fmt.Sprint(corso.ID)
		db.AddDispensa(c, "TitoloMateriale 1"+" "+corso.NomeCorso, "Dispensa 1"+" "+corso.NomeCorso+lorem)
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(5) + 2
		for i := 2; i < x; i++ {
			db.AddDispensa(c, "TitoloMateriale "+strconv.Itoa(i)+" "+corso.NomeCorso, "Dispensa "+strconv.Itoa(i)+" "+corso.NomeCorso+lorem)
		}
	}

	for i, v := range popoStuVoti {
		ll.Info("Ciclo n %v", i)
		corsiNomi := strings.Split(v[3], ",")
		var corsiDB []db.Corso
		var votiDB []db.Voto
		for _, c := range corsiNomi {
			co, err := db.GetCorsoFromName(c)
			if err != nil {
				ll.Error("%s", err)
			}
			corsiDB = append(corsiDB, co)

			rand.Seed(time.Now().UnixNano())
			x := rand.Intn(31)
			if x > 18 {
				votiDB = append(votiDB, db.Voto{Voto: x, Corso: co.ID})
			}
		}

		studente := db.Studente{Cognome: v[0], Nome: v[1], Matricola: v[2], Corsi: corsiDB}
		db.AddStudente(studente)
		studente, err := db.GetStudenteFromNome(v[0], v[1])

		if err != nil {
			ll.Error("%s", err)
		}
		for _, v := range votiDB {
			v.Studente = studente.ID
			db.AddVoto(v)
		}
	}
}

func main() {
	db.InitDB("database")
	populeDB()
	//db.AddVoto(db.Voto{Voto: 28, Corso: 2, Studente: 1})
	//s, err := db.GetStudenteFromMatricola("S960228")
	//if err != nil {
	//	log.Error(err.Error())
	//}-
	//
	//fmt.Println(s)
	//voti2, err := db.GetAPIVotiFromStudente(s)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, v := range voti2 {
	//	fmt.Println("Voto", v)
	//}

	//fmt.Println(s)

	//corsi, err := db.GetCorsiSeguiti(s)
	//if err != nil {
	//	panic(err)
	//}
	//for _, v := range corsi {
	//	fmt.Println("Corso", v)
	//}

	//c, err := db.GetCorso("Crittografia")
	//if err != nil {
	//	panic(err)
	//}
	//voti, err := db.GetAPIVotiFromCorso(c)
	//if err != nil {
	//	panic(err)
	//}
	//for _, v := range voti {
	//	fmt.Println("Voto", v)
	//}
	//
	//m, err := db.GetDispensa(4, 8)
	//if err != nil {
	//	panic(err)
	//}

	//db.RemoveDispensa(m)

	//c, err = db.GetCorsoFromDB("Crittografia")
	//if err != nil {
	//	panic(err)
	//}
	//materiali, err := db.GetAPIDispense(c)
	//if err != nil {
	//	panic(err)
	//}
	//for _, v := range materiali {
	//	fmt.Println("Dispensa", v)
	//}

}
