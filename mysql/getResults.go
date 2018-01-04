package mysql

import (
	"errors"
)

func (self *MysqlDB) GetResults(qtype int, query string, param ...interface{}) ([]map[string]string, error) {
	lastQuery = query
	if qtype == Statement {
		return self.stmtQuery(query, param...)
	} else {
		return self.query(query)
	}
}

func (self *MysqlDB) query(query string) ([]map[string]string, error) {
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
		errs := errors.New(errorSetField)
		printErrors(errs)
		return nil, errs
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
	if self.fieldlist != nil && len(columns) != len(self.fieldlist) {
		return nil, errors.New(errorSetField)
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
		err := rows.Scan(colbuff...)
		printErrors(err)
		rowData := make(map[string]string, len(columns))
		for k, column := range columnName {
			rowData[clos[k]] = string(column)
		}
		result = append(result, rowData)
	}
	self.fieldlist = nil
	return result, nil
}
