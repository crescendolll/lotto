package lottoapi

import (
	"fmt"
	"lotto/database"
	"lotto/lottojson"
	"lotto/lottologic"

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
	default:
		response = lottojson.Errorresponse{
			Errormessage: "Unbekannte Methode",
		}
	}

	return response

}

func CreateRegistrationResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}
	var neuerNutzer lottologic.Nutzer

	if apirequest.Param["name"] == "" || apirequest.Param["passwort"] == "" {
		response = lottojson.Errorresponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		neuerNutzer = lottologic.Nutzer{
			Benutzername: apirequest.Param["name"],
			Pw_hash:      apirequest.Param["passwort"],
			Ist_spieler:  true,
		}

		databasehandle := database.OpenLottoConnection()
		insertError := database.InsertIntoSpieler(databasehandle, neuerNutzer.Benutzername, neuerNutzer.Pw_hash)
		database.CloseLottoConnection(databasehandle)

		if insertError != nil {
			response = lottojson.Errorresponse{
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
		response = lottojson.Errorresponse{
			Errormessage: "Name und Passwort sind nicht belegt",
		}
	} else {

		databasehandle := database.OpenLottoConnection()
		nutzer, err := database.SelectFromSpielerByName(databasehandle, apirequest.Param["name"])
		database.CloseLottoConnection(databasehandle)

		if err != nil {
			response = lottojson.Errorresponse{
				Errormessage: err.Error(),
			}
		} else {

			if !lottologic.CheckPasswordHash(apirequest.Param["passwort"], nutzer.Pw_hash) {
				response = lottojson.Errorresponse{
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

func CreateLogoutResponse(apirequest lottojson.LottoRequest) interface{} {

	var response interface{}

	auth := apirequest.Auth

	_, found := aktiveNutzer[auth]

	if !found {
		response = lottojson.Errorresponse{
			Errormessage: "Unauthorisierter Zugriff",
		}
	} else {
		LogoutNutzer(auth)
		response = lottojson.LogoutResponse{
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
