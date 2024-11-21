package main

import (
	"fmt"
	"goIM/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime&loc=Local", "huluwa", "vegagame", "192.168.1.180", 3306, "goIM")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	/* err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Relation{})
	if err != nil {
		panic(err)
	} */

	err = db.AutoMigrate(&models.Group{})
	if err != nil {
		panic(err)
	}
}
