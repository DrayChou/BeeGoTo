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
	Conffiletype string
	OauthClient  *oauth.Client
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

	tconf, err := config.NewConfig(this.Conffiletype, this.Conf)
	if err != nil {
		return TwitterError{"Auth", "配置文件加载失败"}, ""
	}
	beego.Debug("TwitterAPI:User:tconf:", tconf)

	b, _ := json.Marshal(tempCred)
	ioutil.WriteFile(tconf.String("cacheDir")+uid+".tmp"+".json", b, 0644)

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

	if this.Conffiletype == "" {
		this.Conffiletype = "ini"
	}

	tconf, err := config.NewConfig(this.Conffiletype, this.Conf)
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

	tokenconf, err := config.NewConfig("json", tconf.String("cacheDir")+uid+".json")

	if err == nil {
		beego.Debug("TwitterAPI:Auth:tokenconf:", tokenconf)

		this.OauthClient.Credentials = oauth.Credentials{
			Token:  tokenconf.String("Token"),
			Secret: tokenconf.String("Secret"),
		}
	} else {
		beego.Debug("TwitterAPI:Auth:Code:", code)
		if code == "" {
			return TwitterError{"Auth", "code 为空，无法授权"}
		}

		tmpconf, err := config.NewConfig("json", tconf.String("cacheDir")+uid+".tmp"+".json")
		beego.Debug("TwitterAPI:Auth:tmpconf:", tmpconf)

		if err != nil {
			return TwitterError{"Auth", "请先取得授权地址"}
		}

		tempCred := &oauth.Credentials{
			Token:  tmpconf.String("Token"),
			Secret: tmpconf.String("Secret"),
		}
		beego.Debug("TwitterAPI:Auth:this.TempCred:", tempCred)

		tokenCred, _, err := this.OauthClient.RequestToken(http.DefaultClient, tempCred, code)
		if err != nil {
			log.Fatal(err)
			return TwitterError{"Auth", "RequestTemporaryCredentials:" + err.Error()}
		}

		beego.Debug("TwitterAPI:Auth:tokenCred:", tokenCred)

		b, _ := json.Marshal(tokenCred)
		ioutil.WriteFile(tconf.String("cacheDir")+uid+".json", b, 0644)
		this.OauthClient.Credentials.Secret = tokenCred.Secret
		this.OauthClient.Credentials.Token = tokenCred.Token
	}

	return nil
}

//func (this *Twitter) Refresh() error {
//	return this.Transport.Refresh()
//}

//func (this *Twitter) User(uid string) (error, User) {

//	u := User{}
//	if uid == "" {
//		if this.Transport.Token == nil {
//			return TwitterError{"User", "未授权，请先授权"}, u
//		}
//		uid = "~me"
//	}
//	request_url := "https://api.douban.com/v2/user/" + uid

//	r, err := this.Transport.Client().Get(request_url)
//	if err != nil {
//		return TwitterError{"User", "请求失败:" + err.Error()}, u
//	}
//	defer r.Body.Close()

//	body, _ := ioutil.ReadAll(r.Body)

//	beego.Debug("TwitterAPI:User:Request StatusCode:", r.StatusCode)
//	beego.Debug("TwitterAPI:User:Request Body:", string(body))

//	err = json.Unmarshal(body, &u)
//	if err != nil {
//		return TwitterError{"User", "JSON解析失败:" + err.Error()}, u
//	}

//	return nil, u
//}

func (this *Twitter) UserTimeLine(uid string, count int64, since_id int64) error {
	r, err := this.OauthClient.Get(http.DefaultClient, &this.OauthClient.Credentials,
		"http://api.twitter.com/1.1/statuses/home_timeline.json", nil)
	if err != nil {
		log.Fatal(err)
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
