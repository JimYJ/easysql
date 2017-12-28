[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)]

simpleDB encapsulated database operation, simplifying the use of the database.(include Mysql database,more database is coming soon)

# How to get

```
go get github.com/JimYJ/simpleDB/simplemysql
```

# Usage

import:

```go
import "github.com/JimYJ/simpleDB"
```

get value:

```go
mysql := simpleDB.GetMysqlConn("0.0.0.0", 3306, "DBname", "root", "123456")
value,err := mysql.GetVal("SELECT count(*) FROM users")
```

get single row data:

```go
mysql := simpleDB.GetMysqlConn("0.0.0.0", 3306, "DBname", "root", "123456")
value,err := mysql.GetRow(simplemysql.Normal,"SELECT name,email FROM users WHERE id = 2")
or
value,err := mysql.GetRow(simplemysql.Normal,"SELECT name,email FROM users WHERE id = ?",2)
```

get single row data with statement:
```go
mysql := simpleDB.GetMysqlConn("0.0.0.0", 3306, "DBname", "root", "123456")
row,err := mysql.GetRow(simplemysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
```

