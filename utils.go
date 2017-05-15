package miorm

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func isZero(k interface{}) bool {
	switch k.(type) {
	case int:
		return k.(int) == 0
	case uint:
		return k.(uint) == 0
	case int16:
		return k.(int16) == 0
	case uint16:
		return k.(uint16) == 0
	case int64:
		return k.(int64) == 0
	case uint64:
		return k.(int64) == 0
	case int8:
		return k.(int8) == 0
	case uint8:
		return k.(uint8) == 0
	case int32:
		return k.(int32) == 0
	case uint32:
		return k.(uint32) == 0
	case string:
		return k.(string) == ""
	case bool:
		return k.(bool) == false
	case time.Time:
		return k.(time.Time).IsZero()
	}
	return false
}

func isPKZero(bean interface{}, tableMeta *dbTable) bool {
	beanv := reflect.ValueOf(bean).Elem()
	for _, v := range tableMeta.PrimaryKeys {
		fv := beanv.FieldByName(v.Property).Interface()
		if isZero(fv) {
			return true
		}
	}
	return false
}

func wherePK(bean interface{}, tableMeta *dbTable) string {
	beanv := reflect.ValueOf(bean).Elem()
	strWhere := ""
	for _, v := range tableMeta.PrimaryKeys {
		fvalue := beanv.FieldByName(v.Property)
		if strWhere == "" {
			strWhere = v.Column + " = " + reflect2SqlWhereValue(&fvalue)
		} else {
			strWhere += " and " + v.Column + " = " + reflect2SqlWhereValue(&fvalue)

		}
	}

	return strWhere
}

func reflect2SqlWhereValue(rawValue *reflect.Value) (str string) {
	aa := (*rawValue).Type()
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt((*rawValue).Int(), 10)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		str = strconv.FormatUint((*rawValue).Uint(), 10)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat((*rawValue).Float(), 'f', -1, 64)
	case reflect.String:
		str = "\"" + (*rawValue).String() + "\""
	case reflect.Bool:
		str = strconv.FormatBool((*rawValue).Bool())

	}
	return
}

func table2Struct(tableMeta *dbTable, bean interface{}, tableDataPtr *map[string]interface{}) bool {
	for k, v := range *tableDataPtr {
		prop := tableMeta.Field2PropertyMapping[k]
		//fmt.Println("k = ", k, " prop = ", prop)
		if prop != "" {
			bv := reflect.ValueOf(bean).Elem()
			fv := bv.FieldByName(prop)
			ft := fv.Type().Kind()
			switch ft {
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				switch v.(type) {
				case []uint8:
					uintv, _ := strconv.ParseUint(string(v.([]byte)), 10, 64)
					fv.SetUint(uintv)
				case int64:
					fv.SetUint(uint64(v.(int64)))
				default:
					panic("unsupported sql scan interface : " + reflect.TypeOf(v).String())

				}

			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
				switch v.(type) {
				case []uint8:
					intv, _ := strconv.ParseInt(string(v.([]byte)), 10, 64)
					fv.SetInt(intv)
				case int64:
					fv.SetInt(v.(int64))
				default:
					panic("unsupported sql san interface : " + reflect.TypeOf(v).String())
				}

			case reflect.Float32, reflect.Float64:
				floatv, _ := strconv.ParseFloat(string(v.([]byte)), 64)
				fv.SetFloat(floatv)
			case reflect.String:
				fv.SetString(string(v.([]byte)))
			case reflect.Bool:
				if v.([]byte)[0] == 48 {
					fv.SetBool(false)
				} else {
					fv.SetBool(true)
				}
			case reflect.Ptr:
			case reflect.Struct:
				if fv.Type() == reflect.TypeOf(t_TIME) {
					timedata := string(v.([]byte))
					if len(timedata) == date_Layout_len {
						t, _ := time.Parse(date_Layout, timedata)

						fv.Set(reflect.ValueOf(t))

					} else if len(timedata) == datetime_Layout_len {

						t, _ := time.Parse(datetime_Layout, timedata)

						fv.Set(reflect.ValueOf(t))
					}
					//fv.SetBytes(v)

				} else {
					fmt.Println("unsupported type :", fv.Type().String())

				}

			}

		}
	}
	return true
}
