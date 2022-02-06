package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"task/pkg/cashe"
)

func main() {

	inmemory := cashe.New()

	err := GetDataFromDB(inmemory)
	if err != nil {
		log.Println(err)
	}

	sc, err := ConnectNatsStream()
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	err = SubscribeMsg(sc, inmemory)
	if err != nil {
		log.Fatal(err)
	}
	HttpHandlersStart(inmemory)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
