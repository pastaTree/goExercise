package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

// ini配置文件解析器

// MySQL config 配置结构体
type MySQLConfig struct {
	Address string `ini:"address"`
	Port int `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

// Redis config 配置结构体
type RedisConfig struct {
	Host string `ini:"HOST"`
	Port int `ini:"port"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

// Config 配置结构体
type Config struct {
	MySQLConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

func loadIni(fileName string, data interface{}) (err error) {
	// 0. 参数的校验
	// 传入的 data 必须是指针类型(需要赋值)
	t := reflect.TypeOf(data)
	fmt.Println(t, t.Kind(), t.Elem().Kind()) // *main.MySQLConfig, ptr, struct
	if t.Kind() != reflect.Ptr {
		err = errors.New("input should be a pointer")  // 创建一个 error 类型的错误
		return
	}
	// 传入的 data 必须是结构体类型的指针(需要把各种键值对赋值给结构体)
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("input should be a struct")
		return
	}
	// 1. 读取文件, 获得字节类型的数据
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	//string(b): 将字节类型的文件内容转换成字符串
	lineSlice := strings.Split(string(b), "\n")
	// 2. 一行行分析数据
	var structName string
	for index, line := range lineSlice {
		// 去掉每行首位的空格, 避免情况如" [redis]"
		line = strings.TrimSpace(line)
		// 2.1 如果是注释就跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 2.2 如果 [ 开头就是结构体的类型标识, mysql or redis
		if strings.HasPrefix(line, "[") {
			// 处理边界情况 "["
			if !strings.HasSuffix(line, "]") {
				err = fmt.Errorf("line %d: syntax error, incorrect value - \"%s\"", index + 1, line)
				return
			}
			// 处理边界情况 "[    ]"
			sectionName := strings.TrimSpace(line[1:len(line) - 1])
			if len(sectionName) == 0 {
				err = fmt.Errorf("line %d: syntax error, incorrect value - \"%s\"", index + 1, line)
				return
			}
			// 接下来从 sectionName 根据反射找到对应的结构体
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				if sectionName == field.Tag.Get("ini") {
					structName = field.Name
					//fmt.Printf("找到%s对应的嵌套结构体%s了\n", sectionName, structName)
					break
				}
			}
		} else {
			// 2.3 如果不以 [ 开头, 这行表示的是 = 分隔的键值对
		}
	}
	return
}

func main() {
	var cfg Config
	err := loadIni("./src/config.ini", &cfg)
	if err != nil {
		fmt.Printf("load config ini failed, error: %v\n", err)
		return
	}
}