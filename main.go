package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

var database *sql.DB

type Nutzer struct {
	benutzername string
	pw_hash      string
	ist_spieler  int8
}

type Tipp struct {
	id      int64
	datum   time.Time
	ziehung string
}

type Ziehung struct {
	datum   time.Time
	ziehung string
}

type Auszahlung struct {
	datum      time.Time
	klasse     int8
	budget     float64
	auszahlung float64
}

func main() {
	// Verbindungseigenschaften festhalten
	configuration := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "lotto",
		ParseTime: true,
	}
	// Get a database handle
	var accessError error
	database, accessError = sql.Open("mysql", configuration.FormatDSN())
	if accessError != nil {
		log.Fatal(accessError)
	}

	pingError := database.Ping()
	if pingError != nil {
		log.Fatal(pingError)
	}
	fmt.Println("Verbindung mit Datenbank lotto aufgebaut")

	ablaufTest()

}

func ablaufTest() {

	var neuerNutzer Nutzer

	neuerNutzer = Nutzer{
		benutzername: "Hans",
		pw_hash:      "23bonobo42",
		ist_spieler:  1,
	}

	var andererNutzer Nutzer

	andererNutzer = Nutzer{
		benutzername: "Sofia",
		pw_hash:      "P3psiman",
		ist_spieler:  1,
	}

	var insertError error

	fmt.Println("\nSpieler Hans anlegen")
	insertError = addSpieler(neuerNutzer.benutzername, neuerNutzer.pw_hash)
	if insertError != nil {
		log.Println(insertError)
	}
	fmt.Println("\nSpieler Hans nochmal anlegen")
	insertError = addSpieler(neuerNutzer.benutzername, neuerNutzer.pw_hash)
	if insertError != nil {
		log.Println(insertError)
	}
	fmt.Println("\nSpielerin Sofia anlegen")
	insertError = addSpieler(andererNutzer.benutzername, andererNutzer.pw_hash)
	if insertError != nil {
		log.Println(insertError)
	}

	var neuerTipp Tipp

	neuerTipp = Tipp{
		datum:   time.Now(),
		ziehung: randomZiehung(),
	}

	for tippindex := 1; tippindex <= 100; tippindex++ {

		neuerTipp.ziehung = randomZiehung()

		fmt.Println("\nSofia gibt einen Tipp ab")
		insertError = addTipp(neuerTipp, andererNutzer)
		if insertError != nil {
			log.Println(insertError)
		}

	}

	var nichtNutzer Nutzer

	nichtNutzer = Nutzer{
		benutzername: "Lea",
		pw_hash:      "B0bd3rB4um3ister",
		ist_spieler:  1,
	}

	neuerTipp.ziehung = randomZiehung()

	fmt.Println("\nnicht registrierte Nutzerin Lea gibt einen Tipp ab")
	insertError = addTipp(neuerTipp, nichtNutzer)
	if insertError != nil {
		log.Println(insertError)
	}

	var neuesPasswort string
	var updateError error

	neuesPasswort = "Apoc4lyps3"

	fmt.Println("\nPasswort von Hans ändern")
	updateError = updateSpieler(neuerNutzer.benutzername, neuerNutzer.benutzername, neuesPasswort)
	if updateError != nil {
		log.Println(updateError)
	} else {
		neuerNutzer.pw_hash = neuesPasswort
	}

	neuesPasswort = "Banan3nbr0t"

	var neuerBenutzername = "Katharina"

	fmt.Println("\nPasswort und Name von Sofia ändern")
	updateError = updateSpieler(andererNutzer.benutzername, neuerBenutzername, neuesPasswort)
	if updateError != nil {
		log.Println(updateError)
	} else {
		andererNutzer.pw_hash = neuesPasswort
		andererNutzer.benutzername = neuerBenutzername
	}

	neuerBenutzername = neuerNutzer.benutzername

	fmt.Println("\nName auf einen bereits belegten Namen ändern")
	updateError = updateSpieler(andererNutzer.benutzername, neuerBenutzername, andererNutzer.pw_hash)
	if updateError != nil {
		log.Println(updateError)
	} else {
		andererNutzer.benutzername = neuerBenutzername
	}

	fmt.Println("\nZeige Katharinas Tipps an")
	tipps, selectError := holeTippsVonSpieler(andererNutzer.benutzername)
	if selectError != nil {
		log.Println(selectError)
	} else {
		for _, tipp := range tipps {
			fmt.Printf("%s für die Ziehung am %s\n", tipp.ziehung, tipp.datum.Format("2006-01-02"))
		}
	}

	var ziehung Ziehung

	ziehung = Ziehung{
		datum:   time.Now(),
		ziehung: randomZiehung(),
	}

	fmt.Println("\nFüge Ziehung ein")
	insertError = addZiehung(ziehung)
	if insertError != nil {
		log.Println(insertError)
	}

	fmt.Println("\nZeige alle Ziehungen an")
	ziehungen, selectError := getZiehungen()
	if selectError != nil {
		log.Println(selectError)
	} else {
		for _, ziehung := range ziehungen {
			fmt.Printf("%s ist die Ziehung am %s\n", ziehung.ziehung, ziehung.datum.Format("2006-01-02"))
		}
	}

	fmt.Println("\nBerechne Auszahlungen")
	auszahlungen, _ := berechneBudgets(ziehung.datum)

	fmt.Println("\nVeröffentliche Auszahlungen")
	veroeffentlicheAuszahlungen(auszahlungen)

	var deleteError error

	fmt.Println("\nLösche Spieler Hans")
	deleteError = loescheSpieler(neuerNutzer.benutzername)
	if deleteError != nil {
		log.Println(deleteError)
	}

	fmt.Println("\nLösche Spielerin Katharina")
	deleteError = loescheSpieler(andererNutzer.benutzername)
	if deleteError != nil {
		log.Println(deleteError)
	}

	var arbeitenderNutzer Nutzer

	arbeitenderNutzer = Nutzer{
		benutzername: "Ingo",
		pw_hash:      "123",
		ist_spieler:  0,
	}

	fmt.Println("\nLösche Mitarbeiter Ingo")
	deleteError = loescheSpieler(arbeitenderNutzer.benutzername)
	if deleteError != nil {
		log.Println(deleteError)
	}

}

