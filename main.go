package main

import (
	"fmt"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Git struct {
	gorm.Model
	OperatorID string `gorm:"operator_id"`
	Resource   string `gorm:"resource"`
	ResourceID string `gorm:"resource_id"`
	Action     string `gorm:"action"`
	Changes    string `gorm:"changes"`
	RemoteIP   string `gorm:"remote_ip"`
	Comment    string `gorm:"comment"`
}

type Resource struct {
	gorm.Model
	Content string `gorm:"content"`
}

var IgnoredColumns []string

func main() {
	db, err := gorm.Open("mysql", "root:@/git?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		panic("can not connect to db")
	}
	db.AutoMigrate(&Git{}, &Resource{})

	updated := func(scope *gorm.Scope) {
		r := &Resource{}
		for _, f := range scope.Fields() {
			log.Printf("dbname: %s, value: %v", f.DBName, f.Field.Interface())
		}
		db.First(r, 1)
		t := reflect.TypeOf(*r)
		v := reflect.ValueOf(*r)
		for i := 0; i < t.NumField(); i++ {
			log.Printf("%s:%v", t.Field(i).Name, v.Field(i))
		}
	}

	db.Callback().Update().Register("log_create", updated)
	db.Create(&Resource{Content: "content"})

	fmt.Println("end")
}
