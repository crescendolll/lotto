package lottologic

import (
	"database/sql"
	"errors"
	"lotto/database"
	"lotto/lottolog"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type Tippauszahlung struct {
	Id         int64
	Datum      time.Time
	Ziehung    string
	Klasse     int8
	Auszahlung float64
}

type Ziehungsstatistik struct {
	Datum        time.Time
	Ziehung      null.String
	Auszahlungen []Auszahlungsstatistik
}

type Auszahlungsstatistik struct {
	Klasse   int8
	Gewinner int
	Gewinn   float64
}

var ErstelleAuszahlungsstatistiken = func(ziehung database.Ziehung) ([]Auszahlungsstatistik, error) {

	var auszahlungsstatistiken []Auszahlungsstatistik
	var auszahlungsstatistik Auszahlungsstatistik
	var auszahlungen []database.Auszahlung
	var gewinneranzahl [10]int
	var fehler error

	auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

	auszahlungen, fehler = database.HoleAuszahlungenZumDatum(ziehung.Datum)

	if fehler != nil {
		return auszahlungsstatistiken, fehler
	} else {
		gewinneranzahl, fehler = BerechneGewinneranzahl(ziehung.Datum)
	}

	if fehler != nil {
		return auszahlungsstatistiken, fehler
	}

	for _, auszahlung := range auszahlungen {
		auszahlungsstatistik.Klasse = auszahlung.Klasse
		auszahlungsstatistik.Gewinn = auszahlung.Auszahlung
		auszahlungsstatistik.Gewinner = gewinneranzahl[auszahlung.Klasse]

		auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
	}

	return auszahlungsstatistiken, fehler

}

var ErstelleZiehungsstatistiken = func(ziehungen []database.Ziehung) ([]Ziehungsstatistik, error) {

	var ziehungsstatistiken []Ziehungsstatistik
	var ziehungsstatistik Ziehungsstatistik
	var auszahlungsstatistiken []Auszahlungsstatistik
	var fehler error

	ziehungsstatistiken = make([]Ziehungsstatistik, 0)

	for _, ziehung := range ziehungen {

		if ziehung.Ziehung.Ptr() != nil {
			auszahlungsstatistiken, fehler = ErstelleAuszahlungsstatistiken(ziehung)

			if fehler != nil {
				return ziehungsstatistiken, fehler
			} else {
				ziehungsstatistik.Datum = ziehung.Datum
				ziehungsstatistik.Ziehung = ziehung.Ziehung
				ziehungsstatistik.Auszahlungen = auszahlungsstatistiken
				ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)
			}

		}

	}

	return ziehungsstatistiken, fehler
}

func ErstelleZiehungsstatistikenFuerZeitraum(startdatum time.Time, enddatum time.Time) ([]Ziehungsstatistik, error) {

	var ziehungsstatistiken []Ziehungsstatistik
	var ziehungen []database.Ziehung
	var fehler error

	ziehungen, fehler = database.HoleZiehungenInnerhalbEinesZeitraums(startdatum, enddatum)

	if fehler != nil {
		return ziehungsstatistiken, fehler
	} else {
		ziehungsstatistiken, fehler = ErstelleZiehungsstatistiken(ziehungen)
	}

	return ziehungsstatistiken, fehler
}

func ErstelleTippauszahlungen(tipps []database.Tipp) ([]Tippauszahlung, error) {

	var tippauszahlungen []Tippauszahlung
	var tippauszahlung Tippauszahlung
	var ziehung database.Ziehung
	var auszahlung database.Auszahlung
	var fehler error

	tippauszahlungen = make([]Tippauszahlung, 0)

	for _, tipp := range tipps {

		tippauszahlung.Id = tipp.Id
		tippauszahlung.Datum = tipp.Datum
		tippauszahlung.Ziehung = tipp.Ziehung

		ziehung, fehler = database.HoleZiehungZumDatum(tipp.Datum)
		if fehler != nil {
			lottolog.FehlerLogger.Println(fehler)
			break
		} else {
			tippauszahlung.Klasse = BerechneGewinnklasse(tipp.Ziehung, ziehung.Ziehung.String)
		}

		auszahlung.Auszahlung = 0.0

		if tippauszahlung.Klasse >= 0 {
			auszahlung, fehler = database.HoleAuszahlungZuDatumUndKlasse(tipp.Datum, tippauszahlung.Klasse)
		}
		if fehler != nil {
			lottolog.FehlerLogger.Println(fehler)
			break
		} else {
			tippauszahlung.Auszahlung = auszahlung.Auszahlung
			tippauszahlungen = append(tippauszahlungen, tippauszahlung)
		}

	}

	return tippauszahlungen, fehler
}

