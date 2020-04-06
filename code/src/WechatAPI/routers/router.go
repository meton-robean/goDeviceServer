package routers

/*
	beego框架的路由实现方式
*/
import (
	"WechatAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// beego.Router("/", &controllers.MainController{})
	//开门接口
	beego.Router("/v1/token", &controllers.WechatController{}, "get:GetToken")
	beego.Router("/v1/open-door", &controllers.WechatController{}, "get:DoorCtrlOpen")
	beego.Router("/v1/get-roominfo", &controllers.WechatController{}, "get:GetRoomInfo")
	beego.Router("/v1/setting-card-password", &controllers.WechatController{}, "get:SettingCardPassword")
	beego.Router("/v1/cancel-card-password", &controllers.WechatController{}, "get:CancleCardPassword")
	beego.Router("/v1/sync-room-info", &controllers.WechatController{}, "post:SyncAllRooms")
	beego.Router("/v1/add-room-info", &controllers.WechatController{}, "get:AddRoomInfo")
	beego.Router("/v1/del-room-info", &controllers.WechatController{}, "get:DelRoomInfo")

	//APP扫描绑定接口
	beego.Router("/v1/login", &controllers.AppController{}, "get:AppLogin")
	beego.Router("/v1/add-gateway", &controllers.AppController{}, "get:AddGateway")
	beego.Router("/v1/bind-room", &controllers.AppController{}, "get:BindDeviceRoom")
	beego.Router("/v1/get-room-info", &controllers.AppController{}, "get:GetAllRoomInfos")

	//模拟推送接收接口
	beego.Router("/test/token", &controllers.TestPushServerController{}, "get:TestToken")
	beego.Router("/test/push", &controllers.TestPushServerController{}, "post:TestPush")

	//接收设备服务的状态上报接口
	beego.Router("/report/door-ctrl-rsp", &controllers.DevStatusController{}, "get:DoorCtrlRsp")
	beego.Router("/report/dev-setting-password-status", &controllers.DevStatusController{}, "get:SettingCardlRsp")
	beego.Router("/report/dev-cancel-password-status", &controllers.DevStatusController{}, "get:CancelCardlRsp")
	beego.Router("/report/card-openlock-record", &controllers.DevStatusController{}, "get:CardDoorOpenlRsp")
}
