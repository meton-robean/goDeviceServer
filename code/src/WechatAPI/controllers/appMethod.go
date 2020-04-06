package controllers

/*
	该模块主要用来接收APP的请求，包括登陆，添加网关，设备绑定
*/
import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

//AppController .
type AppController struct {
	beego.Controller
}

//AppLogin APP登陆
func (c *AppController) AppLogin() {
	username := c.GetString("username")
	pwd := c.GetString("pwd")

	log.Debug("username:", username)
	log.Debug("pwd:", pwd)
	userInfo, err := DBOpt.GetDataOpt().GetUserPwd(username)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if userInfo.UserType != 3 {
		log.Error("只有管理用户有权限登陆")
		c.Data["json"] = common.GetErrCodeJSON(10014)
		c.ServeJSON()
		return
	}

	// m := md5.New()
	// m.Write([]byte(pwd))
	// pwdMd5 := hex.EncodeToString(m.Sum(nil))

	if userInfo.UserPwd != pwd {
		log.Error("登陆密码不匹配:name=", username, ",pwd=", pwd, ",pwdMd5=", userInfo.UserPwd)
		c.Data["json"] = common.GetErrCodeJSON(10010)
		c.ServeJSON()
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["appid"] = userInfo.AppID
	dataMap["secret"] = userInfo.Secret
	dataMap["userid"] = userInfo.UserID
	dataMap["code"] = 0
	c.Data["json"] = dataMap
	c.ServeJSON()
}

//AddGateway 添加网关
func (c *AppController) AddGateway() {
	gwid := c.GetString("gwid")
	gwid = strings.ToUpper(gwid)
	gwname := c.GetString("gwname")
	token := c.GetString("token")
	userid, err := c.GetInt("userid")
	if err != nil {
		log.Error("err:", err)
	}

	if gwid == "" || gwname == "" || userid == 0 || token == "" {
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

	//检查用户ID的合法性
	status, err = DBOpt.GetDataOpt().CheckUserID(userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("用户ID不存在:", userid)
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	//添加网关到数据库
	err = DBOpt.GetDataOpt().AddGatewayInfo(userid, gwid, gwname)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}

//BindDeviceRoom 绑定房间与设备
func (c *AppController) BindDeviceRoom() {
	gwid := c.GetString("gwid")
	gwid = strings.ToUpper(gwid)
	deviceid := c.GetString("deviceid")
	deviceid = strings.ToUpper(deviceid)
	roomnu := c.GetString("roomnu")
	token := c.GetString("token")
	userid, err := c.GetInt("userid")
	if err != nil {
		log.Error("err:", err)
	}

	if deviceid == "" || gwid == "" || roomnu == "" || userid == 0 || token == "" {
		log.Error("参数出错:", "gwid=", gwid, "deviceid=", deviceid, ",roomnu=", roomnu, ",userid=", userid, ",token=", token)
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

	//检查用户ID的合法性
	status, err = DBOpt.GetDataOpt().CheckUserID(userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("用户ID不存在:", userid)
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	gid, err := DBOpt.GetDataOpt().CheckGatewayExist(gwid, userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if gid == 0 {
		log.Error("用户下该网关不存在:", gwid, ",userid:", userid)
		c.Data["json"] = common.GetErrCodeJSON(10013)
		c.ServeJSON()
	}

	//检查该用户ID下的设备ID与房间号的绑定情况,是否已经被绑定
	status, err = DBOpt.GetDataOpt().CheckRoomBeenBind(roomnu, userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if status {
		log.Info("该用户的该房间号已经被绑定过了:", userid, ",", roomnu)
		c.Data["json"] = common.GetErrCodeJSON(10015)
		c.ServeJSON()
		return
	}

	//检查该用户ID下的设备ID与房间号的绑定情况,是否已经被绑定
	status, err = DBOpt.GetDataOpt().CheckDeviceBeenBind(deviceid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if status {
		log.Info("设备ID已经被绑定过了:", deviceid, ",", roomnu)
		c.Data["json"] = common.GetErrCodeJSON(10011)
		c.ServeJSON()
		return
	}

	//添加设备的绑定信息
	err = DBOpt.GetDataOpt().AddDeviceAndRoomBind(userid, gid, deviceid, roomnu)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
	return
}

//SyncRoomInfo 同步所有的房间
func (c *AppController) GetAllRoomInfos() {
	username := c.GetString("username")
	token := c.GetString("token")

	if username == "" || token == "" {
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

	dataRoomInfo, err := DBOpt.GetDataOpt().GetAllRoomInfos(username)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["username"] = username
	dataMap["data"] = dataRoomInfo
	c.Data["json"] = dataMap
	c.ServeJSON()
	return
}
