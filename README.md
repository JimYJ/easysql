[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)

simpleDB encapsulated database operation, simplifying the use of the database.(include Mysql database,more database is coming soon)

# How to get

```
go get github.com/JimYJ/simpleDB/mysql
```

# Usage

**import:**

```go
import "github.com/JimYJ/simpleDB/mysql"
```

**conn db:**
```go
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
mysqldb, err := mysql.GetMysqlConn()
```

**get value:**

```go
value,err := mysqldb.GetVal("SELECT count(*) FROM users")
```

**get single row data:**
```go
row,err := mysqldb.GetRow(mysql.Normal,"SELECT name,email FROM users WHERE id = 2")
```

**get single row data with statement:**
```go
row,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
```

**get multi-rows data:**
```go
rows,err := mysqldb.GetResults(mysql.Normal,"SELECT name,email FROM users where type = 'public'")
```

**get multi-rows data with statement:**
```go
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```


**If you do not want to expose the database field name,you can set field name:**
```go
mysql.SetFields([]string{"username", "useremail"})
rows,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
or
mysql.SetFields([]string{"username", "useremail"})
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```


