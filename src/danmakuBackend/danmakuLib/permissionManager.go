package danmakuLib

import (

	"database/sql"
	"fmt"
)

var permissionMap map[string]int

func QueryPermission(userId string) int {
	if permissionMap == nil{
		permissionMap = make(map[string]int)
	}

	permission, userFound := permissionMap[userId];
	if userFound {
		return permission
	} else {
		config := GetConfig()
		db, err := sql.Open("mysql", config.DBsource)
		if err != nil {
			println("failed to connect database.", err)
			db.Close()
			return -1
		}
		var permission int
		dbQuery := fmt.Sprintf("SELECT permission FROM users WHERE reg_code='%s';",
			userId)
		rows, err := db.Query(dbQuery)
		if rows.Next() {
			rows.Scan(&permission)
			permissionMap[userId] = permission
			return permission
		} else {
			return -1
		}
	}
	return -1
}

func SetPermission(userId string, permission int) bool {
	if permissionMap == nil{
		permissionMap = make(map[string]int)
	}

	permissionMap[userId] = permission

	config := GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.", err)
		db.Close()
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE users SET permission=? WHERE reg_code=?")
	result, err := stmt.Exec(permission, userId)
	affect, err := result.RowsAffected()
	defer stmt.Close()

	if err != nil {
		return false
	}

	if affect == 1 {
		return true
	} else {
		return false
	}
}
