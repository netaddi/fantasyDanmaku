package danmakuLib

import (

	"database/sql"
	"fmt"
)

var permissionMap map[string]int

func InitPermissionMap() {
	permissionMap = make(map[string]int)
}

func QueryPermission(userId string) int {

	permission, userFound := permissionMap[userId]
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
		defer db.Close()

		var permission int
		dbQuery := fmt.Sprintf("SELECT permission FROM users WHERE reg_code='%s';",
			userId)
		rows, err := db.Query(dbQuery)
		defer rows.Close()
		if err != nil {
			println("failed to connect database.")
			return 1
		}
		if rows.Next() {
			_ = rows.Scan(&permission)
			permissionMap[userId] = permission
			return permission
		} else {
			return -1
		}
	}
}

func SetPermission(userId string, permission int) bool {
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

func SetMemoryPermission(userId string, permission int) {
	permissionMap[userId] = permission
}