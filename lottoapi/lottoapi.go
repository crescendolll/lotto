package lottoapi

import (
	"lotto/database"
	"lotto/lottojson"
	"lotto/lottolog"
	"lotto/lottologic"
	"time"

	"github.com/guregu/null"
	"github.com/segmentio/ksuid"
)

var aktiveNutzer map[string]database.Nutzer

func BearbeiteRequest(request lottojson.LottoRequest) interface{} {

	var response interface{}

	switch request.Methode {
	case "login":
		response = ErstelleResponseAufLogin(request)
	case "logout":
		response = ErstelleResponseAufLogout(request)
	case "registriere":
		response = ErstelleResponseAufRegistrierung(request)
	case "aendereKontodaten":
		response = ErstelleResponseAufKontoaenderung(request)
	case "loescheKontodaten":
		response = ErstelleResponseAufKontoloeschung(request)
	case "neuerTipp":
		response = ErstelleResponseAufTippabgabe(request)
	case "neueZiehung":
		response = ErstelleResponseAufZiehungseroeffnung(request)
	case "beendeZiehung":
		response = ErstelleResponseAufZiehungsschliessung(request)
	case "zeigeAktuelleSpiele":
		response = ErstelleResponseAufAnfrageNachLaufendenZiehungen(request)
	case "zeigeTipps":
		response = ErstelleResponseAufAnfrageNachAbgegebenenTipps(request)
	case "holeZiehungen":
		response = ErstelleResponseAufAnfrageNachGespieltenZiehungen(request)
	default:
		response = lottojson.ErrorResponse{
			Errormessage: "Unbekannte Methode",
		}
	}

	return response

}

func ErstelleResponseAufAnfrageNachGespieltenZiehungen(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var startdatum time.Time
	var enddatum time.Time
	var fehler error
	var ziehungsauszahlungen []lottologic.Ziehungsstatistik

	if request.Param["bis"] == "" {
		enddatum = time.Now()
	} else {
		enddatum, fehler = time.Parse("2006-01-02", request.Param["bis"])
	}

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
		return response
	}

	if request.Param["von"] != "" {
		startdatum, fehler = time.Parse("2006-01-02", request.Param["von"])
	}

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
		return response
	}

	ziehungsauszahlungen, fehler = lottologic.ErstelleZiehungsstatistikenFuerZeitraum(startdatum, enddatum)

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
	} else {
		response = lottojson.GetZiehungenMitAuszahlungenResponse{
			Errormessage:         "",
			Ziehungsauszahlungen: ziehungsauszahlungen,
		}
	}

	return response
}

func ErstelleResponseAufAnfrageNachAbgegebenenTipps(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var startdatum time.Time
	var enddatum time.Time
	var nutzer database.Nutzer
	var fehler error
	var tippauszahlungen []lottologic.Tippauszahlung

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if request.Param["bis"] == "" {
		enddatum = time.Now()
	} else {
		enddatum, fehler = time.Parse("2006-01-02", request.Param["bis"])
	}

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
		return response
	}

	if request.Param["von"] != "" {
		startdatum, fehler = time.Parse("2006-01-02", request.Param["von"])
	}

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
		return response
	}

	nutzer = aktiveNutzer[request.Auth]

	tippauszahlungen, fehler = lottologic.ErstelleTippauszahlungenFuerSpieler(nutzer.Benutzername, startdatum, enddatum)

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
	} else {
		response = lottojson.GetTippsResponse{
			Errormessage:     "",
			Tippauszahlungen: tippauszahlungen,
		}
	}

	return response
}

func ErstelleResponseAufZiehungsschliessung(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var mitarbeiter database.Nutzer
	var ziehung database.Ziehung

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if request.Param["datum"] == "" || request.Param["ziehung"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum und Ziehung angeben",
		}
		return response
	} else {

		var parserFehler error

		ziehung.Datum, parserFehler = time.Parse("2006-01-02", request.Param["datum"])
		ziehung.Ziehung = null.StringFrom(request.Param["ziehung"])

		if parserFehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: parserFehler.Error(),
			}
			return response
		}

		mitarbeiter = aktiveNutzer[request.Auth]

		updateFehler := lottologic.SchliesseZiehung(ziehung, mitarbeiter)

		if updateFehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: updateFehler.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}
	}

	return response

}

func ErstelleResponseAufAnfrageNachLaufendenZiehungen(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var ziehungen []database.Ziehung
	var fehler error
	var ziehungstage []string

	ziehungen, fehler = database.HoleOffeneZiehungen()

	if fehler != nil {
		response = lottojson.ErrorResponse{
			Errormessage: fehler.Error(),
		}
	} else {
		ziehungstage = make([]string, 0)

		for _, ziehung := range ziehungen {
			ziehungstage = append(ziehungstage, ziehung.Datum.Format("2006-01-02"))
		}
		response = lottojson.GetZiehungenResponse{
			Errormessage: "",
			Ziehungen:    ziehungstage,
		}
	}

	return response

}

