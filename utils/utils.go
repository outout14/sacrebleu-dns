package utils

import (
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalf("%s\n ", err.Error())
		panic(err)
	}
}

func InitLogger(conf *Conf) {
	var logLevel = logrus.InfoLevel

	if conf.App_mode != "production" {
		logLevel = logrus.DebugLevel
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
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

	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	if conf.App.Logfile {
		logrus.AddHook(rotateFileHook)
	}

	log.WithFields(log.Fields{"app_mode": conf.App_mode}).Info("Application mode")
	log.WithFields(log.Fields{"logLevel": logLevel}).Debug("Log level")
}

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
