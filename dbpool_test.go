package miorm

import (
	"fmt"
	"testing"
)

func TestDBPool(t *testing.T) {
	pool := NewDBPool("mysql", "root:@/myorm?charset=utf8", "miorm", 10, 20)
	pool.initPool()
	defer pool.Close()
	records, err := pool.GetDB().Query("select * from user limit 0,1")
	if err != nil {
		t.Fail()
	} else {
		fmt.Println("records len = ", len(records))
		for i, tmp := range records {
			fmt.Println("tmp = ", tmp)
			for k, v := range tmp {
				fmt.Println("[", i, "] , key = ", k, ", value = ", v)
			}
		}
	}

}