func ErstelleTippauszahlungenFuerSpieler(spielername string, startdatum time.Time, enddatum time.Time) ([]Tippauszahlung, error) {

	var tippauszahlungen []Tippauszahlung
	var tipps []database.Tipp
	var fehler error

	tipps, fehler = database.HoleTippsVonSpielerImZeitraum(spielername, startdatum, enddatum)

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler)
		return nil, fehler
	}

	tippauszahlungen, fehler = ErstelleTippauszahlungen(tipps)

	return tippauszahlungen, fehler

}

func SchliesseZiehung(ziehung database.Ziehung, mitarbeiter database.Nutzer) error {

	var updateFehler error
	var ziehungTransaktion *sql.Tx
	var auszahlungen []database.Auszahlung

	if !IstValideZiehung(ziehung.Ziehung.String) {
		updateFehler = errors.New(ziehung.Ziehung.String + " ist keine gültige Ziehung")
		return updateFehler
	}

	if !IstSchliessbaresZiehungsdatum(ziehung.Datum) {
		updateFehler = errors.New(ziehung.Datum.Format("2006-01-02") + " ist kein Datum einer abschliessbaren Ziehung")
		return updateFehler
	}

	lottolog.InfoLogger.Println("Transaktion zur Schliessung einer Ziehung beginnt")
	ziehungTransaktion, updateFehler = database.HoleVerbindung().Begin()

	if updateFehler == nil {
		updateFehler = database.AendereZiehung(ziehung, ziehungTransaktion)
	}

	if updateFehler == nil {
		updateFehler = database.FuegeMitarbeiterZiehungVerknuepfungEin(ziehung.Datum, mitarbeiter.Benutzername, true, ziehungTransaktion)
	}

	if updateFehler == nil {
		auszahlungen, updateFehler = BerechneAuszahlungen(ziehungTransaktion, ziehung.Datum)
	}

	if updateFehler == nil {
		for _, auszahlung := range auszahlungen {
			updateFehler = database.FuegeAuszahlungEin(ziehungTransaktion, auszahlung)
			if updateFehler != nil {
				break
			}
		}
	}

	if updateFehler != nil {
		lottolog.FehlerLogger.Println(updateFehler.Error())
		lottolog.WarnungLogger.Println("Transaktion zur Schliessung einer Ziehung abgebrochen")
		ziehungTransaktion.Rollback()
	} else {
		lottolog.InfoLogger.Println("Transaktion zur Schliessung einer Ziehung erfolgt")
		ziehungTransaktion.Commit()
	}

	return updateFehler

}

func FuegeSpielerNachPruefungEin(nutzer database.Nutzer, passwort string) error {

	var hash string

	ist_verfuegbar, fehler := database.HoleVerfuegbarkeitEinesBenutzernamens(nutzer.Benutzername)

	if fehler == nil {
		hash, fehler = HashePasswort(passwort)
	}

	if fehler == nil {
		if ist_verfuegbar {
			nutzer.Pw_hash = hash
			fehler = database.FuegeNutzerEin(nutzer)
		} else {
			fehler = errors.New("Name " + nutzer.Benutzername + " bereits vergeben")
		}
	}

	return fehler
}

func FuegeMitarbeiterNachPruefungEin(benutzername string, passwort string) error {

	var hash string
	var mitarbeiter database.Nutzer

	ist_verfuegbar, fehler := database.HoleVerfuegbarkeitEinesBenutzernamens(benutzername)

	if fehler == nil {
		if ist_verfuegbar {
			hash, fehler = HashePasswort(passwort)
		} else {
			fehler = errors.New("Name " + benutzername + " bereits vergeben")
		}
	} else {
		return fehler
	}

	if fehler != nil {
		return fehler
	} else {
		mitarbeiter = database.Nutzer{
			Benutzername: benutzername,
			Pw_hash:      hash,
			Ist_spieler:  false,
		}
		fehler = database.FuegeNutzerEin(mitarbeiter)
	}

	return fehler
}

