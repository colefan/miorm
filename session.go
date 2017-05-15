package miorm

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
	//"strings"
	//"strconv"
)

type Session struct {
	db           *Orm
	modelBean    interface{}
	table        *dbTable
	sqlbuilder   *SqlBuilder
	sqlDesc      string
	lastSql      string
	params       []interface{}
	insertAutoId string
}

func (s *Session) init() {
	//TODO
	s.params = make([]interface{}, 0, 0)
	s.sqlbuilder = NewSqlBuilder()
}

func (s *Session) assembleSelectSql(tableBean interface{}) string {
	sql := ""
	if s.sqlDesc != "" {
		sql = s.sqlDesc
	} else {
		if s.sqlbuilder.strSelect == "" {
			s.sqlbuilder.SELECT("*")
		}

		if s.sqlbuilder.strFrom == "" {
			s.sqlbuilder.FROM(s.table.TableName)
		}

		if s.sqlbuilder.strWhere == "" && len(s.table.PrimaryKeys) > 0 {
			//如果有－pk的值，则获取他的值为条件
			if false == isPKZero(tableBean, s.table) {
				s.sqlbuilder.WHERE(wherePK(tableBean, s.table))
			}

		}

		sql = s.sqlbuilder.SQL()

	}
	return sql
}

func (s *Session) assembleSqlParams(rawsql string) string {
	//#{paramname}
	//
	if s.modelBean == nil {
		return rawsql
	}

	reg := regexp.MustCompile(`#\{[0-9A-Za-z_]+\}`)
	params := reg.FindAllString(rawsql, -1)
	rawsql = reg.ReplaceAllString(rawsql, "?")

	beanv := reflect.ValueOf(s.modelBean).Elem()
	for _, param := range params {

		param = param[2 : len(param)-1]
		fmt.Println("param ,", param)
		fieldV := beanv.FieldByName(param)
		if fieldV.IsValid() {
			s.params = append(s.params, beanv.FieldByName(param).Interface())
		} else {

			panic("unknow property name: " + param)
		}

	}

	return rawsql

}

func (s *Session) RawResultSelectOne(tableBean interface{}) (records []map[string]interface{}, err error) {
	logstr := ""
	startms := int64(0)
	endms := int64(0)
	s.table = s.db.mappingTable(tableBean)
	s.modelBean = tableBean
	s.sqlbuilder.LIMIT(1, 0)
	rawsql := s.assembleSelectSql(tableBean)
	s.lastSql = s.assembleSqlParams(rawsql)
	//exec

	defer func() {
		if logstr != "" {
			endms = time.Now().UnixNano() / 1000
			s.db.logger.Debug(logstr, " exec[", endms-startms, "]us")
		}
	}()

	startms = time.Now().UnixNano() / 1000

	if s.db.ShowSQL {
		logstr = s.lastSql
		paramslog := ""
		for i, _ := range s.params {
			paramslog = fmt.Sprint(s.params[i]) + ","

		}

		if paramslog != "" {
			logstr += " -- params[" + paramslog + "]"
		}

	}
	//starttime = time.Now()
	//startms = time.Now().UnixNano() / 1000
	records, err = s.db.dbpool.GetDB().Query(s.lastSql, s.params...)
	return

}

func (s *Session) SelectOne(tableBean interface{}) (result interface{}, err error) {
	records, err := s.RawResultSelectOne(tableBean)
	if err != nil {
		//	endms = time.Now().UnixNano() / 1000
		return tableBean, err
	}

	if records == nil || len(records) == 0 {
		//	endms = time.Now().UnixNano()
		return tableBean, nil
	}
	//endms = time.Now().UnixNano()

	table2Struct(s.table, tableBean, &records[0])
	//endms = time.Now().UnixNano()
	return tableBean, nil
}

