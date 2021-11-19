package lottologic

import (
	"database/sql"
	"lotto/database"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func TestIstValideZiehung(test *testing.T) {

	ziehung := "0102030405067"

	ergebnis := IstValideZiehung(ziehung)

	assert.True(test, ergebnis)

}

func TestIstKeineValideZiehung(test *testing.T) {

	ziehung := "kigcrui gegbrg"

	ergebnis := IstValideZiehung(ziehung)

	assert.False(test, ergebnis)

}

func TestKeinGewinn(test *testing.T) {

	tipp := "1324354642319"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, -4, ergebnis)
}

func TestZweiRichtige(test *testing.T) {

	tipp := "1122354642319"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 0, ergebnis)
}

func TestZweiRichtigeMitSuperzahl(test *testing.T) {

	tipp := "1122354642312"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 1, ergebnis)
}

func TestDreiRichtige(test *testing.T) {

	tipp := "1122334642313"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 2, ergebnis)
}

func TestDreiRichtigeMitSuperzahl(test *testing.T) {

	tipp := "1122334642312"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 3, ergebnis)
}

func TestVierRichtige(test *testing.T) {

	tipp := "1122334442313"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 4, ergebnis)
}

func TestVierRichtigeMiitSuperzahl(test *testing.T) {

	tipp := "1122334442313"

	ziehung := "1122334410193"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 5, ergebnis)
}

func TestFuenfRichtige(test *testing.T) {

	tipp := "1122334410313"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 6, ergebnis)
}

func TestFuenfRichtigeMitSuperzahl(test *testing.T) {

	tipp := "1122334410312"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 7, ergebnis)
}

func TestSechsRichtige(test *testing.T) {

	tipp := "1122334410193"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 8, ergebnis)
}

func TestSechsRichtigeMitSuperzahll(test *testing.T) {

	tipp := "1122334410192"

	ziehung := "1122334410192"

	ergebnis := BerechneGewinnklasse(tipp, ziehung)

	assert.EqualValues(test, 9, ergebnis)
}

func TestBerechneJackpotsteigerung(test *testing.T) {

	jackpotAuszahlung := database.Auszahlung{

		Budget: 100.00,
	}

	originalHoleLetztenJackpot := database.HoleLetztenJackpot

	database.HoleLetztenJackpot = func(transaktion *sql.Tx) (database.Auszahlung, error) {
		letzterJackpot := database.Auszahlung{

			Budget:     100.00,
			Auszahlung: 0,
		}

		return letzterJackpot, nil
	}

	jackpotAuszahlung, fehler := BerechneJackpotsteigerung(nil, jackpotAuszahlung)

	database.HoleLetztenJackpot = originalHoleLetztenJackpot

	assert.NoError(test, fehler)
	assert.EqualValues(test, 200.00, jackpotAuszahlung.Budget)

}

func TestErstelleAuszahlungsstatistiken(test *testing.T) {

	ziehung := database.Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0102030405067"),
	}

	originalHoleAuszahlungenZumDatum := database.HoleAuszahlungenZumDatum

	database.HoleAuszahlungenZumDatum = func(datum time.Time) ([]database.Auszahlung, error) {

		var auszahlungen []database.Auszahlung

		auszahlungen = make([]database.Auszahlung, 0)

		for klasse := 0; klasse < 10; klasse++ {
			auszahlung := database.Auszahlung{
				Datum:      time.Now(),
				Klasse:     int8(klasse),
				Budget:     100.00,
				Auszahlung: 100.00,
			}
			auszahlungen = append(auszahlungen, auszahlung)
		}

		return auszahlungen, nil

	}

	originalBerechneGewinneranzahl := BerechneGewinneranzahl

	BerechneGewinneranzahl = func(ziehungsdatum time.Time) ([10]int, error) {

		gewinner := [10]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

		return gewinner, nil
	}

	var auszahlungsstatistiken []Auszahlungsstatistik

	auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

	for klasse := 0; klasse < 10; klasse++ {
		auszahlungsstatistik := Auszahlungsstatistik{
			Klasse:   int8(klasse),
			Gewinner: 1,
			Gewinn:   100.00,
		}
		auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
	}

	ergebnis, fehler := ErstelleAuszahlungsstatistiken(ziehung)

	database.HoleAuszahlungenZumDatum = originalHoleAuszahlungenZumDatum
	BerechneGewinneranzahl = originalBerechneGewinneranzahl

	assert.NoError(test, fehler)
	assert.EqualValues(test, ergebnis, auszahlungsstatistiken)
}

