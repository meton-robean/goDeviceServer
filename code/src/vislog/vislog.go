package vislog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	Mb             int = 1024000
	DefaultMaxSize     = 50
	_depth             = 7
)

type VisFormatter struct {
	TimestampFormat string
}

var strScanID = "none"

func (f *VisFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = logrus.DefaultTimestampFormat
	}

	_, file, line, ok := runtime.Caller(_depth)
	var strFile string
	if !ok {
		strFile = "???:?"
	} else {
		strFile = fmt.Sprint(filepath.Base(file), ":", line)
	}
	strTime := time.Now().Format("2006-01-02 15:04:05")
	serialized := fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"filename\":\"%s\"}",
		strTime, entry.Level.String(), entry.Message, strFile)

	return append([]byte(serialized), '\n'), nil
}

type FileWriter struct {
	lock     sync.Mutex
	fileName string
	fp       *os.File
}

func NewFileWriter(fileName string) (*FileWriter, error) {
	w := &FileWriter{fileName: fileName}
	var err error
	w.fp, err = os.OpenFile(w.fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	return w, err
}

func (w *FileWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	//var fileInfo os.FileInfo
	//fileInfo, err = os.Stat(w.fileName)
	_, err = os.Stat(w.fileName)
	if err == nil {
		strTime := time.Now().AddDate(0, 0, -1).Format("20060102")
		err = os.Rename(w.fileName, w.fileName+"."+strTime)
		//err = os.Rename(w.fileName, w.fileName+"."+fileInfo.ModTime().Format("2006-01-02-15-04-05"))
		if err != nil {
			return
		}
	}
	w.fp, err = os.OpenFile(w.fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	return
}

type VislogHook struct {
	fileWriter *FileWriter
	maxSize    int
}

func SetScanID(scanID string) {
	strScanID = scanID
}

func NewVislogHook(fileName string, size ...int) (*VislogHook, error) {
	// strTime := time.Now().Format("20060102_150405")
	// fileName += "."
	// fileName += strTime
	fw, err := NewFileWriter(fileName)
	if err != nil {
		return nil, err
	}
	maxSize := DefaultMaxSize
	if len(size) > 0 {
		maxSize = size[0]
	}
	vislogHook := &VislogHook{
		fileWriter: fw,
		maxSize:    maxSize,
	}
	return vislogHook, nil
}

func (hook *VislogHook) Fire(entry *logrus.Entry) error {
	formatter := VisFormatter{}
	line, err := formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	if hook.fileWriter != nil {
		fileInfo, err := hook.fileWriter.fp.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get file info, %v", err)
			return err
		}
		if fileInfo.Size()+int64(len(line)) > int64(hook.maxSize*Mb) || fileInfo.ModTime().Day() != time.Now().Day() {
			if err := hook.fileWriter.Rotate(); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to rotate, %v", err)
				return err
			}
		}
		if _, err = hook.fileWriter.fp.Write(line); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write data, %v", err)
			return err
		}
	}

	return nil

}

func (hook *VislogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
