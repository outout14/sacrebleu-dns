package utils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//DB SQL database as global var
var db *gorm.DB

//SQLDatabase Initialize the SQL Database
//Requires a conf struct
func SQLDatabase(conf *Conf) *gorm.DB {
	logrus.WithFields(logrus.Fields{"database": conf.Database.Db, "driver": conf.Database.Type}).Infof("SQL : Connection to DB")
	var gormLogLevel logger.LogLevel

	//Set GORM log level based on conf AppMode
	if conf.AppMode != "production" {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Silent
	}

	//Connect to the Database
	if conf.Database.Type == "postgresql" {
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable", conf.Database.Username, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Db)
		DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(gormLogLevel),
		})
		CheckErr(err)
		db = DB
		return DB
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Database.Username, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Db)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	CheckErr(err)
	db = DB

	return DB
}

//SQLMigrate : Launch the database migration (creation of tables)
func SQLMigrate() {
	logrus.Info("SQL : Database migration launched")
	db.AutoMigrate(&Record{})
}

//Check for a record in the SQL database
func sqlCheckForRecord(redisKey string, dKey string, entry Record) ([]Record, bool) {
	var records []Record

	rows, err := db.Where("fqdn = ? AND type = ?", dKey, entry.Qtype).Model(&Record{}).Rows()
	if err != nil {
		return records, true
	}
	defer rows.Close()
	for rows.Next() {
		var entry Record
		db.ScanRows(rows, &entry)

		if entry.Content != "" { //If Record content not empty
			records = append(records, entry)
		}
	}
	//Cache the request in Redis if any result
	_ = redisSet(redisDb, redisKey, 30*time.Second, records) //Set it in the Redis database for 30sec
	return records, false
}

//Check for a wildcard record in the SQL database
func sqlCheckForReverse6Wildcard(redisKey string, dKey string, entry Record) []Record {
	returnedEntry := entry

	rows, err := db.Table("records").Select("id", "content", "fqdn").Where("fqdn LIKE ?", "*%.ip6.arpa.").Rows()

	DbgErr(err) //Check for empty row or non important error

	var records []Record

	//For each result check if it match the reverse IP
	for rows.Next() {
		rows.Scan(&returnedEntry.ID, &returnedEntry.Content, &returnedEntry.Fqdn)
		CheckErr(err)

		//Check if the record is matching the reversed IP
		if checkReverse6(entry, returnedEntry) {
			logrus.Debug("REVERSE : Correct wildcard reverse.")
			//Cache the request in Redis if any result
			_ = redisSet(redisDb, redisKey, 10*time.Second, returnedEntry)
			records = append(records, entry)
		}
		logrus.Debug("REVERSE : WRONG wildcard reverse .")
	}

	return records

}
