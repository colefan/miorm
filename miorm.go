package miorm

import (
	"os"
	"reflect"
	"strings"
)

var logger *DefaultLogger

const (
	orm_dbpool_idle_size = 10
	orm_dbpool_max_size  = 50
)

type Orm struct {
	ShowSQL     bool
	TxSupport   bool
	tableMapper map[string]*dbTable
	logger      *DefaultLogger

	dbpool *DBPool
}

func NewOrm(dbtype string, dblinkurl string, dbname string, supporttx bool) *Orm {
	o := &Orm{}
	o.ShowSQL = true
	o.TxSupport = supporttx
	o.tableMapper = make(map[string]*dbTable)
	o.dbpool = NewDBPool(dbtype, dblinkurl, dbname, orm_dbpool_idle_size, orm_dbpool_max_size)
	o.logger = NewDefaultLogger(os.Stdout)
	o.logger.SetLevel(LOG_DEBUG)
	err := o.dbpool.initPool()
	if err != nil {
		o.logger.Err(err)
		return nil
	}
	return o

}

func (o *Orm) SetDbPoolCache(idleSize int, maxSize int) {
	if maxSize < idleSize {
		o.logger.Warning("unsuppoted size params ,idlesize ", idleSize, " is smaller than maxsize ", maxSize, " ")
	}
	o.dbpool.ResizePool(idleSize, maxSize)
}

// create a session for user
func (o *Orm) NewSession() *Session {
	s := &Session{db: o}
	s.init()
	return s
}

//select * from table limit 1
func (o *Orm) SelectOne(tableBean interface{}) (interface{}, error) {
	s := o.NewSession()
	return s.SelectOne(tableBean)
}

//select * from table
func (o *Orm) Select(tableBean interface{}) ([]interface{}, error) {
	s := o.NewSession()
	return s.Select(tableBean)
}

func (o *Orm) RawSelect(sql string) ([]map[string]interface{}, error) {
	s := o.NewSession()
	return s.RawSelect(sql)
}

func (o *Orm) Insert(tableBean interface{}, fields ...interface{}) (affecteRows int64, err error) {
	s := o.NewSession()
	return s.Insert(tableBean, fields...)
}

func (o *Orm) MultiInsert(tableBeanSlice []interface{}) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.MultiInsert(tableBeanSlice)
}

func (o *Orm) RawInsert(sql string) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.RawInsert(sql)
}

func (o *Orm) Delete(tableBean interface{}) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.Delete(tableBean)
}

func (o *Orm) MultiDelete(tableBeanSlice []interface{}) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.MultiDelete(tableBeanSlice)

}

func (o *Orm) RawDelete(sql string) (int64, error) {
	s := o.NewSession()
	return s.RawDelete(sql)
}

func (o *Orm) Update(tableBean interface{}) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.Update(tableBean)
}

func (o *Orm) RawUpdate(sql string) (affectedRows int64, err error) {
	s := o.NewSession()
	return s.RawUpdate(sql)
}

//sqlprovider
//
func (o *Orm) SqlProvider(sqldesc string) *Session {
	s := o.NewSession()
	return s.SqlProvider(sqldesc)

}

func (o *Orm) ResultAutoID(autoid string) *Session {
	s := o.NewSession()
	return s.ResultAutoID(autoid)

}

func (o *Orm) Where(cond string) *Session {
	s := o.NewSession()
	s.Where(cond)
	return s
}

func (o *Orm) OrderBy(field string) *Session {
	s := o.NewSession()
	s.OrderBy(field)
	return s

}

func (o *Orm) GroupBy(field string) *Session {
	s := o.NewSession()
	s.GroupBy(field)
	return s
}

func (o *Orm) Limit(limit int, start int) *Session {
	s := o.NewSession()
	s.Limit(limit, start)
	return s
}

func (o *Orm) Having(field string) *Session {
	s := o.NewSession()
	s.Having(field)
	return s
}

func (o *Orm) SelectFields(selFields string) *Session {
	s := o.NewSession()
	s.SelectFields(selFields)
	return s
}

func (o *Orm) mappingTable(tableBean interface{}) *dbTable {
	table := o.tableMapper[reflect.TypeOf(tableBean).String()]
	if table == nil {
		table = newDbTable()
		t := reflect.TypeOf(tableBean)
		tfields := t.Elem()
		for i := 0; i < tfields.NumField(); i++ {
			fname := strings.ToLower(tfields.Field(i).Name)
			ftag := tfields.Field(i).Tag.Get("orm")
			if fname == "orm_table_name" {
				//parse table name
				tablename := tfields.Field(i).Tag.Get("table")
				if tablename != "" {
					table.TableName = tablename
				}
			} else {
				//parse field
				column := newDbColumn()
				column.Property = tfields.Field(i).Name
				column.Type = tfields.Field(i).Type.Kind()

				if ftag == "" {
					column.Column = strings.ToLower(column.Property)
				} else {

					o.parseOrmFieldTag(ftag, column)

				}

				if false == column.IsExclude {

					if column.IsPrimaryKey {
						table.PrimaryKeys = append(table.PrimaryKeys, column)
					}

					if column.IsAutoId {
						table.AutoId = column
					}

					table.Fields = append(table.Fields, column)
					table.Property2ColumnMapping[column.Property] = column.Column
					table.Field2PropertyMapping[column.Column] = column.Property
				}

			}

		}

		if table.TableName == "" {
			if t.Kind() == reflect.Ptr {
				table.TableName = t.Elem().Name()
			} else {
				table.TableName = t.Name()
			}
		}

		o.tableMapper[reflect.TypeOf(tableBean).String()] = table
	}

	return table

}

// parse field tag formatter
func (o *Orm) parseOrmFieldTag(tag string, column *dbColumn) {
	splittag := strings.Split(tag, ";")
	for _, s := range splittag {
		if s == "auto" {
			column.IsAutoId = true
		} else if s == "pk" {
			column.IsPrimaryKey = true
		} else if strings.HasPrefix(s, "column=") && len(s) > 7 {
			column.Column = s[7:]
		} else if s == "exclude" {
			column.IsExclude = true
			break
		}
	}

	if column.Column == "" {
		column.Column = strings.ToLower(column.Property)
	}

}

// update include fields, split with ','
func (o *Orm) UpdateIncludeFields(fields string) *Session {
	s := o.NewSession()
	s.UpdateInclude(fields)
	return s
}

//update exclude fields, split with ','
func (o *Orm) UpdateExcludeFields(fields string) *Session {
	s := o.NewSession()
	s.UpdateExclude(fields)
	return s
}

func (o *Orm) Close() {
	if o.dbpool != nil {
		o.dbpool.Close()
	}
}

func init() {
	logger = NewDefaultLogger(os.Stdout)

}
