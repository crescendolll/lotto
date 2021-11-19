package database

import (
	"database/sql"
	"errors"
	"log"
	"lotto/lottolog"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func NeuerMock() (*sql.DB, sqlmock.Sqlmock) {
	mockVerbindung, mockDatenbank, fehler := sqlmock.New()

	if fehler != nil {
		log.Fatalf("Fehler '%s' beim Erzeugen des Mocks", fehler)
	}

	return mockVerbindung, mockDatenbank
}

func TestFuegeSpielerEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Otto",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,\\?\\)"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeNutzerEin(nutzer)

	assert.NoError(test, fehler)
}

func TestFuegeMitarbeiterEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Hans",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  false,
	}

	abfrage := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,\\?\\)"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeNutzerEin(nutzer)

	assert.NoError(test, fehler)
}

func TestFuegeNutzerDoppeltEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Viktoria",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,\\?\\)"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 0))

	fehler := FuegeNutzerEin(nutzer)

	assert.NoError(test, fehler)

	statement = mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnError(errors.New("Error 1062: Duplicate entry 'Viktoria' for key 'nutzer.PRIMARY'"))

	fehler = FuegeNutzerEin(nutzer)

	assert.Error(test, fehler)
}

func TestHoleNutzerdatenEinesSpielers(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Otto",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "SELECT \\* from nutzer WHERE benutzername = \\?"

	selektion := sqlmock.NewRows([]string{"benutzername", "pw_hash", "ist_spieler"}).AddRow(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler)
	mockDatenbank.ExpectQuery(abfrage).WithArgs(nutzer.Benutzername).WillReturnRows(selektion)

	ergebnis, fehler := HoleNutzerdatenZumNamen(nutzer.Benutzername)

	assert.NoError(test, fehler)
	assert.EqualValues(test, nutzer, ergebnis)

}

func TestHoleNutzerdatenEinesMitarbeiters(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Hans",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  false,
	}

	abfrage := "SELECT \\* from nutzer WHERE benutzername = \\?"

	selektion := sqlmock.NewRows([]string{"benutzername", "pw_hash", "ist_spieler"}).AddRow(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler)
	mockDatenbank.ExpectQuery(abfrage).WithArgs(nutzer.Benutzername).WillReturnRows(selektion)

	ergebnis, fehler := HoleNutzerdatenZumNamen(nutzer.Benutzername)

	assert.NoError(test, fehler)
	assert.EqualValues(test, nutzer, ergebnis)

}

func TestHoleNutzerdatenEinesNichtVorhandenenNutzers(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	abfrage := "SELECT \\* from nutzer WHERE benutzername = \\?"

	selektion := sqlmock.NewRows([]string{"benutzername", "pw_hash", "ist_spieler"})
	mockDatenbank.ExpectQuery(abfrage).WithArgs("Otto").WillReturnRows(selektion)

	ergebnis, fehler := HoleNutzerdatenZumNamen("Otto")

	nutzer := Nutzer{}

	assert.NoError(test, fehler)
	assert.EqualValues(test, nutzer, ergebnis)

}

func TestNameVerfuegbar(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Hans",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  false,
	}

	abfrage := "SELECT \\(COUNT\\(\\*\\) = 0\\) as verfuegbar FROM nutzer WHERE benutzername = \\?"

	selektion := sqlmock.NewRows([]string{"verfuegbar"}).AddRow(true)
	mockDatenbank.ExpectQuery(abfrage).WithArgs(nutzer.Benutzername).WillReturnRows(selektion)

	ergebnis, fehler := HoleVerfuegbarkeitEinesBenutzernamens(nutzer.Benutzername)

	assert.NoError(test, fehler)
	assert.True(test, ergebnis)

}

func TestNameNichtVerfuegbar(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Hans",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  false,
	}

	abfrage := "SELECT \\(COUNT\\(\\*\\) = 0\\) as verfuegbar FROM nutzer WHERE benutzername = \\?"

	selektion := sqlmock.NewRows([]string{"verfuegbar"}).AddRow(false)
	mockDatenbank.ExpectQuery(abfrage).WithArgs(nutzer.Benutzername).WillReturnRows(selektion)

	ergebnis, fehler := HoleVerfuegbarkeitEinesBenutzernamens(nutzer.Benutzername)

	assert.NoError(test, fehler)
	assert.False(test, ergebnis)

}

func TestAendereNutzerdaten(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Otto",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "UPDATE nutzer set pw_hash = \\?, benutzername = \\? WHERE benutzername = \\?"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Pw_hash, nutzer.Benutzername, "Otto").WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := AendereNutzerdaten(nutzer, "Otto")

	assert.NoError(test, fehler)
}

