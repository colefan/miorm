package miorm

import (
	"reflect"
)

// table info struct
type dbTable struct {
	TableName              string
	PrimaryKeys            []*dbColumn
	Bean                   interface{}
	Fields                 []*dbColumn
	Property2ColumnMapping map[string]string
	Field2PropertyMapping  map[string]string
	AutoId                 *dbColumn
}

type dbColumn struct {
	Column       string
	DbType       reflect.Kind
	Property     string
	Type         reflect.Kind
	IsPrimaryKey bool
	IsAutoId     bool
	IsExclude    bool
}

func newDbTable() *dbTable {
	t := &dbTable{}
	t.Fields = make([]*dbColumn, 0, 0)
	t.PrimaryKeys = make([]*dbColumn, 0, 0)
	t.Property2ColumnMapping = make(map[string]string)
	t.Field2PropertyMapping = make(map[string]string)

	return t
}

func newDbColumn() *dbColumn {
	c := &dbColumn{}
	return c
}
