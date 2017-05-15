package miorm

import (
	"strconv"
	"strings"
)

type SqlBuilder struct {
	strSelect     string
	strFrom       string
	strWhere      string
	strOrderBy    string
	strLimit      string
	strWhereLink  string
	strHaving     string
	strGroupBy    string
	updateInclude []string
	updateExclude []string
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{strWhereLink: " and "}
}

func (s *SqlBuilder) SELECT(param string) *SqlBuilder {
	s.strSelect = param
	return s
}

func (s *SqlBuilder) FROM(param string) *SqlBuilder {
	s.strFrom = param
	return s
}

func (s *SqlBuilder) WHERE(param string) *SqlBuilder {
	if s.strWhere == "" {
		s.strWhere = param
	} else {
		s.strWhere += s.strWhereLink + param

	}
	s.strWhereLink = " and "
	return s
}

func (s *SqlBuilder) OR() *SqlBuilder {
	s.strWhereLink = " or "
	return s
}

func (s *SqlBuilder) ORDERBY(param string) *SqlBuilder {
	if s.strOrderBy == "" {
		s.strOrderBy = param
	} else {
		s.strOrderBy += "," + param
	}
	return s
}

func (s *SqlBuilder) LIMIT(limit int, start int) *SqlBuilder {
	s.strLimit = strconv.Itoa(start) + "," + strconv.Itoa(limit)
	return s
}

func (s *SqlBuilder) GROUPBY(field string) *SqlBuilder {
	if s.strGroupBy == "" {
		s.strGroupBy = field
	} else {
		s.strGroupBy = s.strGroupBy + "," + field
	}
	return s
}

func (s *SqlBuilder) HAVING(field string) {
	if s.strHaving == "" {
		s.strHaving = field
	} else {
		s.strHaving = s.strHaving + "," + field
	}

}

func (s *SqlBuilder) SQL() string {
	sql := ""
	sql += "select " + s.strSelect
	sql += " from " + s.strFrom
	if s.strWhere != "" {
		sql += " where (" + s.strWhere + ")"
	}

	if s.strGroupBy != "" {
		sql += " group by " + s.strGroupBy + " "
		if s.strHaving != "" {
			sql += " having " + s.strHaving + " "
		}
	}

	if s.strOrderBy != "" {
		sql += " order by " + s.strOrderBy
	}

	if s.strLimit != "" {
		sql += " limit " + s.strLimit
	}

	return sql
}

func (s *SqlBuilder) UpdateInclude(fields string) {
	if fields != "" {
		v := strings.Split(fields, ",")
		for _, tmp := range v {
			s.updateInclude = append(s.updateInclude, strings.ToLower(tmp))
		}
	}
}

func (s *SqlBuilder) UpdateExclude(fields string) {
	if fields != "" {
		v := strings.Split(fields, ",")
		for _, tmp := range v {
			s.updateExclude = append(s.updateExclude, strings.ToLower(tmp))
		}
	}
}

func (s *SqlBuilder) IsUpdateField(field string) bool {
	if len(s.updateInclude) == 0 && len(s.updateExclude) == 0 {
		return true
	}

	field = strings.ToLower(field)
	if len(s.updateExclude) > 0 {
		for _, v := range s.updateExclude {
			if v == field {
				return false
			}
		}

	}

	if len(s.updateInclude) > 0 {
		for _, v := range s.updateInclude {
			if v == field {
				return true
			}
		}
		return false
	}

	return true
}
