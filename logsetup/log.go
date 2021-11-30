package logsetup

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func LogFile() (*os.File, error) {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //opening a log file
	if err != nil {
		return nil, err
	}
	formater := new(log.TextFormatter)
	formater.TimestampFormat = "02-01-2006 15:04:05"
	formater.FullTimestamp = true
	log.SetFormatter(formater)
	return file, nil
}
