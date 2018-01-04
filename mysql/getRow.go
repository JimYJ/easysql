package mysql

import (
	"errors"
)

func (self *MysqlDB) GetRow(qtype int, query string, param ...interface{}) (map[string]string, error) {
	lastQuery = query
	if qtype == Statement {
		return self.stmtQueryRow(query, param...)
	} else {
		return self.queryRow(query)
	}
}

func (self *MysqlDB) queryRow(query string) (map[string]string, error) {
	rows, err := self.dbConn.Query(query)
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
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		err := errors.New(errorSetField)
		printErrors(err)
		return nil, err
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
		err := rows.Scan(colbuff...)
		printErrors(err)
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
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		err := errors.New(errorSetField)
		printErrors(err)
		return nil, err
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
		err := rows.Scan(colbuff...)
		printErrors(err)
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		break
	}
	self.fieldlist = nil
	return rowData, nil
}
