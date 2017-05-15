package miorm

import (
	"time"
)

//const values for reflect.TypeOf
var (
	t_STRING string
	t_TIME   time.Time
)

var (
	datetime_Layout     = "2006-01-02 15:04:05"
	datetime_Layout_len = len("2006-01-02 15:04:05")
	date_Layout         = "2006-01-02"
	date_Layout_len     = len("2006-01-02")
)
