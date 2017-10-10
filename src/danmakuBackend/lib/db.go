package lib

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//func queryDB(sqlQuery string) bool {
//
//	config := GetConfig()
//	db, err := sql.Open("mysql", config.DBsource)
//	if err != nil {
//		println("failed to connect database.")
//		return false
//	}
//
//
//}

//func checkErr(err error) {
//	if err != nil {
//		panic(err)
//	}
//}