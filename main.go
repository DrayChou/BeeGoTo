package main

import (
	"BeeGoTo/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.BeeLogger.SetLogger("file", `{"filename":"logs/logs.log"}`)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/oauth/google", &controllers.OauthGoogleController{})
	beego.Run()
}
