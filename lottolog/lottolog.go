package lottolog

import (
	"log"
	"os"
)

var (
	WarnungLogger *log.Logger
	InfoLogger    *log.Logger
	FehlerLogger  *log.Logger
)

func OeffneLogdatei() {
	logdatei, fehler := os.OpenFile(os.Getenv("LOGPATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fehler != nil {
		log.Fatal(fehler)
	}

	InitialisiereLogLevel(logdatei)
}

func OeffneTestLogdatei() {
	logdatei, fehler := os.OpenFile(os.Getenv("TESTLOGPATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fehler != nil {
		log.Fatal(fehler)
	}
	InitialisiereLogLevel(logdatei)

}

func InitialisiereLogLevel(logdatei *os.File) {

	InfoLogger = log.New(logdatei, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnungLogger = log.New(logdatei, "WARNUNG: ", log.Ldate|log.Ltime|log.Lshortfile)
	FehlerLogger = log.New(logdatei, "FEHLER: ", log.Ldate|log.Ltime|log.Lshortfile)

}
