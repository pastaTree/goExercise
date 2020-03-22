package main

import mylogger "myLogger"

// 声明一个全局的接口变量
var log mylogger.Logger

func main() {
	// console log test
	log = mylogger.NewConsoleLog("debug")  // 终端log实例
	id := 2
	name := "George"
	log.Debug("这是一条Debug日志, id: %d, name: %s", id, name)
	log.Trace("这是一条Trace日志")
	log.Info("这是一条Info日志")
	log.Warning("这是一条Warning日志")
	log.Error("这是一条Error日志, id: %d, name: %s", id, name)
	log.Fatal("这是一条Fatal日志")

	// file log test
	log = mylogger.NewFileLog("debug", "./", "file.log", 10 * 1024 * 1024) // 文件log实例
	for {
		log.Debug("这是一条Debug日志, id: %d, name: %s", id, name)
		log.Trace("这是一条Trace日志")
		log.Info("这是一条Info日志")
		log.Warning("这是一条Warning日志")
		log.Error("这是一条Error日志, id: %d, name: %s", id, name)
		log.Fatal("这是一条Fatal日志")
	}
}
