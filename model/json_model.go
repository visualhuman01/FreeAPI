package model

type DBSource_Param struct {
	Source_name 		string `json:"source_name"`
	Source_type			string `json:"source_type"`
	Source_ipaddr 		string `json:"source_ipaddr"`
	Source_port 		string `json:"source_port"`
	Source_database 	string `json:"source_database"`
	Source_uid 			string `json:"source_uid"`
	Source_pwd 			string `json:"source_pwd"`
	Source_des			string `json:"source_des"`
}
type Result_Data struct {
	Code	int
	Msg		string
}
type DBBuild_Param struct {
	Id	int `json:"id"`
}

