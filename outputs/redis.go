package outputs

import (
	"gopkg.in/redis.v3"
)

var RedisClient *redis.Client
var LogKey string

func SetupRedisClient(uri string, key string) (err error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})

	Log.Debug(RedisClient)
	_, err = RedisClient.Ping().Result()
	if err != nil {
		RedisClient = nil
	}

	LogKey = key
	return
}

func ShipRedis(event []byte) (err error) {
	RedisClient.RPush(LogKey, string(event))
	return
}
