// Copyright 2011 The goauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program makes a call to the specified API, authenticated with OAuth2.
// a list of example APIs can be found at https://code.google.com/oauthplayground/
package controllers

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type GoogleUser struct {
	id          string
	name        string
	given_name  string
	family_name string
	link        string
	picture     string
	gender      string
	locale      string
}

type OauthGoogleController struct {
	beego.Controller
}

func (this *OauthGoogleController) Get() {
	this.Layout = "layout.tpl"

	flag.Parse()

	var (
		clientId     = flag.String("id", beego.AppConfig.String("google_clientId"), "Client ID")
		clientSecret = flag.String("secret", beego.AppConfig.String("google_clientSecret"), "Client Secret")
		scope        = flag.String("scope", beego.AppConfig.String("google_scope"), "OAuth scope")
		redirectURL  = flag.String("redirect_url", beego.AppConfig.String("google_redirectURL"), "Redirect URL")
		authURL      = flag.String("auth_url", beego.AppConfig.String("google_authURL"), "Authentication URL")
		tokenURL     = flag.String("token_url", beego.AppConfig.String("google_tokenURL"), "Token URL")
		requestURL   = flag.String("request_url", beego.AppConfig.String("google_requestURL"), "API request")
		code         = flag.String("code", beego.AppConfig.String("google_code"), "Authorization Code")
		cachefile    = flag.String("cache", beego.AppConfig.String("google_cachefile"), "Token cache file")
	)

	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scope:        *scope,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			this.Data["url"] = url
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			return
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token

	// Make the request.
	r, err := transport.Client().Get(*requestURL)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer r.Body.Close()
	fmt.Println("r.Body:", r.Body)

	if err != nil {
		log.Fatal("Request Error:", err)
	}

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body:", body)
	var userInfo []GoogleUser
	if json.Unmarshal(body, &userInfo) != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("userInfo:", userInfo)

	// Write the response to standard output.
	io.Copy(os.Stdout, r.Body)

	// Send final carriage return, just to be neat.
	fmt.Println()
}

func (this *OauthGoogleController) Post() {

}