func TestErstelleZiehungsstatistiken(test *testing.T) {

	originalErstelleAuszahlungsstatistiken := ErstelleAuszahlungsstatistiken

	ErstelleAuszahlungsstatistiken = func(ziehung database.Ziehung) ([]Auszahlungsstatistik, error) {
		var auszahlungsstatistiken []Auszahlungsstatistik

		auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

		for klasse := 0; klasse < 10; klasse++ {
			auszahlungsstatistik := Auszahlungsstatistik{
				Klasse:   int8(klasse),
				Gewinner: 1,
				Gewinn:   100.00,
			}
			auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
		}

		return auszahlungsstatistiken, nil
	}

	var auszahlungsstatistiken []Auszahlungsstatistik

	auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

	for klasse := 0; klasse < 10; klasse++ {
		auszahlungsstatistik := Auszahlungsstatistik{
			Klasse:   int8(klasse),
			Gewinner: 1,
			Gewinn:   100.00,
		}
		auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
	}

	var ziehungsstatistiken []Ziehungsstatistik

	var ziehungen []database.Ziehung

	ziehungen = make([]database.Ziehung, 0)

	ziehung := database.Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0102030405067"),
	}

	ziehungen = append(ziehungen, ziehung)
	ziehungen = append(ziehungen, ziehung)

	ziehungsstatistik := Ziehungsstatistik{
		Datum:        ziehung.Datum.Format("2006-01-02"),
		Ziehung:      ziehung.Ziehung,
		Auszahlungen: auszahlungsstatistiken,
	}

	ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)
	ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)

	ergebnis, fehler := ErstelleZiehungsstatistiken(ziehungen)

	ErstelleAuszahlungsstatistiken = originalErstelleAuszahlungsstatistiken

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehungsstatistiken, ergebnis)
}

func TestErstelleZiehungsstatistikenFuerZeitraum(test *testing.T) {

	originalHoleZiehungenInnerhalbEinesZeitraums := database.HoleZiehungenInnerhalbEinesZeitraums
	originalErstelleZiehungsstatistiken := ErstelleZiehungsstatistiken

	database.HoleZiehungenInnerhalbEinesZeitraums = func(startdatum, enddatum time.Time) ([]database.Ziehung, error) {
		var ziehungen []database.Ziehung

		ziehungen = make([]database.Ziehung, 0)

		ziehung := database.Ziehung{
			Datum:   time.Now(),
			Ziehung: null.StringFrom("0102030405067"),
		}

		ziehungen = append(ziehungen, ziehung)
		ziehungen = append(ziehungen, ziehung)

		return ziehungen, nil
	}

	ErstelleZiehungsstatistiken = func(ziehungen []database.Ziehung) ([]Ziehungsstatistik, error) {
		var auszahlungsstatistiken []Auszahlungsstatistik

		auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

		for klasse := 0; klasse < 10; klasse++ {
			auszahlungsstatistik := Auszahlungsstatistik{
				Klasse:   int8(klasse),
				Gewinner: 1,
				Gewinn:   100.00,
			}
			auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
		}

		var ziehungsstatistiken []Ziehungsstatistik

		var testziehungen []database.Ziehung

		testziehungen = make([]database.Ziehung, 0)

		ziehung := database.Ziehung{
			Datum:   time.Now(),
			Ziehung: null.StringFrom("0102030405067"),
		}

		testziehungen = append(testziehungen, ziehung)
		testziehungen = append(testziehungen, ziehung)

		ziehungsstatistik := Ziehungsstatistik{
			Datum:        ziehung.Datum.Format("2006-01-02"),
			Ziehung:      ziehung.Ziehung,
			Auszahlungen: auszahlungsstatistiken,
		}

		ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)
		ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)

		return ziehungsstatistiken, nil
	}

	var auszahlungsstatistiken []Auszahlungsstatistik

	auszahlungsstatistiken = make([]Auszahlungsstatistik, 0)

	for klasse := 0; klasse < 10; klasse++ {
		auszahlungsstatistik := Auszahlungsstatistik{
			Klasse:   int8(klasse),
			Gewinner: 1,
			Gewinn:   100.00,
		}
		auszahlungsstatistiken = append(auszahlungsstatistiken, auszahlungsstatistik)
	}

	var ziehungsstatistiken []Ziehungsstatistik

	var testziehungen []database.Ziehung

	testziehungen = make([]database.Ziehung, 0)

	ziehung := database.Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0102030405067"),
	}

	testziehungen = append(testziehungen, ziehung)
	testziehungen = append(testziehungen, ziehung)

	ziehungsstatistik := Ziehungsstatistik{
		Datum:        ziehung.Datum.Format("2006-01-02"),
		Ziehung:      ziehung.Ziehung,
		Auszahlungen: auszahlungsstatistiken,
	}

	ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)
	ziehungsstatistiken = append(ziehungsstatistiken, ziehungsstatistik)

	ergebnis, fehler := ErstelleZiehungsstatistikenFuerZeitraum(time.Now(), time.Now())

	database.HoleZiehungenInnerhalbEinesZeitraums = originalHoleZiehungenInnerhalbEinesZeitraums
	ErstelleZiehungsstatistiken = originalErstelleZiehungsstatistiken

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehungsstatistiken, ergebnis)
}
