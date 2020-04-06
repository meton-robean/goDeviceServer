package main

import (
	"DeviceServer/Common"
	"DeviceServer/Config"
	"DeviceServer/DBOpt"
	"DeviceServer/HTTPServer"
	"DeviceServer/Handle"
	"DeviceServer/ThirdPush"
	"fmt"
	"gotcp"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vislog"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

var Srv *gotcp.Server

var (
	version         = "1.1.3.3"
	versionTime     = "20180827"
	versionFunction = "1 增加发卡功,重新１分钟内掉线的不需要短信通知\n"
)

func usage() bool {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v" ||
		args[1] == "version") {
		fmt.Println("version:", version)
		fmt.Println("build time:", versionTime)
		fmt.Println("function:", versionFunction)
		return true
	} else if len(args) == 2 && (args[1] == "--help" || args[1] == "-h" ||
		args[1] == "help") {
		fmt.Println("1 --version -v version ")
		fmt.Println("2 kill -USR1 pid  open debug info")
		fmt.Println("3 kill -USR2 pid  close debug info")
		return true
	}

	return false
}

func main() {
	Common.ServerStarTime = time.Now().Unix()
	//ThirdPush.PushEmail("wenzhongjian@suanier.com", "公司", "aaaaaaa")
	if usage() {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("main:", err)
		}
	}()

	start()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	stop()
	log.Info("DeviceServer server quit")
}

func start() {
	Config.InitConfig()
	config := Config.GetConfig()
	ThirdPush.SendPhoneMessage("13723450181", "服务重启")
	//初始化日志
	initLog(config.LogFile, config.LogLevel, config.SysLogAddr)

	//初始化公共组件
	err := Common.InitCommon()
	if err != nil {
		log.Error("err:", err)
		return
	}

	//初始化数据库
	DBOpt.GetDataOpt().InitDatabase(config.Database)

	log.Info("DeviceServer server is starting.....version:", version, ",port:", config.Addr)
	//初始化网关监听服务
	Srv = gotcp.NewServer(&Handle.CallBack{})
	go Srv.StartServer(config.Addr, "GateWayServer", DBOpt.GetDataOpt().SetGatwayOffline)

	//初始化HTTP服务，接收WechatAPI的消息
	go HTTPServer.HTTPInit(config.HTTPServer)

}

func stop() {
}

func initLog(logfile string, loglevel string, syslogAddr string) {
	if logfile == "" {
		logfile = "DevStatusServer.log"
	}
	hook, err := vislog.NewVislogHook(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(hook)

	syshook, err := logrus_syslog.NewSyslogHook("udp", syslogAddr, syslog.LOG_DEBUG, os.Args[0])
	if err == nil {
		log.AddHook(syshook)
	}

	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})

}
