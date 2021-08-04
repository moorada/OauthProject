package db

import (
	"github.com/evilsocket/islazy/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type Studente struct {
	gorm.Model
	Matricola string  `gorm:"unique;not null"`
	Nome      string  `gorm:"index:idx_nc,unique"`
	Cognome   string  `gorm:"index:idx_nc,unique"`
	Corsi     []Corso `gorm:"many2many:corsi_seguiti;"`
}

type APIStudente struct {
	Matricola string `gorm:"unique;not null"`
	Nome      string `gorm:"index:idx_nc,unique"`
	Cognome   string `gorm:"index:idx_nc,unique"`
}

type Voto struct {
	gorm.Model
	Corso    uint `gorm:"index:idx_cs,unique"`
	Studente uint `gorm:"index:idx_cs,unique"`
	Voto     int
}

//type APIVoto struct {
//	Corso    uint `gorm:"index:idx_cs,unique"`
//	Studente uint `gorm:"index:idx_cs,unique"`
//	Voto     int
//}

type APIVoto struct {
	Voto      int
	NomeCorso string
	Matricola string
}

type Docente struct {
	gorm.Model
	Matricola string `gorm:"unique;not null"`
	Nome      string `gorm:"index:idx_ncp,unique"`
	Cognome   string `gorm:"index:idx_ncp,unique"`
}
type APIDocente struct {
	Matricola string `gorm:"unique;not null"`
	Nome      string `gorm:"index:idx_ncp,unique"`
	Cognome   string `gorm:"index:idx_ncp,unique"`
}

type Corso struct {
	gorm.Model
	NomeCorso    string `gorm:"unique;not null"`
	DocenteID    uint
	Informazioni string
	Dispensa     []Dispensa
}
type APICorso struct {
	NomeCorso    string `gorm:"unique;not null"`
	Informazioni string
}

type Dispensa struct {
	gorm.Model
	CorsoID   uint   `gorm:"index:idx_ncc,unique"`
	Capitolo  int    `gorm:"index:idx_ncc,unique"`
	Titolo    string `gorm:"unique;not null"`
	Contenuto string
}

type APIDispensa struct {
	CorsoID   uint
	Capitolo  int
	Titolo    string `gorm:"unique;not null"`
	Contenuto string
}

/*Initialize the DataBase*/
func InitDB(nameDB string) {
	var err error
	db, err = gorm.Open(sqlite.Open("./"+nameDB+".db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect database %s", err.Error())
	}
	// Migrate the schema
	db.AutoMigrate(&Studente{})
	db.AutoMigrate(&Docente{})
	db.AutoMigrate(&Corso{})
	db.AutoMigrate(&Dispensa{})
	db.AutoMigrate(&Voto{})

}

/*Add a Docente in the DataBase*/
func AddDocente(d Docente) {
	db.Create(&d)
}

func getVoti(s Studente) ([]Voto, error) {
	var voti []Voto
	err := db.Where(map[string]interface{}{"studente": s.ID}).Find(&voti).Error
	return voti, err
}

func GetAPIVotiFromStudente(s Studente) ([]APIVoto, error) {
	var voti []APIVoto
	err := db.Table("votos").Select("votos.voto, corsos.nome_corso, studentes.matricola").Where(map[string]interface{}{"studente": s.ID}).Joins("left join studentes on votos.studente = studentes.id").Joins("left join corsos on votos.corso = corsos.id").Scan(&voti).Error
	return voti, err
}

func GetAPIVotiFromCorso(c Corso) ([]APIVoto, error) {
	var voti []APIVoto
	err := db.Table("votos").Select("votos.voto, corsos.nome_corso, studentes.matricola").Where(map[string]interface{}{"corso": c.ID}).Joins("left join studentes on votos.studente = studentes.id").Joins("left join corsos on votos.corso = corsos.id").Scan(&voti).Error
	return voti, err
}

func GetAPIVotoFromCorso(c Corso, s Studente) (APIVoto, error) {
	var voti []APIVoto
	var voto APIVoto
	err := db.Table("votos").Select("votos.voto, corsos.nome_corso, studentes.matricola").Where(map[string]interface{}{"corso": c.ID, "studente": s.ID}).Joins("left join studentes on votos.studente = studentes.id").Joins("left join corsos on votos.corso = corsos.id").Scan(&voti).First(&voto).Error

	return voto, err

}
func getDispense(c Corso) ([]Dispensa, error) {
	var dispensa []Dispensa
	err := db.Where(map[string]interface{}{"corso_id": c.ID}).Find(&dispensa).Error
	return dispensa, err
}

func GetDispensa(cId int, mId int) (Dispensa, error) {
	var dispensa Dispensa
	err := db.Where(map[string]interface{}{"corso_id": cId, "capitolo": mId}).First(&dispensa).Error
	return dispensa, err
}

