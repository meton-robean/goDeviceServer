package ThirdPush

import (
	"DeviceServer/Config"
	"bytes"
	"fmt"
	"net/smtp"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

/*
模块说明： 第三方推送接口，邮件推送，短信推送
*/

//PushEmail 邮件推送接口
func PushEmail(toPerson, gatewayName, gatewayID string) {
	auth := smtp.PlainAuth("", "563951092@qq.com", "psibjzctuwspbcag", "smtp.qq.com")
	to := []string{toPerson}
	nickname := "测试"
	user := "563951092@qq.com"
	subject := "网关掉线通知"
	contentType := "Content-Type: text/plain; charset=UTF-8"
	body := fmt.Sprintf("网关名: %s,网关ID:%s ,该设备掉线了，请及时处理!", gatewayName, gatewayID)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		log.Error("掉线通知失败: ", err)
		return
	}
	log.Info("掉线通知发送成功")
}

//SendPhoneMessage 发送短信
func SendPhoneMessage(phone, gatewayID string) {
	cmdStr := "python " + Config.GetConfig().EmailPythonPath + " " + phone + " " + gatewayID + " " + time.Now().Format("2006-01-02@15:04:05")
	log.Info("cmdStr:", cmdStr)
	result := execshell(cmdStr)
	log.Info("短信发送状态:", result)
}

func execshell(s string) string {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error("err:", err)
	}
	result := out.String()
	return result[:len(result)-1]
}
