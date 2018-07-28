package v2

import (
	"errors"
)

//GetRow get single row data
func (mdb *MysqlDB) GetRow(query string, param ...interface{}) (map[string]interface{}, error) {
	lastQuery = getQuery(query, param...)
	var rs map[string]interface{}
	var err error
	if cacheMode {
		value, found := checkCache()
		if found {
			return value.(map[string]interface{}), nil
		}
	}
	rs, err = mdb.stmtQueryRow(query, param...)
	setCache(rs)
	return rs, err
}

func (mdb *MysqlDB) stmtQueryRow(query string, param ...interface{}) (map[string]interface{}, error) {
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
	mdb.fieldlist = nil
	return rowData, nil
}
