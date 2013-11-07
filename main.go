package main

import (
	"BeeGoTo/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.BeeLogger.SetLogger("file", `{"filename":"cache/logs/logs.log"}`)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/oauth/google", &controllers.OauthGoogleController{})
	beego.Router("/oauth/douban", &controllers.OauthDoubanController{})
	beego.Run()
}
