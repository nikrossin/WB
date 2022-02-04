package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-query-store3"
	durableID = "service-json"
)

func ConnectNatsStream() (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(1, 3),
		stan.SetConnectionLostHandler(func(con stan.Conn, err error) {
			fmt.Printf("Connection nats lost: %s", err)
		}))
	if err != nil {
		return sc, err
	}

	return sc, nil
}
func SubscribeMsg(sc stan.Conn, data jsonFile, db *sql.DB, cash Cashe) error {

	handler := func(msg *stan.Msg) {
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			err = fmt.Errorf("incorrect messange from NATS %s", err)
		}

		err = SendDB(data, db)
		if err != nil {
			panic(err)
		} else {
			cash[data.OrderUid] = data
		}

	}
	_, err := sc.Subscribe("foo", handler, stan.DurableName(durableID))
	if err != nil {
		return fmt.Errorf("error subcribe: %s", err)
	}
	return nil

}
