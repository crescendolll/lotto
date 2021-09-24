package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"lotto/lottologic"
	"os"
	"strconv"
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
	fmt.Println(configuration.FormatDSN())
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

func HoleZiehungenMitAuszahlungen(databasehandle *sql.DB, von time.Time, bis time.Time) ([]lottologic.Ziehungauszahlung, error) {

	var ziehungsauszahlungen []lottologic.Ziehungauszahlung
	var err error

	var ziehungen []lottologic.Ziehung

	var ziehungsauszahlung lottologic.Ziehungauszahlung

	var auszahlungen []lottologic.Auszahlung

	var gewinneranzahl [10]int

	var auszahlungsstatistiken []lottologic.Auszahlungsstatistik
	var auszahlungsstatistik lottologic.Auszahlungsstatistik

	ziehungen, err = SelectFromZiehungen(databasehandle, von, bis)

	if err != nil {
		return ziehungsauszahlungen, err
	}

	for _, ziehung := range ziehungen {
		ziehungsauszahlung.Datum = ziehung.Datum
		ziehungsauszahlung.Ziehung = ziehung.Ziehung

		auszahlungsstatistiken = make([]lottologic.Auszahlungsstatistik, 0)

		if ziehung.Ziehung.Ptr() != nil {
			berechnungstransaktion, _ := databasehandle.Begin()

			auszahlungen, err = SelectAuszahlungenByDate(databasehandle, ziehung.Datum)

			if err != nil {
				berechnungstransaktion.Rollback()
				return ziehungsauszahlungen, err
			}

			gewinneranzahl = lottologic.BerechneGewinneranzahl(berechnungstransaktion, ziehung.Datum)

			for _, auszahlung := range auszahlungen {
				auszahlungsstatistik.Klasse = auszahlung.Klasse
				auszahlungsstatistik.Gewinn = auszahlung.Auszahlung
				auszahlungsstatistik.Gewinner = gewinneranzahl[auszahlung.Klasse]

				auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
			}

			berechnungstransaktion.Commit()
		}

		ziehungsauszahlung.Auszahlungen = auszahlungsstatistiken

		ziehungsauszahlungen = append(ziehungsauszahlungen, ziehungsauszahlung)
	}

	return ziehungsauszahlungen, err
}

