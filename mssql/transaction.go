package mssql

import (
	"context"
	"database/sql"
	"errors"
)

// TxConn tx obj
type TxConn struct {
	tx        *sql.Tx
	fieldlist []string
}

//Begin transaction begin with default isolation level is dependent
func (mdb *MsSQL) Begin() (*TxConn, error) {
	var err error
	txConn := &TxConn{}
	txConn.tx, err = mdb.dbConn.Begin()
	printErrors(err)
	if err != nil {
		return nil, err
	}
	return txConn, nil
}

// BeginWithIsol transaction begin with custom isolation level is dependent
// LevelDefault 默认级别
// LevelReadUncommitted 读未提交
// LevelReadCommitted 读已提交
// LevelWriteCommitted 写已提交
// LevelRepeatableRead 可重复读
// LevelSnapshot 可读快照
// LevelSerializable 可串行化
// LevelLinearizable 可线性化
func (mdb *MsSQL) BeginWithIsol(isolLevel sql.IsolationLevel, readOnly bool) (*TxConn, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	// defer cancel()
	var err error
	txConn := &TxConn{}
	txConn.tx, err = mdb.dbConn.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: isolLevel,
		ReadOnly:  readOnly,
	})
	printErrors(err)
	if err != nil {
		return nil, err
	}
	return txConn, nil
}

//SetFields  set field name
func (txConn *TxConn) SetFields(fieldlist []string) {
	txConn.fieldlist = fieldlist
}

//Commit transaction commit
func (txConn *TxConn) Commit() error {
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		printErrors(err)
		return err
	}
	err := txConn.tx.Commit()
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

//Rollback transaction rollback
func (txConn *TxConn) Rollback() error {
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		printErrors(err)
		return err
	}
	err := txConn.tx.Rollback()
	printErrors(err)
	if err != nil {
		return err
	}
	return nil
}

func (txConn *TxConn) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		return 0, err
	}
	stmt, err := txConn.tx.Prepare(query)
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

//Update Update operation
func (txConn *TxConn) Update(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return txConn.stmtExec(query, update, args...)
}

//Insert Insert operation
func (txConn *TxConn) Insert(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return txConn.stmtExec(query, insert, args...)
}

//Delete Delete operation
func (txConn *TxConn) Delete(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return txConn.stmtExec(query, delete, args...)
}

// GetVal get single value by transaction
func (txConn *TxConn) GetVal(query string, args ...interface{}) (interface{}, error) {
	lastQuery = getQuery(query, args...)
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		return nil, err
	}
	stmt, err := txConn.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	var err2 error
	printErrors(err2)
	row := stmt.QueryRow(args...)
	var value interface{}
	err2 = row.Scan(&value)
	b, ok := value.([]byte)
	if ok {
		value = string(b)
	}
	return value, err2
}

// GetRow get single row data by transaction
func (txConn *TxConn) GetRow(query string, args ...interface{}) (map[string]interface{}, error) {
	lastQuery = getQuery(query, args...)
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		return nil, err
	}
	stmt, err := txConn.tx.Prepare(query)
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
	if txConn.fieldlist != nil && len(columns) != len(txConn.fieldlist) {
		err := errors.New(errorSetField)
		printErrors(err)
		return nil, err
	}
	var clos []string
	if txConn.fieldlist == nil {
		clos = columns
	} else {
		clos = txConn.fieldlist
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
				b, ok := column.([]byte)
				if ok {
					rowData[clos[k]] = string(b)
				} else {
					rowData[clos[k]] = column
				}
			} else {
				rowData[clos[k]] = nil
			}
		}
		break
	}
	txConn.fieldlist = nil
	return rowData, nil
}

// GetResults get multiple rows data by transaction
func (txConn *TxConn) GetResults(query string, args ...interface{}) ([]map[string]interface{}, error) {
	lastQuery = getQuery(query, args...)
	if txConn.tx == nil {
		err := errors.New(errorTxInit)
		return nil, err
	}
	stmt, err := txConn.tx.Prepare(query)
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
	if txConn.fieldlist != nil && len(columns) != len(txConn.fieldlist) {
		return nil, errors.New(errorSetField)
	}
	var clos []string
	if txConn.fieldlist == nil {
		clos = columns
	} else {
		clos = txConn.fieldlist
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
				b, ok := column.([]byte)
				if ok {
					rowData[clos[k]] = string(b)
				} else {
					rowData[clos[k]] = column
				}
			} else {
				rowData[clos[k]] = nil
			}
		}
		result = append(result, rowData)
	}
	txConn.fieldlist = nil
	return result, nil
}
