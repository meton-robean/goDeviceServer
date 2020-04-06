package Redis

/**
 * 文档说明：
 * 该接口将Redis当数据使用，使用Redis的HMap的结构，将Redis格式化为DB的数据库表
 * 主要用途为使用Redis作为Mysql的数据缓冲，提高Mysql的访问速度
 *
 */

/**
 * 函数描述： 初始化数据库名
 * 参数说明：
 *		dbName(string) 初始化数据库名字
 * 返回参数说明：
 *		nil
 */
func (r *RedisOpt) InitRedisDBName(dbName string) {

}

/**
 * 函数描述： 写入指定数据库和扫描ID的表结构
 * 参数说明：
 *		dbName(string) 初始化数据库名字
 *		scanID(string) 扫描ID
 *		dataStruct(map[string]string) 存储数据结构格式
 *		timeSecond(int) 超时时间,默认为30天
 * 返回参数说明：
 *		nil
 */
func (r *RedisOpt) SetRedisDBData(dbName, scanID string, dataStrcut map[string]string, timeSecond int) error {
	keyString := "DB@" + dbName + ":" + scanID
	if timeSecond == 0 {
		timeSecond = 3600 * 24 * 30
	}
	return r.HMSet(keyString, dataStrcut, timeSecond)
}

/**
 * 函数描述： 写入指定数据库和扫描ID,指定字段写入Redis
 * 参数说明：
 *		dbName(string) 初始化数据库名字
 *		scanID(string) 扫描ID
 *		field(string) 表字段
 *		value(interface{}) 值
 *		timeSecond(int) 超时时间,默认为30天
 * 返回参数说明：
 *		nil
 */
func (r *RedisOpt) SetRedisDBDataField(dbName, scanID, field string, value interface{}, timeSecond int) error {
	keyString := "DB@" + dbName + ":" + scanID
	if timeSecond == 0 {
		timeSecond = 3600 * 24 * 30
	}
	return r.HSet(keyString, field, value, timeSecond)
}

/**
 * 函数描述： 获取指定数据表与扫描ID的数据结构
 * 参数说明：
 *		dbName(string) 初始化数据库名字
 *		scanID(string) 扫描ID
 * 返回参数说明：
 *		dataResult(map[string]string) 请求结果
 *		status(bool) 如果err=nil,则需要判断status,status为false说明找不到key,为true，说明成功
 *		err(error) 状态
 */
func (r *RedisOpt) GetRedisDBDataAll(dbName, scanID string) (dataResult map[string]string, status bool, err error) {
	keyString := "DB@" + dbName + ":" + scanID

	return r.HGetAll(keyString)
}

/**
 * 函数描述： 获取指定数据表与扫描ID的指定字段的值
 * 参数说明：
 *		dbName(string) 初始化数据库名字
 *		scanID(string) 扫描ID
 *		field(string) 字段
 * 返回参数说明：
 *		value(string) 请求结果
 *		status(bool) 如果err=nil,则需要判断status,status为false说明找不到key,为true，说明成功
 *		err(error) 状态
 */
func (r *RedisOpt) GetRedisDBDataFiled(dbName, scanID, field string) (value string, status bool, err error) {
	keyString := "DB@" + dbName + ":" + scanID

	return r.HGet(keyString, field)
}
