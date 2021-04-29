package sqlite

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
)

func List() {
	CreateSchema()
	InsertStake()
	db := OpenTheDB()
	defer db.Close()
	ListRows(db)
}

func ListRows(db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("select provider,ssh_key from stakes"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	for rows.Next() {
		var s1 string
		var s2 string
		rows.Scan(&s1, &s2)
		decodedBytes, _ := base64.StdEncoding.DecodeString(s2)
		shhh := decrypt(decodedBytes, phrase)
		fmt.Println(s1, string(shhh))
	}
}