func TestAendereNutzerdatenZuEinemVergebenenNamen(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Otto",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,\\?\\)"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 1))

	_ = FuegeNutzerEin(nutzer)

	nutzer.Benutzername = "Hans"

	statement = mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 1))

	_ = FuegeNutzerEin(nutzer)

	nutzer.Benutzername = "Otto"
	nutzer.Pw_hash = "24banaba43"

	abfrage = "UPDATE nutzer set pw_hash = \\?, benutzername = \\? WHERE benutzername = \\?"

	statement = mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Pw_hash, nutzer.Benutzername, "Hans").WillReturnError(errors.New("Error 1062: Duplicate entry 'Otto' for key 'nutzer.PRIMARY'"))

	fehler := AendereNutzerdaten(nutzer, "Hans")

	assert.Error(test, fehler)
}

func TestLoescheNutzerdaten(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	nutzer := Nutzer{
		Benutzername: "Otto",
		Pw_hash:      "23bonobo42",
		Ist_spieler:  true,
	}

	abfrage := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,\\?\\)"

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername, nutzer.Pw_hash, nutzer.Ist_spieler).WillReturnResult(sqlmock.NewResult(0, 1))

	_ = FuegeNutzerEin(nutzer)

	abfrage = "DELETE from nutzer WHERE benutzername = \\?"

	statement = mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(nutzer.Benutzername).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := LoescheNutzerdaten(nutzer.Benutzername)

	assert.NoError(test, fehler)
}

func TestFuegeTippEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	tipp := Tipp{
		Id:      1,
		Datum:   time.Now(),
		Ziehung: "0102030405067",
	}

	abfrage := "INSERT INTO tipps \\(id, datum, ziehung\\) VALUES \\(\\?,\\?,\\?\\)"

	mockDatenbank.ExpectBegin()

	mockTransaktion, _ := mockVerbindung.Begin()

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(tipp.Id, tipp.Datum, tipp.Ziehung).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeTippEin(tipp, mockTransaktion)

	assert.NoError(test, fehler)
}

func TestHoleTippdatenEinesSpielers(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	var tipps []Tipp

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	tipp := Tipp{
		Id:      1,
		Datum:   time.Now(),
		Ziehung: "0102030405067",
	}

	tipps = append(tipps, tipp)
	selektion := sqlmock.NewRows([]string{"tipps.id", "tipps.datum", "tipps.ziehung"}).AddRow(tipp.Id, tipp.Datum, tipp.Ziehung)

	tipp.Id = 2
	tipp.Ziehung = "1533324509271"

	tipps = append(tipps, tipp)
	selektion.AddRow(tipp.Id, tipp.Datum, tipp.Ziehung)

	abfrage := "SELECT tipps.id, tipps.datum, tipps.ziehung FROM spieler_tipps, tipps " +
		"WHERE spieler_tipps.spielername = \\? AND spieler_tipps.id = tipps.id AND tipps.datum >= \\? AND tipps.datum <= \\?"

	mockDatenbank.ExpectQuery(abfrage).WithArgs("Luisa", time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02")).WillReturnRows(selektion)

	ergebnis, fehler := HoleTippsVonSpielerImZeitraum("Luisa", time.Now(), time.Now())

	assert.NoError(test, fehler)
	assert.EqualValues(test, tipps, ergebnis)

}

func TestHoleTippsZumDatum(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	var tipps []Tipp

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	tipp := Tipp{
		Ziehung: "0102030405067",
	}

	tipps = append(tipps, tipp)
	selektion := sqlmock.NewRows([]string{"ziehung"}).AddRow(tipp.Ziehung)

	tipp.Ziehung = "1533324509271"

	tipps = append(tipps, tipp)
	selektion.AddRow(tipp.Ziehung)

	abfrage := "SELECT ziehung FROM tipps WHERE datum = \\'" + time.Now().Format("2006-01-02") + "\\'"

	mockDatenbank.ExpectQuery(abfrage).WithArgs().WillReturnRows(selektion)

	ergebnis, fehler := HoleTippsZumDatum(time.Now())

	assert.NoError(test, fehler)
	assert.EqualValues(test, tipps, ergebnis)

}

func TestFuegeSpielerTippVerknuepfungEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	abfrage := "INSERT INTO spieler_tipps \\(id, spielername\\) VALUES \\(\\?,\\?\\)"

	mockDatenbank.ExpectBegin()

	mockTransaktion, _ := mockVerbindung.Begin()

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(1, "Horst").WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeSpielerTippVerknuepfungEin(1, "Horst", mockTransaktion)

	assert.NoError(test, fehler)
}

func TestFuegeZiehungEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum: time.Now(),
	}

	abfrage := "INSERT INTO ziehungen \\(datum\\) VALUES \\(\\?\\)"

	mockDatenbank.ExpectBegin()
	mockTransaktion, _ := mockVerbindung.Begin()

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(ziehung.Datum.Format("2006-01-02")).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeZiehungEin(ziehung, mockTransaktion)

	assert.NoError(test, fehler)
}

func TestHoleOffeneZiehungen(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	var ziehungen []Ziehung

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum: time.Now(),
	}

	ziehungen = append(ziehungen, ziehung)
	selektion := sqlmock.NewRows([]string{"datum", "ziehung"}).AddRow(ziehung.Datum, ziehung.Ziehung)

	ziehung.Datum = time.Now().AddDate(0, 0, 1)

	ziehungen = append(ziehungen, ziehung)
	selektion.AddRow(ziehung.Datum, ziehung.Ziehung)

	abfrage := "SELECT \\* from ziehungen WHERE ziehung IS NULL"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleOffeneZiehungen()

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehungen, ergebnis)

}

func TestHoleLetztesZiehungsdatum(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum: time.Now(),
	}

	selektion := sqlmock.NewRows([]string{"letzteZiehung"}).AddRow(ziehung.Datum)

	abfrage := "SELECT MAX\\(datum\\) as letzteZiehung FROM ziehungen WHERE ziehung IS NOT NULL"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleLetztesZiehungsdatum()

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehung.Datum, ergebnis)

}

func TestZiehungVerfuegbar(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum: time.Now(),
	}

	selektion := sqlmock.NewRows([]string{"verfuegbar"}).AddRow(true)

	abfrage := "SELECT \\(COUNT\\(\\*\\) > 0\\) as verfuegbar FROM ziehungen WHERE datum = \\'" + ziehung.Datum.Format("2006-01-02") + "\\' AND ziehung IS NULL"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleVerfuegbarkeitEinerZiehungZumDatum(ziehung.Datum)

	assert.NoError(test, fehler)
	assert.True(test, ergebnis)

}

func TestZiehungNichtVerfuegbar(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum: time.Now(),
	}

	selektion := sqlmock.NewRows([]string{"verfuegbar"}).AddRow(false)

	abfrage := "SELECT \\(COUNT\\(\\*\\) > 0\\) as verfuegbar FROM ziehungen WHERE datum = \\'" + ziehung.Datum.Format("2006-01-02") + "\\' AND ziehung IS NULL"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleVerfuegbarkeitEinerZiehungZumDatum(ziehung.Datum)

	assert.NoError(test, fehler)
	assert.False(test, ergebnis)

}

func TestHoleZiehungZumDatum(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0114253741292"),
	}

	selektion := sqlmock.NewRows([]string{"datum", "ziehung"}).AddRow(ziehung.Datum, ziehung.Ziehung)

	abfrage := "SELECT \\* from ziehungen WHERE datum = \\'" + ziehung.Datum.Format("2006-01-02") + "\\'"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleZiehungZumDatum(ziehung.Datum)

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehung, ergebnis)

}

func TestHoleZiehungenInnerhalbEinesZeitraums(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	var ziehungen []Ziehung

	ziehung := Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0114253741292"),
	}

	ziehungen = append(ziehungen, ziehung)
	selektion := sqlmock.NewRows([]string{"datum", "ziehung"}).AddRow(ziehung.Datum, ziehung.Ziehung)

	ziehung.Datum = time.Now().AddDate(0, 0, 1)
	ziehung.Ziehung = null.StringFrom("1324354631428")

	ziehungen = append(ziehungen, ziehung)
	selektion.AddRow(ziehung.Datum, ziehung.Ziehung)

	abfrage := "SELECT \\* from ziehungen WHERE datum >= \\'" + time.Now().Format("2006-01-02") + "\\' AND datum <= \\'" + time.Now().AddDate(0, 0, 1).Format("2006-01-02") + "\\'"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleZiehungenInnerhalbEinesZeitraums(time.Now(), time.Now().AddDate(0, 0, 1))

	assert.NoError(test, fehler)
	assert.EqualValues(test, ziehungen, ergebnis)

}

func TestAendereZiehung(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	ziehung := Ziehung{
		Datum:   time.Now(),
		Ziehung: null.StringFrom("0114253741292"),
	}

	abfrage := "UPDATE ziehungen SET ziehung = \\? WHERE datum = \\?"

	mockDatenbank.ExpectBegin()
	mockTransaktion, _ := mockVerbindung.Begin()

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(ziehung.Ziehung, ziehung.Datum.Format("2006-01-02")).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := AendereZiehung(ziehung, mockTransaktion)

	assert.NoError(test, fehler)
}

