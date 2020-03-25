package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
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
	Test bool `ini:"test"`
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
		// 跳过空行
		if len(line) == 0 {
			continue
		}
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
			// 2.3.1 以 = 分隔这一行, 左 key 由 value
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line: %d, syntax error, incorrect value - \"%s\"", index + 1, line)
				return
			}
			index := strings.Index(line, "=")
			// "port =3306"
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index + 1:])
			// 2.3.2 根据 structName, 去 data 中把对应的嵌套结构体取出来
			v := reflect.ValueOf(data)
			sValue := v.Elem().FieldByName(structName) //拿到嵌套结构体的值信息
			sType := sValue.Type() // 拿到嵌套结构体的类型信息
			// 判断 config 中的字段是否是个结构体
			if sType.Kind() != reflect.Struct {
				err = fmt.Errorf("data 中的%s字段应该是个结构体", structName)
			}
			//2.3.3 遍历嵌套结构体每个字段, 判断 tag 是不是等于 key
			var fieldName string
			var fieldType reflect.StructField
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i) // tag 是存储在类型信息中的
				fieldType = field
				if field.Tag.Get("ini") == key {
					// 找到对应字段
					fieldName = field.Name
					break
				}
			}
			// 2.3.4 如果 key = tag, 给字段赋值
			// 根据 fieldName, 取出这个字段
			if len(fieldName) == 0 {
				// 在结构体中找不到对应的字段
				continue
			}
			fieldObj := sValue.FieldByName(fieldName)
			// 对其赋值
			//fmt.Println(fieldName, fieldType.Type.Kind(), value)
			switch fieldType.Type.Kind() {
			case reflect.String:
				fieldObj.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var valueInt int64
				valueInt, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("line: %d, syntax error, incorrect value - \"%s\"", index + 1, line)
					return
				}
				fieldObj.SetInt(valueInt)
			case reflect.Bool:
				var valueBool bool
				valueBool, err = strconv.ParseBool(value)
				if err != nil {
					err = fmt.Errorf("line: %d, syntax error, incorrect value - \"%s\"", index + 1, line)
					return
				}
				fieldObj.SetBool(valueBool)
			case reflect.Float32, reflect.Float64:
				var valueFloat float64
				valueFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					err = fmt.Errorf("line: %d, syntax error, incorrect value - \"%s\"", index + 1, line)
					return
				}
				fieldObj.SetFloat(valueFloat)
			}
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
	fmt.Printf("%#v\n", cfg)
}