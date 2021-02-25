package main

import (
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-api/api/types"
	"github.com/outout14/sacrebleu-dns/core"
	"github.com/outout14/sacrebleu-dns/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//Global vars
//conf Configuration
var conf *utils.Conf

//DB the SQL database
var DB *sql.DB

//Main loop
func main() {
	configPatch := flag.String("config", "config.ini", "the patch to the config file")  //Get the config patch from --config flag
	sqlMigration := flag.Bool("sqlmigrate", false, "initialize / migrate the database") //Detect if migration asked
	flag.Parse()

	//Load the INI configuration file
	conf = new(utils.Conf)
	err := ini.MapTo(conf, *configPatch)
	utils.CheckErr(err)

	//Set up the Logrus logger
	utils.InitLogger(conf)

	//Attach DNS request handler func for all domains
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) { core.HandleDNSRequest(w, r, conf) })

	//Initialize the redis database
	utils.RedisDatabase(conf)

	//Initialize the sql database
	db := utils.SQLDatabase(conf)
	if *sqlMigration {
		types.SQLMigrate(db)
	}

	//Start the UDP listener
	go func() {
		server := &dns.Server{Addr: conf.App.IP + ":" + strconv.Itoa(conf.App.Port), Net: "udp"}                          //define the server
		logrus.WithFields(logrus.Fields{"ip": conf.App.IP, "port": conf.App.Port}).Infof("SERVER : Started UDP listener") //log
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	//Start the TCP listener
	go func() {
		server := &dns.Server{Addr: conf.App.IP + ":" + strconv.Itoa(conf.App.Port), Net: "tcp"}                          //define the server
		logrus.WithFields(logrus.Fields{"ip": conf.App.IP, "port": conf.App.Port}).Infof("SERVER : Started TCP listener") //log
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	logrus.Infof("Signal (%v) received, stopping\n", s)

}
