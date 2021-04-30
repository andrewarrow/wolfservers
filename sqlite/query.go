package sqlite

import (
	"fmt"
)

func NameExists(name string) bool {
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select name from stakes where name=?", name)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var s1 string
		rows.Scan(&s1)
		if s1 == name {
			return true
		}
	}
	return false
}
