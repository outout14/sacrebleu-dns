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
	//Get config patch
	configPatch := flag.String("config", "extra/config.ini.example", "the patch to the config file")
	flag.Parse()

	//Load Configuration
	conf = new(utils.Conf)
	err := ini.MapTo(conf, *configPatch)
	utils.CheckErr(err)

	utils.InitLogger(conf)

	// attach request handler func
	dns.HandleFunc(".", core.HandleDnsRequest)

	//Init redis database
	utils.RedisDatabase(conf)

	//Init sql database
	utils.SqlDatabase(conf)

	// start server
	server := &dns.Server{Addr: conf.App.Ip + strconv.Itoa(conf.App.Port), Net: "udp"} //define the server
	logrus.WithFields(logrus.Fields{"ip": conf.App.Ip, "port": conf.App.Port}).Infof("SERVER : Started")
	err = server.ListenAndServe() //start it
	utils.CheckErr(err)

	defer server.Shutdown() //shut down on application closing
}
