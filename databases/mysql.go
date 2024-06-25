package sdb

import (
	"GinProject12/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"sync"
)

var (
	db       *gorm.DB
	once     sync.Once
	sqlUser  = "root"
	passwd   = "sgj123456"
	ip       = "127.0.0.1"
	sport    = "3306"
	database = "ginproject12"
)

// InitMysql 单例模式
func InitMysql() {
	once.Do(func() {
		var err error
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sqlUser, passwd, ip, sport, database))

		if err != nil {
			panic(err.Error())
		}
	})

}

func GetMysqlDB() *gorm.DB {
	if db == nil {
		InitMysql()
	}
	db.AutoMigrate(model.User{})
	return db
}

func CloseMysql() {
	if db != nil {
		panic(db.Close())
		log.Println("databases close!")
	}
}

func CreateInfo(info any) {
	if err := db.Create(info).Error; err != nil {
		panic(err.Error())
	}
}

func IsExistInfoByString(query string, data string) bool {
	if err := db.Where(fmt.Sprintf("%s=?", query), data).Error; err != nil {
		return true
	}
	return false
}

func FindInfoByString(query string, data string) *model.User {
	user := &model.User{}
	db.Where(fmt.Sprintf("%s=?", query), data).First(&user)
	return user
}

func FindInfoById(id int) *model.User {
	var user = &model.User{}
	db.First(&user, id)
	return user
}
