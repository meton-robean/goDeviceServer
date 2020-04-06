package controllers

/*
	该模块主要用来接收微信端的请求设置
*/

import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"WechatAPI/config"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"crypto/sha256"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

//WechatController .
type WechatController struct {
	beego.Controller
}

//GetToken 通过APPID+Secrete生成token
func (c *WechatController) GetToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//判断appid,secret是否存在,权限判断
	status, err := DBOpt.GetDataOpt().CheckAppIDSecret(appid, secret)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Debug("appid,secrete不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	//组织应答报文,生成token的方法
	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	//将token保存到Redis缓存
	common.RedisTokenOpt.Set(token, 1, config.GetConfig().RedisTokenTimeOut)

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTokenTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

//GetRoomInfo 通过设备ID获取房间信息
func (c *WechatController) GetRoomInfo() {
	DeviceID := c.GetString("deviceid")
	Token := c.GetString("token")
	log.Info("DeviceID=", DeviceID, ",Token=", Token)
	if DeviceID == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	log.Debug("token:", Token)
	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(Token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	//通过设备ID获取房间信息与对应的酒店appid
	roomnu, appid, err := DBOpt.GetDataOpt().GetRoomInfo(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	data := make(map[string]interface{})
	if len(roomnu) != 0 {
		data["roomnu"] = roomnu
		data["appid"] = appid
		data["code"] = 0
		c.Data["json"] = data
	} else {
		c.Data["json"] = common.GetErrCodeJSON(10005)
	}
	c.ServeJSON()
	return
}

//DoorCtrlOpen 开门
func (c *WechatController) DoorCtrlOpen() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	method, err := c.GetInt("method")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")
	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//门禁开门
	if method == 3 {

	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, method)
	if !status {
		return
	}

	//向设备服务请求开门
	httpServerIP := fmt.Sprintf("http://%s/dev-ctrl?gwid=%s&deviceid=%s&requestid=%s", serverIP, gatewayID, DeviceID, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

	if method == 1 {
		//保存开门信息
		err = DBOpt.GetDataOpt().WechatOpenMethod(DeviceID)
	} else {
		err = DBOpt.GetDataOpt().WechatOpenMethod(DeviceID)
	}
	if err != nil {
		log.Error("err:", err)
	}

	return
}

//SettingCardPassword 发卡、密码
func (c *WechatController) SettingCardPassword() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	keyvalue := c.GetString("keyvalue")
	keytype, err := c.GetInt("keytype")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	expireDate, err := c.GetInt64("expire-date")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")
	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" ||
		keyvalue == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//向设备服务请求发卡
	httpServerIP := fmt.Sprintf("http://%s/setting-card-password?gwid=%s&deviceid=%s&keyvalue=%s&keytype=%d&expire-date=%d&requestid=%s",
		serverIP, gatewayID, DeviceID, keyvalue, keytype, expireDate, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

//CancleCardPassword 发卡、密码
func (c *WechatController) CancleCardPassword() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	keyvalue := c.GetString("keyvalue")
	keytype, err := c.GetInt("keytype")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")

	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" ||
		keyvalue == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//向设备服务请求发卡
	httpServerIP := fmt.Sprintf("http://%s/cancel-card-password?gwid=%s&deviceid=%s&keyvalue=%s&keytype=%d&requestid=%s",
		serverIP, gatewayID, DeviceID, keyvalue, keytype, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

func (c *WechatController) checkAppidUser(roomnu, appid, token string, method int) (serverIP, gatewayID, DeviceID string, gwOnline bool) {
	var devOnline bool

	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return "", "", "", false
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return "", "", "", false
	}

	if method == 3 {
		gatewayID, DeviceID, gwOnline, err = DBOpt.GetDataOpt().GetDoorCardInfo(roomnu, appid)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
	} else {
		//通过房间号与酒店appid获取设备id信息
		DeviceID, err = DBOpt.GetDataOpt().GetDeviceID(roomnu, appid)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
		if len(DeviceID) == 0 {
			log.Error("房间数据不存在:", roomnu, ",userid:", appid)
			c.Data["json"] = common.GetErrCodeJSON(10004)
			c.ServeJSON()
			return "", "", "", false
		}
		log.Debug("DeviceID:", DeviceID)

		//通过设备ID获取网关ID与在线状态
		gatewayID, gwOnline, devOnline, err = DBOpt.GetDataOpt().CheckGatewayOnline(DeviceID)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
	}

	//用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
	dataBuf, isExist, err := common.RedisServerListOpt.Get(gatewayID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return "", "", "", false
	}
	if !isExist {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10008)
		c.ServeJSON()
		return "", "", "", false
	}
	serverIP = string(dataBuf)

	log.Debug("gatewayID:", gatewayID, ",gwOnline:", gwOnline)

	//目前网关心跳只是一人空包，没有网关ＩＤ，无法做到网关是否线
	devOnline = true
	var errcode int
	if gwOnline {
		//网关在线
		errcode = 0
		if devOnline {
			errcode = 0
		} else {
			//设备不在线
			errcode = 10009
		}
	} else {
		//网关不在线
		errcode = 10008
	}
	c.Data["json"] = common.GetErrCodeJSON(errcode)
	c.ServeJSON()
	if errcode != 0 {
		log.Info("网关或者设备不在线：gw=", gatewayID, ",deviceID=", DeviceID)
		return "", "", "", false
	}
	return serverIP, gatewayID, DeviceID, true
}

//SyncAllRooms 同步所有房间
func (c *WechatController) SyncAllRooms() {
	dataMap := make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataMap); err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	Token := dataMap["token"]
	if Token == nil {
		log.Error("token字段不存在")
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	log.Debug("token:", Token)
	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(Token.(string))
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	appid := dataMap["appid"]
	if appid == nil {
		log.Error("appid字段不存在")
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	userid, err := DBOpt.GetDataOpt().GetUserID(appid.(string))
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	dataInfo := dataMap["data"].([]interface{})
	if dataInfo == nil {
		log.Error("data字段不存在")
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	dataRoomInfos := make([]common.RoomInfo, len(dataInfo))
	for k, v := range dataInfo {
		vMap := v.(map[string]interface{})
		dataRoomInfos[k].RName = vMap["rname"].(string)
		dataRoomInfos[k].Roomnu = vMap["roomnu"].(string)
	}
	log.Debug("roomInfo:", dataRoomInfos)

	if err := DBOpt.GetDataOpt().SyncRoomInfos(dataRoomInfos, userid); err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}

//AddRoomInfo 添加一个房间
func (c *WechatController) AddRoomInfo() {
	appid := c.GetString("appid")
	token := c.GetString("token")
	rname := c.GetString("rname")
	roomnu := c.GetString("roomnu")

	log.Debug("roomnu:", roomnu, ",rname:", rname)
	if appid == "" || token == "" || rname == "" ||
		roomnu == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}
	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	userid, err := DBOpt.GetDataOpt().GetUserID(appid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if userid < 0 {
		log.Error("用户不存在")
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	dataRoomInfos := make([]common.RoomInfo, 1)
	dataRoomInfos[0].RName = rname
	dataRoomInfos[0].Roomnu = roomnu
	if err := DBOpt.GetDataOpt().SyncRoomInfos(dataRoomInfos, userid); err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}

//DelRoomInfo 添加一个房间
func (c *WechatController) DelRoomInfo() {
	appid := c.GetString("appid")
	token := c.GetString("token")
	roomnu := c.GetString("roomnu")

	if appid == "" || token == "" || roomnu == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	userid, err := DBOpt.GetDataOpt().GetUserID(appid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if userid < 0 {
		log.Error("用户不存在")
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	if err := DBOpt.GetDataOpt().DelRoomInfo(roomnu, userid); err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}
