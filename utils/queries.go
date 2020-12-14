package utils

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func GetRecord(entry Record) Record {
	//Check for strict record in Redis cache
	redisKey := entry.Fqdn + "--" + fmt.Sprint(entry.Qtype)

	result, redisErr := redisCheckForRecord(redisKey, entry)

	var sqlErr int

	//If reverse DNS
	reverseCheck := IsReverse(entry.Fqdn)
	if reverseCheck > 0 {

		if redisErr == redis.Nil {
			logrus.Debug("QUERIES : Check for strict reverse in MySQL")
			result, sqlErr = sqlCheckForRecord(redisKey, entry.Fqdn, entry)
			if sqlErr == 1 {
				logrus.Debug("QUERIES : Check for wildcard reverse in MySQL")
				result, _ = sqlCheckForReverse6Wildcard(redisKey, entry.Fqdn, entry)
			}
		}

		//For dynamic reverse dns
		if strings.Contains(result.Content, "%s") {
			record := ExtractAddressFromReverse(entry.Fqdn)
			var recordFormated string
			if reverseCheck == 1 {
				recordFormated = strings.ReplaceAll(record, ".", "-")
			} else {
				recordFormated = strings.ReplaceAll(record, ":", "-")
			}
			result.Content = fmt.Sprintf(result.Content, recordFormated)
		}
	} else if redisErr == redis.Nil { //If strict record NOT in Redis cache & not Reverse
		//Check for wildcard in Redis cache
		logrus.Debug("QUERIES : Check for wildcard in redis cache")
		mainDomainKey := fmt.Sprintf("*%s", entry.Fqdn[strings.Index(entry.Fqdn, "."):]) //Remove the last subdomain
		redismdKey := fmt.Sprintf("%s--%v", mainDomainKey, entry.Qtype)
		result, redisErr = redisCheckForRecord(redismdKey, entry)
		//If none of both check in mysql
		if redisErr == redis.Nil {
			//Check for strict record in mysql
			logrus.Debug("QUERIES : Check for strict record in MSQL")
			result, sqlErr = sqlCheckForRecord(redisKey, entry.Fqdn, entry)
			if sqlErr == 1 {
				//Check for wildcard record in mysql
				logrus.Debug("QUERIES : Check for wildcard in MSQL")
				result, _ = sqlCheckForRecord(redismdKey, fmt.Sprint(mainDomainKey), entry)
			}
		}
	}

	return result
}
