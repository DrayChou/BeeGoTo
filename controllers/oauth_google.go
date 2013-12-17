// Copyright 2011 The goauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program makes a call to the specified API, authenticated with OAuth2.
// a list of example APIs can be found at https://code.google.com/oauthplayground/
package controllers

import (
	"code.google.com/p/goauth2/oauth"
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"os"
)

type OauthGoogleController struct {
	beego.Controller
}

var GOConfig = &oauth.Config{
	ClientId:     beego.AppConfig.String("clientId"),
	ClientSecret: beego.AppConfig.String("clientSecret"),
	RedirectURL:  beego.AppConfig.String("redirectURL"),
	Scope:        beego.AppConfig.String("scope"),
	AuthURL:      beego.AppConfig.String("authURL"),
	TokenURL:     beego.AppConfig.String("tokenURL"),
	TokenCache:   oauth.CacheFile(beego.AppConfig.String("cachefile")),
}

var GOtransport = &oauth.Transport{Config: GOConfig}

func (this *OauthGoogleController) Prepare() {

}

func (this *OauthGoogleController) Get() {
	this.TplNames = "oauth_google/get.tpl"

	// Try to pull the token from the cache; if this fails, we need to get one.
	_, err := GOConfig.TokenCache.Token()
	if err != nil {
		if GOConfig.ClientId == "" || GOConfig.ClientSecret == "" {
			fmt.Fprint(os.Stderr)
			os.Exit(2)
		}

		url := GOConfig.AuthCodeURL("")
		this.Data["AuthCodeURL"] = url
	}

}

func (this *OauthGoogleController) Post() {
	this.TplNames = "oauth_google/post.tpl"

	t_code := this.GetString("code")

	token, err := GOtransport.Exchange(t_code)
	if err != nil {
		fmt.Printf("Exchange:", err)
	}
	// (The Exchange method will automatically cache the token.)
	fmt.Printf("Token is cached in %v\n", GOConfig.TokenCache)

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	GOtransport.Token = token

	// Make the request.
	r, err := GOtransport.Client().Get(GOConfig.RedirectURL)
	if err != nil {
		fmt.Printf("Get:", err)
	}
	defer r.Body.Close()
	fmt.Println("r.Body:", r.Body)

	if err != nil {
		fmt.Printf("Request Error:", err)
	}

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body:", body)

	// Write the response to standard output.
	io.Copy(os.Stdout, r.Body)

	// Send final carriage return, just to be neat.
	fmt.Println()
}
