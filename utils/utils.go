package utils

import (
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

//If fatal error, log it and panic
func CheckErr(err error) {
	if err != nil {
		log.Fatalf("%s\n ", err.Error())
	}
}

//If basic error, log it as classic error but don't panic and keep kalm
func DbgErr(err error) {
	if err != nil {
		log.Errorf("%s\n ", err.Error())
	}
}

//Init the logrus logger with rotateFileHook.
//Conf struct passed to get informations about the logger (debug or not)
func InitLogger(conf *Conf) {
	var logLevel = logrus.InfoLevel //By default the level is Info.

	if conf.App_mode != "production" { //If the configuration contains anything different than "production"; the level is set to Debug
		logLevel = logrus.DebugLevel
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{ //Rotate file hook, By default 50Mb max and 28 days retention
		Filename:   conf.App.Logdir + "/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)                        //Set the log level
	logrus.SetOutput(colorable.NewColorableStdout()) //Force colors in the Stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	if conf.App.Logfile { //If file logging is enabled
		logrus.AddHook(rotateFileHook)
	}

	log.WithFields(log.Fields{"app_mode": conf.App_mode}).Info("Application mode")
	log.WithFields(log.Fields{"logLevel": logLevel}).Debug("Log level")
}

//Check if a reverse wildcard correspond to a record using strings.Contains
//Return bool
func checkReverse6(entry Record, result Record) bool {
	check := strings.Replace(entry.Fqdn, result.Fqdn[1:], "", 1)
	logrus.WithFields(logrus.Fields{"entry": entry.Fqdn, "result": result.Fqdn[1:]}).Debug("REVERSE checkReverse6 :")
	logrus.Debugf("REVERSE checkReverse6 : %s", check)
	if strings.Contains(check, IP6arpa) {
		return false
	} else {
		return true
	}
}
