package sqlite

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
)

func List() {
	db := OpenTheDB()
	defer db.Close()
	ListRows(db)
}

func ListRows(db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("select name,provider,ssh_key,producer_ip,relay_ip from stakes"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	for rows.Next() {
		var s1 string
		var s2 string
		var s3 string
		var s4 string
		var s5 string
		rows.Scan(&s1, &s2, &s3, &s4, &s5)
		decodedBytes, _ := base64.StdEncoding.DecodeString(s3)
		shhh := decrypt(decodedBytes, phrase)
		fmt.Println(s1, s2, len(shhh), s4, s5)
	}
}

func MakeIpMap(db *sql.DB) map[string]string {
	m := map[string]string{}
	rows, err := db.Query(fmt.Sprintf("select name,producer_ip,relay_ip from stakes"))
	if err != nil {
		fmt.Println(err)
		return m
	}
	defer rows.Close()
	for rows.Next() {
		var s1 string
		var s2 string
		var s3 string
		rows.Scan(&s1, &s2, &s3)
		m[s2] = s1
		m[s3] = s1
	}
	return m
}

func MakeIpToId(db *sql.DB) map[string]string {
	m := map[string]string{}
	rows, err := db.Query(fmt.Sprintf("select producer_ip,producer_id,relay_ip,relay_id from stakes"))
	if err != nil {
		fmt.Println(err)
		return m
	}
	defer rows.Close()
	for rows.Next() {
		var s1 string
		var s2 string
		var s3 string
		var s4 string
		rows.Scan(&s1, &s2, &s3, &s4)
		m[s1] = s2
		m[s3] = s4
	}
	return m
}