func addSpieler(benutzername string, pw_hash string) error {

	var insertError error

	if name_verfuegbar(benutzername) {
		_, insertError = database.Exec("INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (?,?,1)", benutzername, pw_hash)
		fmt.Printf("INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (%s,%s,1)\n", benutzername, pw_hash)
	} else {
		insertError = errors.New("Name " + benutzername + " bereits vergeben")
	}

	return insertError
}

func updateSpieler(alterBenutzername string, neuerBenutzername string, neuesPasswort string) error {

	var updateError error

	if alterBenutzername != neuerBenutzername {
		if !name_verfuegbar(neuerBenutzername) {
			updateError = errors.New("Name " + neuerBenutzername + " bereits vergeben")
			return updateError
		}
	}

	_, updateError = database.Exec("UPDATE nutzer set pw_hash = ?, benutzername = ? WHERE benutzername = ?", neuesPasswort, neuerBenutzername, alterBenutzername)
	fmt.Printf("UPDATE nutzer set pw_hash = %s, benutzername = %s WHERE benutzername = %s\n", neuesPasswort, neuerBenutzername, alterBenutzername)

	return updateError

}

func loescheSpieler(benutzername string) error {

	var deleteError error

	_, deleteError = database.Exec("DELETE from nutzer WHERE benutzername = ? AND ist_spieler = 1", benutzername)
	fmt.Printf("DELETE from nutzer WHERE benutzername = %s AND ist_spieler = 1\n", benutzername)

	return deleteError

}

func name_verfuegbar(benutzername string) bool {
	var ist_verfuegbar bool

	var anzahl_vorkommen int

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

	return ist_verfuegbar
}

func addTipp(tipp Tipp, spieler Nutzer) error {

	var insertError error
	var readError error
	var result sql.Result

	if !ist_valide_ziehung(tipp.ziehung) {
		insertError = errors.New("Keine gültige Lottoziehung: " + tipp.ziehung)
		return insertError
	}

	fmt.Println("Transaktion beginnt")
	tippTransaktion, _ := database.Begin()

	result, insertError = tippTransaktion.Exec("INSERT INTO tipps (datum, ziehung) VALUES (?,?)", tipp.datum.Format("2006-01-02"), tipp.ziehung)
	fmt.Printf("INSERT INTO tipps (datum, ziehung) VALUES (%s,%s)\n", tipp.datum.Format("2006-01-02"), tipp.ziehung)
	if insertError != nil {
		fmt.Println("Transaktion abgebrochen")
		tippTransaktion.Rollback()
		return insertError
	}

	var id int64

	id, readError = result.LastInsertId()
	if readError != nil {
		fmt.Println("Transaktion abgebrochen")
		tippTransaktion.Rollback()
		return readError
	} else {
		tipp.id = id
	}

	insertError = verknuepfeSpielerMitTipp(tipp.id, spieler.benutzername, tippTransaktion)

	if insertError != nil {
		fmt.Println("Transaktion abgebrochen")
		tippTransaktion.Rollback()
	} else {
		fmt.Println("Transaktion erfolgt")
		tippTransaktion.Commit()
	}

	return insertError

}

