[![Build Status](https://travis-ci.org/JimYJ/easysql.svg?branch=master)](https://travis-ci.org/JimYJ/easysql)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)
[English](https://github.com/JimYJ/easysql/blob/master/README.md) 

easysql 封装了数据库操作，简化了数据库的使用，目前已支持Mysql，今后将会支持更多数据库

## 更新:
1.2 增加缓存功能，应对高并发请求下的重复请求<br>
1.1 增加debug模式和debug的功能<br>
1.0 完成基础功能<br>

# 获取和安装

```
go get github.com/JimYJ/easysql/mysql
```

# Mysql驱动（使用go get获取安装）
```go
github.com/Go-SQL-Driver/MySQL
```

# 使用

**引用本库:**

```go
import "github.com/JimYJ/easysql/mysql"
```

**连接数据库:** 
Init()只需要初始化一次，Getmysqldb()为并发安全的单例模式，可以多次调用，推荐使用，考虑到多数据库的连接，Newmysqldb()没有做保护，请谨慎调用，建议再封装一层
```go
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", MaxIdleConns, MaxOpenConns)//数据库ip，端口，数据库名，用户，密码，最大空闲数，最大连接数
mysqldb, err := mysql.Getmysqldb()//singleton pattern
or
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123", MaxIdleConns, MaxOpenConns)
```

**关闭连接:**
```go
mysqldb.Close()
```

**设置缓存失效时间**
```go
mysql.SetCacheTimeout(5 * time.Second)//默认超时时间为5秒，设置缓存超时必须在开启缓存之前，不然设置的时间不会生效
```
**启用缓存:**
```go
mysql.UseCache()
```

**关闭缓存:**
```go
mysql.CloseCache()
```



**获取值:**
```go
value,err := mysqldb.GetVal(mysql.Normal,"SELECT count(*) FROM users")
```

**Stmt模式获取值（推荐）:**
```go
value,err := mysqldb.GetVal(mysql.Statement,"SELECT count(*) FROM users")
or
value,err := mysqldb.GetVal(mysql.Statement,"SELECT count(*) FROM users where type = ?","public")
```

**获取单行数据:**
```go
row,err := mysqldb.GetRow(mysql.Normal,"SELECT name,email FROM users WHERE id = 2")
```

**Stmt模式获取单行数据（推荐）:**
```go
row,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
```

**获取多行数据:**
```go
rows,err := mysqldb.GetResults(mysql.Normal,"SELECT name,email FROM users where type = 'public'")
```

**Stmt模式获取多行数据（推荐）:**
```go
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```


**如果是开发接口，不希望暴露数据库的字段名，可以返回的自定义字段名，要和数据库取值顺序一致:**
```go
mysql.SetFields([]string{"username", "useremail"})
row,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
or
mysql.SetFields([]string{"username", "useremail"})
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```

**插入数据:**
```go
insertId, err := mysqldb.Insert(mysql.Normal, "insert into users set name = ?", "jim")
```


**Stmt模式插入数据（推荐）:**
```go
insertId, err := mysqldb.Insert(mysql.Statement, "insert into users set name = ?", "jim")
```

**更新数据:**
```go
rowsAffected, err := mysqldb.Update(mysql.Normal, "update users set name = ? where id =?", "jim", 1)
```

**Stmt模式更新数据（推荐）:**
```go
rowsAffected, err := mysqldb.Update(mysql.Statement, "update users set name = ? where id =?", "jim", 1)
```

**删除数据:**
```go
rowsAffected, err := mysqldb.Delete(mysql.Normal, "delete from users where id =?", 453)
```

**Stmt模式删除数据（推荐）:**
```go
rowsAffected, err := mysqldb.Delete(mysql.Statement, "delete from users where id =?", 453)
```

**事务操作:**
```go
mysqldb.TxBegin()
insertId, err := mysqldb.TxInsert(mysql.Normal, "insert into users set name = ?", "jim")
rowsAffected, err := mysqldb.TxUpdate(mysql.Normal, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqldb.TxDelete(mysql.Normal, "delete from users where id =?", 453)
mysqldb.TxRollback()
or
mysqldb.TxCommit()
```

**Stmt模式事务操作（推荐）:**
```go
mysqldb.TxBegin()
insertId, err := mysqldb.TxInsert(mysql.Statement, "insert into users set name = ?", "jim")
rowsAffected, err := mysqldb.TxUpdate(mysql.Statement, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqldb.TxDelete(mysql.Statement, "delete from users where id =?", 453)
mysqldb.TxRollback()
or
mysqldb.TxCommit()
```

**调试:**
Debug()会打印最新操作的SQL
```go
mysql.Debug()
```

**打印所有错误信息**
```go
mysql.ShowErrors()
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
mysqldb, err := mysql.Getmysqldb()
...
or
mysql.ShowErrors()
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
```

**关闭打印错误信息**
```go
mysql.HideErrors()
```



