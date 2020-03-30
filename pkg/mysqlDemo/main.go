package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	id int
	name string
	age int
}

// 初始化数据库
func initDB() (err error) {
	// 连接数据库
	dsn := "root:Bristol123395@tcp(127.0.0.1:3306)/goDB"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(10)
	db.SetMaxIdleConns(5)
	fmt.Println("连接数据库成功!")
	return
}

// 查询单条
func queryOne(id int) {
	// 1. 写一条 SQL 语句
	sqlStr := `select id, name, age from user where id = ?;`
	// 2. 执行 SQL, 并为结构体赋值
	var u1 user
	// 从连接池拿一个连接去做单条查询, 注意必须调用 Scan 方法从而关闭连接, 为防止忘记直接写在后面
	db.QueryRow(sqlStr, id).Scan(&u1.id, &u1.name, &u1.age)
	fmt.Printf("%#v\n", u1)
}

// 查询多条
func queryMore(n int) {
	// 1. SQL 语句
	sqlStr := `select id, name, age from user where id > ?;`
	// 2. 执行 SQL 语句
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Printf("execute query %s, failed with err: %v\n ", sqlStr, err)
		return
	}
	// 3. 关闭 rows
	defer rows.Close()
	// 4. 循环取值
	for rows.Next() {
		var u1 user
		err := rows.Scan(&u1.id, &u1.name, &u1.age)
		if err != nil {
			fmt.Printf("scan failed, err: %v\n", err)
		}
		fmt.Printf("u1: %#v\n", u1)
	}
}

// 插入数据
func insert() {
	// 1. 写 SQL 语句
	sqlStr := `insert into user(name, id) values("douding", 10);`
	// 2. 执行 SQL 语句
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("insert failed, err: %v", err)
		return
	}
	// 如果是插入数据的操作, 会拿到插入数据的 ID 值
	id, err := ret.LastInsertId()
	 if err != nil {
	 	fmt.Printf("get id failed, err: %v\n", err)
	 	return
	 }
	 fmt.Printf("ID: %d\n", id)
}

// 更新数据
func update(newAge, id int) {
	// 1. 写 SQL 语句
	sqlStr := `update user set age = ? where id = ?;`
	// 2. 执行 SQL 语句
	ret, err := db.Exec(sqlStr, newAge, id)
	if err != nil {
		fmt.Printf("insert failed, err: %v", err)
		return
	}
	// 如果是更新数据的操作, 会返回更新行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed, err: %v\n", err)
		return
	}
	fmt.Printf("更新了%d行\n", n)
}

// 删除数据
func delete(id int) {
	// 1. 写 SQL 语句
	sqlStr := `delete from user where id = ?;`
	// 2. 执行 SQL 语句
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		fmt.Printf("insert failed, err: %v", err)
		return
	}
	// 如果是删除数据的操作, 会返回删除行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed, err: %v\n", err)
		return
	}
	fmt.Printf("删除了%d行\n", n)
}

// 预处理方式插入多条语句
func prepareInsert() {
	sqlStr := `insert into user(name, age) values(?, ?);`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err: %v\n", err)
		return
	}
	defer stmt.Close()
	// 后序只需要拿到 stmt 去操作
	var m = map[string]int{
		"豆丁": 10,
		"晓雪": 29,
	}
	for k, v := range m {
		stmt.Exec(k, v)
		fmt.Printf("插入成功! name: %s, age: %d\n", k, v)
	}
}

// 事物操作
func transaction() {
	// 1. 开始事务
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("begin failed, err: %v\n", err)
		return
	}
	// 2. 执行多个 SQL 操作
	sqlStr1 := `update user set age = age+5 where id = 1`
	sqlStr2 := `update xxx set age = age+3 where id = 12`
	// 执行 SQL 1
	_, err = tx.Exec(sqlStr1)
	if err != nil {
		tx.Rollback()
		fmt.Println("执行 SQL 1 失败了, 需要回滚")
		return
	}
	// 执行 SQL 2
	_, err = tx.Exec(sqlStr2)
	if err != nil {
		tx.Rollback()
		fmt.Println("执行 SQL 2 失败了, 需要回滚")
		return
	}
	// 都执行成功, 提交本次操作
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Println("提交失败了, 需要回滚")
		return
	}
	fmt.Println("事物执行成功!")
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init DB failed, error: %#v\n", err)
	}
	//queryOne(1)
	//queryMore(0)
	//insert()
	//update(99, 10)
	//delete(10)
	//prepareInsert()
	transaction()
}
