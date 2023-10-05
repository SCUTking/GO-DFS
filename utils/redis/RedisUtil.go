package redis

import (
	"GO-DFS/config"
	"GO-DFS/model"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/goccy/go-json"
)

var redisdb *redis.Client

// Setup redis链接池
func Setup(config *config.GlobalConfig) (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr, // 指定
		Password: config.RedisPassword,
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err = redisdb.Ping().Result()
	return nil
}

func SetHash(key string, data interface{}, time int) error {

	redisdb.Set(key, data, 0)

	//value, err := json.Marshal(data)
	//if err != nil {
	//	return err
	//}

	return nil
}

func SetForever(key string, data interface{}) error {

	redisdb.Set(key, data, 0)
	//value, err := json.Marshal(data)
	//if err != nil {
	//	return err
	//}

	return nil
}

func GetFileInfoByHash(hashStr string) model.FileInfo {
	result, err := redisdb.Get(hashStr).Result()

	if err != nil {
		fmt.Println(err)
	}

	resultBytes := []byte(result)
	var fileInfoTemp model.FileInfo
	json.Unmarshal(resultBytes, &fileInfoTemp)

	return fileInfoTemp

}
