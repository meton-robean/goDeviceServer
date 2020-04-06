package common

import (
	Redis "RedisOpt"
)

var errCodeMap map[int]string

//RedisTokenOpt 操作
var RedisTokenOpt *Redis.RedisOpt

//RedisServerListOpt 服务列表
var RedisServerListOpt *Redis.RedisOpt

func init() {
	RedisTokenOpt = &Redis.RedisOpt{}
	RedisServerListOpt = &Redis.RedisOpt{}

	errCodeMap = make(map[int]string)
	errCodeMap[0] = "成功"
	errCodeMap[10000] = "其他"
	errCodeMap[10001] = "secret 不正确"
	errCodeMap[10002] = "接口凭证过期"
	errCodeMap[10003] = "参数出错"
	errCodeMap[10004] = "房间号不存在"
	errCodeMap[10005] = "设备ID不存在"
	errCodeMap[10006] = "数据库服务异常"
	errCodeMap[10007] = "Redis服务异常"
	errCodeMap[10008] = "网关不在线"
	errCodeMap[10009] = "网关已存在"
	errCodeMap[10010] = "用户数据不匹配"
	errCodeMap[10011] = "设备已被绑定"
	errCodeMap[10012] = "用户ID不存在"
	errCodeMap[10013] = "网关不存在"
	errCodeMap[10014] = "用户没有登陆权限"
	errCodeMap[10015] = "房间号已绑定"
}

//GetErrCodeJSON 获取错误信息
func GetErrCodeJSON(code int) map[string]interface{} {
	data := make(map[string]interface{})
	data["code"] = code
	data["errmsg"] = errCodeMap[code]
	return data
}