func (s *Session) RawResultSelect(tableBean interface{}) (records []map[string]interface{}, err error) {
	startms := int64(0)
	endms := int64(0)
	sqllog := ""
	s.table = s.db.mappingTable(tableBean)
	s.modelBean = tableBean
	rawsql := s.assembleSelectSql(tableBean)
	s.lastSql = s.assembleSqlParams(rawsql)
	//exec
	if s.db.ShowSQL {
		sqllog = s.lastSql
		if len(s.params) > 0 {
			sqllog += " params[" + fmt.Sprint(s.params...) + "]"
		}
	}

	defer func() {
		if sqllog != "" {
			endms = time.Now().UnixNano() / 1000
			s.db.logger.Debug(sqllog, " exec[", endms-startms, "]us")
		}
	}()

	startms = time.Now().UnixNano() / 1000
	records, err = s.db.dbpool.GetDB().Query(s.lastSql, s.params...)
	return

}

func (s *Session) Select(tableBean interface{}) (retslice []interface{}, err error) {
	records, err := s.RawResultSelect(tableBean)
	if err != nil {
		return
	}
	if records == nil || len(records) == 0 {
		return
	}

	retslice = make([]interface{}, len(records))
	beant := reflect.TypeOf(tableBean)
	if beant.Kind() == reflect.Ptr {
		beant = beant.Elem()
	}

	for i, _ := range records {
		newbean := reflect.New(beant).Interface()
		//logger.Debugf("newbean -> ", newbean)
		table2Struct(s.table, newbean, &records[i])
		retslice[i] = newbean
	}

	return
}

func (s *Session) RawSelect(sql string) (records []map[string]interface{}, err error) {
	//TODO
	startms := time.Now().UnixNano() / 1000
	records, err = s.db.dbpool.GetDB().Query(sql)
	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		s.db.logger.Debug(sql, " exec[", endms-startms, "]")

	}
	return
}

func (s *Session) Insert(tableBean interface{}) (affectedRows int64, err error) {
	//TODO
	s.table = s.db.mappingTable(tableBean)
	s.modelBean = tableBean
	if s.sqlDesc != "" {
		s.lastSql = s.assembleSqlParams(s.sqlDesc)
	} else {
		beanv := reflect.ValueOf(tableBean).Elem()

		strInsert := "insert into " + s.table.TableName + " ("
		strValues := ""
		for _, field := range s.table.Fields {
			v := beanv.FieldByName(field.Property).Interface()
			if field.IsAutoId {
				continue
			}
			strInsert += field.Column + ","
			strValues += "?,"
			s.params = append(s.params, v)
		}

		if strValues == "" {
			return 0, errors.New("Insert " + reflect.TypeOf(tableBean).String() + " is empty ,no data")
		}

		strInsert = strInsert[0:len(strInsert)-1] + ")"
		strValues = "(" + strValues[0:len(strValues)-1] + ")"

		s.lastSql = strInsert + " values " + strValues

	}
	//exec
	startms := time.Now().UnixNano() / 1000

	if s.db.TxSupport {
		if s.table.AutoId == nil {
			affectedRows, err = s.db.dbpool.GetDB().TxInsert(s.lastSql, s.params...)
		} else {
			id, rows, err1 := s.db.dbpool.GetDB().TxInsertWithAutoID(s.lastSql, s.params...)
			if err1 == nil {
				beanv := reflect.ValueOf(tableBean).Elem()
				beanv.FieldByName(s.table.AutoId.Property).SetInt(id)
			}
			affectedRows = rows
			err = err1
		}

	} else {
		if s.table.AutoId == nil {
			affectedRows, err = s.db.dbpool.GetDB().Insert(s.lastSql, s.params...)
		} else {
			id, rows, err1 := s.db.dbpool.GetDB().InsertWithAutoID(s.lastSql, s.params...)
			if err1 == nil {
				beanv := reflect.ValueOf(tableBean).Elem()
				beanv.FieldByName(s.table.AutoId.Property).SetInt(id)
			}
			affectedRows = rows
			err = err1
		}

	}

	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		if len(s.params) > 0 {
			paramstr := ""
			for _, p := range s.params {
				paramstr += fmt.Sprint(p, ",")
			}

			s.db.logger.Debug(s.lastSql, " params[", paramstr, "] exec[", endms-startms, "]us")

		} else {
			s.db.logger.Debug(s.lastSql, " exec[", endms-startms, "]us")
		}

	}

	return
}

