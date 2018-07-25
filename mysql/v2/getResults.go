package v2

import (
	"errors"
)

//GetResults get multiple rows data
func (mdb *MysqlDB) GetResults(qtype int, query string, param ...interface{}) ([]map[string]interface{}, error) {
	lastQuery = getQuery(query, param...)
	var rs []map[string]interface{}
	var err error
	if cacheMode {
		value, found := checkCache()
		if found {
			return value.([]map[string]interface{}), nil
		}
	}
	rs, err = mdb.stmtQuery(query, param...)
	setCache(rs)
	return rs, err
}

func (mdb *MysqlDB) stmtQuery(query string, param ...interface{}) ([]map[string]interface{}, error) {
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
