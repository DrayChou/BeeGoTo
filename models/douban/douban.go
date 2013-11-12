package douban

import (
	"code.google.com/p/goauth2/oauth"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"io/ioutil"
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

type UserTime struct {
	Id          int
	Layout      int
	Title       string
	Text        string
	Created_at  string
	User        User
	Attachments Attachments
}

type Return struct {
	Code    int
	Msg     string
	Request string
}

type Douban struct {
	User     User
	TimeLine []UserTime
	Return   Return
}

var DBConfig = &oauth.Config{
	ClientId:     dbconf.String("douban_clientId"),
	ClientSecret: dbconf.String("douban_clientSecret"),
	RedirectURL:  dbconf.String("douban_redirectURL"),
	Scope:        dbconf.String("douban_scope"),
	AuthURL:      dbconf.String("douban_authURL"),
	TokenURL:     dbconf.String("douban_tokenURL"),
	TokenCache:   oauth.CacheFile(dbconf.String("douban_cachefile") + md5_name + ".json"),
}
var DBtransport = &oauth.Transport{Config: DBConfig}

func (this *Douban) Auth(uid string, code string) (status bool, url string, msg string) {
	if uid == "" {
		return false, "", "內網用戶ID不能爲空"
	}

	h := md5.New()
	h.Write([]byte(uid)) // 需要加密的字符串为 sharejs.com
	md5_name := hex.EncodeToString(h.Sum(nil))

	fmt.Println("User_id:", uid)
	fmt.Println("MD5:", md5_name)
	fmt.Println("Code:", code)

	dbconf, err := config.NewConfig("ini", "conf/douban.conf")
	if err != nil {
		fmt.Println("NewConfig Error:", err)
		return false, "", "讀取配置文件出錯"
	}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := DBConfig.TokenCache.Token()
	if err != nil {

		if token == nil {
			if code != "" {

				token, err = DBtransport.Exchange(code)
				fmt.Println("token:", token)

				if err != nil {
					fmt.Println("Exchange:", err)
					return false, "", "令牌換取失敗"
				}

				// (The Exchange method will automatically cache the token.)
				fmt.Println("Token is cached in %v\n", DBConfig.TokenCache)

				return true, "", "令牌更新成功"
			}
			return false, "", ""
		} else {
			url := DBConfig.AuthCodeURL("")
			return true, url, "成功"
		}
	}
	DBtransport.Token = token

	return true, "", "令牌加載成功"
}

func User(uid string) (status bool, user User) {
	var user douban.User

	if uid == "" {
		uid = "~me"
	}
	url := "https://api.douban.com/v2/user/" + uid

	// Make the request.
	r, err := DBtransport.Client().Get(url)
	if err != nil {
		fmt.Println("Request Error:", err)
		return false, user
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	fmt.Println("Request Return:", string(body))

	err = json.Unmarshal(body, &line)
	if err != nil {
		fmt.Println("Unmarshal Error:", err)
	}
	fmt.Printf("%+v", user)

	return true, user
}
