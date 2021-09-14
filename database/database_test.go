package database

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	databasehandle, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("Fehler '%s' beim Erzeugen des Mocks", err)
	}

	return databasehandle, mock
}

func TestInsertIntoSpieler(test *testing.T) {
	databasehandle, mock := NewMock()

	benutzername := "Ingo"
	pw_hash := "23bonobo42"

	selectQuery := "SELECT COUNT\\(\\*\\) as anzahl FROM nutzer WHERE benutzername = \\?"

	rows := sqlmock.NewRows([]string{"anzahl"}).AddRow(0)

	mock.ExpectQuery(selectQuery).WithArgs(benutzername).WillReturnRows(rows)

	insertQuery := "INSERT INTO nutzer \\(benutzername, pw_hash, ist_spieler\\) VALUES \\(\\?,\\?,1\\)"

	prep := mock.ExpectPrepare(insertQuery)
	prep.ExpectExec().WithArgs(benutzername, pw_hash).WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertIntoSpieler(databasehandle, benutzername, pw_hash)

	assert.NoError(test, err)
}
