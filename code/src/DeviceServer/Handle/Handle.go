package Handle

import (
	"DeviceServer/Common"
	"bytes"
	"encoding/json"
	"errors"
	"gotcp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//gatewayID,conn
var ConnInfo map[string]*gotcp.Conn = make(map[string]*gotcp.Conn)

type CallBack struct {
}

func (cb *CallBack) Close() {
}

func (cb *CallBack) HandleMsg(conn *gotcp.Conn, MsgBody []byte) error {
	defer func() {
		if e := recover(); e != nil {
			log.Error("HandleMsg:", e)
			return
		}
	}()

	if len(MsgBody) < 10 {
		baseSendMsg(conn, []byte("abc"))
		//log.Debug("错误包:", string(MsgBody))
		return nil
	}

	if len(MsgBody) < len(Common.DefaultHead)+10 {
		baseSendMsg(conn, []byte("abc"))
		log.Error("wrong pack")
		return errors.New("wrong pack")
	}

	log.Debug("接收到消息msg:", string(MsgBody))

	if !strings.Contains(string(MsgBody), Common.DefaultHead) {
		log.Debug("head err")
		return errors.New("head err")
	}
	jsonData := MsgBody[len(Common.DefaultHead)+5:]
	data := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("err: ", err)
		return nil
	}
	val, exist := data["cmd"]
	if !exist {
		log.Error("cmd not exist:", string(MsgBody))
		return nil
	}

	//log.Info("data:", data)
	cmd := val.(string)
	switch cmd {
	case "gw_register": //心跳
		gatewayRegisterRsp(conn, cmd, data)
	case "d2s_status": //开门返回来的状态
		doorCtrlDealRsp(conn, cmd, data)
	case "d2s_request_devices": //网关请求所有节点信息
		requestDeviceListRsp(conn, cmd, data)
	case "d2s_battery": //上报电量
		doorReportBarryRsp(conn, cmd, data)
	case "dev_single_password_setting":
		devSettingPasswordRsp(conn, cmd, data)
	case "dev_single_password_cancel":
		devCancelPasswordRsp(conn, cmd, data)
	case "openlock_record_return":
		cardOpenLockRecord(conn, cmd, data)
	default:
		baseSendMsg(conn, []byte("abc"))
		log.Error("cmd invalid:", cmd)
	}

	return nil
}

//响应网关程序
func ackGateway(conn *gotcp.Conn, dataMap map[string]interface{}) {
	dataBuf, err := json.Marshal(dataMap)
	if err != nil {
		log.Error("err:", err)
		return
	}

	//获取打包的协议
	protoclBuf := getPackage(dataBuf)
	log.Debug("ackmsg:", string(protoclBuf))
	baseSendMsg(conn, protoclBuf)
}

func baseSendMsg(conn *gotcp.Conn, msg []byte) {
	//log.Debug("send msg:", string(msg))
	conn.SendChan <- msg
}

/*
获取网关通信协议的组装包格式
*/
func getPackage(msg []byte) []byte {
	var crc int
	len := len(msg)

	var dataBuff bytes.Buffer
	//包头
	dataBuff.WriteString(Common.DefaultHead)
	//状态
	dataBuff.WriteByte(0x23)
	//包体长度
	dataBuff.WriteByte(byte(len >> 24))
	dataBuff.WriteByte(byte(len >> 16))
	dataBuff.WriteByte(byte(len >> 8))
	dataBuff.WriteByte(byte(len))
	//包内容
	dataBuff.Write(msg)
	//crc，目前不需要用到
	dataBuff.WriteByte(byte(crc >> 24))
	dataBuff.WriteByte(byte(crc >> 16))
	dataBuff.WriteByte(byte(crc >> 8))
	dataBuff.WriteByte(byte(crc))
	dataBuff.WriteString("\n")

	return dataBuff.Bytes()
}
