package controllers

import (
	"BeeGoTo/models/twitter"
	//"fmt"
	"github.com/astaxie/beego"
)

type OauthTwitterController struct {
	beego.Controller
}

func (this *OauthTwitterController) Prepare() {

}

func (this *OauthTwitterController) Get() {

	t_token := this.GetString("oauth_token")
	t_token_secret := this.GetString("oauth_token_secret")
	db := &twitter.Twitter{Conf: "conf/twitter.conf"}

	if ok := db.Auth("we_get", t_token, t_token_secret); ok != nil {
		if ok2, url := db.AuthUrl("we_get"); ok2 != nil {

		} else {
			this.Data["AuthCodeURL"] = url
		}
	}

	//_, dbu := db.User("we_get")
	//fmt.Println("dbu:", dbu)

	//this.Data["dbu"] = dbu

	//err := db.UserTimeLine("we_get", 20, 0)
	//fmt.Println("err:", err)
	//fmt.Println("dbtl:", dbtl)
	//_, dbs := db.Shuo("test")
	//fmt.Println("dbs:", dbs)
}

func (this *OauthTwitterController) Post() {

}
