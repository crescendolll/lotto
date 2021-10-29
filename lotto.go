package main

import (
	"flag"
	"fmt"
	"lotto/database"
	"lotto/lottohttp"
	"lotto/lottolog"
	"lotto/lottologic"
)

func main() {

	lottolog.OeffneLogdatei()

	database.OeffneVerbindungZurLottoDatenbank()

	if parseAdminFlag() {
		MitarbeiterEintragen()
	} else {
		lottohttp.StarteLottoServer()
	}

	database.SchliesseVerbindung()

}

func MitarbeiterEintragen() {

	var name string
	var passwort string

	fmt.Println("Willkommen bei der Mitarbeiter-Eingabe")

	fmt.Println("Bitte Namen des Mitarbeiters eingeben")
	fmt.Scanln(&name)

	fmt.Println("Bitte Passwort des Mitarbeiters eingeben")
	fmt.Scanln(&passwort)

	fehler := lottologic.FuegeMitarbeiterNachPruefungEin(name, passwort)

	if fehler != nil {
		fmt.Println(fehler.Error())
	} else {
		fmt.Println("Mitarbeiter erfolgreich eingetragen")
	}

}

func parseAdminFlag() bool {

	boolPtr := flag.Bool("admin", false, "Mitarbeiterverwaltung starten")

	flag.Parse()

	return *boolPtr
}
