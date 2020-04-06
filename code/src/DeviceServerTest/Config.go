package main

import (
	"flag"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Option struct {
	LogFile    string `yaml:"LogFile"`
	LogLevel   string `yaml:"LogLevel"`
	SysLogAddr string `yaml:"SysLogAddr"`

	UpdateAddr        string `yaml:"UpdateAddr"`
	HardwareAddr      string `yaml:"HardwareAddr"`
	DisplayServerAddr string `yaml:"DisplayServerAddr"`
	ScanServerAddr    string `yaml:"ScanServerAddr"`
	BodyEvalAddr      string `yaml:"BodyEvalAddr"`
	ScanDealAddr      string `yaml:"ScanDealAddr"`

	TaskServiceAddr        string `yaml:"TaskServiceAddr"`
	QrcodeServiceAddr      string `yaml:"QrcodeServiceAddr"`
	DeviceID               string `yaml:"DeviceID"`
	GGGetRFIDShellPath     string `yaml:"GGGetRFIDShellPath"`
	GGCardIDReqUserInfoURL string `yaml:"GGCardIDReqUserInfoURL"`
	DebugFlag              bool   `yaml:"DebugFlag"`

	NetWorkdTime int64 `yaml:"NetWorkdTime"`
	SaveModeFile bool  `yaml:"SaveModeFile"`

	TestTime int    `yaml:"TestTime"`
	TestFile string `yaml:"TestFile"`
	TestNum  int    `yaml:"TestNum"`
	TestJson string `yaml:"TestJson"`

	Mode                int  `yaml:"Mode"`
	EnableCloudObject   bool `yaml:"EnableCloudObject"`
	CloudeObjectSelect  int  `yaml:"CloudeObjectSelect"`
	SupportNewAlgorithm bool `yaml:"SupportNewAlgorithm"`
}

func LoadConfig() (p Option) {
	var confName string
	flag.StringVar(&confName, "f", "config.yml", "config file of monitor")
	flag.Parse()
	f, err := os.Open(confName)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("read config file err: " + err.Error())
	}
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		log.Fatal("unmarshal yaml config err: " + err.Error())
	}
	return p
}
