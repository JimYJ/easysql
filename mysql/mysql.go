package simpleDB

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
)

var (
	Statement     int      = 1
	Normal        int      = 0
	Insert        int      = 0
	Update        int      = 1
	Delete        int      = 2
	customColumns []string = nil
)

type MysqlDB struct {
	host, user, dbname, pass string
	port                     int
	dbConn                   *sql.DB
	fieldlist                []string
	tx                       *sql.Tx
}

func (self *MysqlDB) Conn(MaxIdleConns int, MaxOpenConns int) error {
	if self.host == "" || self.pass == "" || self.user == "" || self.dbname == "" {
		errs := errors.New("DB param is not initialize!")
		return errs
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", self.user, self.pass, self.host, self.port, self.dbname))
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	err2 := db.Ping()
	if err2 != nil {
		return err2
	}
	self.dbConn = db
	return nil
}

func (self *MysqlDB) Close() {
	self.dbConn.Close()
}

func (self *MysqlDB) GetResults(qtype int, query string, param ...interface{}) ([]map[string]string, error) {
	if qtype == Statement {
		return self.stmtQuery(query, param...)
	} else {
		return self.query(query)
	}
}

func (self *MysqlDB) SetFields(fieldlist []string) {
	self.fieldlist = fieldlist
}

func (self *MysqlDB) query(query string) ([]map[string]string, error) {
	rows, err := self.dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		return nil, errors.New("Field List is Error!")
	}
	var clos []string
	if self.fieldlist == nil {
		clos = columns
	} else {
		clos = self.fieldlist
	}
	/* check custom field end*/
	columnName := make([]string, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	var result []map[string]string
	for rows.Next() {
		_ = rows.Scan(colbuff...)
		// common.ErrHendle("Scan Warning:", err)
		rowData := make(map[string]string, len(columns))
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		result = append(result, rowData)
	}
	self.fieldlist = nil
	return result, nil
}

func (self *MysqlDB) stmtQuery(query string, param ...interface{}) ([]map[string]string, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(param...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		return nil, errors.New("Field List is Error!")
	}
	var clos []string
	if self.fieldlist == nil {
		clos = columns
	} else {
		clos = self.fieldlist
	}
	/* check custom field end*/
	columnName := make([]string, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	var result []map[string]string
	for rows.Next() {
		_ = rows.Scan(colbuff...)
		// common.ErrHendle("Scan Warning:", err)
		rowData := make(map[string]string, len(columns))
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		result = append(result, rowData)
	}
	self.fieldlist = nil
	return result, nil
}

func (self *MysqlDB) getVal(query string) (string, error) {
	row := self.dbConn.QueryRow(query)
	var str string
	err := row.Scan(&str)
	// common.ErrHendle("Scan Warning:", err)
	return str, err
}

func (self *MysqlDB) getValByStmt(query string, param ...interface{}) (string, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	var str string
	err2 := row.Scan(&str)
	return str, err2
}

func (self *MysqlDB) GetVal(qtype int, query string, param ...interface{}) (string, error) {
	if qtype == Statement {
		return self.getValByStmt(query, param...)
	} else {
		return self.getVal(query)
	}
}

func (self *MysqlDB) GetRow(qtype int, query string, param ...interface{}) (map[string]string, error) {
	if qtype == Statement {
		return self.stmtQueryRow(query, param...)
	} else {
		return self.queryRow(query)
	}
}

func (self *MysqlDB) queryRow(query string) (map[string]string, error) {
	rows, err := self.dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		return nil, errors.New("Field List is Error!")
	}
	var clos []string
	if self.fieldlist == nil {
		clos = columns
	} else {
		clos = self.fieldlist
	}
	/* check custom field end*/
	columnName := make([]string, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	rowData := make(map[string]string, len(columns))
	for rows.Next() {
		_ = rows.Scan(colbuff...)
		// common.ErrHendle("Scan Warning:", err)
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		break
	}
	self.fieldlist = nil
	return rowData, nil
}

func (self *MysqlDB) stmtQueryRow(query string, param ...interface{}) (map[string]string, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(param...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		return nil, errors.New("Field List is Error!")
	}
	var clos []string
	if self.fieldlist == nil {
		clos = columns
	} else {
		clos = self.fieldlist
	}
	/* check custom field end*/
	columnName := make([]string, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	rowData := make(map[string]string, len(columns))
	for rows.Next() {
		_ = rows.Scan(colbuff...)
		// common.ErrHendle("Scan Warning:", err)
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		break
	}
	self.fieldlist = nil
	return rowData, nil
}

func (self *MysqlDB) exec(query string, qtype int, args ...interface{}) (int64, error) {
	rs, err := self.dbConn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update || qtype == Delete {
		result, err2 = rs.RowsAffected()
		log.Println(result)
	}
	self.fieldlist = nil
	return result, err2
}

func (self *MysqlDB) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return 0, err
	}
	rs, err := stmt.Exec(args...)
	if err != nil {
		if qtype == Insert {
			log.Println("Insert Fail:", err)
			return 0, err
		} else if qtype == Update {
			log.Println("Update Fail:", err)
			return 0, err
		} else if qtype == Delete {
			log.Println("Delete Fail:", err)
			return 0, err
		}
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update || qtype == Delete {
		result, err2 = rs.RowsAffected()
		log.Println(result)
	}
	return result, err2
}

func (self *MysqlDB) TxBegin() error {
	var err error
	self.tx, err = self.dbConn.Begin()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) TxCommit() error {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return err
	}
	err := self.tx.Commit()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) TxRollback() error {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return err
	}
	err := self.tx.Rollback()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) txExec(query string, qtype int, args ...interface{}) (int64, error) {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return 0, err
	}
	rs, err := self.tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update {
		result, err2 = rs.RowsAffected()
		log.Println(result)
	}
	self.fieldlist = nil
	return result, err2
}

func (self *MysqlDB) stmtTxExec(query string, qtype int, args ...interface{}) (int64, error) {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return 0, err
	}
	stmt, err := self.tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	rs, err := stmt.Exec(args...)
	if err != nil {
		if qtype == Insert {
			log.Println("Insert Fail:", err)
			return 0, err
		} else if qtype == Update {
			log.Println("Update Fail:", err)
			return 0, err
		} else if qtype == Delete {
			log.Println("Delete Fail:", err)
			return 0, err
		}
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update || qtype == Delete {
		result, err2 = rs.RowsAffected()
		log.Println(result)
	}
	return result, err2
}

func (self *MysqlDB) Update(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Update, args...)
	} else {
		return self.exec(query, Update, args...)
	}
}

func (self *MysqlDB) Insert(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Insert, args...)
	} else {
		return self.exec(query, Insert, args...)
	}
}

func (self *MysqlDB) Delete(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Delete, args...)
	} else {
		return self.exec(query, Delete, args...)
	}
}

func (self *MysqlDB) TxUpdate(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Update, args...)
	} else {
		return self.txExec(query, Update, args...)
	}
}

func (self *MysqlDB) TxInsert(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Insert, args...)
	} else {
		return self.txExec(query, Insert, args...)
	}
}

func (self *MysqlDB) TxDelete(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Delete, args...)
	} else {
		return self.txExec(query, Delete, args...)
	}
}
