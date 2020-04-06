package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	// "sync"

	LogOpt "LogOpt"
	"gotcp"
	"syscall"
	"vislog"

	"runtime"

	log "github.com/Sirupsen/logrus"
)

var gConfigOpt Option
var ScanSrv *gotcp.Server
var DisplaySrv *gotcp.Server
var UpdateSrv *gotcp.Server
var BodyEvalSrv *gotcp.Server
var ScanDealSrv *gotcp.Server

//业务打点
var gLogJog *LogOpt.LogOpt

var (
	serverName  string
	version     string
	versionTime string
)

func init() {
	version = "1.0.9"
	versionTime = "20170831"
	serverName = "DeviceServerTest"
}

func usage() bool {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v" ||
		args[1] == "version") {
		fmt.Println("version:", version)
		fmt.Println("build time:", versionTime)
		fmt.Println("funciont:", "bodyEval,taskServer remove upload file")
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

func testFunc() {
	personData, err := ioutil.ReadFile(os.Args[4])
	if err != nil {
		log.Error("read json error:", err)
		return
	}
	//fmt.Println(string(personData))
	lineStr := strings.Split(string(personData), "\n")
	fmt.Println("len:", len(lineStr))
	for i, val := range lineStr {
		result := strings.Split(string(val), " ")
		for ii, vali := range result {
			fmt.Println(i*(ii+1), ",", vali)
		}
	}
}

func main() {
	// testFunc()
	// return
	// if usage() {
	// 	return
	// }

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	stop()
	log.Info("控制服务退出....")
}

func start() {
	gConfigOpt = LoadConfig()

	if gConfigOpt.Mode == 0 {
		gConfigOpt.Mode = 2
	}
	log.Info("File:", gConfigOpt.TestFile, ",Json:", gConfigOpt.TestJson)
	log.Info("DeviceID111:", gConfigOpt.DeviceID)

	initLog(gConfigOpt.LogFile, gConfigOpt.LogLevel, gConfigOpt.SysLogAddr)
	log.Info("###########版本:", version)
	log.Info("控制服务开始启动....")

	go qrcodeClient.ConnectQrcodeServer(gConfigOpt.QrcodeServiceAddr)

}

func stop() {

}

func initLog(logfile string, loglevel string, syslogAddr string) {
	if logfile == "" {
		logfile = "DeviceServerTest.log"
	}
	hook, err := vislog.NewVislogHook(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(hook)

	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})
}