func verknuepfeSpielerMitTipp(tippID int64, spielername string, tippTransaktion *sql.Tx) error {

	var insertError error

	_, insertError = tippTransaktion.Exec("INSERT INTO spieler_tipps (id, spielername) VALUES (?,?)", tippID, spielername)
	fmt.Printf("INSERT INTO spieler_tipps (id, spielername) VALUES (%d,%s)\n", tippID, spielername)

	return insertError

}

func randomZiehung() string {

	var ziehung string

	ziehung = ""

	rand.Seed(time.Now().UnixNano())

	var zufallszahl int
	var ziehungszahlen []int

	ziehungszahlen = []int{}

	for index := 1; index < 7; index++ {

		zufallszahl = rand.Intn(49) + 1

		for enthaelt(ziehungszahlen, zufallszahl) {
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

func ist_valide_ziehung(ziehung string) bool {

	var valid bool

	valid, _ = regexp.MatchString("(0[1-9]|[1-4][0-9]){6}[0-9]", ziehung)

	valid = valid && (len(ziehung) == 13)

	if valid {

		var zahl int
		var ziehungszahlen []int

		ziehungszahlen = []int{}

		for index := 0; index < 6; index++ {

			zahl, _ = strconv.Atoi(ziehung[2*index : 2*index+2])

			if enthaelt(ziehungszahlen, zahl) {
				valid = false
				break
			} else {
				ziehungszahlen = append(ziehungszahlen, zahl)
			}

		}

	}

	return valid

}

func enthaelt(zahlenraum []int, gesucht int) bool {
	for _, zahl := range zahlenraum {
		if zahl == gesucht {
			return true
		}
	}
	return false
}

func holeTippsVonSpieler(spielername string) ([]Tipp, error) {

	var tipps []Tipp

	rows, selectError := database.Query("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = ? AND spieler_tipps.id = tipps.id", spielername)
	fmt.Printf("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = %s AND spieler_tipps.id = tipps.id\n", spielername)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	for rows.Next() {
		var tipp Tipp
		if scanError := rows.Scan(&tipp.id, &tipp.datum, &tipp.ziehung); scanError != nil {
			return nil, scanError
		}
		tipps = append(tipps, tipp)
	}

	if sqlError := rows.Err(); sqlError != nil {
		return nil, sqlError
	}
	return tipps, nil
}

func addZiehung(ziehung Ziehung) error {

	var insertError error

	_, insertError = database.Exec("INSERT INTO ziehungen (datum, ziehung) VALUES (?,?)", ziehung.datum.Format("2006-01-02"), ziehung.ziehung)
	fmt.Printf("INSERT INTO ziehungen (datum, ziehung) VALUES (%s,%s)\n", ziehung.datum.Format("2006-01-02"), ziehung.ziehung)

	return insertError

}

func getZiehungen() ([]Ziehung, error) {

	var ziehungen []Ziehung
	var selectError error

	rows, selectError := database.Query("SELECT * from ziehungen")
	fmt.Printf("SELECT * FROM ziehungen\n")
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	for rows.Next() {
		var ziehung Ziehung
		if scanError := rows.Scan(&ziehung.datum, &ziehung.ziehung); scanError != nil {
			return nil, scanError
		}
		ziehungen = append(ziehungen, ziehung)
	}

	if sqlError := rows.Err(); sqlError != nil {
		return nil, sqlError
	}

	return ziehungen, selectError

}

func berechneBudgets(datum time.Time) ([]Auszahlung, error) {

	var auszahlungen []Auszahlung

	var gesamteinsatz float64
	var gesamttipps int

	einsatz := 5.0

	query := "SELECT COUNT(*) as anzahl FROM tipps WHERE datum = '" + datum.Format("2006-01-02") + "'"

	rows, err := database.Query(query)
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

	gewinnanzahl = berechneGewinnanzahl(datum)

	fmt.Printf("Startbudget: %.2f\n", gesamteinsatz)

	zweirichtige := Auszahlung{
		datum:      datum,
		klasse:     0,
		auszahlung: 2.0,
	}

	zweirichtige.budget = zweirichtige.auszahlung * float64(gewinnanzahl[zweirichtige.klasse])

	auszahlungen = append(auszahlungen, zweirichtige)

	gesamteinsatz = gesamteinsatz - zweirichtige.budget

	fmt.Printf("Budget nach Abzug vom Budget für 2 Richtige: %.2f\n", gesamteinsatz)

	zweirichtigesuper := Auszahlung{
		datum:      datum,
		klasse:     1,
		auszahlung: 5.0,
	}

	zweirichtigesuper.budget = zweirichtigesuper.auszahlung * float64(gewinnanzahl[zweirichtigesuper.klasse])

	auszahlungen = append(auszahlungen, zweirichtigesuper)

	gesamteinsatz = gesamteinsatz - zweirichtigesuper.budget

	fmt.Printf("Budget nach Abzug vom Budget für 2 Richtige mit Superzahl: %.2f\n", gesamteinsatz)

	sechsrichtigesuper := Auszahlung{
		datum:  datum,
		klasse: 9,
	}

	sechsrichtigesuper.budget = gesamteinsatz * 0.128
	sechsrichtigesuper.budget = math.Round(sechsrichtigesuper.budget*100) / 100

	sechsrichtigesuper.auszahlung = 0.0

	if gewinnanzahl[sechsrichtigesuper.klasse] > 0 {
		sechsrichtigesuper.auszahlung = sechsrichtigesuper.budget / float64(gewinnanzahl[sechsrichtigesuper.klasse])
		sechsrichtigesuper.auszahlung = math.Round(sechsrichtigesuper.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, sechsrichtigesuper)

	gesamteinsatz = gesamteinsatz - sechsrichtigesuper.budget

	fmt.Printf("Budget nach Abzug vom Budget für 6 Richtige mit Superzahl: %.2f\n", gesamteinsatz)

	sechsrichtige := Auszahlung{
		datum:  datum,
		klasse: 8,
	}

	sechsrichtige.budget = gesamteinsatz * 0.1
	sechsrichtige.budget = math.Round(sechsrichtige.budget*100) / 100

	sechsrichtige.auszahlung = 0.0

	if gewinnanzahl[sechsrichtige.klasse] > 0 {
		sechsrichtige.auszahlung = sechsrichtige.budget / float64(gewinnanzahl[sechsrichtige.klasse])
		sechsrichtige.auszahlung = math.Round(sechsrichtige.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, sechsrichtige)

	fuenfrichtige := Auszahlung{
		datum:  datum,
		klasse: 6,
	}

	fuenfrichtige.budget = gesamteinsatz * 0.15
	fuenfrichtige.budget = math.Round(fuenfrichtige.budget*100) / 100

	fuenfrichtige.auszahlung = 0.0

	if gewinnanzahl[fuenfrichtige.klasse] > 0 {
		fuenfrichtige.auszahlung = fuenfrichtige.budget / float64(gewinnanzahl[fuenfrichtige.klasse])
		fuenfrichtige.auszahlung = math.Round(fuenfrichtige.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, fuenfrichtige)

	fuenfrichtigesuper := Auszahlung{
		datum:  datum,
		klasse: 7,
	}

	fuenfrichtigesuper.budget = gesamteinsatz * 0.05
	fuenfrichtigesuper.budget = math.Round(fuenfrichtigesuper.budget*100) / 100

	fuenfrichtigesuper.auszahlung = 0.0

	if gewinnanzahl[fuenfrichtigesuper.klasse] > 0 {
		fuenfrichtigesuper.auszahlung = fuenfrichtigesuper.budget / float64(gewinnanzahl[fuenfrichtigesuper.klasse])
		fuenfrichtigesuper.auszahlung = math.Round(fuenfrichtigesuper.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, fuenfrichtigesuper)

	vierrichtige := Auszahlung{
		datum:  datum,
		klasse: 4,
	}

	vierrichtige.budget = gesamteinsatz * 0.1
	vierrichtige.budget = math.Round(vierrichtige.budget*100) / 100

	vierrichtige.auszahlung = 0.0

	if gewinnanzahl[vierrichtige.klasse] > 0 {
		vierrichtige.auszahlung = vierrichtige.budget / float64(gewinnanzahl[vierrichtige.klasse])
		vierrichtige.auszahlung = math.Round(vierrichtige.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, vierrichtige)

	vierrichtigesuper := Auszahlung{
		datum:  datum,
		klasse: 5,
	}

	vierrichtigesuper.budget = gesamteinsatz * 0.05
	vierrichtigesuper.budget = math.Round(vierrichtigesuper.budget*100) / 100

	vierrichtigesuper.auszahlung = 0.0

	if gewinnanzahl[vierrichtigesuper.klasse] > 0 {
		vierrichtigesuper.auszahlung = vierrichtigesuper.budget / float64(gewinnanzahl[vierrichtigesuper.klasse])
		vierrichtigesuper.auszahlung = math.Round(vierrichtigesuper.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, vierrichtigesuper)

	dreirichtige := Auszahlung{
		datum:  datum,
		klasse: 2,
	}

	dreirichtige.budget = gesamteinsatz * 0.45
	dreirichtige.budget = math.Round(dreirichtige.budget*100) / 100

	dreirichtige.auszahlung = 0.0

	if gewinnanzahl[dreirichtige.klasse] > 0 {
		dreirichtige.auszahlung = dreirichtige.budget / float64(gewinnanzahl[dreirichtige.klasse])
		dreirichtige.auszahlung = math.Round(dreirichtige.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, dreirichtige)

	dreirichtigesuper := Auszahlung{
		datum:  datum,
		klasse: 3,
	}

	dreirichtigesuper.budget = gesamteinsatz * 0.1
	dreirichtigesuper.budget = math.Round(dreirichtigesuper.budget*100) / 100

	dreirichtigesuper.auszahlung = 0.0

	if gewinnanzahl[dreirichtigesuper.klasse] > 0 {
		dreirichtigesuper.auszahlung = dreirichtigesuper.budget / float64(gewinnanzahl[dreirichtigesuper.klasse])
		dreirichtigesuper.auszahlung = math.Round(dreirichtigesuper.auszahlung*100) / 100
	}

	auszahlungen = append(auszahlungen, dreirichtigesuper)

	return auszahlungen, nil

}

func berechneGewinnanzahl(datum time.Time) [10]int {

	var gewinnanzahl [10]int

	var ziehung string

	var tipp string

	for klasse := 0; klasse <= 9; klasse++ {

		gewinnanzahl[klasse] = 0

	}

	query := "SELECT ziehung FROM ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'"

	ziehungsdaten, err := database.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for ziehungsdaten.Next() {
			ziehungsdaten.Scan(&ziehung)
		}
	}

	query = "SELECT ziehung FROM tipps WHERE datum = '" + datum.Format("2006-01-02") + "'"

	tippdaten, err := database.Query(query)
	fmt.Println(query)

	if err != nil {
		log.Println(err)
	} else {
		for tippdaten.Next() {
			tippdaten.Scan(&tipp)
			gewinnklasse := berechneGewinnklasse(tipp, ziehung)
			fmt.Printf("Tipp %s erzielt Gewinnklasse %d bei Ziehung %s\n", tipp, gewinnklasse, ziehung)
			if gewinnklasse >= 0 {
				gewinnanzahl[gewinnklasse]++
			}
		}
	}

	for klasse := 0; klasse <= 9; klasse++ {

		fmt.Printf("Gewinner in Klasse %d: %d\n", klasse+1, gewinnanzahl[klasse])

	}

	return gewinnanzahl

}

func berechneGewinnklasse(tipp string, ziehung string) int {

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

func veroeffentlicheAuszahlungen(auszahlungen []Auszahlung) error {

	var insertError error

	for _, auszahlung := range auszahlungen {
		_, insertError = database.Exec("INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (?,?,?,?)", auszahlung.datum.Format("2006-01-02"), auszahlung.klasse, auszahlung.budget, auszahlung.auszahlung)
		fmt.Printf("INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (%s,%d,%.2f,%.2f)\n", auszahlung.datum.Format("2006-01-02"), auszahlung.klasse, auszahlung.budget, auszahlung.auszahlung)
	}

	return insertError

}
