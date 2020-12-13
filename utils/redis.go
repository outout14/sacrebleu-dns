package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

var redisDb *redis.Client

func RedisDatabase(conf *Conf) *redis.Client {
	logrus.WithFields(logrus.Fields{"ip": conf.Redis.Ip, "port": conf.Redis.Port}).Infof("REDIS : Connection to DB")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", conf.Redis.Ip, conf.Redis.Port),
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})

	//Test Redis connection
	err := rdb.Set(ctx, "alive", 1, 0).Err()
	alive, err := rdb.Get(ctx, "alive").Result()
	CheckErr(err)
	if alive != "1" {
		logrus.WithFields(logrus.Fields{"alive": alive}).Panic("REDIS : Test not passed. alive != 1")
	}
	CheckErr(err)
	logrus.WithFields(logrus.Fields{"db": conf.Redis.Db}).Info("REDIS : Successfull connection")

	redisDb = rdb

	return rdb
}

func redisCheckForRecord(redisKey string, entry Record) (Record, error) {
	val, err := redisDb.Get(ctx, redisKey).Result()

	//If Record in Redis cache
	if err == nil {
		err := json.Unmarshal([]byte(val), &entry)
		logrus.Debugf("REDIS : %s => %s", redisKey, entry.Content)
		return entry, err
	} else {
		//Else return nil
		return entry, redis.Nil
	}
}

func redisSet(c *redis.Client, key string, ttl time.Duration, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, p, ttl).Err()
}
