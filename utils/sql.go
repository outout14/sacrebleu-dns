package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql" //MySQL driver
	"github.com/sirupsen/logrus"
)

//DB SQL database as global var
var DB *sql.DB

//SQLDatabase Initialize the (My)SQL Database
//Requires a conf struct
func SQLDatabase(conf *Conf) {
	logrus.WithFields(logrus.Fields{"database": conf.Database.Db}).Infof("SQL : Connection to DB")
	//Connect to the Database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Database.Username, conf.Database.Password, conf.Database.IP, conf.Database.Port, conf.Database.Db))
	CheckErr(err)
	DB = db
	SQLTest() //Test SQL connexion
}

//SQLTest : Test the SQL connexion by selecting all records from the database
func SQLTest() {
	_, err := DB.Query("SELECT name, content FROM records")
	CheckErr(err) //Panic if any error
}

//Check for a record in the SQL database
func sqlCheckForRecord(redisKey string, dKey string, entry Record) (Record, int) {
	dbg := DB.QueryRow(
		"SELECT id, content, ttl FROM records WHERE `name` = ? AND `type` = ?;", dKey, entry.Qtype).Scan(
		&entry.ID,
		&entry.Content,
		&entry.TTL,
	)

	if dbg != nil { //If any err
		logrus.Debugf("SQL : %v", dbg)
	}

	logrus.Debugf("SQL : %s => %s", entry.Fqdn, entry.Content) //log the result

	if entry.Content != "" { //If Record content not empty
		//Cache the request in Redis if any result
		logrus.Debugf("REDIS : Set entry for %s", redisKey)
		_ = redisSet(redisDb, redisKey, 30*time.Second, entry) //Set it in the Redis database for 30sec
		return entry, 0
	}
	//Else return 1 for err
	return entry, 1

}

//Check for a wildcard record in the SQL database
func sqlCheckForReverse6Wildcard(redisKey string, dKey string, entry Record) (Record, error) {
	returnedEntry := entry

	results, err := DB.Query("SELECT id, content, name FROM records WHERE name LIKE '*%.ip6.arpa.';") //Get ALL reverse IPs
	DbgErr(err)                                                                                       //Check for empty row or non important error

	//For each result check if it match the reverse IP
	for results.Next() {
		err = results.Scan(&returnedEntry.ID, &returnedEntry.Content, &returnedEntry.Fqdn)
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
