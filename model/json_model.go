package model

type DBSource_Param struct {
	Name 		string 	`json:"name"`
	Ipaddr 		string 	`json:"ipaddr"`
	Port 		string 	`json:"port"`
	Database 	string 	`json:"database"`
	Uid 		string 	`json:"uid"`
	Pwd 		string 	`json:"pwd"`
}
type Result_Data struct {
	Code	int
	Msg		string
}

