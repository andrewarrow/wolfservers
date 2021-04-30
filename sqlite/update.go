package sqlite

func UpdateIps(name, producer, relay string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `update stakes set producer_ip=?, relay_ip=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(producer, relay, name)
	tx.Commit()
}
func UpdateIds(name, producer, relay string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `update stakes set producer_id=?, relay_id=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(producer, relay, name)
	tx.Commit()
}
