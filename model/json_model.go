package model

type DBSource_Param struct {
	Source_name     string `json:"source_name"`
	Source_type     string `json:"source_type"`
	Source_ipaddr   string `json:"source_ipaddr"`
	Source_port     string `json:"source_port"`
	Source_database string `json:"source_database"`
	Source_uid      string `json:"source_uid"`
	Source_pwd      string `json:"source_pwd"`
	Source_des      string `json:"source_des"`
}
type Result_Data struct {
	Code int
	Msg  string
}
type BuildDB_Param struct {
	Source_id int `json:"source_id"`
}
type BuildTable_Param struct {
	Source_id int `json:"source_id"`
	Table_id  int `json:"table_id"`
}
type Table_Param struct {
	Source_id  string        `json:"source_id"`
	Table_name string        `json:"table_name"`
	Table_des  string        `json:"table_des"`
	Field      []Field_Param `json:"field"`
}
type Field_Param struct {
	Field_name                 string `json:"field_name"`
	Datatype_id                string `json:"datatype_id"`
	Datatype_name              string `json:"datatype_name"`
	Datatype_is_fixed          int    `json:"datatype_is_fixed"`
	Datatype_is_quotation_mark int    `json:"datatype_is_quotation_mark"`
	Field_len                  string `json:"field_len"`
	Field_default              string `json:"field_default"`
	Field_pk                   int    `json:"field_pk"`
	Field_null                 int    `json:"field_null"`
	Field_auto                 int    `json:"field_auto"`
	Field_unsigned             int    `json:"field_unsigned"`
	Field_zero                 int    `json:"field_zero"`
}
