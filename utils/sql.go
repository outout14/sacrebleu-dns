package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func SqlDatabase(conf *Conf) {
	logrus.WithFields(logrus.Fields{"database": conf.Database.Db}).Infof("SQL : Connection to DB")
	//db conn
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Database.Username, conf.Database.Password, conf.Database.Ip, conf.Database.Port, conf.Database.Db))
	CheckErr(err)
	DB = db

	// if there is an error opening the connection, handle it
	CheckErr(err)
}

func SqlTest() {
	results, err := DB.Query("SELECT name, content FROM records")
	CheckErr(err)

	for results.Next() {
		var record Record
		err = results.Scan(&record.Fqdn, &record.Content)
		CheckErr(err)
		logrus.Debugf(record.Content)
	}
}

func sqlCheckForRecord(redisKey string, dKey string, entry Record) (Record, int) {
	dbg := DB.QueryRow(
		"SELECT id, content, ttl FROM lighthouse.records WHERE `name` = ? AND `type` = ?;", dKey, entry.Qtype).Scan(
		&entry.Id,
		&entry.Content,
		&entry.TTL,
	)

	//logrus.WithFields(logrus.Fields{"name": dKey, "type": entry.Qtype}).Debugf("SQL : ")

	if dbg != nil {
		logrus.Debugf("SQL : %v", dbg)
	}

	logrus.Debugf("SQL : %s => %s", entry.Fqdn, entry.Content)

	if entry.Content != "" {
		//Cache the request in Redis if any result
		logrus.Debugf("REDIS : Set entry for %s", redisKey)
		logrus.Warningf("REDIS : %s", redisKey)
		_ = redisSet(redisDb, redisKey, 30*time.Second, entry)
		return entry, 0
	} else {
		//Else return nil
		return entry, 1
	}
}

func sqlCheckForReverse6Wildcard(redisKey string, dKey string, entry Record) (Record, error) {
	returnedEntry := entry

	results, err := DB.Query("SELECT id, content, name FROM lighthouse.records WHERE name LIKE '*%.ip6.arpa.';")

	for results.Next() {
		err = results.Scan(&returnedEntry.Id, &returnedEntry.Content, &returnedEntry.Fqdn)
		CheckErr(err)

		if checkReverse6(entry, returnedEntry) {
			logrus.Debug("REVERSE : Correct wildcard reverse.")
			//Cache the request in Redis if any result
			_ = redisSet(redisDb, redisKey, 10*time.Second, returnedEntry)
			return returnedEntry, err
		} else {
			logrus.Debug("REVERSE : WRONG wildcard reverse .")
		}
	}

	return entry, redis.Nil

}
