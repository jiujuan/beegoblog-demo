package controllers

import (
	"github.com/astaxie/beego"
	"mybeegodemo/models"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["IsLogin"] = CheckAccount(this.Ctx)

	cate := this.Input().Get("cate")
	topics, err := models.GetAllTopics(cate, this.Input().Get("label"), true)
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Topics"] = topics
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	this.Data["Categories"] = categories
}
