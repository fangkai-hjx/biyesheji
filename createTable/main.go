package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ModelTrain struct {
	gorm.Model
	Environment string // 环境
	Version string // 版本
	URL string  //访问地址
	Status string //训练环境状态
}
func main() {
	db, err := gorm.Open("mysql", "root:root@(localhost)/mall_admin?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()


	// 创建
	create := db.Create(&ModelTrain{
		Model:       gorm.Model{},
		Environment: "tensorflow",
		Version:     "1.8",
		URL:         "http://222.201.187.166:8081",
		Status:      "Success",
	})
	fmt.Println(create)
	//// 读取
	//var modelTrain ModelTrain
	//db.First(&modelTrain, 1) // 查询id为1的product
	//db.First(&modelTrain, "code = ?", "L1212") // 查询code为l1212的product
	//
	//// 更新 - 更新product的price为2000
	//db.Model(&modelTrain).Update("Status", "Success")
	//
	//// 删除 - 删除product
	//db.Delete(&modelTrain)
}