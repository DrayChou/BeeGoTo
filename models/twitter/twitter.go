package twitter

import (
	"encoding/json"
	//"flag"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
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
	OauthClient  *oauth.Client
}

func (this *Twitter) LoadToken(uid string) (err_r error, token oauth.Credentials) {

	tconf, err := config.NewConfig(this.ConfFileType, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "twitter 配置文件加载失败"}, token
	}

	tokenconf, err := config.NewConfig("json", tconf.String("cacheDir")+uid+".json")
	if err != nil {
		return TwitterError{"Auth", "user 配置文件加载失败"}, token
	}

	if tokenconf.String("Token") == "" || tokenconf.String("Secret") == "" {
		return TwitterError{"Auth", "user 配置文件有误"}, token
	}

	token = oauth.Credentials{
		Token:  tokenconf.String("Token"),
		Secret: tokenconf.String("Secret"),
	}

	return nil, token
}

func (this *Twitter) SaveToken(uid string, token oauth.Credentials) (err_r error) {
	tconf, err := config.NewConfig(this.ConfFileType, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "twitter 配置文件加载失败"}
	}

	b, err := json.Marshal(token)
	if err != nil {
		return TwitterError{"Auth", "token 编码失败"}
	}
	ioutil.WriteFile(tconf.String("cacheDir")+uid+".json", b, 0644)

	return nil
}

func (this *Twitter) AuthUrl(uid string) (error, string) {
	if this.OauthClient == nil {
		return TwitterError{"AuthUrl", "推特配置异常"}, ""
	}

	tempCred, err := this.OauthClient.RequestTemporaryCredentials(http.DefaultClient, "oob", nil)
	if err != nil {
		return TwitterError{"AuthUrl", "RequestTemporaryCredentials:" + err.Error()}, ""
	}
	beego.Debug("TwitterAPI:AuthUrl:tempCred:", tempCred)

	this.SaveToken(uid+".tmp", *tempCred)

	url := this.OauthClient.AuthorizationURL(tempCred, nil)
	fmt.Println("url:", url)
	return nil, url
}

func (this *Twitter) Auth(uid string, code string) error {
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

	tconf, err := config.NewConfig(this.ConfFileType, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "配置文件加载失败"}
	}

	beego.Debug("TwitterAPI:User:tconf:", tconf)

	this.OauthClient = &oauth.Client{
		TemporaryCredentialRequestURI: tconf.String("redirectURL"),
		ResourceOwnerAuthorizationURI: tconf.String("authURL"),
		TokenRequestURI:               tconf.String("tokenURL"),
		Credentials: oauth.Credentials{
			Token:  tconf.String("clientId"),
			Secret: tconf.String("clientSecret"),
		},
	}

	beego.Debug("TwitterAPI:Auth:this.OauthClient:", this.OauthClient)

	err, token := this.LoadToken(uid)
	beego.Debug("TwitterAPI:Auth:token:", token)

	if err == nil {
		this.OauthClient.Credentials = oauth.Credentials{
			Token:  token.Token,
			Secret: token.Secret,
		}
	} else {
		beego.Debug("TwitterAPI:Auth:Code:", code)
		if code == "" {
			return TwitterError{"Auth", "code 为空，无法授权"}
		}

		err, tempCred := this.LoadToken(uid + ".tmp")
		beego.Debug("TwitterAPI:Auth:tempCred:", tempCred)

		if err != nil {
			return TwitterError{"Auth", "请先取得授权地址"}
		}

		tokenCred, _, err := this.OauthClient.RequestToken(http.DefaultClient, &tempCred, code)
		if err != nil {
			log.Fatal(err)
			return TwitterError{"Auth", "RequestTemporaryCredentials:" + err.Error()}
		}

		beego.Debug("TwitterAPI:Auth:tokenCred:", tokenCred)

		this.SaveToken(uid, *tokenCred)
		this.OauthClient.Credentials = oauth.Credentials{
			Token:  tokenCred.Token,
			Secret: tokenCred.Secret,
		}
	}

	return nil
}

//func (this *Twitter) Refresh() error {
//	return this.Transport.Refresh()
//}

func (this *Twitter) User(uid string) error {

	request_url := "https://api.douban.com/v2/user/" + uid
	if uid == "" {
		if this.Transport.Token == nil {
			return TwitterError{"User", "未授权，请先授权"}
		}
		request_url = "https://api.twitter.com/1.1/account/verify_credentials.json"
	}

	r, err := this.Transport.Client().Get(request_url)
	if err != nil {
		return TwitterError{"User", "请求失败:" + err.Error()}
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	beego.Debug("TwitterAPI:User:Request StatusCode:", r.StatusCode)
	beego.Debug("TwitterAPI:User:Request Body:", string(body))

	err = json.Unmarshal(body, &u)
	if err != nil {
		return TwitterError{"User", "JSON解析失败:" + err.Error()}
	}

	return nil
}

func (this *Twitter) UserTimeLine(uid string, count int64, since_id int64) error {
	_, token := this.LoadToken(uid)
	beego.Debug("TwitterAPI:UserTimeLine:token:", token)

	r, err := this.OauthClient.Get(http.DefaultClient, &token,
		"http://api.twitter.com/1.1/statuses/home_timeline.json", nil)
	if err != nil {
		return TwitterError{"Auth", "RequestTemporaryCredentials:" + err.Error()}
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	beego.Debug("TwitterAPI:UserTimeLine:Request StatusCode:", r.StatusCode)
	beego.Debug("TwitterAPI:UserTimeLine:Request Body:", string(body))

	//err = json.Unmarshal(body, &tl)
	//if err != nil {
	//	fmt.Println("Unmarshal Error:", err)
	//	return TwitterError{"UserTimeLine", "JSON解析失败:" + err.Error()}, tl
	//}

	return nil
}

//func (this *Twitter) Shuo(text string) (error, TimeLine) {
//	tl := TimeLine{}

//	if this.Transport.Token == nil {
//		return TwitterError{"Shuo", "未授权，请先授权"}, tl
//	}

//	if text == "" {
//		return TwitterError{"Shuo", "参数错误"}, tl
//	}
//	request_url := "https://api.douban.com/shuo/v2/statuses/"

//	v := url.Values{}
//	v.Set("source", this.Transport.ClientId)
//	v.Set("text", text)
//	r, err := this.Transport.Client().PostForm(request_url, v)
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
