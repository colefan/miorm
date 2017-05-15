package miorm

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type MyUser struct {
	orm_table_name string `table:"user"`
	Id             int    `orm:"pk;auto"`
	Username       string
	D1             time.Time
	D2             time.Time
	F              float32
	IsBoy          bool
	Dd             time.Time
}

func TestType(t *testing.T) {
	t1, err := time.Parse(date_Layout, "2105-11-07")
	fmt.Println("t1 = ", t1)
	t2, err := time.Parse(datetime_Layout, "2015-11-01 12:12:12")
	fmt.Println("t2 = ", t2, ", error = ", err)
}

func TestReflect(t *testing.T) {
	u1 := new(MyUser)
	u1.Id = 1
	fmt.Println(u1)
	fmt.Println("u1-type =", reflect.TypeOf(u1))

	ut := reflect.TypeOf(u1).Elem()
	fmt.Println("ut = ", ut.String())

	newu := reflect.New(ut).Interface()

	fmt.Println("newu : ", newu)
	fmt.Println("newu:type ", reflect.TypeOf(newu))
	fmt.Println("u1:", u1)
}

func TestMyOrm(t *testing.T) {
	o := NewOrm("mysql", "root:@/myorm?charset=utf8", "myorm", true)
	defer o.Close()

	u := &MyUser{Id: 2}
	u2, _ := o.SelectOne(u)
	fmt.Println("u2 :", u2)
	fmt.Println("u :", u)
	list, err := o.Select(&MyUser{})
	if err != nil {
		fmt.Println("list error ", err.Error())
	} else {
		fmt.Println("list len ", len(list))
		for i, item := range list {
			newu := item.(*MyUser)
			fmt.Println("list[", i, "] ", newu.Id, newu.Username, newu.D1, newu.D2, newu.Dd, newu.IsBoy)
		}
	}
}
