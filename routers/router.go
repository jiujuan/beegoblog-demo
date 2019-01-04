package routers

import (
	"github.com/astaxie/beego"
	"mybeegodemo/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/category", &controllers.CategoryController{})

	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})

	beego.Router("/reply", &controllers.ReplyController{})
	beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")

	beego.Router("/attachment/:all", &controllers.AttachmentController{})

	beego.Router("/login", &controllers.LoginController{})
}
