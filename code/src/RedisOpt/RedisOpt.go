package Redis

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

//RedisOpt .
type RedisOpt struct {
	redisClusterClient *redis.ClusterClient
	clusterAddr        []string

	redisSingleClient *redis.Client
	signleAddr        string
	dbNum             int

	redisPasswd string
	//redis mode: 1 single mode, 2 cluster mode
	redisMothed int
}

//InitSingle .k
//Init single client
func (opt *RedisOpt) InitSingle(redisAddr, redisPasswd string, DBNum int) error {
	opt.redisMothed = 1
	opt.signleAddr = redisAddr
	opt.redisPasswd = redisPasswd
	opt.dbNum = DBNum
	opt.redisSingleClient = redis.NewClient(&redis.Options{
		Addr:     opt.signleAddr,
		Password: opt.redisPasswd,
		DB:       DBNum,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := opt.redisSingleClient.Ping().Result()
	return err
}

//InitCluster .
//Init cluster client
func (opt *RedisOpt) InitCluster(ClusterAddr []string, redisPasswd string) error {
	opt.redisMothed = 2
	opt.clusterAddr = ClusterAddr
	opt.redisPasswd = redisPasswd
	opt.redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    opt.clusterAddr,
		Password: opt.redisPasswd,
	})

	statusCmd := opt.redisClusterClient.Ping()
	_, err := statusCmd.Result()
	return err
}

func (opt *RedisOpt) reconnectRedis() {
	log.Error("reconnect RedisServer")
	var err error
	if opt.redisMothed == 2 {
		err = opt.InitCluster(opt.clusterAddr, opt.redisPasswd)
	} else {
		err = opt.InitSingle(opt.signleAddr, opt.redisPasswd, opt.dbNum)
	}
	if err != nil {
		log.Error("reconnect Redis failed")
		return
	}
	log.Info("reconnect Redis success")
}

//Keys get all key
func (opt *RedisOpt) Keys(pattern string) ([]string, error) {
	var result *redis.StringSliceCmd
	if opt.redisMothed == 1 {
		result = opt.redisSingleClient.Keys(pattern)
		//return result.Val()
	} else {
		return opt.ClusterKeys(pattern)
		//return result.Val()
	}
	//var result []string
	return result.Val(), result.Err()
}

//ClusterKeys 集群获取所有的key
func (opt *RedisOpt) ClusterKeys(pattern string) ([]string, error) {
	if opt.redisMothed != 2 {
		return nil, nil
	}
	result := []string{}
	for _, addr := range opt.clusterAddr {
		IPPort := strings.Split(addr, ":")
		log.Debug("redis:", IPPort[0], ",", IPPort[1])
		cmd := exec.Command("redis-cli", "-h", IPPort[0], "-p", IPPort[1], "keys", pattern)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			continue
		}
		retArray := strings.Split(out.String(), "\n")
		result = append(result, retArray...)
	}
	return result, nil
}

//Exist check if key exist
func (opt *RedisOpt) Exist(key string) (bool, error) {
	var result *redis.BoolCmd
	if opt.redisMothed == 2 {
		result = opt.redisClusterClient.Exists(key)
	} else {
		result = opt.redisSingleClient.Exists(key)
	}
	return result.Val(), result.Err()
}

//Set .
func (opt *RedisOpt) Set(key string, data interface{}, sec int) error {
	var cmd *redis.StatusCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Set(key, data, time.Duration(sec)*time.Second)
	} else {
		cmd = opt.redisSingleClient.Set(key, data, time.Duration(sec)*time.Second)
	}
	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	return err
}

//Get  返回值说明：返回结果，是否找到，错误码
func (opt *RedisOpt) Get(key string) ([]byte, bool, error) {
	var cmd *redis.StringCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Get(key)
	} else {
		cmd = opt.redisSingleClient.Get(key)
	}
	data, err := cmd.Bytes()
	if err != nil {
		if err.Error() != "redis: nil" {
			opt.reconnectRedis()
		} else {
			return data, false, nil
		}
		return data, false, err
	}
	return data, true, nil
}

//Delete .
func (opt *RedisOpt) Delete(key string) error {
	var cmd *redis.IntCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Del(key)
	} else {
		cmd = opt.redisSingleClient.Del(key)
	}
	if cmd.Err() != nil {
		opt.reconnectRedis()
		return cmd.Err()
	}
	return nil
}

//-------------------------------
//HExists .
func (opt *RedisOpt) HExists(key, field string) bool {
	var cmd *redis.BoolCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HExists(key, field)
	} else {
		cmd = opt.redisSingleClient.HExists(key, field)
	}

	status, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	return status
}

//HSet .
func (opt *RedisOpt) HSet(key, field string, value interface{}, second int) error {
	var cmd *redis.BoolCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HSet(key, field, value)
	} else {
		cmd = opt.redisSingleClient.HSet(key, field, value)
	}

	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	if second > 0 {
		opt.redisClusterClient.Expire(key, time.Duration(second)*time.Second)
	}
	return err
}

//HMSet .
func (opt *RedisOpt) HMSet(key string, fields map[string]string, second int) error {
	var cmd *redis.StatusCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HMSet(key, fields)
	} else {
		cmd = opt.redisSingleClient.HMSet(key, fields)
	}
	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	if second > 0 {
		opt.redisClusterClient.Expire(key, time.Duration(second)*time.Second)
	}
	return err
}

//HGet 返回值说明：返回结果，是否找到，错误码
func (opt *RedisOpt) HGet(key, field string) (string, bool, error) {
	var cmd *redis.StringCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HGet(key, field)
	} else {
		cmd = opt.redisSingleClient.HGet(key, field)
	}

	data, err := cmd.Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			opt.reconnectRedis()
		} else {
			return data, false, nil
		}
		return data, false, err
	}
	return data, true, nil
}

//HGetAll 返回值说明：返回结果，是否找到，错误码
func (opt *RedisOpt) HGetAll(key string) (map[string]string, bool, error) {
	var cmd *redis.StringStringMapCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HGetAll(key)
	} else {
		cmd = opt.redisSingleClient.HGetAll(key)
	}

	data, err := cmd.Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			opt.reconnectRedis()
		} else {
			return data, false, nil
		}
		return data, false, err
	}
	return data, true, nil
}

//HDelete .
func (opt *RedisOpt) HDelete(key string) error {
	var cmd *redis.IntCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.HDel(key)
	} else {
		cmd = opt.redisSingleClient.HDel(key)
	}

	if cmd.Err() != nil {
		if cmd.Err().Error() != "redis: nil" {
			opt.reconnectRedis()
		}
		return cmd.Err()
	}
	return nil
}
