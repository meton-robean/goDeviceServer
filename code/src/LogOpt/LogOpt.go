package LogOpt

import (
	"fmt"
	"os"
	"time"
)

type LogOpt struct {
	mFileName string
	mDate     string
	mFile     *os.File
}

//与ReleaseLogOpt对应关系，需要调用ReleaseLogOpt释放对应的句柄
func CreateLogOpt(strFileName string) (*LogOpt, error) {
	logOpt := new(LogOpt)

	err := logOpt.InitLogOpt(strFileName)
	if err != nil {
		return nil, err
	}

	return logOpt, nil
}

func (opt *LogOpt) check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
* 判断文件是否存在 存在返回 true 不存在返回false
 */
func (opt *LogOpt) checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (opt *LogOpt) InitLogOpt(strFileName string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Panic InitLogOpt:", r)
		}
	}()

	opt.mFileName = strFileName
	err := opt.OpenFile()
	return err
}

func (opt *LogOpt) OpenFile() error {
	var err error
	if opt.checkFileIsExist(opt.mFileName) { //如果文件存在
		opt.mFile, err = os.OpenFile(opt.mFileName, os.O_WRONLY, 0644) //打开文件
		if err != nil {
			return err
		}

		_, err := opt.mFile.Seek(0, os.SEEK_END)
		if err != nil {
			return err
		}
		//fmt.Println("文件已经存在，打开:", n)
	} else {
		opt.mFile, err = os.Create(opt.mFileName) //创建文件
		if err != nil {
			return err
		}
		fmt.Println("文件不存在，创建:")
	}

	//opt.mDate = time.Now().Format("2010-10-10")
	return nil
}

func (opt *LogOpt) PrintMsg(scanID, deviceID, dotType, serverName, jobMsg string) error {
	// nowTime := time.Now().Format("2010-10-10")
	// if nowTime != opt.mDate {
	// 	os.Remove(opt.mFileName)
	// 	opt.OpenFile()
	// }
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	serialized := fmt.Sprintf("{\"time\":\"%s\",\"scanID\":\"%s\",\"deviceID\":\"%s\",\"dotType\":\"%s\",\"serverName\":\"%s\",\"jobMsg\":\"%s\"}",
		timeStr, scanID, deviceID, dotType, serverName, jobMsg)
	_, err := opt.mFile.WriteString(serialized)
	_, err = opt.mFile.WriteString("\n")
	if err != nil {
		opt.OpenFile()
	}
	return err
}

func (opt *LogOpt) PrintMsgTime(scanID string, useTime int64) error {
	serialized := fmt.Sprintf("scanID:%s,time=%d", scanID, useTime)
	_, err := opt.mFile.WriteString(serialized)
	_, err = opt.mFile.WriteString("\n")
	if err != nil {
		opt.OpenFile()
	}
	return err
}

func ReleaseLogOpt(opt *LogOpt) {
	if opt.mFile != nil {
		opt.mFile.Close()
	}
}
