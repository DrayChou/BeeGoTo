package twitter

import (
	"encoding/json"
	//"flag"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"net/url"
	//"os"
	//"strings"
	//"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

type TwitterError struct {
	prefix string
	msg    string
}

func (oe TwitterError) Error() string {
	fmt.Println("TwitterAPI: " + oe.prefix + ": " + oe.msg)
	beego.Error("TwitterAPI: " + oe.prefix + ": " + oe.msg)
	return "TwitterError: " + oe.prefix + ": " + oe.msg
}

type Twitter struct {
	Conf         string
	ConfFileType string
	Service      *oauth1a.Service
	UserConfig   *oauth1a.UserConfig
}

func (this *Twitter) LoadToken(uid string) (err_r error, client *twittergo.Client) {

	tconf, err := config.NewConfig(this.ConfFileType, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "twitter 配置文件加载失败"}, client
	}
	beego.Debug("TwitterAPI:LoadToken:tconf:", tconf)

	this.Service.ClientConfig = &oauth1a.ClientConfig{
		ConsumerKey:    tconf.String("clientId"),
		ConsumerSecret: tconf.String("clientSecret"),
		CallbackURL:    tconf.String("callbackURL"),
	}

	tokenconf, err := config.NewConfig("json", tconf.String("cacheDir")+uid+".json")
	if err != nil {
		return TwitterError{"Auth", "user 配置文件加载失败"}, client
	}
	beego.Debug("TwitterAPI:LoadToken:tokenconf:", tokenconf)

	if tokenconf.String("Token") == "" || tokenconf.String("Secret") == "" {
		return TwitterError{"Auth", "user 配置文件有误"}, client
	}

	params, err := url.ParseQuery(tokenconf.String("AccessValues"))
	this.UserConfig = &oauth1a.UserConfig{
		RequestTokenSecret: tokenconf.String("RequestTokenSecret"),
		RequestTokenKey:    tokenconf.String("RequestTokenKey"),
		AccessTokenSecret:  tokenconf.String("AccessTokenSecret"),
		AccessTokenKey:     tokenconf.String("AccessTokenKey"),
		Verifier:           tokenconf.String("Verifier"),
		AccessValues:       params,
	}
	beego.Debug("TwitterAPI:LoadToken:params:", params)
	beego.Debug("TwitterAPI:LoadToken:this.UserConfig:", this.UserConfig)

	if this.UserConfig.AccessTokenSecret == "" || this.UserConfig.AccessTokenKey == "" {
		return TwitterError{"Auth", "请先取得推特授权"}, client
	}

	user := oauth1a.NewAuthorizedConfig(this.UserConfig.AccessTokenKey, this.UserConfig.AccessTokenSecret)
	client = twittergo.NewClient(this.Service.ClientConfig, user)
	beego.Debug("TwitterAPI:LoadToken:user:", user)
	beego.Debug("TwitterAPI:LoadToken:client:", client)

	return nil, client
}

func (this *Twitter) SaveToken(uid string, UserConfig *oauth1a.UserConfig) (err_r error) {
	tconf, err := config.NewConfig(this.ConfFileType, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "twitter 配置文件加载失败"}
	}
	beego.Debug("TwitterAPI:SaveToken:tconf:", tconf)

	b, err := json.Marshal(UserConfig)
	if err != nil {
		return TwitterError{"Auth", "token 编码失败"}
	}
	beego.Debug("TwitterAPI:SaveToken:b:", b)

	ioutil.WriteFile(tconf.String("cacheDir")+uid+".json", b, 0644)

	return nil
}

func (this *Twitter) AuthUrl(uid string) (err error, url string) {
	if this.Service == nil {
		return TwitterError{"AuthUrl", "推特配置异常"}, url
	}

	httpClient := new(http.Client)
	if err = this.UserConfig.GetRequestToken(this.Service, httpClient); err != nil {
		return TwitterError{"AuthUrl", "RequestTemporaryCredentials:" + err.Error()}, url
	}
	beego.Debug("TwitterAPI:AuthUrl:userConfig:", this.UserConfig)

	if url, err = this.UserConfig.GetAuthorizeURL(this.Service); err != nil {
		return TwitterError{"AuthUrl", "RequestTemporaryCredentials:" + err.Error()}, url
	}
	fmt.Println("url:", url)
	beego.Debug("TwitterAPI:AuthUrl:url:", url)

	//this.SaveToken(uid+".tmp", userConfig)
	return nil, url
}

