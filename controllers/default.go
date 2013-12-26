package controllers

import (
	"BeeGoTo/models"
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {

	//models.MgoTest()

	this.SetSession("test", "test")
	v := this.GetSession("asta")
	if v == nil {
		this.SetSession("asta", int(1))
		this.Data["num"] = 0
	} else {
		this.SetSession("asta", v.(int)+1)
		this.Data["num"] = v.(int)
	}

	UserInfo := this.GetSession("UserInfo")
	if UserInfo == nil {
		this.SetSession("asta", int(1))
		this.Data["num"] = 0
	}

	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"

	fmt.Println(this.Data, this.GetSession("asta"))
}
