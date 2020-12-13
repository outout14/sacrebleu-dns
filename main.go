package main

import (
	"database/sql"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/miekg/dns"
	"github.com/outout14/lighthouse/core"
	"github.com/outout14/lighthouse/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//Global vars
var conf *utils.Conf
var DB *sql.DB
var redisDb *redis.Client

//Main loop
func main() {
	//Load Configuration
	conf = new(utils.Conf)
	err := ini.MapTo(conf, "./config.ini")
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
