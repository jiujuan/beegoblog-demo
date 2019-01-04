package controllers

import (
	"github.com/astaxie/beego"
	"mybeegodemo/models"
)

type ReplyController struct {
	beego.Controller
}

func (this *ReplyController) Add() {
	tid := this.Input().Get("tid")
	err := models.AddReply(tid, this.Input().Get("nickname"), this.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}

	err = models.ModifyReplyCount(tid)
	if err != nil {
		beego.Error("replycount error")
	}
	this.Redirect("/topic/view/"+tid, 302)
}

func (this *ReplyController) Delete() {
	tid := this.Input().Get("tid")

	if !CheckAccount(this.Ctx) {
		beego.Error("not login")
	} else {
		err := models.DelReply(this.Input().Get("rid"))
		if err != nil {
			beego.Error(err)
		}

		err = models.ModifyReplyCount(tid)
		if err != nil {
			beego.Error("replycount error")
		}
	}
	this.Redirect("/topic/view/"+tid, 302)
}
