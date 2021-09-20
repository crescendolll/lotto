package lottologic

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Nutzer struct {
	Benutzername string
	Pw_hash      string
	Ist_spieler  bool
}

type Tipp struct {
	Id      int64
	Datum   time.Time
	Ziehung string
}

type Ziehung struct {
	Datum   time.Time
	Ziehung string
}

type Auszahlung struct {
	Datum      time.Time
	Klasse     int8
	Budget     float64
	Auszahlung float64
}

func IstNameVerfuegbar(database *sql.DB, benutzername string) error {

	var ist_verfuegbar bool

	var anzahl_vorkommen int

	var belegtError error

	rows, err := database.Query("SELECT COUNT(*) as anzahl FROM nutzer WHERE benutzername = ?", benutzername)
	fmt.Printf("SELECT COUNT(*) as anzahl FROM nutzer WHERE benutzername = %s\n", benutzername)

	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			rows.Scan(&anzahl_vorkommen)
		}
		ist_verfuegbar = (anzahl_vorkommen == 0)
	}

	if !ist_verfuegbar {
		belegtError = errors.New("Name " + benutzername + " vergeben")
	}

	return belegtError
}

func ErzeugeZufallsziehung() string {

	var ziehung string

	ziehung = ""

	rand.Seed(time.Now().UnixNano())

	var zufallszahl int
	var ziehungszahlen []int

	ziehungszahlen = []int{}

	for index := 1; index < 7; index++ {

		zufallszahl = rand.Intn(49) + 1

		for IstZahlInSlice(ziehungszahlen, zufallszahl) {
			zufallszahl = rand.Intn(49) + 1
		}

		if zufallszahl < 10 {
			ziehung = ziehung + "0"
		}
		ziehung = ziehung + strconv.Itoa(zufallszahl)
		ziehungszahlen = append(ziehungszahlen, zufallszahl)
	}

	zufallszahl = rand.Intn(10)
	ziehung = ziehung + strconv.Itoa(zufallszahl)

	return ziehung

}

func IstValideZiehung(ziehung string) bool {

	var valid bool

	valid, _ = regexp.MatchString("(0[1-9]|[1-4][0-9]){6}[0-9]", ziehung)

	valid = valid && (len(ziehung) == 13)

	if valid {

		var zahl int
		var ziehungszahlen []int

		ziehungszahlen = []int{}

		for index := 0; index < 6; index++ {

			zahl, _ = strconv.Atoi(ziehung[2*index : 2*index+2])

			if IstZahlInSlice(ziehungszahlen, zahl) {
				valid = false
				break
			} else {
				ziehungszahlen = append(ziehungszahlen, zahl)
			}

		}

	}

	return valid

}

func IstValidesTippdatum(databasehandle *sql.DB, tippdatum time.Time) bool {

	var valid bool

	valid = false

	var letzteZiehung time.Time

	query := "SELECT MAX(datum) as letzteZiehung FROM ziehungen"
	fmt.Println(query)

	rows, err := databasehandle.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
		return false
	} else {
		for rows.Next() {
			rows.Scan(&letzteZiehung)
			fmt.Printf("letzte Ziehung am %s\n", letzteZiehung.Format("2006-01-02"))
		}
	}

	valid = tippdatum.Truncate(24 * time.Hour).After(letzteZiehung.Truncate(24 * time.Hour))

	return valid

}

func IstZahlInSlice(slice []int, gesucht int) bool {
	for _, zahl := range slice {
		if zahl == gesucht {
			return true
		}
	}
	return false
}

func BerechneAuszahlungen(databasehandle *sql.DB, datum time.Time) ([]Auszahlung, error) {

	var auszahlungen []Auszahlung
	var gesamteinsatz float64
	var gesamttipps int

	einsatz := 5.0

	query := "SELECT COUNT(*) as anzahl FROM tipps WHERE datum = '" + datum.Format("2006-01-02") + "'"

	rows, err := databasehandle.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			rows.Scan(&gesamttipps)
			fmt.Printf("Abgegebene Tipps: %d\n", gesamttipps)
		}
	}

	gesamteinsatz = float64(gesamttipps) * einsatz

	var gewinnanzahl [10]int

	gewinnanzahl = BerechneGewinnanzahl(databasehandle, datum)

	fmt.Printf("Startbudget: %.2f\n", gesamteinsatz)

	zweirichtige := Auszahlung{
		Datum:      datum,
		Klasse:     0,
		Auszahlung: 2.0,
	}

	zweirichtige.Budget = zweirichtige.Auszahlung * float64(gewinnanzahl[zweirichtige.Klasse])
	auszahlungen = append(auszahlungen, zweirichtige)
	gesamteinsatz = gesamteinsatz - zweirichtige.Budget

	fmt.Printf("Budget nach Abzug vom Budget für 2 Richtige: %.2f\n", gesamteinsatz)

	zweirichtigesuper := Auszahlung{
		Datum:      datum,
		Klasse:     1,
		Auszahlung: 5.0,
	}

	zweirichtigesuper.Budget = zweirichtigesuper.Auszahlung * float64(gewinnanzahl[zweirichtigesuper.Klasse])
	auszahlungen = append(auszahlungen, zweirichtigesuper)
	gesamteinsatz = gesamteinsatz - zweirichtigesuper.Budget

	fmt.Printf("Budget nach Abzug vom Budget für 2 Richtige mit Superzahl: %.2f\n", gesamteinsatz)

	sechsrichtigesuper := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.128, datum, 9)
	gesamteinsatz = gesamteinsatz - sechsrichtigesuper.Budget
	sechsrichtigesuper = BerechneJackpot(databasehandle, sechsrichtigesuper)
	auszahlungen = append(auszahlungen, sechsrichtigesuper)

	fmt.Printf("Budget nach Abzug vom Budget für 6 Richtige mit Superzahl: %.2f\n", gesamteinsatz)

	sechsrichtige := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.1, datum, 8)
	auszahlungen = append(auszahlungen, sechsrichtige)

	fuenfrichtige := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.15, datum, 6)
	auszahlungen = append(auszahlungen, fuenfrichtige)

	fuenfrichtigesuper := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.05, datum, 7)
	auszahlungen = append(auszahlungen, fuenfrichtigesuper)

	vierrichtige := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.1, datum, 4)
	auszahlungen = append(auszahlungen, vierrichtige)

	vierrichtigesuper := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.05, datum, 5)
	auszahlungen = append(auszahlungen, vierrichtigesuper)

	dreirichtige := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.45, datum, 2)
	auszahlungen = append(auszahlungen, dreirichtige)

	dreirichtigesuper := BerechneAuszahlung(gesamteinsatz, gewinnanzahl, 0.1, datum, 3)
	auszahlungen = append(auszahlungen, dreirichtigesuper)

	return auszahlungen, nil

}

