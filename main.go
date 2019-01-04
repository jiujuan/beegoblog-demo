package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "mybeegodemo/controllers"
	"mybeegodemo/models"
	_ "mybeegodemo/routers"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true)

	// beego.Run("/", &controllers.MainController{})
	beego.Run()
}
