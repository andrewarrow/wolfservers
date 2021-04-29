package sqlite

func UpdateIps(name, producer, relay string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `update stakes set producer=?, relay=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(producer, relay, name)
	tx.Commit()
}
