package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	var data jsonFile
	cash := Cashe(make(map[string]jsonFile))
	connStr := "user=myuser password=qwerty dbname=taskdb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = GetDataFromDB(db, cash)
	if err != nil {
		fmt.Println(err)
	}

	sc, err := ConnectNatsStream()
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	err = SubscribeMsg(sc, data, db, cash)
	if err != nil {
		panic(err)
	}
	HttpServerStart(cash)

}
