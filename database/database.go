package database

import (
	"context"
	"database/sql"
	"log"
	"lotto/lottolog"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/guregu/null"
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
	Ziehung null.String
}

type Auszahlung struct {
	Datum      time.Time
	Klasse     int8
	Budget     float64
	Auszahlung float64
}

var datenbankVerbindung *sql.DB

func HoleVerbindung() *sql.DB {
	return datenbankVerbindung
}

func SetzeVerbindung(verbindung *sql.DB) {
	datenbankVerbindung = verbindung
}

func SchliesseVerbindung() {
	fehler := datenbankVerbindung.Close()
	if fehler != nil {
		lottolog.FehlerLogger.Fatal(fehler)
	}
}

func OeffneVerbindungZurLottoDatenbank() {

	// Verbindung konfigurieren, die Zugangsdaten zur Datenbank sind als Umgebungsvariablen zu setzen
	konfiguration := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "127.0.0.1:" + os.Getenv("DBPORT"),
		DBName:    "lotto",
		ParseTime: true,
	}

	// Verbindung öffnen
	verbindung, zugriffsfehler := sql.Open("mysql", konfiguration.FormatDSN())

	if zugriffsfehler != nil {
		lottolog.FehlerLogger.Fatal(zugriffsfehler)
	}

	// Verbindungstest
	verbindungsfehler := verbindung.Ping()
	if verbindungsfehler != nil {
		lottolog.FehlerLogger.Fatal(verbindungsfehler)
	} else {
		datenbankVerbindung = verbindung
	}

}

func FuegeNutzerEin(nutzer Nutzer) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (?,?,?)"
	lottolog.InfoLogger.Printf("INSERT INTO nutzer (benutzername, pw_hash, ist_spieler) VALUES (%s,%s,%t)", nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler)
	statement, fehler := datenbankVerbindung.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler)
	return fehler

}

func HoleNutzerdatenZumNamen(benutzername string) (Nutzer, error) {

	var nutzer Nutzer
	var abfrageFehler error

	datensaetze, abfrageFehler := datenbankVerbindung.Query("SELECT * from nutzer WHERE benutzername = ?", benutzername)
	lottolog.InfoLogger.Printf("SELECT * from nutzer WHERE benutzername = %s\n", benutzername)
	if abfrageFehler != nil {
		return nutzer, abfrageFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		if scanFehler := datensaetze.Scan(&nutzer.Benutzername, &nutzer.Pw_hash, &nutzer.Ist_spieler); scanFehler != nil {
			return nutzer, scanFehler
		}
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return nutzer, sqlFehler
	}

	return nutzer, abfrageFehler

}

func HoleVerfuegbarkeitEinesBenutzernamens(benutzername string) (bool, error) {

	var ist_verfuegbar bool

	datensaetze, fehler := datenbankVerbindung.Query("SELECT (COUNT(*) = 0) as verfuegbar FROM nutzer WHERE benutzername = ?", benutzername)
	lottolog.InfoLogger.Printf("SELECT (COUNT(*) = 0) as verfuegbar FROM nutzer WHERE benutzername = %s\n", benutzername)

	if fehler == nil {
		for datensaetze.Next() {
			datensaetze.Scan(&ist_verfuegbar)
		}
	}

	return ist_verfuegbar, fehler
}

func AendereNutzerdaten(neueNutzerdaten Nutzer, alterBenutzername string) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "UPDATE nutzer set pw_hash = ?, benutzername = ? WHERE benutzername = ?"
	lottolog.InfoLogger.Printf("UPDATE nutzer set pw_hash = %s, benutzername = %s WHERE benutzername = %s", neueNutzerdaten.Pw_hash, neueNutzerdaten.Benutzername, alterBenutzername)
	statement, fehler := datenbankVerbindung.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, neueNutzerdaten.Pw_hash, neueNutzerdaten.Benutzername, alterBenutzername)
	return fehler

}

func LoescheNutzerdaten(benutzername string) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "DELETE from nutzer WHERE benutzername = ?"
	lottolog.InfoLogger.Printf("DELETE from nutzer WHERE benutzername = %s", benutzername)
	statement, fehler := datenbankVerbindung.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, benutzername)
	return fehler

}