func TestFuegeMitarbeiterZiehungVerknuepfungEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	abfrage := "INSERT INTO mitarbeiter_ziehungen \\(datum, mitarbeitername, aktion\\) VALUES \\(\\?,\\?,\\?\\)"

	mockDatenbank.ExpectBegin()

	mockTransaktion, _ := mockVerbindung.Begin()

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(time.Now().Format("2006-01-02"), "Horst", false).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeMitarbeiterZiehungVerknuepfungEin(time.Now(), "Horst", false, mockTransaktion)

	assert.NoError(test, fehler)
}

func TestFuegeAuszahlungEin(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	abfrage := "INSERT INTO auszahlungen \\(datum, klasse, budget, auszahlung\\) VALUES \\(\\?,\\?,\\?,\\?\\)"

	mockDatenbank.ExpectBegin()

	mockTransaktion, _ := mockVerbindung.Begin()

	auszahlung := Auszahlung{
		Datum:      time.Now(),
		Klasse:     9,
		Budget:     1000.00,
		Auszahlung: 32.50,
	}

	statement := mockDatenbank.ExpectPrepare(abfrage)
	statement.ExpectExec().WithArgs(auszahlung.Datum.Format("2006-01-02"), auszahlung.Klasse, auszahlung.Budget, auszahlung.Auszahlung).WillReturnResult(sqlmock.NewResult(0, 1))

	fehler := FuegeAuszahlungEin(mockTransaktion, auszahlung)

	assert.NoError(test, fehler)
}

func TestHoleAuszahlungenZumDatum(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	var auszahlungen []Auszahlung

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	auszahlung := Auszahlung{
		Klasse:     0,
		Auszahlung: 32.50,
	}

	auszahlungen = append(auszahlungen, auszahlung)
	selektion := sqlmock.NewRows([]string{"klasse", "auszahlung"}).AddRow(auszahlung.Klasse, auszahlung.Auszahlung)

	for testindex := 1; testindex < 10; testindex++ {
		auszahlung.Klasse = int8(testindex)
		auszahlung.Auszahlung *= 2

		auszahlungen = append(auszahlungen, auszahlung)
		selektion.AddRow(auszahlung.Klasse, auszahlung.Auszahlung)
	}

	abfrage := "SELECT klasse, auszahlung from auszahlungen WHERE datum = \\'" + time.Now().Format("2006-01-02") + "\\'"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleAuszahlungenZumDatum(time.Now())

	assert.NoError(test, fehler)
	assert.EqualValues(test, auszahlungen, ergebnis)

}

func TestHoleAuszahlungZuDatumUndKlasse(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	auszahlung := Auszahlung{
		Auszahlung: 32.50,
	}

	selektion := sqlmock.NewRows([]string{"auszahlung"}).AddRow(auszahlung.Auszahlung)

	abfrage := "SELECT auszahlung from auszahlungen WHERE datum = \\'" + time.Now().Format("2006-01-02") + "\\' AND klasse = 9"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleAuszahlungZuDatumUndKlasse(time.Now(), 9)

	assert.NoError(test, fehler)
	assert.EqualValues(test, auszahlung, ergebnis)

}

func TestHoleGroessteTippID(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	selektion := sqlmock.NewRows([]string{"maxID"}).AddRow(0)

	abfrage := "SELECT MAX\\(id\\) as maxID FROM tipps"

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleGroessteTippID()

	assert.NoError(test, fehler)
	assert.EqualValues(test, 0, ergebnis)

}

func TestHoleLetztenJackpot(test *testing.T) {

	lottolog.OeffneTestLogdatei(Lottokonfig.Testlogdateipfad)

	mockVerbindung, mockDatenbank := NeuerMock()

	SetzeVerbindung(mockVerbindung)

	auszahlung := Auszahlung{
		Budget:     100.00,
		Auszahlung: 100.00,
	}

	selektion := sqlmock.NewRows([]string{"budget", "auszahlung"}).AddRow(auszahlung.Budget, auszahlung.Auszahlung)

	abfrage := "SELECT budget, auszahlung FROM auszahlungen WHERE klasse = 9 AND datum = \\(SELECT MAX\\(datum\\) FROM auszahlungen\\)"

	mockDatenbank.ExpectBegin()
	mockTransaktion, _ := mockVerbindung.Begin()

	mockDatenbank.ExpectQuery(abfrage).WillReturnRows(selektion)

	ergebnis, fehler := HoleLetztenJackpot(mockTransaktion)

	assert.NoError(test, fehler)
	assert.EqualValues(test, auszahlung, ergebnis)

}