func ErstelleResponseAufZiehungseroeffnung(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var mitarbeiter database.Nutzer
	var neueZiehung database.Ziehung

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if request.Param["datum"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum der Ziehung angeben",
		}
		return response
	} else {

		var timeErr error

		neueZiehung.Datum, timeErr = time.Parse("2006-01-02", request.Param["datum"])

		if timeErr != nil {
			response = lottojson.ErrorResponse{
				Errormessage: timeErr.Error(),
			}
			return response
		}

		mitarbeiter = aktiveNutzer[request.Auth]

		insertFehler := lottologic.EroeffneZiehung(neueZiehung, mitarbeiter)

		if insertFehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: insertFehler.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}
	}

	return response

}

func ErstelleResponseAufTippabgabe(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var neuerTipp database.Tipp
	var nutzer database.Nutzer

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
		return response
	}

	if request.Param["tipp"] == "" || request.Param["datum"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Datum und Tipp angeben",
		}
		return response
	} else {

		neuerTipp = database.Tipp{
			Ziehung: request.Param["tipp"],
		}

		var parserFehler error

		neuerTipp.Datum, parserFehler = time.Parse("2006-01-02", request.Param["datum"])

		if parserFehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: parserFehler.Error(),
			}
			return response
		}

		nutzer = aktiveNutzer[request.Auth]

		insertError := lottologic.FuegeTippNachPruefungEin(neuerTipp, nutzer)

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

func ErstelleResponseAufKontoloeschung(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var fehler error

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {

		nutzer := aktiveNutzer[request.Auth]

		fehler = database.LoescheNutzerdaten(nutzer.Benutzername)

		if fehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: fehler.Error(),
			}
		} else {
			response = lottojson.ErrorResponse{
				Errormessage: "",
			}
		}

	}

	return response

}

func ErstelleResponseAufKontoaenderung(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var fehler error

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {

		if request.Param["neuespasswort"] == "" {
			response = lottojson.ErrorResponse{
				Errormessage: "kein neues Passwort uebergeben",
			}
		} else {

			nutzer := aktiveNutzer[request.Auth]
			fehler = lottologic.AendereSpielerdatenNachPruefung(nutzer.Benutzername, nutzer.Benutzername, request.Param["neuespasswort"])

			if fehler != nil {
				response = lottojson.ErrorResponse{
					Errormessage: fehler.Error(),
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

func ErstelleResponseAufRegistrierung(request lottojson.LottoRequest) interface{} {

	var response interface{}
	var neuerNutzer database.Nutzer

	if request.Param["name"] == "" || request.Param["passwort"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		neuerNutzer = database.Nutzer{
			Benutzername: request.Param["name"],
			Ist_spieler:  true,
		}

		fehler := lottologic.FuegeSpielerNachPruefungEin(neuerNutzer, request.Param["passwort"])

		if fehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: fehler.Error(),
			}
		} else {

			auth := LoginNutzer(neuerNutzer)
			response = lottojson.RegistrationResponse{
				Errormessage: "",
				Auth:         auth,
			}
		}
	}

	return response

}

func ErstelleResponseAufLogin(request lottojson.LottoRequest) interface{} {

	var response interface{}

	if request.Param["name"] == "" || request.Param["passwort"] == "" {
		response = lottojson.ErrorResponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		nutzer, fehler := database.HoleNutzerdatenZumNamen(request.Param["name"])

		if fehler != nil {
			response = lottojson.ErrorResponse{
				Errormessage: fehler.Error(),
			}
		} else {

			if !lottologic.PruefePasswortHash(request.Param["passwort"], nutzer.Pw_hash) {
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

	return response

}

func ValidiereAuthToken(authToken string) bool {

	_, gueltig := aktiveNutzer[authToken]

	return gueltig

}

func ErstelleResponseAufLogout(request lottojson.LottoRequest) interface{} {

	var response interface{}

	if !ValidiereAuthToken(request.Auth) {
		response = lottojson.ErrorResponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {
		LogoutNutzer(request.Auth)
		response = lottojson.ErrorResponse{
			Errormessage: "",
		}
	}

	return response

}

func InitialisiereNutzerliste() {

	aktiveNutzer = make(map[string]database.Nutzer)

}

func LoginNutzer(nutzer database.Nutzer) string {

	authToken := ksuid.New().String()

	aktiveNutzer[authToken] = nutzer

	lottolog.InfoLogger.Printf("Nutzer %s hat sich eingeloggt", nutzer.Benutzername)

	return authToken
}

func LogoutNutzer(authToken string) {

	lottolog.InfoLogger.Printf("Nutzer %s hat sich ausgeloggt", aktiveNutzer[authToken].Benutzername)

	delete(aktiveNutzer, authToken)

}
