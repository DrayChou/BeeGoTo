package controllers

import (
	"BeeGoTo/models/douban"
	"github.com/astaxie/beego"
)

type TestController struct {
	beego.Controller
}

func (this *TestController) Get() {

	db = new douban.Douban
	&douban.Auth("we_get")
	&douban.Userinfo("we_get")
}
