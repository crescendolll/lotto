package lottoapi

import (
	"fmt"
	"lotto/database"
	"lotto/lottojson"
	"lotto/lottologic"
	"time"

	"github.com/guregu/null"
	"github.com/segmentio/ksuid"
)

var aktiveNutzer map[string]lottologic.Nutzer

// Ausgabe ist interface{} - beliebig, weil jede Response erzeugt werden kann
func HandleRequest(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	switch apirequest.Methode {
	case "login":
		response = CreateLoginResponse(apirequest)
	case "logout":
		response = CreateLogoutResponse(apirequest)
	case "registriere":
		response = CreateRegistrationResponse(apirequest)
	case "aendereKontodaten":
		response = CreateUpdateResponse(apirequest)
	case "loescheKontodaten":
		response = CreateDeleteResponse(apirequest)
	case "neuerTipp":
		response = CreateTippResponse(apirequest)
	case "neueZiehung":
		response = CreateNewGameResponse(apirequest)
	case "beendeZiehung":
		response = CreateCloseGameResponse(apirequest)
	case "zeigeAktuelleSpiele":
		response = CreateCurrentGameResponse(apirequest)
	case "zeigeTipps":
		response = CreateCurrentTippsResponse(apirequest)
	case "holeZiehungen":
		response = CreateGetZiehungMitAuszahlungResponse(apirequest)
	default:
		response = lottojson.ErrorResponse{
			Errormessage: "Unbekannte Methode",
		}
	}

	return response

}

func CreateGetZiehungMitAuszahlungResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var von time.Time
	var bis time.Time
	var err error
	var ziehungsauszahlungen []lottologic.Ziehungauszahlung

	if apirequest.Param["bis"] == "" {
		bis = time.Now()
	} else {
		bis, err = time.Parse("2006-01-02", apirequest.Param["bis"])
	}

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
		return response
	}

	if apirequest.Param["von"] != "" {
		von, err = time.Parse("2006-01-02", apirequest.Param["von"])
	}

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
		return response
	}

	databasehandle := database.OpenLottoConnection()
	ziehungsauszahlungen, err = database.HoleZiehungenMitAuszahlungen(databasehandle, von, bis)
	database.CloseLottoConnection(databasehandle)

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
	} else {
		response = lottojson.GetZiehungenMitAuszahlungenResponse{
			Errormessage:         "",
			Ziehungsauszahlungen: ziehungsauszahlungen,
		}
	}

	return response
}

func CreateCurrentTippsResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var von time.Time
	var bis time.Time
	var nutzer lottologic.Nutzer
	var err error
	var tippauszahlungen []lottologic.Tippauszahlung

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if apirequest.Param["bis"] == "" {
		bis = time.Now()
	} else {
		bis, err = time.Parse("2006-01-02", apirequest.Param["bis"])
	}

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
		return response
	}

	if apirequest.Param["von"] != "" {
		von, err = time.Parse("2006-01-02", apirequest.Param["von"])
	}

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
		return response
	}

	nutzer = aktiveNutzer[apirequest.Auth]

	databasehandle := database.OpenLottoConnection()
	tippauszahlungen, err = database.HoleTippauszahlungenFuerSpieler(databasehandle, nutzer.Benutzername, von, bis)
	database.CloseLottoConnection(databasehandle)

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
	} else {
		response = lottojson.GetTippsResponse{
			Errormessage:     "",
			Tippauszahlungen: tippauszahlungen,
		}
	}

	for _, tippausz := range tippauszahlungen {
		fmt.Printf("| %d | %s | %s | %d | %.2f\n", tippausz.Id, tippausz.Datum.Format("2006-01-02"), tippausz.Ziehung, tippausz.Klasse, tippausz.Auszahlung)
	}

	return response
}

func CreateCloseGameResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var mitarbeiter lottologic.Nutzer
	var ziehung lottologic.Ziehung

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if apirequest.Param["datum"] == "" || apirequest.Param["ziehung"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum und Ziehung angeben",
		}
		return response
	} else {

		var timeErr error

		ziehung.Datum, timeErr = time.Parse("2006-01-02", apirequest.Param["datum"])
		ziehung.Ziehung = null.StringFrom(apirequest.Param["ziehung"])

		if timeErr != nil {
			response = lottojson.ErrorResponse{
				Errormessage: timeErr.Error(),
			}
			return response
		}

		mitarbeiter = aktiveNutzer[apirequest.Auth]

		databasehandle := database.OpenLottoConnection()
		updateError := database.UpdateZiehungen(databasehandle, ziehung, mitarbeiter)
		database.CloseLottoConnection(databasehandle)

		if updateError != nil {
			response = lottojson.ErrorResponse{
				Errormessage: updateError.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}
	}

	return response

}

func CreateCurrentGameResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var ziehungen []lottologic.Ziehung
	var err error

	databasehandle := database.OpenLottoConnection()
	ziehungen, err = database.SelectLaufendeZiehungen(databasehandle)

	var ziehungstage []string
	ziehungstage = make([]string, 0)

	for _, ziehung := range ziehungen {
		ziehungstage = append(ziehungstage, ziehung.Datum.Format("2006-01-02"))
	}

	if err != nil {
		response = lottojson.ErrorResponse{
			Errormessage: err.Error(),
		}
	} else {
		response = lottojson.GetZiehungenResponse{
			Errormessage: "",
			Ziehungen:    ziehungstage,
		}
	}

	database.CloseLottoConnection(databasehandle)

	return response

}

func CreateNewGameResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var mitarbeiter lottologic.Nutzer
	var neueZiehung lottologic.Ziehung

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if apirequest.Param["datum"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum der Ziehung angeben",
		}
		return response
	} else {

		var timeErr error

		neueZiehung.Datum, timeErr = time.Parse("2006-01-02", apirequest.Param["datum"])

		if timeErr != nil {
			response = lottojson.ErrorResponse{
				Errormessage: timeErr.Error(),
			}
			return response
		}

		mitarbeiter = aktiveNutzer[apirequest.Auth]

		databasehandle := database.OpenLottoConnection()
		insertError := database.InsertIntoZiehungen(databasehandle, neueZiehung, mitarbeiter)
		database.CloseLottoConnection(databasehandle)

		if insertError != nil {
			response = lottojson.ErrorResponse{
				Errormessage: insertError.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}
	}

	return response

}

func CreateTippResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var neuerTipp lottologic.Tipp
	var nutzer lottologic.Nutzer

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if apirequest.Param["tipp"] == "" || apirequest.Param["datum"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum und Tipp angeben",
		}
		return response
	} else {

		neuerTipp = lottologic.Tipp{
			Ziehung: apirequest.Param["tipp"],
		}

		var timeErr error

		neuerTipp.Datum, timeErr = time.Parse("2006-01-02", apirequest.Param["datum"])

		if timeErr != nil {
			response = lottojson.ErrorResponse{
				Errormessage: timeErr.Error(),
			}
			return response
		}

		nutzer = aktiveNutzer[apirequest.Auth]

		databasehandle := database.OpenLottoConnection()
		insertError := database.InsertIntoTipps(databasehandle, neuerTipp, nutzer)
		database.CloseLottoConnection(databasehandle)

		if insertError != nil {
			response = lottojson.ErrorResponse{
				Errormessage: insertError.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}
	}

	return response

}

func CreateDeleteResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var deleteError error

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {

		nutzer := aktiveNutzer[apirequest.Auth]
		databasehandle := database.OpenLottoConnection()
		deleteError = database.DeleteFromSpieler(databasehandle, nutzer.Benutzername)
		database.CloseLottoConnection(databasehandle)
		if deleteError != nil {
			response = lottojson.ErrorResponse{
				Errormessage: deleteError.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}

	}

	return response

}

func CreateUpdateResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var updateError error

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {

		if apirequest.Param["neuespasswort"] == "" {
			response = lottojson.ErrorResponse{
				Errormessage: "kein neues Passwort uebergeben",
			}
		} else {
			nutzer := aktiveNutzer[apirequest.Auth]
			databasehandle := database.OpenLottoConnection()
			updateError = database.UpdateSpieler(databasehandle, nutzer.Benutzername, nutzer.Benutzername, apirequest.Param["neuespasswort"])
			database.CloseLottoConnection(databasehandle)
			if updateError != nil {
				response = lottojson.ErrorResponse{
					Errormessage: updateError.Error(),
				}
			} else {
				response = lottojson.ErrorResponse{
					Errormessage: "",
				}
			}
		}

	}

	return response

}

func CreateRegistrationResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var neuerNutzer lottologic.Nutzer

	if apirequest.Param["name"] == "" || apirequest.Param["passwort"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		neuerNutzer = lottologic.Nutzer{
			Benutzername: apirequest.Param["name"],
			Pw_hash:      apirequest.Param["passwort"],
			Ist_spieler:  true,
		}

		databasehandle := database.OpenLottoConnection()
		insertError := database.InsertSpielerIntoNutzer(databasehandle, neuerNutzer.Benutzername, neuerNutzer.Pw_hash)
		database.CloseLottoConnection(databasehandle)

		if insertError != nil {
			response = lottojson.ErrorResponse{
				Errormessage: insertError.Error(),
			}
		} else {

			auth := LoginNutzer(neuerNutzer)
			response = lottojson.RegistrationResponse{
				Errormessage: "",
				Auth:         auth,
			}
		}
	}

	fmt.Println("Aktuelle Nutzer:")
	for key, value := range aktiveNutzer {
		fmt.Printf("Nutzer %s mit PW %s unter Token %s aktiv\n", value.Benutzername, value.Pw_hash, key)
	}

	return response

}

func CreateLoginResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	if apirequest.Param["name"] == "" || apirequest.Param["passwort"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		databasehandle := database.OpenLottoConnection()
		nutzer, err := database.SelectFromSpielerByName(databasehandle, apirequest.Param["name"])
		database.CloseLottoConnection(databasehandle)

		if err != nil {
			response = lottojson.ErrorResponse{
				Errormessage: err.Error(),
			}
		} else {

			if !lottologic.CheckPasswordHash(apirequest.Param["passwort"], nutzer.Pw_hash) {
				response = lottojson.ErrorResponse{
					Errormessage: "Passwort oder Benutzername sind nicht korrekt",
				}
			} else {
				auth := LoginNutzer(nutzer)
				response = lottojson.LoginResponse{
					Errormessage: "",
					IstSpieler:   nutzer.Ist_spieler,
					Auth:         auth,
				}
			}
		}
	}

	fmt.Println("Aktuelle Nutzer:")
	for key, value := range aktiveNutzer {
		fmt.Printf("Nutzer %s mit PW %s unter Token %s aktiv\n", value.Benutzername, value.Pw_hash, key)
	}

	return response

}

func ValidateAuth(auth string) bool {

	_, found := aktiveNutzer[auth]

	return found

}

func CreateLogoutResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	if !ValidateAuth(apirequest.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {
		LogoutNutzer(apirequest.Auth)
		response = lottojson.ErrorResponse{
			Errormessage: "",
		}
	}

	for key, value := range aktiveNutzer {
		fmt.Printf("Nutzer %s mit PW %s unter Token %s aktiv\n", value.Benutzername, value.Pw_hash, key)
	}

	return response

}

func InitNutzer() {

	aktiveNutzer = make(map[string]lottologic.Nutzer)

}

func LoginNutzer(nutzer lottologic.Nutzer) string {

	auth := ksuid.New().String()

	aktiveNutzer[auth] = nutzer

	return auth
}

func LogoutNutzer(auth string) {

	delete(aktiveNutzer, auth)

}
