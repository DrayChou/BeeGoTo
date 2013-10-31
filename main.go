package main

import (
	"BeeGoTo/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Run()
}
