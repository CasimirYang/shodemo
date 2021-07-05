package log

import (
	"github.com/cihub/seelog"
	"log"
)

var Logger seelog.LoggerInterface

func init() {
	var err error
	Logger, err = seelog.LoggerFromConfigAsFile("../seelog.xml")

	if err != nil {
		log.Fatal(err)
	}

	if seelog.ReplaceLogger(Logger) != nil {
		log.Fatal(err)
	}
}
