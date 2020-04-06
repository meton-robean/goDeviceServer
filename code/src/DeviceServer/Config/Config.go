package Config

import (
	"flag"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var configOpt Option
var onceDataOpt sync.Once

//InitConfig 初始化配置文件，只能加载一次
func InitConfig() {
	onceDataOpt.Do(func() {
		configOpt = loadConfig()
	})
}

type Option struct {
	Addr     string `yaml:"Addr"`
	Database string `yaml:"Database"`

	LogFile    string `yaml:"LogFile"`
	LogLevel   string `yaml:"LogLevel"`
	SysLogAddr string `yaml:"SysLogAddr"`

	ReportHTTPAddr string `yaml:"ReportHTTPAddr"`
	HTTPServer     string `yaml:"HTTPServer"`

	RedisAddr      string `yaml:"RedisAddr"`
	RedisPwd       string `yaml:"RedisPwd"`
	RedisTimeOut   int    `yaml:"RedisTimeOut"`
	RedisServerNum int    `yaml:"RedisServerNum"`

	EmailPythonPath string `yaml:"EmailPythonPath"`
}

func loadConfig() (p Option) {
	var confName string
	flag.StringVar(&confName, "f", "config.yml", "config file")
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

//GetConfig get global config
func GetConfig() *Option {
	return &configOpt
}
