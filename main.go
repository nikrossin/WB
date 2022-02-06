package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"task/pkg/cashe"
)

func main() {

	inmemory := cashe.New()

	err := GetDataFromDB(inmemory)
	if err != nil {
		fmt.Println(err)
	}

	sc, err := ConnectNatsStream()
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	err = SubscribeMsg(sc, inmemory)
	if err != nil {
		panic(err)
	}

	HttpHandlersStart(inmemory)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
