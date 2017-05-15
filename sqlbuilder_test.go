package miorm

import (
	"fmt"
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	builder := NewSqlBuilder()
	str := builder.SELECT("*").FROM("tabename").WHERE("a=#{a}").WHERE("b=#{b}").OR().WHERE("c=#{c}").ORDERBY("a").ORDERBY("b").LIMIT(1, 20).SQL()
	fmt.Println(str)
}
