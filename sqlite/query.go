package sqlite

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
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

func CreateNodeKeysOnDisk(name string) {
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select counter, vkey, skey from nodes where name=?", name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	phrase := os.Getenv("WOLF_PHRASE")
	rows.Next()
	var s1 string
	var s2 string
	var s3 string
	rows.Scan(&s1, &s2, &s3)
	decodedBytes, _ := base64.StdEncoding.DecodeString(s1)
	shhh := decrypt(decodedBytes, phrase)
	ioutil.WriteFile("node.counter", shhh, 0755)

	decodedBytes, _ = base64.StdEncoding.DecodeString(s2)
	shhh = decrypt(decodedBytes, phrase)
	ioutil.WriteFile("node.vkey", shhh, 0755)

	decodedBytes, _ = base64.StdEncoding.DecodeString(s3)
	shhh = decrypt(decodedBytes, phrase)
	ioutil.WriteFile("node.skey", shhh, 0755)
}
