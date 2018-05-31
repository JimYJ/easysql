package mysql

import (
	"errors"
)

//GetRow get single row data
func (mdb *MysqlDB) GetRow(qtype int, query string, param ...interface{}) (map[string]string, error) {
	lastQuery = getQuery(query, param...)
	var rs map[string]string
	var err error
	if cacheMode {
		value, found := checkCache()
		if found {
			return value.(map[string]string), nil
		}
	}
	if qtype == Statement {
		rs, err = mdb.stmtQueryRow(query, param...)
		setCache(rs)
		return rs, err
	}
	rs, err = mdb.queryRow(query)
	setCache(rs)
	return rs, err
}

func (mdb *MysqlDB) queryRow(query string) (map[string]string, error) {
	rows, err := mdb.dbConn.Query(query)
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
	rowData := make(map[string]string, len(columns))
	for rows.Next() {
		err := rows.Scan(colbuff...)
		printErrors(err)
		for k, column := range columnName {
			if column != nil {
				rowData[clos[k]] = anyToString(column)
			} else {
				rowData[clos[k]] = ""
			}

		}
		break
	}
	mdb.fieldlist = nil
	return rowData, nil
}

func (mdb *MysqlDB) stmtQueryRow(query string, param ...interface{}) (map[string]string, error) {
	stmt, err := mdb.dbConn.Prepare(query)
	printErrors(err)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(param...)
	printErrors(err)
	if err != nil {
		return nil, err
	}
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
	columnName := make([]string, len(columns))
	colbuff := make([]interface{}, len(columns))
	for i := range colbuff {
		colbuff[i] = &columnName[i]
	}
	rowData := make(map[string]string, len(columns))
	for rows.Next() {
		err := rows.Scan(colbuff...)
		printErrors(err)
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		break
	}
	mdb.fieldlist = nil
	return rowData, nil
}
