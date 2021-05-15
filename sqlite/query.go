package sqlite

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lox/go-touchid"
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
func BioMetricNo() bool {
	ok, err := touchid.Authenticate("access llamas")
	if err != nil {
		log.Fatal(err)
	}

	if ok {
		log.Printf("Authenticated")
		return false
	} else {
		log.Fatal("Failed to authenticate")
	}
	return true
}

func ShowOaths() {
	if BioMetricNo() {
		return
	}
	fmt.Println("")
	fmt.Println("        oathtool --totp -b ''")
	fmt.Println("")
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select name,seed,username,password from oaths")
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
		rows.Scan(&s1, &s2, &s3, &s4)
		decodedBytes, _ := base64.StdEncoding.DecodeString(s1)
		shhh := decrypt(decodedBytes, phrase)
		name := string(shhh)
		decodedBytes, _ = base64.StdEncoding.DecodeString(s2)
		shhh = decrypt(decodedBytes, phrase)
		seed := string(shhh)
		decodedBytes, _ = base64.StdEncoding.DecodeString(s3)
		shhh = decrypt(decodedBytes, phrase)
		username := string(shhh)
		decodedBytes, _ = base64.StdEncoding.DecodeString(s4)
		shhh = decrypt(decodedBytes, phrase)
		password := string(shhh)
		fmt.Println(name, seed, username, password)
	}
}
func LoadPats() map[string]string {
	db := OpenTheDB()
	defer db.Close()
	m := map[string]string{}
	rows, err := db.Query("select provider,pat from pats")
	if err != nil {
		fmt.Println(err)
		return m
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	for rows.Next() {
		var s1 string
		var s2 string
		rows.Scan(&s1, &s2)
		decodedBytes, _ := base64.StdEncoding.DecodeString(s2)
		shhh := decrypt(decodedBytes, phrase)
		m[s1] = string(shhh)
	}
	return m
}

func PaymentAndStakeSigning(name string) (string, string) {
	if BioMetricNo() {
		return "", ""
	}
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select ps,ss from payment where name=?", name)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	rows.Next()
	var s1 string
	var s2 string
	rows.Scan(&s1, &s2)
	if s1 == "" {
		return "", ""
	}
	decodedBytes, _ := base64.StdEncoding.DecodeString(s1)
	shhh1 := decrypt(decodedBytes, phrase)
	decodedBytes, _ = base64.StdEncoding.DecodeString(s2)
	shhh2 := decrypt(decodedBytes, phrase)
	return string(shhh1), string(shhh2)
}
func PaymentStakeV(name string) string {
	if BioMetricNo() {
		return ""
	}
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select sv from payment where name=?", name)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	rows.Next()
	var s1 string
	rows.Scan(&s1)
	if s1 == "" {
		return ""
	}
	decodedBytes, _ := base64.StdEncoding.DecodeString(s1)
	shhh := decrypt(decodedBytes, phrase)
	return string(shhh)
}
func PaymentAddressQuery(name string) (string, string) {
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select pa,sa from payment where name=?", name)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	rows.Next()
	var s1 string
	var s2 string
	rows.Scan(&s1, &s2)
	if s1 == "" {
		return "", ""
	}
	decodedBytes, _ := base64.StdEncoding.DecodeString(s1)
	shhh1 := decrypt(decodedBytes, phrase)
	decodedBytes, _ = base64.StdEncoding.DecodeString(s2)
	shhh2 := decrypt(decodedBytes, phrase)
	return string(shhh1), string(shhh2)
}
func NodeKeysQuery(name string) int {
	db := OpenTheDB()
	defer db.Close()
	rows, err := db.Query("select vkey from nodes where name=?", name)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer rows.Close()

	rows.Next()
	var s1 string
	rows.Scan(&s1)
	return len(s1)
}
func CreateNodeKeysOnDisk(name string) {
	if BioMetricNo() {
		return
	}
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