func SelectAuszahlungenByDate(databasehandle *sql.DB, datum time.Time) ([]lottologic.Auszahlung, error) {

	var auszahlungen []lottologic.Auszahlung
	var auszahlung lottologic.Auszahlung
	var err error

	rows, selectError := databasehandle.Query("SELECT klasse, auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "'")
	fmt.Printf("SELECT klasse, auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "'\n")
	if selectError != nil {
		return auszahlungen, selectError
	}
	defer rows.Close()

	for rows.Next() {
		if scanError := rows.Scan(&auszahlung.Klasse, &auszahlung.Auszahlung); scanError != nil {
			return auszahlungen, scanError
		}
		auszahlungen = append(auszahlungen, auszahlung)
	}

	if sqlError := rows.Err(); sqlError != nil {
		return auszahlungen, sqlError
	}
	return auszahlungen, err

}

func HoleTippauszahlungenFuerSpieler(databasehandle *sql.DB, spielername string, von time.Time, bis time.Time) ([]lottologic.Tippauszahlung, error) {

	var tippauszahlungen []lottologic.Tippauszahlung
	var err error

	var tipps []lottologic.Tipp

	var ziehung lottologic.Ziehung
	var auszahlung lottologic.Auszahlung
	var tippauszahlung lottologic.Tippauszahlung

	tipps, err = SelectFromTippsBySpielername(databasehandle, spielername, von, bis)

	tippauszahlungen = make([]lottologic.Tippauszahlung, 0)

	if err != nil {
		return nil, err
	}

	for _, tipp := range tipps {

		tippauszahlung.Id = tipp.Id
		tippauszahlung.Datum = tipp.Datum
		tippauszahlung.Ziehung = tipp.Ziehung

		ziehung, err = SelectFromZiehungByDate(databasehandle, tipp.Datum)
		if err != nil {
			break
		} else {
			tippauszahlung.Klasse = lottologic.BerechneGewinnklasse(tipp.Ziehung, ziehung.Ziehung.String)
		}

		auszahlung.Auszahlung = 0.0

		if tippauszahlung.Klasse >= 0 {
			auszahlung, err = SelectFromAuszahlungByDateAndClass(databasehandle, tipp.Datum, tippauszahlung.Klasse)
		}
		if err != nil {
			break
		}

		tippauszahlung.Auszahlung = auszahlung.Auszahlung

		tippauszahlungen = append(tippauszahlungen, tippauszahlung)

	}

	return tippauszahlungen, err

}

func SelectFromAuszahlungByDateAndClass(databasehandle *sql.DB, datum time.Time, klasse int8) (lottologic.Auszahlung, error) {

	var auszahlung lottologic.Auszahlung
	var err error

	rows, selectError := databasehandle.Query("SELECT auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "' AND klasse = " + strconv.Itoa(int(klasse)))
	fmt.Printf("SELECT auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "' AND klasse = " + strconv.Itoa(int(klasse)) + "\n")
	if selectError != nil {
		return auszahlung, selectError
	}
	defer rows.Close()

	for rows.Next() {
		if scanError := rows.Scan(&auszahlung.Auszahlung); scanError != nil {
			return auszahlung, scanError
		}
	}

	if sqlError := rows.Err(); sqlError != nil {
		return auszahlung, sqlError
	}
	return auszahlung, err

}

func SelectFromZiehungByDate(databasehandle *sql.DB, datum time.Time) (lottologic.Ziehung, error) {

	var ziehung lottologic.Ziehung
	var selectError error

	rows, selectError := databasehandle.Query("SELECT * from ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'")
	fmt.Printf("SELECT * FROM ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'\n")
	if selectError != nil {
		return ziehung, selectError
	}
	defer rows.Close()

	for rows.Next() {
		if scanError := rows.Scan(&ziehung.Datum, &ziehung.Ziehung); scanError != nil {
			return ziehung, scanError
		}
	}

	if sqlError := rows.Err(); sqlError != nil {
		return ziehung, sqlError
	}

	return ziehung, selectError
}

func InsertSpielerIntoNutzer(databasehandle *sql.DB, benutzername string, password string) error {

	var insertError error

	insertError = lottologic.IstNameVerfuegbar(databasehandle, benutzername)

	if insertError != nil {
		return insertError
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pw_hash, _ := lottologic.HashPassword(password)

	query := "INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (?,?,1)"
	statement, err := databasehandle.PrepareContext(context, query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(context, benutzername, pw_hash)
	return err
}

func InsertMitarbeiterIntoNutzer(databasehandle *sql.DB, benutzername string, password string) error {

	var insertError error

	insertError = lottologic.IstNameVerfuegbar(databasehandle, benutzername)

	if insertError != nil {
		return insertError
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pw_hash, _ := lottologic.HashPassword(password)

	query := "INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (?,?,0)"
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
		neuerHash, _ := lottologic.HashPassword(neuesPasswort)
		_, updateError = databasehandle.Exec("UPDATE nutzer set pw_hash = ?, benutzername = ? WHERE benutzername = ?", neuerHash, neuerBenutzername, alterBenutzername)
		fmt.Printf("UPDATE nutzer set pw_hash = %s, benutzername = %s WHERE benutzername = %s\n", neuerHash, neuerBenutzername, alterBenutzername)
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

	if !lottologic.IstValidesZiehungsdatum(databasehandle, tipp.Datum) {
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

func InsertIntoMitarbeiterZiehungen(ziehungsdatum time.Time, mitarbeitername string, aktion bool, tippTransaktion *sql.Tx) error {

	var insertError error

	_, insertError = tippTransaktion.Exec("INSERT INTO mitarbeiter_ziehungen (datum, mitarbeitername, aktion) VALUES (?,?,?)", ziehungsdatum.Format("2006-01-02"), mitarbeitername, aktion)
	fmt.Printf("INSERT INTO mitarbeiter_ziehungen (datum, mitarbeitername, aktion) VALUES (%s,%s,%t)\n", ziehungsdatum.Format("2006-01-02"), mitarbeitername, aktion)

	return insertError

}

func SelectFromTippsBySpielername(databasehandle *sql.DB, spielername string, von time.Time, bis time.Time) ([]lottologic.Tipp, error) {

	var tipps []lottologic.Tipp

	rows, selectError := databasehandle.Query("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = ? AND spieler_tipps.id = tipps.id AND tipps.datum >= ? AND tipps.datum <= ?", spielername, von, bis)
	fmt.Printf("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = %s AND spieler_tipps.id = tipps.id AND tipps.datum >= %s AND tipps.datum <= %s\n", spielername, von.Format("2006-01-02"), bis.Format("2006-01-02"))
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

func InsertIntoZiehungen(databasehandle *sql.DB, ziehung lottologic.Ziehung, mitarbeiter lottologic.Nutzer) error {

	var insertError error

	if !lottologic.IstValidesZiehungsdatum(databasehandle, ziehung.Datum) {
		insertError = errors.New("Kein gültiges Tippdatum: " + ziehung.Datum.Format("2006-01-02"))
		return insertError
	}

	fmt.Println("Transaktion beginnt")
	ziehungTransaktion, _ := databasehandle.Begin()

	_, insertError = ziehungTransaktion.Exec("INSERT INTO ziehungen (datum) VALUES (?)", ziehung.Datum.Format("2006-01-02"))
	fmt.Printf("INSERT INTO ziehungen (datum) VALUES (%s)\n", ziehung.Datum.Format("2006-01-02"))

	if insertError == nil {
		insertError = InsertIntoMitarbeiterZiehungen(ziehung.Datum, mitarbeiter.Benutzername, false, ziehungTransaktion)
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

func UpdateZiehungen(databasehandle *sql.DB, ziehung lottologic.Ziehung, mitarbeiter lottologic.Nutzer) error {

	var updateError error
	var auszahlungen []lottologic.Auszahlung

	if !lottologic.IstValideZiehung(ziehung.Ziehung.String) {
		updateError = errors.New(ziehung.Ziehung.String + " ist keine gültige Ziehung")
		return updateError
	}

	if !lottologic.IstOffenesZiehungsdatum(databasehandle, ziehung.Datum) {
		updateError = errors.New(ziehung.Datum.Format("2006-01-02") + " ist kein Datum eines offenen Spiels")
		return updateError
	}

	fmt.Println("Transaktion beginnt")
	ziehungTransaktion, _ := databasehandle.Begin()

	_, updateError = ziehungTransaktion.Exec("UPDATE ziehungen SET ziehung = ? WHERE datum = ?", ziehung.Ziehung, ziehung.Datum.Format("2006-01-02"))
	fmt.Printf("UPDATE ziehungen SET ziehung = %s WHERE datum = %s\n", ziehung.Ziehung.String, ziehung.Datum.Format("2006-01-02"))

	if updateError == nil {
		updateError = InsertIntoMitarbeiterZiehungen(ziehung.Datum, mitarbeiter.Benutzername, true, ziehungTransaktion)
	}

	if updateError == nil {
		auszahlungen, updateError = lottologic.BerechneAuszahlungen(ziehungTransaktion, ziehung.Datum)
	}

	if updateError == nil {
		updateError = InsertIntoAuszahlungen(ziehungTransaktion, auszahlungen)
	}

	if updateError != nil {
		fmt.Println("Transaktion abgebrochen")
		ziehungTransaktion.Rollback()
	} else {
		fmt.Println("Transaktion erfolgt")
		ziehungTransaktion.Commit()
	}

	return updateError

}

func SelectLaufendeZiehungen(databasehandle *sql.DB) ([]lottologic.Ziehung, error) {

	var ziehungen []lottologic.Ziehung
	var selectError error

	rows, selectError := databasehandle.Query("SELECT * from ziehungen WHERE ziehung IS NULL")
	fmt.Printf("SELECT * FROM ziehungen WHERE ziehung IS NULL\n")
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

func SelectFromZiehungen(databasehandle *sql.DB, von time.Time, bis time.Time) ([]lottologic.Ziehung, error) {

	var ziehungen []lottologic.Ziehung
	var selectError error

	rows, selectError := databasehandle.Query("SELECT * from ziehungen WHERE datum >= '" + von.Format("2006-01-02") + "' AND datum <= '" + bis.Format("2006-01-02") + "'")
	fmt.Printf("SELECT * from ziehungen WHERE datum >= '" + von.Format("2006-01-02") + "' AND datum <= '" + bis.Format("2006-01-02") + "'\n")
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