func FuegeTippEin(tipp Tipp, transaktion *sql.Tx) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO tipps (id, datum, ziehung) VALUES (?,?,?)"
	lottolog.InfoLogger.Printf("INSERT INTO tipps (id, datum, ziehung) VALUES (%d,%s,%s)\n", tipp.Id, tipp.Datum.Format("2006-01-02"), tipp.Ziehung)
	statement, fehler := transaktion.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, tipp.Id, tipp.Datum, tipp.Ziehung)
	return fehler

}

func HoleGroessteTippID() (int64, error) {

	var maxID int64

	abfrage := "SELECT MAX(id) as maxID FROM tipps"

	datensatz, fehler := datenbankVerbindung.Query(abfrage)
	lottolog.InfoLogger.Println(abfrage)

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler)
	} else {
		for datensatz.Next() {
			datensatz.Scan(&maxID)
		}
	}

	return maxID, fehler
}

func HoleTippanzahlZumDatum(ziehungsdatum time.Time, transaktion *sql.Tx) (int, error) {

	var tippanzahl int

	abfrage := "SELECT COUNT(*) as anzahl FROM tipps WHERE datum = '" + ziehungsdatum.Format("2006-01-02") + "'"

	datensaetze, fehler := transaktion.Query(abfrage)
	lottolog.InfoLogger.Println(abfrage)

	if fehler == nil {
		for datensaetze.Next() {
			datensaetze.Scan(&tippanzahl)
		}
	}

	return tippanzahl, fehler

}

func HoleTippsVonSpielerImZeitraum(spielername string, startdatum time.Time, enddatum time.Time) ([]Tipp, error) {

	var tipps []Tipp

	datensaetze, selektionFehler := datenbankVerbindung.Query("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = ? AND spieler_tipps.id = tipps.id AND tipps.datum >= ? AND tipps.datum <= ?", spielername, startdatum.Format("2006-01-02"), enddatum.Format("2006-01-02"))
	lottolog.InfoLogger.Printf("SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps "+
		"WHERE spieler_tipps.spielername = %s AND spieler_tipps.id = tipps.id AND tipps.datum >= %s AND tipps.datum <= %s\n", spielername, startdatum.Format("2006-01-02"), enddatum.Format("2006-01-02"))
	if selektionFehler != nil {
		return nil, selektionFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		var tipp Tipp
		if scanFehler := datensaetze.Scan(&tipp.Id, &tipp.Datum, &tipp.Ziehung); scanFehler != nil {
			return nil, scanFehler
		}
		tipps = append(tipps, tipp)
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return nil, sqlFehler
	}
	return tipps, nil
}

func HoleTippsZumDatum(ziehungsdatum time.Time) ([]Tipp, error) {

	var tipps []Tipp
	var tipp Tipp

	abfrage := "SELECT ziehung FROM tipps WHERE datum = '" + ziehungsdatum.Format("2006-01-02") + "'"

	datensaetze, fehler := datenbankVerbindung.Query(abfrage)
	lottolog.InfoLogger.Println(abfrage)

	tipps = make([]Tipp, 0)

	if fehler == nil {
		for datensaetze.Next() {
			datensaetze.Scan(&tipp.Ziehung)
			tipps = append(tipps, tipp)
		}
	}

	return tipps, fehler

}

func FuegeSpielerTippVerknuepfungEin(tippID int64, spielername string, transaktion *sql.Tx) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO spieler_tipps (id, spielername) VALUES (?,?)"
	lottolog.InfoLogger.Printf("INSERT INTO spieler_tipps (id, spielername) VALUES (%d,%s)\n", tippID, spielername)
	statement, fehler := transaktion.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, tippID, spielername)
	return fehler

}

func FuegeZiehungEin(ziehung Ziehung, transaktion *sql.Tx) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO ziehungen (datum) VALUES (?)"
	lottolog.InfoLogger.Printf("INSERT INTO ziehungen (datum) VALUES (%s)\n", ziehung.Datum.Format("2006-01-02"))
	statement, fehler := transaktion.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, ziehung.Datum.Format("2006-01-02"))
	return fehler

}