func FuegeTippNachPruefungEin(tipp database.Tipp, spieler database.Nutzer) error {

	var fehler error
	var tippTransaktion *sql.Tx

	if !IstValideZiehung(tipp.Ziehung) {
		fehler = errors.New("Kein gueltiger Lottotipp: " + tipp.Ziehung)
		return fehler
	}

	if !IstOffenesZiehungsdatum(tipp.Datum) {
		fehler = errors.New("Kein gueltiges Tippdatum: " + tipp.Datum.Format("2006-01-02"))
		return fehler
	}

	tipp.Id, fehler = database.HoleGroessteTippID()

	if fehler != nil {
		return fehler
	} else {
		tipp.Id += 1
	}

	lottolog.InfoLogger.Println("Transaktion zur Abgabe eines Tipps beginnt")
	tippTransaktion, fehler = database.HoleVerbindung().Begin()

	if fehler == nil {
		fehler = database.FuegeTippEin(tipp, tippTransaktion)
	}

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler.Error())
		lottolog.WarnungLogger.Println("Transaktion zur Abgabe eines Tipps abgebrochen")
		tippTransaktion.Rollback()
	} else {
		fehler = database.FuegeSpielerTippVerknuepfungEin(tipp.Id, spieler.Benutzername, tippTransaktion)
	}

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler.Error())
		lottolog.WarnungLogger.Println("Transaktion zur Abgabe eines Tipps abgebrochen")
		tippTransaktion.Rollback()
	} else {
		lottolog.InfoLogger.Println("Transaktion zur Abgabe eines Tipps erfolgt")
		tippTransaktion.Commit()
	}

	return fehler

}

func EroeffneZiehung(ziehung database.Ziehung, mitarbeiter database.Nutzer) error {

	var fehler error

	if !IstGueltigesZiehungsdatum(ziehung.Datum) {
		fehler = errors.New("Kein gültiges Ziehungsdatum: " + ziehung.Datum.Format("2006-01-02"))
		return fehler
	}

	lottolog.InfoLogger.Println("Transaktion zum Eroeffnen einer Ziehung beginnt")
	ziehungTransaktion, fehler := database.HoleVerbindung().Begin()

	if fehler == nil {
		fehler = database.FuegeZiehungEin(ziehung, ziehungTransaktion)
	}

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler.Error())
		lottolog.WarnungLogger.Println("Transaktion zum Eroeffnen einer Ziehung abgebrochen")
		ziehungTransaktion.Rollback()
	} else {
		fehler = database.FuegeMitarbeiterZiehungVerknuepfungEin(ziehung.Datum, mitarbeiter.Benutzername, false, ziehungTransaktion)
	}

	if fehler != nil {
		lottolog.FehlerLogger.Println(fehler.Error())
		lottolog.WarnungLogger.Println("Transaktion zum Eroeffnen einer Ziehung abgebrochen")
		ziehungTransaktion.Rollback()
	} else {
		lottolog.InfoLogger.Println("Transaktion erfolgt")
		ziehungTransaktion.Commit()
	}

	return fehler

}

func AendereSpielerdatenNachPruefung(benutzername string, neuesPasswort string) error {

	var fehler error
	var neuerHash string

	neuerHash, fehler = HashePasswort(neuesPasswort)

	if fehler == nil {
		neueNutzerdaten := database.Nutzer{
			Benutzername: benutzername,
			Pw_hash:      neuerHash,
		}
		fehler = database.AendereNutzerdaten(neueNutzerdaten, benutzername)
	}

	return fehler

}