func (s *Session) MultiInsert(tableBeanSlice []interface{}) (affectedRows int64, err error) {
	//TODO
	if tableBeanSlice == nil || len(tableBeanSlice) == 0 {
		err = errors.New("MulitInsert params tableBeanSlice is nil or empty")
		return
	}

	startms := time.Now().UnixNano() / 1000

	s.table = s.db.mappingTable(tableBeanSlice[0])

	paramsSlice := make([][]interface{}, 0)

	defer func() {
		if s.db.ShowSQL {
			endms := time.Now().UnixNano() / 1000
			s.db.logger.Debug(s.lastSql, " exec[", endms-startms, "]us")
		}

	}()

	if s.sqlDesc != "" {
		//assemble sql params
		reg := regexp.MustCompile(`#\{[0-9A-Za-z_]+\}`)
		params := reg.FindAllString(s.sqlDesc, -1)
		s.lastSql = reg.ReplaceAllString(s.sqlDesc, "?")

		if params != nil && len(params) > 0 {
			for _, bean := range tableBeanSlice {
				beanv := reflect.ValueOf(bean).Elem()
				beanParam := make([]interface{}, len(params))
				for i, param := range params {
					param = param[2 : len(param)-1]
					beanParam[i] = beanv.FieldByName(param).Interface()
				}
				paramsSlice = append(paramsSlice, beanParam)
			}
		}

	} else {

		strInsert := "insert into " + s.table.TableName + " ("
		strValues := ""
		paramnames := make([]string, 0)
		for _, field := range s.table.Fields {
			if field.IsAutoId {
				continue
			}
			strInsert += field.Column + ","
			strValues += "?,"
			paramnames = append(paramnames, field.Property)
		}

		if strValues == "" {
			return 0, errors.New("Insert " + s.table.TableName + " , beanslice is empty ,no data")
		}

		strInsert = strInsert[0:len(strInsert)-1] + ")"
		strValues = "(" + strValues[0:len(strValues)-1] + ")"

		s.lastSql = strInsert + " values " + strValues

		for _, bean := range tableBeanSlice {
			beanv := reflect.ValueOf(bean).Elem()
			beanParam := make([]interface{}, len(paramnames))
			for i, param := range paramnames {
				beanParam[i] = beanv.FieldByName(param).Interface()
			}
			paramsSlice = append(paramsSlice, beanParam)
		}

	}
	//exec
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxBatchExec(s.lastSql, paramsSlice)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().BatchExec(s.lastSql, paramsSlice)
	}

	return
}

func (s *Session) RawInsert(sql string) (affectedRows int64, err error) {
	//TODO
	starttime := time.Now().UnixNano() / 1000
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxInsert(sql)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().Insert(sql)
	}

	if s.db.ShowSQL {
		s.db.logger.Debug(sql, " exec[", time.Now().UnixNano()/1000-starttime, "]us")
	}

	return
}

