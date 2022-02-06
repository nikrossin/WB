package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-query-store2"
)

type jsonFile struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Items   `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Items struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

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
func main() {
	file, err := os.Open("model.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)

	var data jsonFile
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	sc, _ := stan.Connect(clusterID, clientID) // Simple Synchronous Publisher
	for i := 0; i < 5; i++ {
		data.OrderUid = RandString(15, true)
		data.Payment.Transaction = RandString(15, true)
		data.Delivery.Name = RandString(10, false)
		msg, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		sc.Publish("foo", msg)
	}

	sc.Close()
}