func ErzeugeZufallsziehung() string {

	var ziehung string
	var zufallszahl int
	var ziehungszahlen []int

	ziehung = ""

	rand.Seed(time.Now().UnixNano())

	ziehungszahlen = []int{}

	for index := 1; index < 7; index++ {

		zufallszahl = rand.Intn(49) + 1

		for IstZahlInSlice(ziehungszahlen, zufallszahl) {
			zufallszahl = rand.Intn(49) + 1
		}

		if zufallszahl < 10 {
			ziehung = ziehung + "0" + strconv.Itoa(zufallszahl)
		} else {
			ziehung = ziehung + strconv.Itoa(zufallszahl)
		}

		ziehungszahlen = append(ziehungszahlen, zufallszahl)
	}

	zufallszahl = rand.Intn(10)
	ziehung = ziehung + strconv.Itoa(zufallszahl)

	return ziehung

}

func IstValideZiehung(ziehung string) bool {

	var valid bool
	var fehler error

	valid, fehler = regexp.MatchString("(0[1-9]|[1-4][0-9]){6}[0-9]", ziehung)

	valid = valid && (len(ziehung) == 13)

	if fehler != nil {
		lottolog.WarnungLogger.Println(fehler.Error())
		valid = false
	}

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

func IstGueltigesZiehungsdatum(ziehungsdatum time.Time) bool {

	var gueltig bool
	var letzteZiehung time.Time

	gueltig = false

	letzteZiehung, fehler := database.HoleLetztesZiehungsdatum()

	if fehler != nil {
		lottolog.WarnungLogger.Println(fehler.Error())
		gueltig = false
	} else {
		gueltig = ziehungsdatum.After(letzteZiehung.AddDate(0, 0, 1)) && ziehungsdatum.After(time.Now())
	}

	return gueltig

}

func IstOffenesZiehungsdatum(tippdatum time.Time) bool {

	offen, fehler := database.HoleVerfuegbarkeitEinerZiehungZumDatum(tippdatum)

	if fehler != nil {
		lottolog.WarnungLogger.Println(fehler.Error())
		offen = false
	}

	if offen {
		offen = tippdatum.After(time.Now())
	}

	return offen

}

func IstSchliessbaresZiehungsdatum(tippdatum time.Time) bool {

	offen, fehler := database.HoleVerfuegbarkeitEinerZiehungZumDatum(tippdatum)

	if fehler != nil {
		lottolog.WarnungLogger.Println(fehler.Error())
		offen = false
	}

	if offen {
		offen = tippdatum.Before(time.Now())
	}

	return offen

}

func IstZahlInSlice(slice []int, gesucht int) bool {
	for _, zahl := range slice {
		if zahl == gesucht {
			return true
		}
	}
	return false
}

func BerechneAuszahlungen(ziehungTransaktion *sql.Tx, ziehungsdatum time.Time) ([]database.Auszahlung, error) {

	var auszahlungen []database.Auszahlung
	var einsatz float64
	var tippanzahl int
	var fehler error

	lospreis := 5.0

	tippanzahl, fehler = database.HoleTippanzahlZumDatum(ziehungsdatum, ziehungTransaktion)

	if fehler != nil {
		return auszahlungen, fehler
	}

	einsatz = float64(tippanzahl) * lospreis

	var gewinneranzahl [10]int

	gewinneranzahl, fehler = BerechneGewinneranzahl(ziehungsdatum)

	if fehler != nil {
		return auszahlungen, fehler
	}

	zweirichtige := database.Auszahlung{
		Datum:      ziehungsdatum,
		Klasse:     0,
		Auszahlung: 2.0,
	}

	zweirichtige.Budget = zweirichtige.Auszahlung * float64(gewinneranzahl[zweirichtige.Klasse])
	auszahlungen = append(auszahlungen, zweirichtige)
	einsatz = einsatz - zweirichtige.Budget

	zweirichtigesuper := database.Auszahlung{
		Datum:      ziehungsdatum,
		Klasse:     1,
		Auszahlung: 5.0,
	}

	zweirichtigesuper.Budget = zweirichtigesuper.Auszahlung * float64(gewinneranzahl[zweirichtigesuper.Klasse])
	auszahlungen = append(auszahlungen, zweirichtigesuper)
	einsatz = einsatz - zweirichtigesuper.Budget

	sechsrichtigesuper := BerechneAuszahlung(einsatz, gewinneranzahl, 0.128, ziehungsdatum, 9)
	einsatz = einsatz - sechsrichtigesuper.Budget

	sechsrichtigesuper, fehler = BerechneJackpotsteigerung(ziehungTransaktion, sechsrichtigesuper)
	auszahlungen = append(auszahlungen, sechsrichtigesuper)

	if fehler != nil {
		return auszahlungen, fehler
	}

	sechsrichtige := BerechneAuszahlung(einsatz, gewinneranzahl, 0.1, ziehungsdatum, 8)
	auszahlungen = append(auszahlungen, sechsrichtige)

	fuenfrichtige := BerechneAuszahlung(einsatz, gewinneranzahl, 0.15, ziehungsdatum, 6)
	auszahlungen = append(auszahlungen, fuenfrichtige)

	fuenfrichtigesuper := BerechneAuszahlung(einsatz, gewinneranzahl, 0.05, ziehungsdatum, 7)
	auszahlungen = append(auszahlungen, fuenfrichtigesuper)

	vierrichtige := BerechneAuszahlung(einsatz, gewinneranzahl, 0.1, ziehungsdatum, 4)
	auszahlungen = append(auszahlungen, vierrichtige)

	vierrichtigesuper := BerechneAuszahlung(einsatz, gewinneranzahl, 0.05, ziehungsdatum, 5)
	auszahlungen = append(auszahlungen, vierrichtigesuper)

	dreirichtige := BerechneAuszahlung(einsatz, gewinneranzahl, 0.45, ziehungsdatum, 2)
	auszahlungen = append(auszahlungen, dreirichtige)

	dreirichtigesuper := BerechneAuszahlung(einsatz, gewinneranzahl, 0.1, ziehungsdatum, 3)
	auszahlungen = append(auszahlungen, dreirichtigesuper)

	return auszahlungen, nil

}

func BerechneAuszahlung(gesamteinsatz float64, gewinnanzahl [10]int, budgetfaktor float64, auszahlungsdatum time.Time, gewinnklasse int8) database.Auszahlung {

	auszahlung := database.Auszahlung{
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

func BerechneJackpotsteigerung(ziehungsTransaktion *sql.Tx, jackpotAuszahlung database.Auszahlung) (database.Auszahlung, error) {

	letzterJackpot, fehler := database.HoleLetztenJackpot(ziehungsTransaktion)

	if letzterJackpot.Auszahlung == 0 {
		jackpotAuszahlung.Budget = jackpotAuszahlung.Budget + letzterJackpot.Budget
	}

	return jackpotAuszahlung, fehler

}

var BerechneGewinneranzahl = func(ziehungsdatum time.Time) ([10]int, error) {

	var gewinneranzahl [10]int
	var ziehung database.Ziehung
	var fehler error
	var tipps []database.Tipp

	for klasse := 0; klasse <= 9; klasse++ {
		gewinneranzahl[klasse] = 0
	}

	ziehung, fehler = database.HoleZiehungZumDatum(ziehungsdatum)

	if fehler != nil {
		return gewinneranzahl, fehler
	}

	tipps, fehler = database.HoleTippsZumDatum(ziehungsdatum)

	if fehler != nil {
		return gewinneranzahl, fehler
	}

	for _, tipp := range tipps {
		gewinnklasse := BerechneGewinnklasse(tipp.Ziehung, ziehung.Ziehung.String)
		if gewinnklasse >= 0 {
			gewinneranzahl[gewinnklasse]++
		}
	}

	return gewinneranzahl, fehler

}

func BerechneGewinnklasse(tipp string, ziehung string) int8 {

	var gewinnklasse int8

	gewinnklasse = -4

	if IstValideZiehung(tipp) && IstValideZiehung(ziehung) {

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

	}

	return gewinnklasse

}

func HashePasswort(passwort string) (string, error) {
	bytes, fehler := bcrypt.GenerateFromPassword([]byte(passwort), 14)
	return string(bytes), fehler
}

func PruefePasswortHash(passwort, hash string) bool {
	fehler := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwort))
	return fehler == nil
}
