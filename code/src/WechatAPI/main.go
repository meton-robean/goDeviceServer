package main

import (
	config "WechatAPI/config"
	"WechatAPI/handle"
	_ "WechatAPI/routers"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"vislog"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

func init() {

}

var (
	version         = "1.1.3.3"
	versionTime     = "20180822"
	versionFunction = "1 增加发卡功\n"
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
	if usage() {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("main:", err)
		}
	}()

	config.InitConfig()
	initLog()
	err := handle.InitServer()
	if err != nil {
		log.Error("err:", err)
		return
	}
	log.Infoln("model server starts success....")
	beego.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func initLog() {
	configOpt := config.GetConfig()

	vislogHook, err := vislog.NewVislogHook(configOpt.LogFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(vislogHook)

	level, err := log.ParseLevel(configOpt.LogLevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})
}
