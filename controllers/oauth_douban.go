// Copyright 2011 The goauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program makes a call to the specified API, authenticated with OAuth2.
// a list of example APIs can be found at https://code.google.com/oauthplayground/
package controllers

import (
	"BeeGoTo/models/douban"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"io/ioutil"
)

type OauthDoubanController struct {
	beego.Controller
}

func (this *OauthDoubanController) Prepare() {

}

func (this *OauthDoubanController) Get() {

	dbconf, err := config.NewConfig("ini", "conf/douban.conf")
	if err != nil {
		fmt.Println(err)
	}

	DBConfig := &oauth.Config{
		ClientId:     dbconf.String("douban_clientId"),
		ClientSecret: dbconf.String("douban_clientSecret"),
		RedirectURL:  dbconf.String("douban_redirectURL"),
		Scope:        dbconf.String("douban_scope"),
		AuthURL:      dbconf.String("douban_authURL"),
		TokenURL:     dbconf.String("douban_tokenURL"),
		TokenCache:   oauth.CacheFile(dbconf.String("douban_cachefile") + "test001" + ".json"),
	}
	DBtransport := &oauth.Transport{Config: DBConfig}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := DBConfig.TokenCache.Token()
	if err != nil {

		url := DBConfig.AuthCodeURL("")
		this.Data["AuthCodeURL"] = url

		t_code := this.GetString("code")
		fmt.Println("Code is %s\n", t_code)
		if token == nil {

			if t_code != "" {

				token, err = DBtransport.Exchange(t_code)
				fmt.Println("token:", token)

				if err != nil {
					fmt.Println("Exchange:", err)
				}
				// (The Exchange method will automatically cache the token.)
				fmt.Println("Token is cached in %v\n", DBConfig.TokenCache)

			}
			return
		}
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	DBtransport.Token = token

	// Make the request.
	r, err := DBtransport.Client().Get(dbconf.String("douban_requestURL"))
	if err != nil {
		fmt.Println("GetErr:", err)
	}
	defer r.Body.Close()

	if err != nil {
		fmt.Println("Request Error:", err)
	}

	body, _ := ioutil.ReadAll(r.Body)

	this.Data["JsonStr"] = string(body)
	fmt.Println("Request Error:", string(body))

	var line douban.User
	err = json.Unmarshal(body, &line)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", line)
}

func (this *OauthDoubanController) Post() {

}
