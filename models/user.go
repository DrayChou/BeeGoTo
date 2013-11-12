package models

import (
	"BeeGoTo/models/douban"
	"fmt"
	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	Douban   douban.DoubanUser
}

func UserCreate() {
	// connect to MongoDB
	mgodb, err := mgo.Dial(beego.AppConfig.String("database_url"))
	if err != nil {
		fmt.Println("connect MongoDB failed...")
		os.Exit(1)
	}
	// close MongoDB
	defer mgodb.Close()
	// set MongoDB model
	mgodb.SetMode(mgo.Monotonic, true)
	// connect DB and Collection in MongoDB
	conn := mgodb.DB("test").C("people")
	// insert data into MongoDB
	err = conn.Insert(&Person{"hahaya", "123456"}, &Person{"sf", "111111"})
	if err != nil {
		fmt.Println("insert into MongoDB failed...")
		os.Exit(1)
	}
	result := Person{}
	//query MongoDB
	err = conn.Find(bson.M{"name": "hahaya"}).One(&result)
	if err != nil {
		fmt.Println("query MongoDB failed...")
		os.Exit(1)
	}
	// output the query result
	fmt.Println("phone", result.Phone)
}
