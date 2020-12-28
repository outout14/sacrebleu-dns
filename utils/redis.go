package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/outout14/sacrebleu-api/api/types"
	"github.com/sirupsen/logrus"
)

//Redis context
var ctx = context.Background()

//Redis client as global var
var redisDb *redis.Client

//RedisDatabase : Initialize the Redis Database
//Requires a conf struct
//Return a *redis.Client
func RedisDatabase(conf *Conf) *redis.Client {
	logrus.WithFields(logrus.Fields{"ip": conf.Redis.IP, "port": conf.Redis.Port}).Infof("REDIS : Connection to DB")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", conf.Redis.IP, conf.Redis.Port),
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	}) //Connect to the DB

	//Test Redis connection
	err := rdb.Set(ctx, "alive", 1, 0).Err()
	CheckErr(err)
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

//Check for a record in the Redis database
//Requires the redis key (as string) and the record to check (struct)
//Return a Record (struct) and error (if any)
func redisCheckForRecord(redisKey string, entry types.Record) ([]types.Record, error) {
	val, err := redisDb.Get(ctx, redisKey).Result()

	var result []types.Record

	//If Record in Redis cache
	if err == nil {
		err := json.Unmarshal([]byte(val), &result)
		return result, err
	}
	return result, redis.Nil
}

//Add a record in the Redis database
//Return an error (if any)
func redisSet(c *redis.Client, key string, ttl time.Duration, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, p, ttl).Err()
}
