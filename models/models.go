package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index;auto_now_add;type(datetime)"`
	Views           int64     `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Labels          string
	Category        string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index;auto_now_add;type(datetime)"`
	Updated         time.Time `orm:"index;auto_now;type(datetime)"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time
	ReplyCount      int
	ReplyLastUserId int64
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"sieze(1000)"`
	Created time.Time `orm:"index;auto_now_add;type(datetime)"`
}

func RegisterDB() {
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/beegoblog?charset=utf8")
}

func AddReply(tid, nickname, content string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Name:    nickname,
		Content: content,
		Tid:     id,
		Created: time.Now(),
	}

	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}
	return nil
}

func GetAllReplies(tid string) (replies []*Comment, err error) {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Comment, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", id).All(&replies)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func DelReply(rid string) error {
	id, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	reply := &Comment{Id: id}
	if o.Read(reply) == nil {
		// tidNum := reply.Id
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	return err
}

//isType 1: true 增加 2：false 减少
func ModifyReplyCount(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		// var count int
		// if isType {
		// 	count = topic.ReplyCount + 1
		// } else if topic.ReplyCount > 0 {
		// 	count = topic.ReplyCount - 1
		// }
		// topic.ReplyCount = count
		// topic.ReplyTime = time.Now() //这个更新时间删除留言会有问题，下面就没有这个问题

		replies := make([]*Comment, 0)
		qs := o.QueryTable("comment")
		_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
		if err != nil {
			return err
		}
		topic.ReplyCount = len(replies)
		topic.ReplyTime = replies[0].Created

		_, err = o.Update(topic)
	}
	return err
}

func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{Title: name}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}

	if o.Read(cate) == nil {
		_, err = o.Delete(cate)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func AddTopic(title, content, category, label, attachment string) error {
	//标签 - 空格作为分隔符
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Labels:     label,
		Content:    content,
		Attachment: attachment,
		Category:   category,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

func GetAllTopics(category string, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()

	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")

	var err error
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.OrderBy("-id").All(&topics)
	}
	if len(topics) > 0 {
		for k, v := range topics {
			if len(v.Labels) > 0 {
				topics[k].Labels = strings.Replace(strings.Replace(v.Labels, "#", " ", -1), "$", "", -1)
			}
		}
	}
	return topics, err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()

	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)

	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, err
}

func ModifyTopic(tid, title, content, category, label, attachment string) error {
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	var oldCate, oldAttach string
	topic := &Topic{Id: id}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Labels = label
		topic.Content = content
		topic.Category = category
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	//删除旧附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	//更新旧分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	//更新新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return nil
}

func DelTopic(tid string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: id}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("category", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			o.Update(cate)
		}
	}

	return nil
}