func (s *Session) Delete(tableBean interface{}) (affectedRows int64, err error) {
	//TODO
	s.table = s.db.mappingTable(tableBean)
	s.modelBean = tableBean
	if s.sqlDesc != "" {
		s.lastSql = s.assembleSqlParams(s.sqlDesc)
	} else {
		if s.sqlbuilder.strWhere != "" {
			//采用builder中的where条件
			//
			rawsql := "delete from " + s.table.TableName + " where " + s.sqlbuilder.strWhere
			s.lastSql = s.assembleSqlParams(rawsql)
		} else {
			//采用主键的值作为where条件

			if isPKZero(tableBean, s.table) {
				return 0, errors.New("Delete " + reflect.TypeOf(tableBean).String() + ", primmary key data is empty,no where condtion")
			}

			wherepk := wherePK(tableBean, s.table)
			s.lastSql = "delete from " + s.table.TableName + " where " + wherepk
		}

	}
	//exec
	startms := time.Now().UnixNano() / 1000
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxExec(s.lastSql, s.params...)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().Exec(s.lastSql, s.params...)
	}

	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		if len(s.params) > 0 {
			paramslog := ""
			for _, p := range s.params {
				paramslog += fmt.Sprint(p, ",")
			}
			s.db.logger.Debug(s.lastSql, " params[", paramslog, "] exec[", endms-startms, "]us")
		} else {
			s.db.logger.Debug(s.lastSql, " exec[", endms-startms, "]us")
		}
	}

	return
}

func (s *Session) MultiDelete(tableBeanSlice []interface{}) (affectedRows int64, err error) {
	if tableBeanSlice == nil || len(tableBeanSlice) == 0 {
		affectedRows = 0
		err = errors.New("MultiDelete param tableBeanSlice is nil or empty")
		return
	}
	s.table = s.db.mappingTable(tableBeanSlice[0])
	paramsSlice := make([][]interface{}, 0)
	if s.sqlDesc != "" {
		reg := regexp.MustCompile(`#\{[0-9A-Za-z_]+\}`)
		params := reg.FindAllString(s.sqlDesc, -1)
		s.lastSql = reg.ReplaceAllString(s.sqlDesc, "?")
		if params != nil && len(params) > 0 {
			for _, bean := range tableBeanSlice {
				beanv := reflect.Indirect(reflect.ValueOf(bean))
				paramlist := make([]interface{}, len(params))
				for i, param := range params {
					param = param[2 : len(param)-1]
					paramlist[i] = beanv.FieldByName(param).Interface()

				}
				paramsSlice = append(paramsSlice, paramlist)
			}

		}

	} else {
		if s.sqlbuilder.strWhere != "" {
			s.lastSql = "delete from " + s.table.TableName + " where " + s.sqlbuilder.strWhere
			reg := regexp.MustCompile(`#\{[0-9A-Za-z_]+\}`)
			params := reg.FindAllString(s.lastSql, -1)
			s.lastSql = reg.ReplaceAllString(s.lastSql, "?")
			if params != nil && len(params) > 0 {
				for _, bean := range tableBeanSlice {
					beanv := reflect.Indirect(reflect.ValueOf(bean))
					paramlist := make([]interface{}, len(params))
					for i, param := range params {
						param = param[2 : len(param)-1]
						paramlist[i] = beanv.FieldByName(param).Interface()

					}
					paramsSlice = append(paramsSlice, paramlist)
				}

			}

		} else {
			s.lastSql = "delete from " + s.table.TableName
			strWhere := ""
			params := make([]string, 0)
			for _, pk := range s.table.PrimaryKeys {
				strWhere += pk.Column + " = ? and"
				params = append(params, pk.Property)
			}

			if strWhere != "" {
				strWhere = strWhere[0 : len(strWhere)-3]
				if params != nil && len(params) > 0 {
					for _, bean := range tableBeanSlice {
						beanv := reflect.Indirect(reflect.ValueOf(bean))
						paramlist := make([]interface{}, len(params))
						for i, param := range params {
							//param = param[2 : len(param)-1]
							paramlist[i] = beanv.FieldByName(param).Interface()

						}
						paramsSlice = append(paramsSlice, paramlist)
					}

				}
				s.lastSql += " where " + strWhere

			} else {
				err = errors.New("delete no where condtion")
				return

			}
		}

	}

	//exec
	//
	startms := time.Now().UnixNano() / 1000

	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxBatchExec(s.lastSql, paramsSlice)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().BatchExec(s.lastSql, paramsSlice)
	}

	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		s.db.logger.Debug(s.lastSql, " exec[", endms-startms, "]us")
	}

	return
}

