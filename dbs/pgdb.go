package dbs

import (
	"database/sql"
	"fmt"
)

/**
这里是postgres数据库的相关逻辑
golang的数据库操作，主要分两类
1.Query():会返回结果，表示查询，它会从数据库获取查询结果（一系列行，可能为空）
2.Exec():不会返回结果，表示执行语句，它不会返回行
golang常见有两种数据库操作模式
QueryRow表示只返回一行的查询，作为Query的一个常见特例。
Prepare表示准备一个需要多次使用，供后续执行用。
 */
 /**
 什么时候用Exec，什么时候用Query，这是一个问题。
 通常DDL和增删改使用Exec，返回结果集的查询使用Query。
 但这不是绝对的，这完全取决于用户是否希望想要获取返回结果。
 例如在PostgreSQL中：INSERT ... RETURNING *;虽然是一条插入语句，但它也有返回结果集，故应当使用Query而不是Exec。
  */

 /**
 statement:需要执行的sql语句
 判断无误后，返回结果
 有错误的话，打印出执行语句，以及错误提示
  */
func Query(db *sql.DB,statement string,args ...interface{})*sql.Rows{
	rows,err := db.Query(statement,args)
	if err != nil {
		fmt.Println("Query执行语句：%v 时报错：%v",statement,err)
	}
	//注意这里的rows不应该defer rows.close(),上层的方便调用，遍历完结果之后才应该关闭。
	return rows
}

/**
单行查询
对于单行查询，Go将没有结果的情况视为错误
l包中定义了一个特殊的错误常量ErrNoRows，当结果为空时，QueryRow().Scan()会返回它
 */
 func QueryRow(db *sql.DB,statement string,args ...interface{})*sql.Row{
 	fmt.Println(args)
 	row:= db.QueryRow(statement,args...)
	return row
 }
 /**
 Exec不需要返回数据集，返回的结果是Result，
 Result接口允许获取执行结果的元数据
  */
func DBExec(db *sql.DB,statement string,args ...interface{})(sql.Result,error){
	resutl,err := db.Exec(statement,args...)
	if err != nil{
		fmt.Println("DBExec执行语句：%v 时报错：%v",statement,err)
		return nil,err
	}
	return resutl,nil
}

/**
stmt:是db.Prepare()返回的结果
 */
 func StmtExec(stmt *sql.Stmt,args ...interface{})sql.Result{
 	res,err := stmt.Exec(args)
 	if err != nil{
		fmt.Println("StmtExec执行语句 时报错：%v",err)
	}
 	return res
 }
