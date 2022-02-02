package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-query-store3"
	durableID = "service-json"
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

type Cashe map[string]jsonFile

func GetDataFromDB(db *sql.DB, cash *Cashe) error {
	rowsOrders, err := db.Query("select * From orders")
	if err != nil {
		return err
	}

	var data jsonFile

	for rowsOrders.Next() {
		var idDelivery int
		var idPayment string

		err := rowsOrders.Scan(&data.OrderUid, &data.TrackNumber, &data.Entry, &idDelivery, &idPayment, &data.Locale, &data.InternalSignature,
			&data.CustomerId, &data.DeliveryService, &data.Shardkey, &data.SmId, &data.DateCreated, &data.OofShard)
		if err != nil {
			return err
		}
		rowsDelivery, err := db.Query("select * From deliveries where id = $1", idDelivery)
		if err != nil {
			return err
		}
		rowsDelivery.Next()
		err = rowsDelivery.Scan(&idDelivery, &data.Delivery.Name, &data.Delivery.Phone, &data.Delivery.Zip, &data.Delivery.City,
			&data.Delivery.Address, &data.Delivery.Region, &data.Delivery.Email)
		if err != nil {
			return err
		}
		if rowsDelivery.Next() {
			return err
		}
		rowsPayment, err := db.Query("select * From payments where transaction = $1", idPayment)
		if err != nil {
			return err
		}
		rowsPayment.Next()
		err = rowsPayment.Scan(&data.Payment.Transaction, &data.Payment.RequestId, &data.Payment.Currency,
			&data.Payment.Provider, &data.Payment.Amount, &data.Payment.PaymentDt, &data.Payment.Bank,
			&data.Payment.DeliveryCost, &data.Payment.GoodsTotal, &data.Payment.CustomFee)
		if err != nil {
			return err
		}
		if rowsPayment.Next() {
			return err
		}

		rowsItems, err := db.Query("select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"+
			" From itemstab where order_id = $1", data.OrderUid)
		if err != nil {
			return err
		}
		for rowsItems.Next() {
			item := Items{}
			err := rowsItems.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			if err != nil {
				return err
			}
			data.Items = append(data.Items, item)
		}
		(*cash)[data.OrderUid] = data
	}
	return nil
}

func SendDB(data jsonFile, db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("select addDatadelivery($1,$2,$3,$4,$5,$6,$7)", data.Delivery.Name,
		data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address,
		data.Delivery.Region, data.Delivery.Email)
	if err != nil {
		fmt.Println("rrr1")
		return err
	}
	fmt.Println(result.RowsAffected())

	result, err = tx.Exec("select addDataPayments($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", data.Payment.Transaction,
		data.Payment.RequestId, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount,
		data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal,
		data.Payment.CustomFee)
	if err != nil {
		fmt.Println("rrr")
		return err
	}
	fmt.Println(result.RowsAffected())

	result, err = tx.Exec("select addDataOrders($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", data.OrderUid,
		data.TrackNumber, data.Entry, data.Payment.Transaction, data.Locale, data.InternalSignature, data.CustomerId,
		data.DeliveryService, data.Shardkey, data.SmId, data.DateCreated, data.OofShard)
	if err != nil {
		fmt.Println("rrr2")
		return err
	}
	fmt.Println(result.RowsAffected())

	for ind, _ := range data.Items {
		result, err = tx.Exec("select addDataItems($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
			data.Items[ind].ChrtId, data.Items[ind].TrackNumber, data.Items[ind].Price, data.Items[ind].Rid,
			data.Items[ind].Name, data.Items[ind].Sale, data.Items[ind].Size, data.Items[ind].TotalPrice,
			data.Items[ind].NmId, data.Items[ind].Brand, data.Items[ind].Status, data.OrderUid)
		if err != nil {
			fmt.Println("rrr3")
			return err
		}
		fmt.Println(result.RowsAffected())
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func main() {
	var data jsonFile
	//ch := make(chan struct{})
	cash := Cashe(make(map[string]jsonFile))
	connStr := "user=myuser password=qwerty dbname=taskdb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = GetDataFromDB(db, &cash)
	if err != nil {
		panic(err)
	}
	fmt.Println(cash)

	sc, err := stan.Connect(clusterID, clientID) // Simple Synchronous Publisher
	if err != nil {
		panic(err)
	}

	handler := func(msg *stan.Msg) {
		err = json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Println("kek")
			panic(err)
		}

		fmt.Println(cash)
		err = SendDB(data, db)
		if err != nil {
			panic(err)
		} else {
			cash[data.OrderUid] = data
		}

	}
	sub, err := sc.Subscribe("foo", handler, stan.DurableName(durableID))
	if err != nil {
		fmt.Println("kk")
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		id, ok := r.URL.Query()["id"]
		var dataResponse string
		if !ok || len(id[0]) < 1 {
			fmt.Fprint(w, "Error ID")
		} else {
			if val, ok := cash[id[0]]; ok {
				jsonData, err := json.Marshal(val)
				if err != nil {
					panic(err)
				} else {
					dataResponse = strings.ReplaceAll(string(jsonData), ",", ",\n")
					dataResponse = strings.ReplaceAll(dataResponse, "{", "{\n")
					dataResponse = strings.ReplaceAll(dataResponse, "},", "\n},")
				}

			} else {
				dataResponse = "Error id"
			}
			fmt.Fprint(w, dataResponse)
		}
	})

	defer sub.Unsubscribe()
	defer sc.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))

	//<-ch
	//runtime.Goexit()

}
