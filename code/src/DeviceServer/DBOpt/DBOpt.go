package DBOpt

import (
	"DeviceServer/Common"
	"DeviceServer/ThirdPush"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type DBOpt struct {
	BaseDB
}

var dataOpt *DBOpt
var onceDataOpt sync.Once

//GetDataOpt .获取数据平台对象
func GetDataOpt() *DBOpt {
	onceDataOpt.Do(func() {
		dataOpt = new(DBOpt)
	})
	return dataOpt
}

//GetDeviceIDList 通过房间号与用户ID获取设备ＩＤ
func (opt *DBOpt) GetDeviceIDList(gatewayID string) (devListMap map[string]bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return devListMap, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select device_id from t_device_info a,t_gateway_info b where a.gw_id=b.id and b.gateway_id=?"
	rows, err := conn.Query(sqlString, gatewayID)
	if err != nil {
		log.Error("err:", err)
		return devListMap, err
	}
	defer rows.Close()

	devListMap = make(map[string]bool)
	var deviceID string
	for rows.Next() {
		err = rows.Scan(&deviceID)
		if err != nil {
			log.Error("err:", err)
			return devListMap, err
		}
		devListMap[deviceID] = true
	}
	return devListMap, err
}

//SetGatwayOnline 设置网关在线
func (opt *DBOpt) SetGatwayOnline(gatewayID string) error {
	if len(gatewayID) == 0 {
		return nil
	}
	log.Debug("SetGatwayOnline:", gatewayID)
	return opt.setGatewayStatus(gatewayID, 1)
}

//SetGatwayOffline 设置网关下线
func (opt *DBOpt) SetGatwayOffline(gatewayID string) error {
	defer func() {
		if e := recover(); e != nil {
			log.Error("HandleMsg:", e)
			return
		}
	}()

	log.Debug("SetGatwayOffline1:", gatewayID)
	err := opt.setGatewayStatus(gatewayID, 0)
	if err != nil {
		log.Error("err:", err)
	}

	timeNow := time.Now().Unix()
	if (timeNow - Common.ServerStarTime) < 60 {
		log.Info("重启服务导致的掉线，不需要通知，１分钟以内")
		return nil
	}
	log.Debug("开始推送掉线通知")
	// email, err := opt.GetAdminEmail()
	// if err != nil {
	// 	log.Error("err:", err)
	// } else {
	// 	log.Debug("email:", email)
	// 	ThirdPush.PushEmail(email, "节点网关", gatewayID)
	// }

	phone, err := opt.GetManagerPhone(gatewayID)
	if err != nil {
		log.Error("err:", err)
	} else {
		log.Debug("phone:", phone)
		if len(phone) < 10 {
			log.Error("错误的手机号:", phone)
			return nil
		}
		ThirdPush.SendPhoneMessage(phone, gatewayID)
	}
	log.Debug("推送掉线通知完毕")
	return err
}

func (opt *DBOpt) setGatewayStatus(gatewayID string, status int) (err error) {
	sqlString := "update t_gateway_info set status=? where gateway_id=?"
	err = opt.exec(nil, sqlString, status, gatewayID)
	if err != nil {
		log.Error("err:", err)
	}
	return
}

//UpdateDeviceBarray 更新电量
func (opt *DBOpt) UpdateDeviceBarray(deviceID string, barray float64) (err error) {
	sqlString := "update t_device_info set barry=? where device_id=?"
	err = opt.exec(nil, sqlString, barray, deviceID)
	if err != nil {
		log.Error("err:", err)
	}
	return
}

//GetAdminEmail 获取超级管理员邮箱
func (opt *DBOpt) GetAdminEmail() (email string, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return email, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select email from t_user_info where user_account='admin'"
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return email, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&email)
		if err != nil {
			log.Error("err:", err)
			return email, err
		}
	}
	return email, err
}

//GetManagerPhone 获取网关对应的管理员电话
func (opt *DBOpt) GetManagerPhone(gatewayID string) (phone string, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return phone, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select user_phone from t_user_info a " +
		"inner join t_gateway_info b on a.id=b.user_id and b.gateway_id=?"
	rows, err := conn.Query(sqlString, gatewayID)
	if err != nil {
		log.Error("err:", err)
		return phone, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&phone)
		if err != nil {
			log.Error("err:", err)
			return phone, err
		}
	}
	return phone, err
}
