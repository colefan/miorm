package miorm

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

//sample of orm bean struct
type User struct {
	orm_table_name string `table:"user_1"`
	Id             uint8  `orm:"pk;auto;column=user_id"`
	Name           string
	fkid           interface{} `orm:"ref(one)"`
}

func TestOrmTags(t *testing.T) {
	var f interface{}
	u := &User{}
	u.Id = 12
	f = u
	stype := reflect.TypeOf(f)

	fmt.Println("stype:", stype.String())
	fmt.Println("stype name:", stype.Name())
	//fmt.Println(stype.Kind())
	stypelem := stype.Elem()
	fmt.Println("stypelem:", stypelem.Name())
	fmt.Println(stypelem)
	for i := 0; i < stypelem.NumField(); i++ {
		fmt.Println("filed-name:", stypelem.Field(i).Name)
		fmt.Println("fieldTag:", stypelem.Field(i).Tag)
		fmt.Println("filed-type:", stypelem.Field(i).Type.Kind())

	}

	//field := &dbColumn{}
	//field.SqlType = reflect.Uint64
	v := reflect.ValueOf(f)
	fmt.Println("v = ", v)

	fmt.Println("v name =", v.Elem().FieldByName("Id").Type())
	vf := v.Elem().FieldByName("Id").Interface()
	fmt.Println("v =", vf)
	v.Elem().FieldByName("Id").SetUint(13)
	fmt.Println("v id = ", u.Id)
}

func TestRegexp(t *testing.T) {
	reg := regexp.MustCompile(`#\{[0-9A-Za-z_]+\}`)
	str := "insert user (a,b,c) values (#{ab_3},#{b_1},#{cded_23}) "
	slist := reg.FindAllString(str, -1)
	for _, s := range slist {
		fmt.Println("find reg :", s)
	}
}
