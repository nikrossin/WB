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
	file, _ := os.Open("static/model.json")
	var testTable []testItem
	var msgData, respData model.Order

	defer file.Close()
	content, err := ioutil.ReadAll(file)

	err = json.Unmarshal(content, &msgData)
	if err != nil {
		t.Errorf("error model struct")
	}
	rand.Seed(time.Now().UnixNano())

	sc, _ := stan.Connect(clusterIDTest, clientIDTest)
	for i := 0; i < 5; i++ {
		msgData.OrderUid = RandString(15, true)
		msgData.Payment.Transaction = RandString(15, true)
		msgData.Delivery.Name = RandString(10, false)
		msg, err := json.Marshal(msgData)
		if err != nil {
			t.Errorf("error marshal msg")
		}
		sc.Publish("foo", msg)
		testTable = append(testTable, testItem{msgData.OrderUid, msgData.Payment.Transaction, msgData.Delivery.Name})
	}
	//wait 2sec while service get msg
	time.Sleep(2 * time.Second)
	for _, val := range testTable {
		resp, _ := http.Get("http://localhost:8080/json?id=" + val.Order)
		if resp.StatusCode != 200 {
			t.Errorf(" id not found")
		}
		msgData.OrderUid, msgData.Payment.Transaction, msgData.Delivery.Name = val.Order, val.Payment, val.Delivery

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(body, &respData)
		if err != nil {
			t.Errorf("error in body response struct")
		}
		if !reflect.DeepEqual(respData, msgData) {
			t.Errorf("Add data from NATs not correct in service")
		}
	}

	sc.Close()
}
