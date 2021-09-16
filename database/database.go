package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"lotto/lottologic"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func OpenLottoConnection() *sql.DB {

	// Verbindung konfigurieren
	configuration := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "lotto",
		ParseTime: true,
	}

	// Datenbankzugriff öffnen
	var database *sql.DB
	var accessError error
	database, accessError = sql.Open("mysql", configuration.FormatDSN())
	if accessError != nil {
		log.Fatal(accessError)
	}

	pingError := database.Ping()
	if pingError != nil {
		log.Fatal(pingError)
	} else {
		fmt.Println("Verbindung mit Datenbank lotto aufgebaut")
	}

	return database
}

func CloseLottoConnection(databasehandle *sql.DB) {

	closingError := databasehandle.Close()
	if closingError != nil {
		log.Fatal(closingError)
	} else {
		fmt.Println("Verbindung mit Datenbank lotto geschlossen")
	}

}

// es dürfen nur neue Spieler eingefügt werden, keine Mitarbeiter
func InsertIntoSpieler(databasehandle *sql.DB, benutzername string, pw_hash string) error {

	var insertError error

	insertError = lottologic.IstNameVerfuegbar(databasehandle, benutzername)

	if insertError != nil {
		return insertError
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (?,?,1)"
	statement, err := databasehandle.PrepareContext(context, query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(context, benutzername, pw_hash)
	return err
}

func UpdateSpieler(databasehandle *sql.DB, alterBenutzername string, neuerBenutzername string, neuesPasswort string) error {

	var updateError error

	if alterBenutzername != neuerBenutzername {
		updateError = lottologic.IstNameVerfuegbar(databasehandle, neuerBenutzername)
	}

	if updateError == nil {
		_, updateError = databasehandle.Exec("UPDATE nutzer set pw_hash = ?, benutzername = ? WHERE benutzername = ?", neuesPasswort, neuerBenutzername, alterBenutzername)
		fmt.Printf("UPDATE nutzer set pw_hash = %s, benutzername = %s WHERE benutzername = %s\n", neuesPasswort, neuerBenutzername, alterBenutzername)
	}

	return updateError

}

// es dürfen nur Spieler gelöscht werden
func DeleteFromSpieler(databasehandle *sql.DB, benutzername string) error {

	var deleteError error

	_, deleteError = databasehandle.Exec("DELETE from nutzer WHERE benutzername = ? AND ist_spieler = 1", benutzername)
	fmt.Printf("DELETE from nutzer WHERE benutzername = %s AND ist_spieler = 1\n", benutzername)

	return deleteError

}

func SelectFromSpielerByName(databasehandle *sql.DB, benutzername string) (lottologic.Nutzer, error) {

	var nutzer lottologic.Nutzer
	var selectError error

	rows, selectError := databasehandle.Query("SELECT * from nutzer WHERE benutzername = ?", benutzername)
	fmt.Printf("SELECT * from nutzer WHERE benutzername = %s\n", benutzername)
	if selectError != nil {
		return nutzer, selectError
	}
	defer rows.Close()

	for rows.Next() {
		if scanError := rows.Scan(&nutzer.Benutzername, &nutzer.Pw_hash, &nutzer.Ist_spieler); scanError != nil {
			return nutzer, scanError
		}
	}

	if sqlError := rows.Err(); sqlError != nil {
		return nutzer, sqlError
	}

	return nutzer, selectError

}

// nur Abgabe eines gültigen Tipps für ein Datum nach der letzten erfolgten Ziehung möglich
func InsertIntoTipps(databasehandle *sql.DB, tipp lottologic.Tipp, spieler lottologic.Nutzer) error {

	var insertError error
	var readError error
	var result sql.Result

	if !lottologic.IstValideZiehung(tipp.Ziehung) {
		insertError = errors.New("Keine gültige Lottoziehung: " + tipp.Ziehung)
		return insertError
	}

	if !lottologic.IstValidesTippdatum(databasehandle, tipp.Datum) {
		insertError = errors.New("Kein gültiges Tippdatum: " + tipp.Datum.Format("2006-01-02"))
		return insertError
	}

	fmt.Println("Transaktion beginnt")
	tippTransaktion, _ := databasehandle.Begin()

	result, insertError = tippTransaktion.Exec("INSERT INTO tipps (datum, ziehung) VALUES (?,?)", tipp.Datum.Format("2006-01-02"), tipp.Ziehung)
	fmt.Printf("INSERT INTO tipps (datum, ziehung) VALUES (%s,%s)\n", tipp.Datum.Format("2006-01-02"), tipp.Ziehung)
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
		tipp.Id = id
	}

	insertError = InsertIntoSpielerTipps(tipp.Id, spieler.Benutzername, tippTransaktion)

	if insertError != nil {
		fmt.Println("Transaktion abgebrochen")
		tippTransaktion.Rollback()
	} else {
		fmt.Println("Transaktion erfolgt")
		tippTransaktion.Commit()
	}

	return insertError

}

func InsertIntoSpielerTipps(tippID int64, spielername string, tippTransaktion *sql.Tx) error {

	var insertError error

	_, insertError = tippTransaktion.Exec("INSERT INTO spieler_tipps (id, spielername) VALUES (?,?)", tippID, spielername)
	fmt.Printf("INSERT INTO spieler_tipps (id, spielername) VALUES (%d,%s)\n", tippID, spielername)

	return insertError

}

func SelectFromTippsBySpielername(databasehandle *sql.DB, spielername string) ([]lottologic.Tipp, error) {

	var tipps []lottologic.Tipp

	rows, selectError := databasehandle.Query("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = ? AND spieler_tipps.id = tipps.id", spielername)
	fmt.Printf("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = %s AND spieler_tipps.id = tipps.id\n", spielername)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	for rows.Next() {
		var tipp lottologic.Tipp
		if scanError := rows.Scan(&tipp.Id, &tipp.Datum, &tipp.Ziehung); scanError != nil {
			return nil, scanError
		}
		tipps = append(tipps, tipp)
	}

	if sqlError := rows.Err(); sqlError != nil {
		return nil, sqlError
	}
	return tipps, nil
}

// es dürfen nur gültige Ziehungen nach der letzten Ziehung eingetragen werden
// wird eine Ziehung eingetragen, werden die Gewinnauszahlungen mitberechnet und gespeichert
func InsertIntoZiehungen(databasehandle *sql.DB, ziehung lottologic.Ziehung) error {

	var insertError error
	var auszahlungen []lottologic.Auszahlung

	if !lottologic.IstValideZiehung(ziehung.Ziehung) {
		insertError = errors.New(ziehung.Ziehung + " ist keine gültige Ziehung")
		return insertError
	}

	if !lottologic.IstValidesTippdatum(databasehandle, ziehung.Datum) {
		insertError = errors.New("Kein gültiges Tippdatum: " + ziehung.Datum.Format("2006-01-02"))
		return insertError
	}

	fmt.Println("Transaktion beginnt")
	ziehungTransaktion, _ := databasehandle.Begin()

	_, insertError = ziehungTransaktion.Exec("INSERT INTO ziehungen (datum, ziehung) VALUES (?,?)", ziehung.Datum.Format("2006-01-02"), ziehung.Ziehung)
	fmt.Printf("INSERT INTO ziehungen (datum, ziehung) VALUES (%s,%s)\n", ziehung.Datum.Format("2006-01-02"), ziehung.Ziehung)

	if insertError == nil {
		auszahlungen, insertError = lottologic.BerechneAuszahlungen(databasehandle, ziehung.Datum)
	}

	if insertError == nil {
		insertError = InsertIntoAuszahlungen(ziehungTransaktion, auszahlungen)
	}

	if insertError != nil {
		fmt.Println("Transaktion abgebrochen")
		ziehungTransaktion.Rollback()
	} else {
		fmt.Println("Transaktion erfolgt")
		ziehungTransaktion.Commit()
	}

	return insertError

}

func SelectFromZiehungen(databasehandle *sql.DB) ([]lottologic.Ziehung, error) {

	var ziehungen []lottologic.Ziehung
	var selectError error

	rows, selectError := databasehandle.Query("SELECT * from ziehungen")
	fmt.Printf("SELECT * FROM ziehungen\n")
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	for rows.Next() {
		var ziehung lottologic.Ziehung
		if scanError := rows.Scan(&ziehung.Datum, &ziehung.Ziehung); scanError != nil {
			return nil, scanError
		}
		ziehungen = append(ziehungen, ziehung)
	}

	if sqlError := rows.Err(); sqlError != nil {
		return nil, sqlError
	}

	return ziehungen, selectError

}

func InsertIntoAuszahlungen(ziehungTransaktion *sql.Tx, auszahlungen []lottologic.Auszahlung) error {

	var insertError error

	for _, auszahlung := range auszahlungen {
		_, insertError = ziehungTransaktion.Exec("INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (?,?,?,?)", auszahlung.Datum.Format("2006-01-02"), auszahlung.Klasse, auszahlung.Budget, auszahlung.Auszahlung)
		fmt.Printf("INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (%s,%d,%.2f,%.2f)\n", auszahlung.Datum.Format("2006-01-02"), auszahlung.Klasse, auszahlung.Budget, auszahlung.Auszahlung)
		if insertError != nil {
			break
		}
	}

	return insertError

}
