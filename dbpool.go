package miorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DEFAULT_MAX_DB_CONNECTIONS = 100
	DEFAUL_IDLE_CONNECTIONS    = 10
)

type DBPool struct {
	idleSize int
	maxSize  int
	dbsource *DBAlias
	rwlock   sync.RWMutex
	dbtype   string
	dblink   string
	database string
}

func NewDBPool(dbtype string, dblink string, dbname string, idleSize int, maxSize int) *DBPool {
	if idleSize <= 0 {
		idleSize = DEFAUL_IDLE_CONNECTIONS
	}

	if maxSize <= 0 {
		maxSize = DEFAULT_MAX_DB_CONNECTIONS
	}

	if maxSize < idleSize {
		maxSize = 2 * idleSize
	}

	dbtype = strings.ToLower(dbtype)

	pool := &DBPool{
		dbtype:   dbtype,
		dblink:   dblink,
		database: dbname,
		idleSize: idleSize,
		maxSize:  maxSize}

	err := pool.initPool()
	if err == nil {
		return pool
	} else {
		fmt.Println(err.Error())
		return nil
	}

}

func (p *DBPool) ResizePool(idlesize int, maxsize int) {
	if maxsize < idlesize {
		maxsize = idlesize
	}
	p.dbsource.db.SetMaxIdleConns(idlesize)
	p.dbsource.db.SetMaxOpenConns(maxsize)

}

func (p *DBPool) initPool() error {
	return p.createDBSource()
}

func (p *DBPool) GetDB() *DBAlias {
	return p.dbsource

}

func (p *DBPool) createDBSource() error {

	switch p.dbtype {
	case "mysql":
		return p.createMysqlDB()
	default:
		return errors.New("unsupport sql type [" + p.dbtype + "]")
	}

}

func (p *DBPool) createMysqlDB() error {
	db, err := sql.Open("mysql", p.dblink)
	if err != nil {
		db.Close()
		return err
	}
	db.SetMaxIdleConns(p.idleSize)
	db.SetMaxOpenConns(p.maxSize)
	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	}

	p.dbsource = NewDBAlias(db)
	return nil
}

func (p *DBPool) Close() {
	//TODO
	if p.dbsource != nil {
		p.dbsource.Close()
	}

}
