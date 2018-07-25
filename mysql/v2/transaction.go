package v2

import (
	"context"
	"database/sql"
	"errors"
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
func (mdb *MysqlDB) TxBeginWithIsol(isolLevel sql.IsolationLevel, readOnly bool) error {
	var err error
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	// defer cancel()
	mdb.tx, err = mdb.dbConn.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: isolLevel,
		ReadOnly:  readOnly,
	})
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

func (mdb *MysqlDB) stmtTxExec(query string, qtype int, args ...interface{}) (int64, error) {
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		return 0, err
	}
	stmt, err := mdb.tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
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
func (mdb *MysqlDB) TxUpdate(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtTxExec(query, update, args...)
}

//TxInsert Insert operation
func (mdb *MysqlDB) TxInsert(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtTxExec(query, insert, args...)
}

//TxDelete Delete operation
func (mdb *MysqlDB) TxDelete(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtTxExec(query, delete, args...)
}

// TxGetVal get single value by transaction
func (mdb *MysqlDB) TxGetVal(query string, args ...interface{}) (interface{}, error) {
	lastQuery = getQuery(query, args...)
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		return "", err
	}
	stmt, err := mdb.tx.Prepare(query)
	if err != nil {
		return "", err
	}
	var err2 error
	printErrors(err2)
	row := stmt.QueryRow(args...)
	var value interface{}
	err2 = row.Scan(&value)
	return value, err2
}

// TxGetRow get single row data by transaction
func (mdb *MysqlDB) TxGetRow(query string, args ...interface{}) (map[string]interface{}, error) {
	lastQuery = getQuery(query, args...)
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		return nil, err
	}
	stmt, err := mdb.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		printErrors(err)
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if mdb.fieldlist != nil && len(columns) != len(mdb.fieldlist) {
		err := errors.New(errorSetField)
		printErrors(err)
		return nil, err
	}
	var clos []string
	if mdb.fieldlist == nil {
		clos = columns
	} else {
		clos = mdb.fieldlist
	}
	/* check custom field end*/
	columnName := make([]interface{}, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	rowData := make(map[string]interface{}, len(columns))
	for rows.Next() {
		err := rows.Scan(colbuff...)
		printErrors(err)
		for k, column := range columnName {
			if column != nil {
				rowData[clos[k]] = column
			} else {
				rowData[clos[k]] = nil
			}
		}
		break
	}
	mdb.fieldlist = nil
	return rowData, nil
}

// TxGetResults get multiple rows data by transaction
func (mdb *MysqlDB) TxGetResults(query string, args ...interface{}) ([]map[string]interface{}, error) {
	lastQuery = getQuery(query, args...)
	if mdb.tx == nil {
		err := errors.New(errorTxInit)
		return nil, err
	}
	stmt, err := mdb.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	printErrors(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	printErrors(err)
	if err != nil {
		return nil, err
	}
	/* check custom field*/
	if mdb.fieldlist != nil && len(columns) != len(mdb.fieldlist) {
		return nil, errors.New(errorSetField)
	}
	var clos []string
	if mdb.fieldlist == nil {
		clos = columns
	} else {
		clos = mdb.fieldlist
	}
	/* check custom field end*/
	columnName := make([]interface{}, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(colbuff...)
		printErrors(err)
		rowData := make(map[string]interface{}, len(columns))
		for k, column := range columnName {
			if column != nil {
				rowData[clos[k]] = column
			} else {
				rowData[clos[k]] = nil
			}
		}
		result = append(result, rowData)
	}
	mdb.fieldlist = nil
	return result, nil
}
