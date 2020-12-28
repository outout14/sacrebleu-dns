package main

import (
	"database/sql"
	"flag"
	"strconv"

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
	dns.HandleFunc(".", core.HandleDNSRequest)

	//Initialize the redis database
	utils.RedisDatabase(conf)

	//Initialize the sql database
	db := utils.SQLDatabase(conf)
	if *sqlMigration {
		types.SQLMigrate(db)
	}

	//Start the DNS server
	server := &dns.Server{Addr: conf.App.IP + ":" + strconv.Itoa(conf.App.Port), Net: "udp"}             //define the server
	logrus.WithFields(logrus.Fields{"ip": conf.App.IP, "port": conf.App.Port}).Infof("SERVER : Started") //log
	err = server.ListenAndServe()                                                                        //start it
	utils.CheckErr(err)

	defer server.Shutdown() //shut down on application closing
}
