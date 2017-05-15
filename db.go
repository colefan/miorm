package miorm

import (
	"database/sql"
	//"fmt"
	//"reflect"
	//	"strconv"
	"strings"
)

type DBAlias struct {
	db *sql.DB
}

func NewDBAlias(db *sql.DB) *DBAlias {
	return &DBAlias{db: db}
}

func (d *DBAlias) Close() {
	d.db.Close()
}

func (d *DBAlias) Query(sql string, args ...interface{}) (records []map[string]interface{}, err error) {
	//TODO

	rows, err := d.db.Query(sql, args...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return
	}

	cols, _ := rows.Columns()
	sanArgs := make([]interface{}, len(cols))
	values := make([]interface{}, len(cols))
	for i, _ := range values {
		sanArgs[i] = &values[i]
	}

	for rows.Next() {
		record := make(map[string]interface{})
		rows.Scan(sanArgs...)
		for i, col := range values {

			if col != nil {
				//fmt.Println("col name = ", cols[i])
				//fmt.Println(" col name = ", cols[i], " , value = ", col, "type =", reflect.TypeOf(col))

				record[strings.ToLower(cols[i])] = col
			}

		}

		if records == nil {
			records = make([]map[string]interface{}, 0)
		}
		records = append(records, record)

	}

	return
}

func (d *DBAlias) InsertWithAutoID(sql string, args ...interface{}) (lastid int64, affectedRows int64, err error) {
	result, err := d.db.Exec(sql, args...)
	if err != nil {
		return
	}

	lastid, _ = result.LastInsertId()
	affectedRows, _ = result.RowsAffected()
	return

}

func (d *DBAlias) Insert(sql string, args ...interface{}) (affectedRows int64, err error) {
	result, err := d.db.Exec(sql, args...)
	if err != nil {
		return
	}
	affectedRows, _ = result.RowsAffected()
	return
}

func (d *DBAlias) TxInsert(sql string, args ...interface{}) (affectedRows int64, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}

	result, err := tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	affectedRows, _ = result.RowsAffected()
	return

}

func (d *DBAlias) TxInsertWithAutoID(sql string, args ...interface{}) (lastid int64, affectedRows int64, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}
	result, err := tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	lastid, _ = result.LastInsertId()
	affectedRows, _ = result.RowsAffected()
	return
}

func (d *DBAlias) BatchExec(sqlbuf string, argsv [][]interface{}) (affectedRows int64, err error) {
	var result sql.Result
	if argsv == nil || len(argsv) == 0 {
		result, err = d.db.Exec(sqlbuf)
		if err == nil {
			affectedRows, err = result.RowsAffected()
		}
		return
	} else {
		var stmt *sql.Stmt
		stmt, err = d.db.Prepare(sqlbuf)
		if err != nil {
			return
		}

		defer stmt.Close()
		affectedRows = 0
		for i, _ := range argsv {
			r, err1 := stmt.Exec(argsv[i]...)
			if err1 != nil {
				err = err1
				return
			}
			num, _ := r.RowsAffected()
			affectedRows += num
		}

	}

	return

}

func (d *DBAlias) TxBatchExec(sqlbuf string, argsv [][]interface{}) (affectedRows int64, err error) {
	var tx *sql.Tx
	tx, err = d.db.Begin()
	if err != nil {
		return
	}

	if argsv == nil || len(argsv) == 0 {
		r, err1 := tx.Exec(sqlbuf)
		if err1 != nil {
			err = err1
			tx.Rollback()
			return
		} else {
			err = tx.Commit()
			affectedRows, _ = r.RowsAffected()
			return
		}

	} else {
		var stmt *sql.Stmt
		stmt, err = tx.Prepare(sqlbuf)
		if err != nil {
			tx.Rollback()
			return
		}
		defer stmt.Close()
		for i, _ := range argsv {
			r, err1 := stmt.Exec(argsv[i]...)
			if err1 != nil {
				tx.Rollback()
				affectedRows = 0
				err = err1
				return
			}

			num, _ := r.RowsAffected()
			affectedRows += num

		}
		err = tx.Commit()

	}
	return

}

func (d *DBAlias) Exec(sql string, args ...interface{}) (affectedRows int64, err error) {
	result, err := d.db.Exec(sql, args...)
	if err != nil {
		return
	}

	affectedRows, err = result.RowsAffected()
	return

}

func (d *DBAlias) TxExec(sqlbuf string, args ...interface{}) (affectedRows int64, err error) {
	var tx *sql.Tx
	tx, err = d.db.Begin()
	if err != nil {
		return
	}

	var result sql.Result
	result, err = tx.Exec(sqlbuf, args...)
	if err != nil {
		err = tx.Rollback()
		return
	}
	err = tx.Commit()

	affectedRows, _ = result.RowsAffected()

	return

}
