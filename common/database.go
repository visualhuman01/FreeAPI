package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type MysqlOperate struct {
	DBtype  string
	ConnStr string
}

func (p *MysqlOperate) InsertData(sqlstr string) int64 {
	db, err := sql.Open("mysql", p.ConnStr)
	defer db.Close()
	checkErr(err)
	res, err := db.Exec(sqlstr)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	return id
}
func (p *MysqlOperate) Exec(sqlstr string) sql.Result {
	db, err := sql.Open("mysql", p.ConnStr)
	defer db.Close()
	checkErr(err)
	res, err := db.Exec(sqlstr)
	checkErr(err)
	return res
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func (p *MysqlOperate) QueryRow(sqlstr string) map[string]interface{} {
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
	columnstype, _ := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(columnstype))
	values := make([]interface{}, len(columnstype))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	rows.Next()
	//将行数据保存到record字典
	err = rows.Scan(scanArgs...)
	row := make(map[string]interface{})
	for j, col := range values {
		row[columnstype[j].Name()] = getData(columnstype[j].DatabaseTypeName(),col)
	}
	return row
}
func (p *MysqlOperate) QueryData(sqlstr string) []map[string]interface{} {
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
	columnstype, _ := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(columnstype))
	values := make([]interface{}, len(columnstype))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[int]map[string]interface{})
	i := 0
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		row := make(map[string]interface{})
		for j, col := range values {
			row[columnstype[j].Name()] = getData(columnstype[j].DatabaseTypeName(),col)
		}
		record[i] = row
		i++
	}
	res := make([]map[string]interface{}, i)
	for k, v := range record {
		res[k] = v
	}
	return res
}
func getData(t string, col interface{}) interface{} {
	var res interface{}
	switch t {
	case "TINYINT","INT":
		c_str := string(col.([]byte))
		c_int, _ := strconv.Atoi(c_str)
		res = c_int
		break
	case "VARCHAR","TEXT","DATETIME":
		if col != nil {
			res = string(col.([]byte))
		} else {
			res = ""
		}
		break
	default:
		res = col
		break
	}
	return res
}
