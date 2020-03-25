package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

var (
	MaxSize = 50000
)

// Logger 日志结构体
type fileLogger struct {
	Level logLevel
	filePath string
	fileName string
	fileObj *os.File
	errFileObj *os.File
	maxFileSize int64
	logChan chan *logMsg
}

type logMsg struct {
	Level logLevel
	msg string
	fileName string
	funcName string
	timestamp string
	line int
}

// NewFileLog 构造函数
func NewFileLog(levelStr, fp, fn string, maxFileSize int64) *fileLogger {
	level, err := parseLogString(levelStr)
	if err != nil {
		panic(err)
	}
	f := &fileLogger{
		Level: level,
		filePath: fp,
		fileName: fn,
		maxFileSize: maxFileSize,
		logChan: make(chan *logMsg, MaxSize),
	}
	err = f.initFile()
	if err != nil {
		panic(err)
	}
	return f
}

func (f *fileLogger) initFile() error {
	fileFullName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fileFullName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file failed, error: %v/n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fileFullName+".err", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open error log file failed, error: %v/n", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	// 开启5个 goroutine 去写日志
	for i := 0; i < 5; i++ {
		go f.writeLogBackground()
	}
	return nil
}

func (f *fileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}

// 判断是否记录日志
func (f *fileLogger) enable(logLevel logLevel) bool {
	return f.Level <= logLevel
}

// 根据文件大小,判断文件是否需要分割
func (f *fileLogger) checkSize(file os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed, error: %v\n", err)
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}

// 根据文件日期,判断文件是否需要切割
//func (f *fileLogger) checkHour(file os.File) bool

// 切割文件
func (f *fileLogger) splitFile(file *os.File) (*os.File, error) {
	// 需要切割日志文件
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed, error: %v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name()) // 拿到当前日志文件完整路径
	nowStr := time.Now().Format("20200322231350")
	newLogName := fmt.Sprintf("%s.back%s", logName, nowStr) // file.log -> file.log.bk20200322
	// 1. 关闭当前日志文件
	file.Close()
	// 2. 备份
	os.Rename(logName, newLogName)
	// 3. 打开一个新的文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file failed, error: %v/n", err)
		return nil, err
	}
	// 4. 将打开的新日志文件对象,赋值给fileObj
	return fileObj, nil
}

func (f *fileLogger) writeLogBackground() {
	for {
		// TODO: 此处存在 Race Condition!
		// TODO: 多个 Goroutine 向同一文件添加内容, 假如文件超过最大限制,
		// TODO: 则在关闭当前 log 创建新 log 的瞬间, 会影响其他正在工作的
		// TODO: Goroutine.
		// 检查是否需要切割
		if f.checkSize(*f.fileObj) {
			newFileObj, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newFileObj
		}

		select {
		case logTmp := <- f.logChan:
			logInfo := fmt.Sprintf("[%s] [%s] [File:%s, Func:%s, Line:%d] %s\n", logTmp.timestamp, parseLogLevel(logTmp.Level), logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)
			fmt.Fprintf(f.fileObj, logInfo)
			if logTmp.Level >= ERROR {
				if f.checkSize(*f.errFileObj) {
					newFileObj, err := f.splitFile(f.errFileObj)
					if err != nil {
						return
					}
					f.errFileObj = newFileObj
				}
				// 如果要记录的日志 >= ERROR级别
				// 还需要在ERROR日志中再记录一遍
				fmt.Fprintf(f.errFileObj, logInfo)
			}
		default:
			// 取不出日志就休息半秒钟
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// 记录日志的方法
func (f *fileLogger) log(level logLevel, format string, a ...interface{}) {
	if f.enable(level) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNum := getInfo(3)
		// 先把日志发送到通道中
		// 1. 造一个logMsg对象
		logTmp := &logMsg{
			Level:     level,
			msg:       msg,
			fileName:  fileName,
			funcName:  funcName,
			timestamp: now.Format("2006-01-01 01:02:03"),
			line:      lineNum,
		}
		// 2. 通过 IO 复用判断 channel 是否可以写入
		select {
		case f.logChan <- logTmp:
		default: // 把日志丢弃保证不出现阻塞
		}
	}
}

func (f *fileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

func (f *fileLogger) Trace(format string, a ...interface{}) {
	f.log(TRACE, format, a...)
}

func (f *fileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

func (f *fileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

func (f *fileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

func (f *fileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}