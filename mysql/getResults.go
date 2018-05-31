package mysql

import (
	"errors"
)

//GetResults get multiple rows data
func (mdb *MysqlDB) GetResults(qtype int, query string, param ...interface{}) ([]map[string]string, error) {
	lastQuery = getQuery(query, param...)
	var rs []map[string]string
	var err error
	if cacheMode {
		value, found := checkCache()
		if found {
			return value.([]map[string]string), nil
		}
	}
	if qtype == Statement {
		rs, err = mdb.stmtQuery(query, param...)
		setCache(rs)
		return rs, err
	}
	rs, err = mdb.query(query)
	setCache(rs)
	return rs, err
}

func (mdb *MysqlDB) query(query string) ([]map[string]string, error) {
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
		errs := errors.New(errorSetField)
		printErrors(errs)
		return nil, errs
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
	mdb.fieldlist = nil
	return result, nil
}

func (mdb *MysqlDB) stmtQuery(query string, param ...interface{}) ([]map[string]string, error) {
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
	var result []map[string]string
	for rows.Next() {
		err := rows.Scan(colbuff...)
		printErrors(err)
		rowData := make(map[string]string, len(columns))
		for k, column := range columnName {
			if column != nil {
				rowData[clos[k]] = column.(string)
			} else {
				rowData[clos[k]] = ""
			}
		}
		result = append(result, rowData)
	}
	mdb.fieldlist = nil
	return result, nil
}
