package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB = getMysqlConnect()

type UserInfo struct {
	gorm.Model
	Name    string `form:"name" json:"name" gorm:"type:varchar(64)"`
	Age     int    `form:"age" json:"age"`
	Address string `form:"address" json:"address" gorm:"type:varchar(255)"`
	Email   string `form:"email" json:"email" binding:"email"`
}

func queryUsers() []UserInfo {
	var users []UserInfo
	DB.Find(&users)
	return users
}

func getUserByName(name string) UserInfo {
	if len(name) == 0 {
		panic("name is empty")
	}
	var user UserInfo
	DB.Where("name = ?", name).First(&user)
	return user
}

func createUser(userInfo *UserInfo) (uint, error) {
	tx := DB.Create(userInfo)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return userInfo.ID, nil
}

func updateUserById(userInfo *UserInfo) error {
	id := userInfo.ID
	var res UserInfo
	result := DB.First(&res, id)
	if result.Error != nil {
		return result.Error
	}

	tx := DB.Save(&userInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func deleteUserByName(name string) error {
	tx := DB.Delete(&UserInfo{}, "name = ?", name)
	//DB.Unscoped().Delete(&UserInfo{}, name)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func getMysqlConnect() *gorm.DB {
	dsn := "root:root@(127.0.0.1:3306)/db_learn?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	err = db.AutoMigrate(&UserInfo{})
	if err != nil {
		panic(err.Error())
	}
	return db
}
