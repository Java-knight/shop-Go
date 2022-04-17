package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"io"
	"log"
	"os"
	"shop_srvs/user_srv/model"
	"time"

	"crypto/md5"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func getMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)

	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(192.168.198.138:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,         // 禁用彩色打印
		},
	)
	
	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	// pbkdf2包名-sha512算法
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("knight%d", i),
			Mobile:   fmt.Sprintf("1776912102%d", i),
			Password: newPassword,
		}
		db.Save(&user)
	}

	////设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	////sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么
	//
	////定义一个表结构， 将表结构直接生成对应的表 - migrations
	//// 迁移 schema
	//_ = db.AutoMigrate(&model.User{}) //此处应该有sql语句

	//fmt.Println(getMd5("123456"))

	// Using the default options
	//salt, encodedPwd := password.Encode("generic password", nil)
	//fmt.Println(salt)
	//fmt.Println(encodedPwd)
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // true

	//// Using custom options
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("generic password", options)
	//// pbkdf2包名-sha512算法
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(newPassword)
	//fmt.Println(len(newPassword))
	//passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	//fmt.Println(check) // true
}
