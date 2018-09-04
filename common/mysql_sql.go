package common

import "strconv"

type MysqlSQL struct {
}

func (p *MysqlSQL)DropTable(tabname string) string {
	sqlstr := "drop table if exists " + tabname + ";"
	return  sqlstr
}
func (p *MysqlSQL) CreateTable(tabname string, field []map[string]interface{}) string {

	sqlstr := "create table " + tabname + "("
	pk := ""
	for _, v := range field {
		sqlstr += v["field_name"].(string) + " " + v["datatype_name"].(string)
		c_int := v["datatype_is_fixed"].(int)
		if c_int == 0 {
			sqlstr += "(" + strconv.Itoa(v["field_len"].(int)) + ")"
		}
		c_int = v["field_null"].(int)
		if c_int == 0 {
			sqlstr += " NOT NULL"
		}
		c_int = v["field_auto"].(int)
		if c_int == 1 {
			sqlstr += " auto_increment"
		}
		c_int = v["field_unsigned"].(int)
		if c_int == 1 {
			sqlstr += " UNSIGNED"
		}
		c_int = v["field_zero"].(int)
		if c_int == 1 {
			sqlstr += " ZEROFILL"
		}
		sqlstr += ","
		c_int = v["field_pk"].(int)
		if c_int == 1 {
			pk += v["field_name"].(string) + ","
		}
	}
	if len([]rune(pk)) > 0 {
		l := len([]rune(pk))
		sqlstr += "primary key (" + string([]byte(pk)[:(l - 1)]) + ")"
	}
	sqlstr += ");"
	return sqlstr
}
