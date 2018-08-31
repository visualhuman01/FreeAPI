package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlOperate struct {
	ConnStr	string
}

func (p *MysqlOperate)Exec(sqlstr string) sql.Result {
	db, err := sql.Open("mysql", p.ConnStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(sqlstr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return res
}

func (p *MysqlOperate)QueryData(sqlstr string) map[int][]interface{} {
	db, err := sql.Open("mysql", p.ConnStr)
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
	record := make(map[int][]interface{})
	i := 0
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record[i] = make([]interface{}, len(columns))
		copy(record[i],values)
		i++
	}
	return  record
}