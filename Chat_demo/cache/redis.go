package cache

import (
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

var (
	RedisClient *redis.Client

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func Init() {
	// 读取配置
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadRedis(file)
	NewRedis()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
}

func NewRedis() {
	db_num, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		DB:       int(db_num),
		Password: RedisPw,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}
