package controllers

import (
	"github.com/astaxie/beego"
	"mybeegodemo/models"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	this.Data["IsTopic"] = true

	var err error
	this.Data["Topics"], err = models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	}
	this.TplName = "topic.html"
}

func (this *TopicController) Add() {
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	this.TplName = "topic_add.html"
}

func (this *TopicController) Post() {
	if !CheckAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	title := this.Input().Get("title")
	content := this.Input().Get("content")
	tid := this.Input().Get("tid")
	category := this.Input().Get("category")
	label := this.Input().Get("label")

	//获取文件
	_, fh, err := this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(title, content, category, label, attachment)
	} else {
		err = models.ModifyTopic(tid, title, content, category, label, attachment)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}

func (this *TopicController) View() {
	this.TplName = "topic_view.html"

	topic, err := models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
	} else {
		this.Data["Topic"] = topic
		this.Data["Tid"] = this.Ctx.Input.Param("0")
		this.Data["Labels"] = strings.Split(topic.Labels, " ")
	}

	replies, err := models.GetAllReplies(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Replies"] = replies
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
}

func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"

	tid := this.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	} else {
		this.Data["Topic"] = topic
		this.Data["Tid"] = tid
	}
}

func (this *TopicController) Delete() {
	if !CheckAccount(this.Ctx) {
		this.Redirect("/", 302)
		return
	}
	err := models.DelTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}
