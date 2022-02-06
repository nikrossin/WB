package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"task/pkg/model"
	"testing"
	"time"
)

const (
	clusterIDTest = "test-cluster"
	clientIDTest  = "order-query-test-app"
)

func RandString(n int, numbers bool) string {
	delta := 0
	if !numbers {
		delta = 10
	}
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)-delta)]
	}
	return string(b)
}

type testItem struct {
	Order    string
	Payment  string
	Delivery string
}

func TestApp(t *testing.T) {
	file, err := os.Open("static/model.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)

	var data model.Order
	var respData model.Order
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	sc, _ := stan.Connect(clusterIDTest, clientIDTest) // Simple Synchronous Publisher
	var testTable []testItem
	for i := 0; i < 5; i++ {
		data.OrderUid = RandString(15, true)
		data.Payment.Transaction = RandString(15, true)
		data.Delivery.Name = RandString(10, false)
		msg, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		sc.Publish("foo", msg)
		testTable = append(testTable, testItem{data.OrderUid, data.Payment.Transaction, data.Delivery.Name})
	}
	time.Sleep(2 * time.Second)

	for _, val := range testTable {

		resp, _ := http.Get("http://localhost:8080/json?id=" + val.Order)
		if resp.StatusCode != 200 {
			t.Errorf(" id not found")
		}
		data.OrderUid, data.Payment.Transaction, data.Delivery.Name = val.Order, val.Payment, val.Delivery
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &respData)
		if !reflect.DeepEqual(respData, data) {
			t.Errorf("Add data from NATs not correct in service")
		}
	}
	sc.Close()
}
