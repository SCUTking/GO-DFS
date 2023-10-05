package main

import (
	"GO-DFS/slaves"
	"GO-DFS/utils"
	"GO-DFS/utils/redis"
	"fmt"
	"os"
	"testing"
)

func TestGetFileHash(t *testing.T) {
	var a utils.FileUtil
	open, err := os.Open("config/config.yaml")
	if err != nil {
	}
	hash := a.GetFileHash(open)

	fmt.Println(hash)
}

func TestUuid(t *testing.T) {
	uuid := utils.Uuid()

	fmt.Println(uuid)
}

func TestMonGoDB(t *testing.T) {

}

func TestRS(t *testing.T) {

}

func TestRecover(t *testing.T) {
	slaves.Decode()

}

func TestRedis(t *testing.T) {
	//初始化
	conf := loadConf()
	redis.Setup(conf)
	err := redis.SetForever("test", "asdeas21312312")
	if err != nil {
		fmt.Println(err.Error())
	}

	redis.GetFileInfoByHash("f31a03f6e6baf85800b3db413acc36341fccd76d900914fc15c3912347eb8f7d")
}
