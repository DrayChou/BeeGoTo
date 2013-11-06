package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type Person struct {
	Name  string
	Phone string
}

func MgoTest() {
	// connect to MongoDB
	session, err := mgo.Dial(beego.AppConfig.String("database_url"))
	if err != nil {
		fmt.Println("connect MongoDB failed...")
		os.Exit(1)
	}
	// close MongoDB
	defer session.Close()
	// set MongoDB model
	session.SetMode(mgo.Monotonic, true)
	// connect DB and Collection in MongoDB
	conn := session.DB("test").C("people")
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