func (s *Session) RawDelete(sql string) (affectedRows int64, err error) {
	startms := time.Now().UnixNano() / 1000
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxExec(sql)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().Exec(sql)
	}

	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		s.db.logger.Debug(sql, " exec[", endms-startms, "]us")
	}
	return
}

func (s *Session) Update(tableBean interface{}) (affectedRows int64, err error) {
	//TODO
	s.table = s.db.mappingTable(tableBean)
	s.modelBean = tableBean
	if s.sqlDesc != "" {
		s.lastSql = s.assembleSqlParams(s.sqlDesc)
	} else {

		strWhere := ""
		if s.sqlbuilder.strWhere != "" {
			strWhere = s.sqlbuilder.strWhere
		} else {
			if isPKZero(tableBean, s.table) {
				return 0, errors.New("Update " + reflect.TypeOf(tableBean).String() + " primary key data is empty,no where condtion")
			}
			strWhere = wherePK(tableBean, s.table)
		}

		rawsql := "update " + s.table.TableName + " set "
		strSetter := ""
		beanv := reflect.ValueOf(tableBean).Elem()
		for _, field := range s.table.Fields {
			if field.IsPrimaryKey {
				continue
			}
			if s.sqlbuilder.IsUpdateField(field.Column) {
				strSetter += field.Column + " = ? ,"
				s.params = append(s.params, beanv.FieldByName(field.Property).Interface())
			}

		}

		if strSetter != "" {
			strSetter = strSetter[0 : len(strSetter)-1]
			s.lastSql = rawsql + strSetter + " where " + strWhere

		} else {
			return 0, errors.New("Update " + reflect.TypeOf(tableBean).String() + " no field can be set")

		}

	}

	startms := time.Now().UnixNano() / 1000
	//exec
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxExec(s.lastSql, s.params...)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().Exec(s.lastSql, s.params...)
	}

	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		if len(s.params) > 0 {
			paramstr := ""
			for _, p := range s.params {
				paramstr += fmt.Sprint(p, ",")
			}
			s.db.logger.Debug(s.lastSql, " params[", paramstr, "] exec[", endms-startms, "]us")
		} else {
			s.db.logger.Debug(s.lastSql, " exec[", endms-startms, "]us")
		}

	}

	return
}

func (s *Session) RawUpdate(sql string) (affectedRows int64, err error) {
	startms := time.Now().UnixNano() / 1000
	if s.db.TxSupport {
		affectedRows, err = s.db.dbpool.GetDB().TxExec(sql)
	} else {
		affectedRows, err = s.db.dbpool.GetDB().Exec(sql)
	}
	if s.db.ShowSQL {
		endms := time.Now().UnixNano() / 1000
		s.db.logger.Debug(sql, " exec[", endms-startms, "]us")
	}
	return
}

func (s *Session) SqlProvider(sqldesc string) *Session {
	s.sqlDesc = sqldesc
	return s
}

func (s *Session) ResultAutoID(autoid string) *Session {
	s.insertAutoId = autoid
	return s
}

func (s *Session) Where(cond string) *Session {
	s.sqlbuilder.WHERE(cond)
	return s
}

func (s *Session) OrderBy(field string) *Session {
	s.sqlbuilder.ORDERBY(field)
	return s
}

func (s *Session) GroupBy(field string) *Session {
	s.sqlbuilder.GROUPBY(field)
	return s

}

func (s *Session) Limit(limit int, start int) *Session {
	s.sqlbuilder.LIMIT(limit, start)
	return s

}

func (s *Session) Having(field string) *Session {
	s.sqlbuilder.HAVING(field)
	return s
}

func (s *Session) SelectFields(selFields string) *Session {
	s.sqlbuilder.SELECT(selFields)
	return s
}

func (s *Session) UpdateInclude(fields string) *Session {
	s.sqlbuilder.UpdateInclude(fields)
	return s
}

func (s *Session) UpdateExclude(fields string) *Session {
	s.sqlbuilder.UpdateExclude(fields)
	return s
}
