package handle

import (
	"WechatAPI/common"
	"WechatAPI/config"

	"WechatAPI/DBOpt"

	log "github.com/Sirupsen/logrus"
)

//InitServer 初始化服务
func InitServer() (err error) {
	config := config.GetConfig()

	if config.ServerStatus {
		return nil
	}

	err = common.RedisTokenOpt.InitSingle(config.RedisAddr, config.RedisPwd, config.RedisTokenNum)
	if err != nil {
		log.Error("err:", err)
		return err
	}

	err = common.RedisServerListOpt.InitSingle(config.RedisAddr, config.RedisPwd, config.RedisServerNum)
	if err != nil {
		log.Error("err:", err)
		return err
	}

	err = DBOpt.GetDataOpt().InitDatabase(config.Database)
	if err != nil {
		log.Error("err:", err)
		return err
	}
	return err
}
