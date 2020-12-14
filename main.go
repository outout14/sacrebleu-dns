package main

import (
	"database/sql"
	"flag"
	"strconv"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-dns/core"
	"github.com/outout14/sacrebleu-dns/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//Global vars
var conf *utils.Conf
var DB *sql.DB

//Main loop
func main() {
	//Get the config patch from --config flag
	configPatch := flag.String("config", "extra/config.ini.example", "the patch to the config file")
	flag.Parse()

	//Load the INI configuration file
	conf = new(utils.Conf)
	err := ini.MapTo(conf, *configPatch)
	utils.CheckErr(err)

	//Set up the Logrus logger
	utils.InitLogger(conf)

	//Attach DNS request handler func for all domains
	dns.HandleFunc(".", core.HandleDnsRequest)

	//Initialize the redis database
	utils.RedisDatabase(conf)

	//Initialize the sql database
	utils.SqlDatabase(conf)

	//Start the DNS server
	server := &dns.Server{Addr: conf.App.Ip + strconv.Itoa(conf.App.Port), Net: "udp"}                   //define the server
	logrus.WithFields(logrus.Fields{"ip": conf.App.Ip, "port": conf.App.Port}).Infof("SERVER : Started") //log
	err = server.ListenAndServe()                                                                        //start it
	utils.CheckErr(err)

	defer server.Shutdown() //shut down on application closing
}