func HoleOffeneZiehungen() ([]Ziehung, error) {

	var ziehungen []Ziehung
	var selektionFehler error

	datensaetze, selektionFehler := datenbankVerbindung.Query("SELECT * from ziehungen WHERE ziehung IS NULL")
	lottolog.InfoLogger.Printf("SELECT * FROM ziehungen WHERE ziehung IS NULL\n")
	if selektionFehler != nil {
		return nil, selektionFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		var ziehung Ziehung
		if scanFehler := datensaetze.Scan(&ziehung.Datum, &ziehung.Ziehung); scanFehler != nil {
			return nil, scanFehler
		}
		ziehungen = append(ziehungen, ziehung)
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return nil, sqlFehler
	}

	return ziehungen, selektionFehler

}

func HoleLetztesZiehungsdatum() (time.Time, error) {

	var letzteZiehung time.Time

	abfrage := "SELECT MAX(datum) as letzteZiehung FROM ziehungen WHERE ziehung IS NOT NULL"
	lottolog.InfoLogger.Println(abfrage)

	datensaetze, fehler := datenbankVerbindung.Query(abfrage)

	if fehler != nil {
		return letzteZiehung, fehler
	} else {
		for datensaetze.Next() {
			datensaetze.Scan(&letzteZiehung)
		}
	}

	return letzteZiehung, fehler
}

func HoleVerfuegbarkeitEinerZiehungZumDatum(tippdatum time.Time) (bool, error) {

	var offen bool

	abfrage := "SELECT (COUNT(*) > 0) as verfuegbar FROM ziehungen WHERE datum = '" + tippdatum.Format("2006-01-02") + "' AND ziehung IS NULL"
	lottolog.InfoLogger.Println(abfrage)

	datensaetze, fehler := datenbankVerbindung.Query(abfrage)

	if fehler == nil {
		for datensaetze.Next() {
			datensaetze.Scan(&offen)
		}
	}

	return offen, fehler
}

func HoleZiehungZumDatum(datum time.Time) (Ziehung, error) {

	var ziehung Ziehung
	var selektionFehler error

	datensaetze, selektionFehler := datenbankVerbindung.Query("SELECT * from ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'")
	lottolog.InfoLogger.Printf("SELECT * FROM ziehungen WHERE datum = '" + datum.Format("2006-01-02") + "'\n")
	if selektionFehler != nil {
		return ziehung, selektionFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		if scanFehler := datensaetze.Scan(&ziehung.Datum, &ziehung.Ziehung); scanFehler != nil {
			return ziehung, scanFehler
		}
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return ziehung, sqlFehler
	}

	return ziehung, selektionFehler
}

func HoleZiehungenInnerhalbEinesZeitraums(startdatum time.Time, enddatum time.Time) ([]Ziehung, error) {

	var ziehungen []Ziehung
	var selektionFehler error

	datensaetze, selektionFehler := datenbankVerbindung.Query("SELECT * from ziehungen WHERE datum >= '" + startdatum.Format("2006-01-02") + "' AND datum <= '" + enddatum.Format("2006-01-02") + "'")
	lottolog.InfoLogger.Printf("SELECT * from ziehungen WHERE datum >= '" + startdatum.Format("2006-01-02") + "' AND datum <= '" + enddatum.Format("2006-01-02") + "'\n")
	if selektionFehler != nil {
		return nil, selektionFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		var ziehung Ziehung
		if scanFehler := datensaetze.Scan(&ziehung.Datum, &ziehung.Ziehung); scanFehler != nil {
			return nil, scanFehler
		}
		ziehungen = append(ziehungen, ziehung)
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return nil, sqlFehler
	}

	return ziehungen, selektionFehler

}

func AendereZiehung(ziehung Ziehung, transaktion *sql.Tx) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "UPDATE ziehungen SET ziehung = ? WHERE datum = ?"
	lottolog.InfoLogger.Printf("UPDATE ziehungen SET ziehung = %s WHERE datum = %s\n", ziehung.Ziehung.String, ziehung.Datum.Format("2006-01-02"))
	statement, fehler := datenbankVerbindung.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, ziehung.Ziehung, ziehung.Datum.Format("2006-01-02"))
	return fehler

}