func GetAPIDispensa(cId int, mId int) (APIDispensa, error) {
	var dispensa APIDispensa
	err := db.Model(&Dispensa{}).Where(map[string]interface{}{"corso_id": cId, "capitolo": mId}).First(&dispensa).Error
	return dispensa, err
}

func UpdateDispensa(dispensa Dispensa) {
	db.Save(dispensa)
}

func GetCorsiSeguiti(s Studente) ([]APICorso, error) {
	var corsi []APICorso
	err := db.Table("corsi_seguiti").Select("corsos.nome_corso, corsos.informazioni").Where(map[string]interface{}{"studente_id": s.ID}).Joins("left join studentes on corsi_seguiti.studente_id = studentes.id").Joins("left join corsos on corsi_seguiti.corso_id = corsos.id").Scan(&corsi).Error
	return corsi, err
}

func GetAPIDispense(c Corso) ([]APIDispensa, error) {
	var dispensa []APIDispensa
	err := db.Model(&Dispensa{}).Where(map[string]interface{}{"corso_id": c.ID}).Find(&dispensa).Error
	return dispensa, err
}
func GetDispense(c Corso) ([]Dispensa, error) {
	var dispensa []Dispensa
	err := db.Model(&Dispensa{}).Where(map[string]interface{}{"corso_id": c.ID}).Find(&dispensa).Error
	return dispensa, err
}

func AddStudente(s Studente) {
	db.Create(&s)
}

func AddVoto(voto Voto) {
	db.Create(&voto)
}

/*Add a Docente in the DataBase*/
func AddCorso(d Corso) {
	db.Create(&d)
}
func AddDispensa(corso string, titolo string, testo string){
	c := GetCorsoFromId(corso)
	dispense, err := getDispense(c)
	if err != nil {
		panic(err)
	}
	maxCap := 0
	for _, v := range dispense {
		if v.Capitolo > maxCap{
			maxCap = v.Capitolo
		}
	}
	m := Dispensa{Titolo: titolo, Capitolo: maxCap+1, Contenuto: testo, CorsoID: c.ID}
	db.Create(&m)
}

func GetVoto(studente Studente, corso Corso) (Voto, error) {
	var voto Voto
	err := db.Where(map[string]interface{}{"corso": corso.ID, "studente": studente.ID}).First(&voto).Error
	return voto, err
}

func UpdateVoto(voto Voto) {
	db.Save(voto)
}

func GetCorsoFromName(nomeCorso string) (Corso, error) {
	var corso Corso
	err := db.Where(map[string]interface{}{"nome_corso": nomeCorso}).First(&corso).Error
	return corso, err
}

func GetCorsoFromId(id string) Corso {
	var corso Corso
	db.First(&corso, id)
	return corso
}
func GetAPICorsoFromId(idCorso string) (APICorso, error) {
	var corso APICorso
	err := db.Model(&Corso{}).Where(map[string]interface{}{"id": idCorso}).First(&corso).Error
	return corso, err
}

func GetAPICorso(nomeCorso string) (APICorso, error) {
	var corso APICorso
	err := db.Model(&Corso{}).Where(map[string]interface{}{"nome_corso": nomeCorso}).First(&corso).Error
	return corso, err
}

/*Get a Docente from the DataBase*/
func GetDocente(idDocente string) (Docente, error) {
	var docente Docente
	err := db.First(&docente, idDocente).Error
	return docente, err
}

func GetAPIDocenteFromMatricola(matricola string) (APIDocente, error) {
	var docente APIDocente
	err := db.Model(&Docente{}).Where(map[string]interface{}{"matricola": matricola}).Scan(&docente).Error
	return docente, err
}

func GetStudenteFromNome(cognome string, nome string) (Studente, error) {
	var studente Studente
	err := db.Where(map[string]interface{}{"nome": nome, "cognome": cognome}).First(&studente).Error
	return studente, err
}

func GetStudenteFromMatricola(matricola string) (Studente, error) {
	var studente Studente
	err := db.Preload("Corsi").First(&studente, "matricola = ?", matricola).Error
	return studente, err
}

func GetAPIStudenteFromMatricola(matricola string) (APIStudente, error) {
	var studente APIStudente
	err := db.Model(&Studente{}).Where(map[string]interface{}{"matricola": matricola}).Scan(&studente).Error
	return studente, err
}

func GetDocenteFromNome(cognome string, nome string) (Docente, error) {
	var docente Docente
	err := db.Where(map[string]interface{}{"nome": nome, "cognome": cognome}).First(&docente).Error
	return docente, err
}

func RemoveDocente(idDocente int) error {
	var docente Docente
	err := db.First(&docente, idDocente).Delete(&docente).Error
	return err
}

func RemoveVoto(voto Voto) error {
	err := db.Unscoped().Delete(&voto).Error
	return err
}

func RemoveDispensa(dispensa Dispensa) error {
	err := db.Unscoped().Delete(&dispensa).Error
	return err
}
