package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"task/pkg/model"
	"testing"
)

func TestSendDB(t *testing.T) {
	file, _ := os.Open("static/model.json")
	var sendData, getData model.Order

	defer file.Close()
	content, _ := ioutil.ReadAll(file)

	err := json.Unmarshal(content, &sendData)
	if err != nil {
		t.Errorf("Structer model not correct")
	}

	//data for test
	sendData.OrderUid = "0000000000000000011"
	sendData.Payment.Transaction = "00000000000000011"
	sendData.Delivery.Name = "Test 011"

	err = SendDB(sendData)
	if err != nil {
		t.Errorf("%s", err)
	}
	var idDelivery int
	var idPayment string
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Errorf("error connect to db")
	}
	rowsOrders, _ := db.Query("select * from orders where order_uid = $1", sendData.OrderUid)
	rowsOrders.Next()
	rowsOrders.Scan(&getData.OrderUid, &getData.TrackNumber, &getData.Entry, &idDelivery, &idPayment, &getData.Locale, &getData.InternalSignature,
		&getData.CustomerId, &getData.DeliveryService, &getData.Shardkey, &getData.SmId, &getData.DateCreated, &getData.OofShard)
	rowsOrders.Close()

	rowsDelivery, _ := db.Query("select * From deliveries where id = $1", idDelivery)
	rowsDelivery.Next()
	rowsDelivery.Scan(&idDelivery, &getData.Delivery.Name, &getData.Delivery.Phone, &getData.Delivery.Zip, &getData.Delivery.City,
		&getData.Delivery.Address, &getData.Delivery.Region, &getData.Delivery.Email)

	if rowsDelivery.Next() {
		t.Errorf("more one row detected in deliveries")
	}
	rowsDelivery.Close()

	rowsPayment, _ := db.Query("select * From payments where transaction = $1", idPayment)
	rowsPayment.Next()
	err = rowsPayment.Scan(&getData.Payment.Transaction, &getData.Payment.RequestId, &getData.Payment.Currency,
		&getData.Payment.Provider, &getData.Payment.Amount, &getData.Payment.PaymentDt, &getData.Payment.Bank,
		&getData.Payment.DeliveryCost, &getData.Payment.GoodsTotal, &getData.Payment.CustomFee)

	if rowsPayment.Next() {
		t.Errorf("more one row detected in payments")
	}
	rowsPayment.Close()

	rowsItems, _ := db.Query("select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"+
		" From itemstab where order_id = $1", getData.OrderUid)
	for rowsItems.Next() {
		item := model.Items{}
		rowsItems.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)

		getData.Items = append(getData.Items, item)
	}
	rowsItems.Close()

	if !reflect.DeepEqual(getData, sendData) {
		fmt.Println(getData, "\n\n\n", sendData)
		t.Errorf("Add data from NATs not correct in service")
	}
}
