package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb" //not use
	"github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

const (
	insert = iota
	update
	delete
)

var (
	//Statement mode
	Statement = 1
	//Normal mode
	Normal            = 0
	customColumns     []string
	msSQL             *MsSQL
	once              sync.Once
	dBHost            string
	dBPort            int
	dBName            string
	dBuser            string
	dBpass            string
	dbcharset         string
	isinit            bool
	maxIdleConns      int
	maxOpenConns      int
	lastQuery         string
	showErrors        = true
	cacheTimeout      = 5 * time.Second
	cacheMode         = false
	checkCacheTimeOut = 10 * time.Minute
	caches            *cache.Cache
)

var (
	errorInit     = "DB param is not initialize"
	errorSetField = "Field List is Error"
	errorTxInit   = "Transaction didn't initializtion"
)

//Init db params
func Init(Host string, Port int, Name, user, pass string, MaxIdleConns, MaxOpenConns int) {
	dBHost = Host
	dBPort = Port
	dBName = Name
	dBuser = user
	dBpass = pass
	isinit = true
	maxIdleConns = MaxIdleConns
	maxOpenConns = MaxOpenConns
}

//GetMsSQLConn create or get a database connection with singleton mode
func GetMsSQLConn() (*MsSQL, error) {
	if isinit == false {
		return nil, errors.New(errorInit)
	}
	var err error
	once.Do(func() {
		msSQL = &MsSQL{dBHost, dBuser, dBName, dBpass, dBPort, nil, nil}
		err = msSQL.conn(maxIdleConns, maxOpenConns)
		if err != nil {
			log.Panicln(err)
		}
	})
	return msSQL, err
}

//NewMysqlConn create a database connection
func NewMysqlConn(Host string, Port int, Name, User, Pass string, MaxIdleConns, MaxOpenConns int) (*MsSQL, error) {
	var err error
	msSQL = &MsSQL{Host, User, Name, Pass, Port, nil, nil}
	err = msSQL.conn(MaxIdleConns, MaxOpenConns)
	if err != nil {
		log.Panicln(err)
	}
	return msSQL, err
}

//MsSQL database
type MsSQL struct {
	host, user, dbname, pass string
	port                     int
	dbConn                   *sql.DB
	fieldlist                []string
}

func (mdb *MsSQL) conn(MaxIdleConns int, MaxOpenConns int) error {
	if mdb.host == "" || mdb.pass == "" || mdb.user == "" || mdb.dbname == "" {
		errs := errors.New(errorInit)
		printErrors(errs)
		return errs
	}
	lastQuery = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", mdb.host, mdb.user, mdb.pass, mdb.port, mdb.dbname)
	db, err := sql.Open("mssql", lastQuery)
	if err != nil {
		printErrors(err)
		return err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	err2 := db.Ping()
	if err2 != nil {
		printErrors(err2)
		return err2
	}
	mdb.dbConn = db
	return nil
}

//Close close db conn
func (mdb *MsSQL) Close() {
	mdb.dbConn.Close()
}

//SetFields  set field name
func (mdb *MsSQL) SetFields(fieldlist []string) {
	mdb.fieldlist = fieldlist
}
