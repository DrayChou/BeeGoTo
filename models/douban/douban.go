package douban

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"io/ioutil"
	"net/url"
)

type User struct {
	Id           string
	Uid          string
	Name         string
	Screen_name  string
	Loc_id       string
	Loc_name     string
	Type         string
	Alt          string
	Signature    string
	Desc         string
	Description  string
	Avatar       string
	Small_avatar string
	Large_avatar string
}

type Attachments struct {
	Type         string
	Title        string
	Description  string
	Href         string
	Original_src string
}

type TimeLine struct {
	Id          int
	Layout      int
	Title       string
	Text        string
	Created_at  string
	User        User
	Attachments []Attachments
}

type Return struct {
	Code    int
	Msg     string
	Request string
}

type DoubanError struct {
	prefix string
	msg    string
}

func (oe DoubanError) Error() string {
	fmt.Println("DoubanAPIError: " + oe.prefix + ": " + oe.msg)
	beego.Error("DoubanAPIError: " + oe.prefix + ": " + oe.msg)
	return "DoubanError: " + oe.prefix + ": " + oe.msg
}

type Douban struct {
	Conf      string
	Transport *oauth.Transport
}

func (this *Douban) AuthUrl() (error, string) {
	if this.Transport.Config == nil {
		return DoubanError{"AuthUrl", "豆瓣配置异常"}, ""
	}
	url := this.Transport.Config.AuthCodeURL("")
	fmt.Println("url:", url)
	return nil, url
}

func (this *Douban) Auth(uid string, code string) error {
	if uid == "" {
		return DoubanError{"Auth", "內網用戶ID不能爲空"}
	}

	beego.Debug("User_id:", uid)
	beego.Debug("Code:", code)

	dbconf, err := config.NewConfig("ini", this.Conf)
	if err != nil {
		return DoubanError{"Auth", "配置文件加载失败"}
	}

	if code == "" {
		code = dbconf.String("code")
	}

	Config := &oauth.Config{
		ClientId:     dbconf.String("clientId"),
		ClientSecret: dbconf.String("clientSecret"),
		RedirectURL:  dbconf.String("redirectURL"),
		Scope:        dbconf.String("scope"),
		AuthURL:      dbconf.String("authURL"),
		TokenURL:     dbconf.String("tokenURL"),
		TokenCache:   oauth.CacheFile(dbconf.String("cacheDir") + uid + ".json"),
	}
	beego.Debug("Config:", Config)

	this.Transport = &oauth.Transport{Config: Config}
	beego.Debug("this.Transport:", this.Transport.Config)

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := Config.TokenCache.Token()
	beego.Debug("token:", token)

	if err != nil {

		if token == nil {
			if code != "" {

				token, err = this.Transport.Exchange(code)
				fmt.Println("token:", token)

				if err != nil {
					return DoubanError{"Auth", "令牌換取失敗:," + err.Error()}
				}

				// (The Exchange method will automatically cache the token.)
				fmt.Println("Token is cached in %v\n", Config.TokenCache)
				return nil
			} else {
				return DoubanError{"Auth", "令牌ID出错"}
			}
		}
	}
	this.Transport.Token = token

	return nil
}

func (this *Douban) Refresh() error {
	return this.Transport.Refresh()
}

func (this *Douban) User(uid string) (error, User) {

	u := User{}
	if uid == "" {
		if this.Transport.Token == nil {
			return DoubanError{"User", "未授权，请先授权"}, u
		}
		uid = "~me"
	}
	request_url := "https://api.douban.com/v2/user/" + uid

	r, err := this.Transport.Client().Get(request_url)
	if err != nil {
		return DoubanError{"User", "请求失败:" + err.Error()}, u
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	beego.Debug("User:Request StatusCode:", r.StatusCode)
	beego.Debug("User:Request Body:", string(body))

	err = json.Unmarshal(body, &u)
	if err != nil {
		return DoubanError{"User", "JSON解析失败:" + err.Error()}, u
	}

	return nil, u
}

func (this *Douban) UserTimeLine(uid string, count int64, since_id int64) (error, []TimeLine) {
	tl := []TimeLine{}

	if uid == "" {
		return DoubanError{"UserTimeLine", "参数错误"}, tl
	}

	request_url := "https://api.douban.com/shuo/v2/statuses/user_timeline/" + uid + "?"
	if count != 0 {
		request_url = request_url + "&count=" + string(count)
	}
	if since_id != 0 {
		request_url = request_url + "&since_id=" + string(since_id)
	}

	// Make the request.
	r, err := this.Transport.Client().Get(request_url)
	if err != nil {
		return DoubanError{"UserTimeLine", "请求失败:" + err.Error()}, tl
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	beego.Debug("UserTimeLine:Request StatusCode:", r.StatusCode)
	beego.Debug("UserTimeLine:Request Body:", string(body))

	err = json.Unmarshal(body, &tl)
	if err != nil {
		fmt.Println("Unmarshal Error:", err)
		return DoubanError{"UserTimeLine", "JSON解析失败:" + err.Error()}, tl
	}

	return nil, tl
}

func (this *Douban) Shuo(text string) (error, TimeLine) {
	tl := TimeLine{}

	if this.Transport.Token == nil {
		return DoubanError{"Shuo", "未授权，请先授权"}, tl
	}

	if text == "" {
		return DoubanError{"Shuo", "参数错误"}, tl
	}
	request_url := "https://api.douban.com/shuo/v2/statuses/"

	v := url.Values{}
	v.Set("source", this.Transport.ClientId)
	v.Set("text", text)
	r, err := this.Transport.Client().PostForm(request_url, v)
	if err != nil {
		return DoubanError{"Shuo", "请求失败:" + err.Error()}, tl
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	beego.Debug("Shuo:Request StatusCode:", r.StatusCode)
	beego.Debug("Shuo:Request Body:", string(body))

	if err = json.Unmarshal(body, &tl); r.StatusCode != 200 || err != nil {
		return DoubanError{"Shuo", "JSON解析失败:" + err.Error()}, tl
	}

	return nil, tl
}
