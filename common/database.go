package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlOperate struct {
	DBtype	string
	ConnStr	string
}

func (p *MysqlOperate)InsertData(sqlstr string) int64 {
	db, err := sql.Open("mysql", p.ConnStr)
	defer db.Close()
	checkErr(err)
	res, err := db.Exec(sqlstr)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	return id
}
func (p *MysqlOperate)Exec(sqlstr string) sql.Result {
	db, err := sql.Open("mysql", p.ConnStr)
	defer db.Close()
	checkErr(err)
	res, err := db.Exec(sqlstr)
	checkErr(err)
	return res
}
func checkErr(err error)  {
	if err != nil{
		panic(err)
	}
}
func (p *MysqlOperate)QueryData(sqlstr string) map[int]map[string] interface{} {
	db, err := sql.Open(p.DBtype, p.ConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstr)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[int]map[string]interface{})
	i := 0
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		row := make(map[string]interface{})
		for j,col := range values{
			row[columns[j]] = col
		}
		record[i] = row
		i++
	}
	return  record
}