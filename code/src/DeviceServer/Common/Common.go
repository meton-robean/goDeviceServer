package Common

import (
	"DeviceServer/Config"
	Redis "RedisOpt"
	"bytes"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

var ServerStarTime int64

//DefaultHead 文件头
const DefaultHead = "HTTP-JSON-BOCHIOT"

//RedisServerOpt 服务列表
var RedisServerOpt *Redis.RedisOpt

func InitCommon() error {
	Config.InitConfig()
	config := Config.GetConfig()

	RedisServerOpt = &Redis.RedisOpt{}
	err := RedisServerOpt.InitSingle(config.RedisAddr, config.RedisPwd, config.RedisServerNum)
	if err != nil {
		log.Error("err:", err)
	}
	return err
}

func execshell(s string) string {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	localIP := out.String()
	return localIP[:len(localIP)-1]
}

//GetLocalIP 获取本地IP地址
func GetLocalIP() string {
	return execshell("ifconfig | grep ^e -A2 | grep 'inet addr' | awk '{print $2}' | awk -F: '{print $2}'")
}