func BerechneAuszahlung(gesamteinsatz float64, gewinnanzahl [10]int, budgetfaktor float64, auszahlungsdatum time.Time, gewinnklasse int8) Auszahlung {

	auszahlung := Auszahlung{
		Datum:  auszahlungsdatum,
		Klasse: gewinnklasse,
	}

	auszahlung.Budget = gesamteinsatz * budgetfaktor
	auszahlung.Budget = math.Round(auszahlung.Budget*100) / 100

	if gewinnanzahl[auszahlung.Klasse] > 0 {
		auszahlung.Auszahlung = auszahlung.Budget / float64(gewinnanzahl[auszahlung.Klasse])
		auszahlung.Auszahlung = math.Round(auszahlung.Auszahlung*100) / 100
	} else {
		auszahlung.Auszahlung = 0.0
	}

	return auszahlung
}

func BerechneJackpot(databasehandle *sql.DB, jackpotAuszahlung Auszahlung) Auszahlung {

	var jackpot float64
	var auszahlung float64

	// Betrachte Jackpotauszahlung der letzten Ziehung
	query := "SELECT budget, auszahlung FROM auszahlungen WHERE klasse = 9 AND datum = (SELECT MAX(datum) FROM auszahlungen)"
	fmt.Println(query)

	rows, err := databasehandle.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			rows.Scan(&jackpot, &auszahlung)
			fmt.Printf("unausgezahlte Jackpotsumme : %.2f\n", jackpot)
		}
	}

	if auszahlung == 0 {
		jackpotAuszahlung.Budget = jackpotAuszahlung.Budget + jackpot
	}

	return jackpotAuszahlung

}

func BerechneGewinnanzahl(databasehandle *sql.DB, datum time.Time) [10]int {

	var gewinnanzahl [10]int

	var ziehung string

	var tipp string

	for klasse := 0; klasse <= 9; klasse++ {

		gewinnanzahl[klasse] = 0

	}

	query := "SELECT ziehung FROM ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'"

	ziehungsdaten, err := databasehandle.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for ziehungsdaten.Next() {
			ziehungsdaten.Scan(&ziehung)
		}
	}

	query = "SELECT ziehung FROM tipps WHERE datum = '" + datum.Format("2006-01-02") + "'"

	tippdaten, err := databasehandle.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for tippdaten.Next() {
			tippdaten.Scan(&tipp)
			gewinnklasse := BerechneGewinnklasse(tipp, ziehung)
			fmt.Printf("Tipp %s erzielt Gewinnklasse %d bei Ziehung %s\n", tipp, gewinnklasse, ziehung)
			if gewinnklasse >= 0 {
				gewinnanzahl[gewinnklasse]++
			}
		}
	}

	for klasse := 0; klasse <= 9; klasse++ {

		fmt.Printf("Gewinner in Klasse %d: %d\n", klasse, gewinnanzahl[klasse])

	}

	return gewinnanzahl

}

func BerechneGewinnklasse(tipp string, ziehung string) int {

	var gewinnklasse int

	gewinnklasse = -4

	for tippindex := 0; tippindex < 6; tippindex++ {

		for ziehungsindex := 0; ziehungsindex < 6; ziehungsindex++ {

			if tipp[2*tippindex:2*tippindex+2] == ziehung[2*ziehungsindex:2*ziehungsindex+2] {
				gewinnklasse = gewinnklasse + 2
			}

		}

	}

	if tipp[12] == ziehung[12] {
		gewinnklasse = gewinnklasse + 1
	}

	return gewinnklasse

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