func (this *Twitter) Auth(uid string, oauth_token string, oauth_verifier string) error {
	if uid == "" {
		return TwitterError{"Auth", "內網用戶ID不能爲空"}
	}

	beego.Debug("TwitterAPI:Auth:User_id:", uid)

	if this.Conf == "" {
		return TwitterError{"Auth", "配置文件参数错误"}
	}

	if this.ConfFileType == "" {
		this.ConfFileType = "ini"
	}

	this.Service = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{},
		Signer:       new(oauth1a.HmacSha1Signer),
	}

	beego.Debug("TwitterAPI:Auth:this.Service:", this.Service)

	err, client := this.LoadToken(uid)
	beego.Debug("TwitterAPI:Auth:client:", client)

	if err != nil {
		beego.Debug("TwitterAPI:Auth:oauth_token:", oauth_token)
		beego.Debug("TwitterAPI:Auth:oauth_verifier:", oauth_verifier)

		httpClient := new(http.Client)
		if err = this.UserConfig.GetAccessToken(oauth_token, oauth_verifier, this.Service, httpClient); err != nil {
			beego.Debug("TwitterAPI:Auth:client:", err.Error())
			return TwitterError{"Auth", err.Error()}
		}
	}
	return nil
}

//func (this *Twitter) Refresh() error {
//	return this.OauthClient.Refresh()
//}

//func (this *Twitter) User(uid string) error {

//	u := User{}
//	request_url := "https://api.twitter.com/1.1/users/show.json?screen_name=" + uid
//	if uid == "" {
//		if this.OauthClient.Token == nil {
//			return TwitterError{"User", "未授权，请先授权"}
//		}
//		request_url = "https://api.twitter.com/1.1/account/verify_credentials.json"
//	}

//	r, err := this.OauthClient.Get(request_url)
//	if err != nil {
//		return TwitterError{"User", "请求失败:" + err.Error()}
//	}
//	defer r.Body.Close()

//	body, _ := ioutil.ReadAll(r.Body)

//	beego.Debug("TwitterAPI:User:Request StatusCode:", r.StatusCode)
//	beego.Debug("TwitterAPI:User:Request Body:", string(body))

//	err = json.Unmarshal(body, &u)
//	if err != nil {
//		return TwitterError{"User", "JSON解析失败:" + err.Error()}
//	}

//	return nil
//}

//func (this *Twitter) UserTimeLine(uid string, count int64, since_id int64) error {
//	_, token := this.LoadToken(uid)
//	beego.Debug("TwitterAPI:UserTimeLine:token:", token)

//	r, err := this.OauthClient.Get(http.DefaultClient, &token,
//		"http://api.twitter.com/1.1/statuses/home_timeline.json", nil)
//	if err != nil {
//		return TwitterError{"Auth", "RequestTemporaryCredentials:" + err.Error()}
//	}
//	defer r.Body.Close()

//	body, _ := ioutil.ReadAll(r.Body)

//	beego.Debug("TwitterAPI:UserTimeLine:Request StatusCode:", r.StatusCode)
//	beego.Debug("TwitterAPI:UserTimeLine:Request Body:", string(body))

//	//err = json.Unmarshal(body, &tl)
//	//if err != nil {
//	//	fmt.Println("Unmarshal Error:", err)
//	//	return TwitterError{"UserTimeLine", "JSON解析失败:" + err.Error()}, tl
//	//}

//	return nil
//}

//func (this *Twitter) Shuo(text string) (error, TimeLine) {
//	tl := TimeLine{}

//	if this.OauthClient.Token == nil {
//		return TwitterError{"Shuo", "未授权，请先授权"}, tl
//	}

//	if text == "" {
//		return TwitterError{"Shuo", "参数错误"}, tl
//	}
//	request_url := "https://api.douban.com/shuo/v2/statuses/"

//	v := url.Values{}
//	v.Set("source", this.OauthClient.ClientId)
//	v.Set("text", text)
//	r, err := this.OauthClient.PostForm(request_url, v)
//	if err != nil {
//		return TwitterError{"Shuo", "请求失败:" + err.Error()}, tl
//	}
//	defer r.Body.Close()

//	body, _ := ioutil.ReadAll(r.Body)

//	beego.Debug("TwitterAPI:Shuo:Request StatusCode:", r.StatusCode)
//	beego.Debug("TwitterAPI:Shuo:Request Body:", string(body))

//	if err = json.Unmarshal(body, &tl); r.StatusCode != 200 || err != nil {
//		return TwitterError{"Shuo", "JSON解析失败:" + err.Error()}, tl
//	}

//	return nil, tl
//}
