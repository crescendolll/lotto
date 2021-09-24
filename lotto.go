package main

import (
	"flag"
	"fmt"
	"lotto/database"
	"lotto/lottohttp"
)

func main() {

	if parseAdminFlag() {
		MitarbeiterEintragen()
	} else {
		lottohttp.OpenLottoServer()
	}

}

func MitarbeiterEintragen() {

	var name string
	var pw string

	fmt.Println("Willkommen bei der Mitarbeiter-Eingabe")

	fmt.Println("Bitte Namen des Mitarbeiters eingeben")
	fmt.Scanln(&name)

	fmt.Println("Bitte Passwort des Mitarbeiters eingeben")
	fmt.Scanln(&pw)

	databasehandle := database.OpenLottoConnection()
	err := database.InsertMitarbeiterIntoNutzer(databasehandle, name, pw)
	database.CloseLottoConnection(databasehandle)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Mitarbeiter erfolgreich eingetragen")
	}

}

func parseAdminFlag() bool {

	boolPtr := flag.Bool("useradmin", false, "activate user administration")

	flag.Parse()

	return *boolPtr
}
