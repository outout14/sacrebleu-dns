package utils

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB SQL database as global var
var db *gorm.DB

//SQLDatabase Initialize the (My)SQL Database
//Requires a conf struct
func SQLDatabase(conf *Conf) {
	logrus.WithFields(logrus.Fields{"database": conf.Database.Db, "driver": conf.Database.Type}).Infof("SQL : Connection to DB")
	//Connect to the Database
	var err error

	if conf.Database.Type == "postgresql" {
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable", conf.Database.Username, conf.Database.Password, conf.Database.IP, conf.Database.Port, conf.Database.Db)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Database.Username, conf.Database.Password, conf.Database.IP, conf.Database.Port, conf.Database.Db)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	CheckErr(err)
}

//Check for a record in the SQL database
func sqlCheckForRecord(redisKey string, dKey string, entry Record) (Record, bool) {
	db.Where("name = ? AND type = ?", dKey, entry.Qtype).First(&entry)

	logrus.Debugf("SQL : %s => %s", entry.Fqdn, entry.Content) //log the result

	if entry.Content != "" { //If Record content not empty
		//Cache the request in Redis if any result
		logrus.Debugf("REDIS : Set entry for %s", redisKey)
		_ = redisSet(redisDb, redisKey, 30*time.Second, entry) //Set it in the Redis database for 30sec
		return entry, false
	}
	//Else return 1 for err
	return entry, true

}

//Check for a wildcard record in the SQL database
func sqlCheckForReverse6Wildcard(redisKey string, dKey string, entry Record) (Record, error) {
	returnedEntry := entry

	rows, err := db.Table("records").Select("id", "content", "name").Where("name LIKE ?", "*%.ip6.arpa.").Rows()

	DbgErr(err) //Check for empty row or non important error

	//For each result check if it match the reverse IP
	for rows.Next() {
		err = rows.Scan(&returnedEntry.ID, &returnedEntry.Content, &returnedEntry.Fqdn)
		CheckErr(err)

		//Check if the record is matching the reversed IP
		if checkReverse6(entry, returnedEntry) {
			logrus.Debug("REVERSE : Correct wildcard reverse.")
			//Cache the request in Redis if any result
			_ = redisSet(redisDb, redisKey, 10*time.Second, returnedEntry)
			return returnedEntry, err
		}
		logrus.Debug("REVERSE : WRONG wildcard reverse .")
	}

	return entry, redis.Nil

}
