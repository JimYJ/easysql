package mysql

import (
	"errors"
)

func (mdb *MysqlDB) TxBegin() error {
	var err error
	mdb.tx, err = mdb.dbConn.Begin()
	printErrors(err)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (mdb *MysqlDB) TxCommit() error {
	if mdb.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		printErrors(err)
		return err
	}
	err := mdb.tx.Commit()
	printErrors(err)
	if err != nil {
		return err
	} else {
		return nil
	}
}

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
	} else {
		return nil
	}
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

func (mdb *MysqlDB) TxUpdate(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtTxExec(query, update, args...)
	} else {
		return mdb.txExec(query, update, args...)
	}
}

func (mdb *MysqlDB) TxInsert(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtTxExec(query, insert, args...)
	} else {
		return mdb.txExec(query, insert, args...)
	}
}

func (mdb *MysqlDB) TxDelete(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtTxExec(query, delete, args...)
	} else {
		return mdb.txExec(query, delete, args...)
	}
}
