package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// Logger 日志结构体
type fileLogger struct {
	Level logLevel
	filePath string
	fileName string
	fileObj *os.File
	errFileObj *os.File
	maxFileSize int64
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

// 记录日志的方法
func (f *fileLogger) log(level logLevel, format string, a ...interface{}) {
	if f.enable(level) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNum := getInfo(3)
		levelStr := parseLogLevel(level)
		if f.checkSize(*f.fileObj) {
			newFileObj, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newFileObj
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [File:%s, Func:%s, Line:%d] %s\n", now.Format("2006-01-01 01:02:03"), levelStr, fileName, funcName, lineNum, msg)
		if level >= ERROR {
			if f.checkSize(*f.errFileObj) {
				newFileObj, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newFileObj
			}
			// 如果要记录的日志 >= ERROR级别
			// 还需要在ERROR日志中再记录一遍
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [File:%s, Func:%s, Line:%d] %s\n", now.Format("2006-01-01 01:02:03"), levelStr, fileName, funcName, lineNum, msg)
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