package miorm

// Copyright 2015 - 2017 The myorm Authors. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.
// Author mail:colefan@126.com

/*
Package myorm is a simple and powful orm framwork for go.

Installation

Make sure you have installed Go 1.1+ and then:

    go get -u github.com/colefan/myorm

1. new orm instance
   o := myorm.NewOrm("mysql","root:@/myorm?charset=utf8","myorm",true)

2. table to struct bean

type User struct{
	orm_table_name string `table:"user"`
	Id int `orm:"auto;pk"`
	Name string `orm:"column=username"`
	Pwd string
}

tag description:
1) orm_table_name string `table:"user"`  //special field and tag for table name,if without this field ,tablename is same as the struct name
2) `orm:"auto;pk;column=aa"`
3)	`orm:"ref(one);column=fid"`
4)`orm:"exclude"` //exclude this field, not a column of table

3. orm operations

user1,err := o.SelectOne(&User{})
user1 ,err:= o.SelectOne(&User{Id:1})

o.SqlProvider(“insert into User(name,pwd) values (#{name},#{pwd})").ResultAutoID("").Insert( &User{name:"a",pwd:"12345"})

o.Where("username=#{Name}").Update(&User{Name:"yjx"});

o.Where("id >100").Delete(&User{});

o.Where("id in (1,2)").Select(&User{});

o.RawSelect(sql)

o.RawUpdate(sql)

o.RawDelete(sql)

o.MultiInsert([]interface{})

o.MultiDelete([]interface{})

4. other select condtions
o.Where(cond).GroupBy(field).Having(field).OrderBy(field).Limit(start,limit).Select(&User{})

*/
