package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

//TxBegin transaction begin with default isolation level is dependent
func (mdb *MysqlDB) TxBegin() error {
	var err error
	mdb.tx, err = mdb.dbConn.Begin()
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

// TxBeginWithIsol transaction begin with custom isolation level is dependent
// LevelDefault 默认级别
// LevelReadUncommitted 读未提交
// LevelReadCommitted 读已提交
// LevelWriteCommitted 写已提交
// LevelRepeatableRead 可重复读
// LevelSnapshot 可读快照
// LevelSerializable 可串行化
// LevelLinearizable 可线性化
func (mdb *MysqlDB) TxBeginWithIsol(opts *sql.TxOptions) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()
	mdb.tx, err = mdb.dbConn.BeginTx(ctx, opts)
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

//TxCommit transaction commit
func (mdb *MysqlDB) TxCommit() error {
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		printErrors(err)
		return err
	}
	err := mdb.tx.Commit()
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

//TxRollback transaction rollback
func (mdb *MysqlDB) TxRollback() error {
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		printErrors(err)
		return err
	}
	err := mdb.tx.Rollback()
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

func (mdb *MysqlDB) txExec(query string, qtype int, args ...interface{}) (int64, error) {
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		printErrors(err)
		return 0, err
	}
	rs, err := mdb.tx.Exec(query, args...)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == update {
		result, err2 = rs.RowsAffected()
	}
	mdb.fieldlist = nil
	printErrors(err2)
	return result, err2
}

func (mdb *MysqlDB) stmtTxExec(query string, qtype int, args ...interface{}) (int64, error) {
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		return 0, err
	}
	stmt, err := mdb.tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	rs, err := stmt.Exec(args...)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == update || qtype == delete {
		result, err2 = rs.RowsAffected()
	}
	printErrors(err2)
	return result, err2
}

//TxUpdate Update operation
func (mdb *MysqlDB) TxUpdate(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	if qtype == Statement {
		return mdb.stmtTxExec(query, update, args...)
	}
	return mdb.txExec(query, update, args...)
}

//TxInsert Insert operation
func (mdb *MysqlDB) TxInsert(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	if qtype == Statement {
		return mdb.stmtTxExec(query, insert, args...)
	}
	return mdb.txExec(query, insert, args...)
}

//TxDelete Delete operation
func (mdb *MysqlDB) TxDelete(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	if qtype == Statement {
		return mdb.stmtTxExec(query, delete, args...)
	}
	return mdb.txExec(query, delete, args...)
}
