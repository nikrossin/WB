package main

import (
	"database/sql"
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
	//используем готовую json модель для формирования сообщения
	content, err := ioutil.ReadAll(file)

	err = json.Unmarshal(content, &msgData)
	if err != nil {
		t.Errorf("error model struct")
	}
	rand.Seed(time.Now().UnixNano())
	sc, _ := stan.Connect(clusterIDTest, clientIDTest)
	//на основе готовой модели генерируем 5 тестовых сообщений, изменяя 3 поля в модели
	for i := 0; i < 5; i++ {
		msgData.OrderUid = RandString(15, true)
		msgData.Payment.Transaction = RandString(15, true)
		msgData.Delivery.Name = RandString(10, false)
		msg, err := json.Marshal(msgData)
		if err != nil {
			t.Errorf("error marshal msg")
		}
		sc.Publish("foo", msg)
		//формируем таблицу тестов
		testTable = append(testTable, testItem{msgData.OrderUid, msgData.Payment.Transaction, msgData.Delivery.Name})
	}
	//ждем некотрое время, пока сервис не получит сообщения
	time.Sleep(2 * time.Second)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, val := range testTable {
		resp, _ := http.Get("http://localhost:8080/json?id=" + val.Order)
		if resp.StatusCode != 200 {
			t.Errorf(" id not found")
		}
		//обновляем поля сообшения, отправленного в nats-streaming для каждого теств
		msgData.OrderUid, msgData.Payment.Transaction, msgData.Delivery.Name = val.Order, val.Payment, val.Delivery

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(body, &respData)
		if err != nil {
			t.Errorf("error in body response struct")
		}
		//проверка добавил ли сервис данные от nats
		if !reflect.DeepEqual(respData, msgData) {
			t.Errorf("Add data from NATs not correct in service")
		}
		rowsOrders, _ := db.Query("select * from orders where order_uid = $1", respData.OrderUid)
		if !rowsOrders.Next() {
			t.Errorf("msg not add in DB")
		}
		rowsOrders.Close()

	}
	sc.Close()
}
