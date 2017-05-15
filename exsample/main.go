package main

import (
	"fmt"
	"github.com/colefan/myorm"
	"time"
)

func main() {
	o := myorm.NewOrm("mysql", "root:@/myorm?charset=utf8", "myorm", true)
	if o == nil {
		return
	}
	user1 := &User{Id: 1, UserName: "yjx"}
	_, err := o.Where("Username=#{UserName}").SelectOne(user1)
	fmt.Println("user1 :", user1, "--", err)

	userlist, _ := o.Select(&User{})
	for i, user := range userlist {
		fmt.Println("list[", i, "] ", user)
	}

	u2 := &User{UserName: "y1", D1: time.Now(), D2: time.Now(), IsBoy: false, Dd: time.Now(), Score: 5.5}
	_, err = o.Insert(u2)
	o.SqlProvider("insert into user(username,isboy) values('yyy',false) ").Insert(&User{})
	fmt.Println("insert u2 :", u2, " err :", err)

	user1.Score = 100.1
	_, err = o.Update(user1)
	fmt.Println("update user1 = ", user1)

	user1.Score = 200.3

	_, err = o.SqlProvider("update user set f=#{Score} where id=#{Id}").Update(user1)

	ul := make([]interface{}, 0)
	for i := 0; i < 10; i++ {
		ul = append(ul, &User{UserName: "yjx00" + fmt.Sprint(i), D1: time.Now(), D2: time.Now(), IsBoy: true, Dd: time.Now(), Score: float32(i)})

	}

	rows, err := o.MultiInsert(ul)
	fmt.Println("multiinsert ul , rows = ", rows, " err = ", err)

	rows, err = o.Where("id in(100,200)").Delete(&User{Id: 11})
	fmt.Println("delete user ,rows = ", rows, " err = ", err)

	rows, err = o.MultiDelete(userlist[0 : len(userlist)/2])
	fmt.Println("multidelete ul ,rows = ", rows, " err = ", err)

	ru, err := o.SelectOne(&User{})
	fmt.Println("ru selectOne :", ru, " err = ", err)
	upuser := ru.(*User)
	upuser.UserName = "updateyjx"
	rows, err = o.Update(upuser)
	fmt.Println("update rows = ", rows, " err = ", err)

	o.Close()

	now := time.Now()
	fmt.Println("ss:", now.Unix())
	fmt.Println("ms:", now.UnixNano()/1000)
	fmt.Println("ns:", now.UnixNano())

	o = myorm.NewOrm("mysql", "root:@/inova?charset=utf8", "inova", true)
	//o.RawSelect("select * from pub_user")

	pu1 := &PubUser{Username: "admin"}
	pu1.Nickname = "admin"
	pu1.Id = 1
	pub2, err := o.Where("id=#{Id}").SelectOne(pu1)
	fmt.Println(pu1)
	fmt.Println(pub2)

}

type PubUser struct {
	orm_table_name string `table:"pub_user"`
	Id             int    `orm:"pk;auto"`
	Username       string
	Nickname       string
	Email          string
	LoginPwd       string `orm:"column=login_pwd"`
	AdminPwd       string `orm:"column=admin_pwd"`
	Status         int16
	Level          int16
	Createtime     time.Time
	Modifytime     time.Time
}

type User struct {
	Id       int `orm:"auto;pk"`
	UserName string
	D1       time.Time
	D2       time.Time
	IsBoy    bool
	Dd       time.Time
	Score    float32 `orm:"column=f"`
}