func FuegeMitarbeiterZiehungVerknuepfungEin(ziehungsdatum time.Time, mitarbeitername string, aktion bool, transaktion *sql.Tx) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO mitarbeiter_ziehungen (datum, mitarbeitername, aktion) VALUES (?,?,?)"
	lottolog.InfoLogger.Printf("INSERT INTO mitarbeiter_ziehungen (datum, mitarbeitername, aktion) VALUES (%s,%s,%t)\n", ziehungsdatum.Format("2006-01-02"), mitarbeitername, aktion)
	statement, fehler := transaktion.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, ziehungsdatum.Format("2006-01-02"), mitarbeitername, aktion)
	return fehler

}

func FuegeAuszahlungEin(transaktion *sql.Tx, auszahlung Auszahlung) error {

	kontext, abbruch := context.WithTimeout(context.Background(), 5*time.Second)
	defer abbruch()

	abfrage := "INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (?,?,?,?)"
	lottolog.InfoLogger.Printf("INSERT INTO auszahlungen (datum, klasse, budget, auszahlung) VALUES (%s,%d,%.2f,%.2f)\n", auszahlung.Datum.Format("2006-01-02"), auszahlung.Klasse, auszahlung.Budget, auszahlung.Auszahlung)
	statement, fehler := transaktion.PrepareContext(kontext, abfrage)
	if fehler != nil {
		return fehler
	}
	defer statement.Close()

	_, fehler = statement.ExecContext(kontext, auszahlung.Datum.Format("2006-01-02"), auszahlung.Klasse, auszahlung.Budget, auszahlung.Auszahlung)
	return fehler

}

func HoleAuszahlungenZumDatum(datum time.Time) ([]Auszahlung, error) {

	var auszahlungen []Auszahlung
	var auszahlung Auszahlung

	datensaetze, selectFehler := datenbankVerbindung.Query("SELECT klasse, auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "'")
	lottolog.InfoLogger.Printf("SELECT klasse, auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "'\n")
	if selectFehler != nil {
		lottolog.FehlerLogger.Println(selectFehler.Error())
		return auszahlungen, selectFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		if scanFehler := datensaetze.Scan(&auszahlung.Klasse, &auszahlung.Auszahlung); scanFehler != nil {
			lottolog.FehlerLogger.Println(scanFehler.Error())
			return auszahlungen, scanFehler
		}
		auszahlungen = append(auszahlungen, auszahlung)
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		log.Println(sqlFehler.Error())
		return auszahlungen, sqlFehler
	}

	return auszahlungen, nil

}

func HoleAuszahlungZuDatumUndKlasse(datum time.Time, klasse int8) (Auszahlung, error) {

	var auszahlung Auszahlung

	datensaetze, selectFehler := datenbankVerbindung.Query("SELECT auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "' AND klasse = " + strconv.Itoa(int(klasse)))
	lottolog.InfoLogger.Printf("SELECT auszahlung from auszahlungen WHERE datum = '" + datum.Format("2006-01-02") + "' AND klasse = " + strconv.Itoa(int(klasse)) + "\n")
	if selectFehler != nil {
		return auszahlung, selectFehler
	}
	defer datensaetze.Close()

	for datensaetze.Next() {
		if scanFehler := datensaetze.Scan(&auszahlung.Auszahlung); scanFehler != nil {
			return auszahlung, scanFehler
		}
	}

	if sqlFehler := datensaetze.Err(); sqlFehler != nil {
		return auszahlung, sqlFehler
	}
	return auszahlung, nil

}

func HoleLetztenJackpot(transaktion *sql.Tx) (Auszahlung, error) {

	var jackpot Auszahlung

	abfrage := "SELECT budget, auszahlung FROM auszahlungen WHERE klasse = 9 AND datum = (SELECT MAX(datum) FROM auszahlungen)"
	lottolog.InfoLogger.Println(abfrage)

	datensaetze, fehler := transaktion.Query(abfrage)

	if fehler == nil {
		for datensaetze.Next() {
			datensaetze.Scan(&jackpot.Budget, &jackpot.Auszahlung)
		}
	}

	return jackpot, fehler
}
