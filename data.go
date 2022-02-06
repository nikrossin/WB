package main

import (
	"database/sql"
	"fmt"
	"task/pkg/memory"
	"task/pkg/model"
)

const (
	connStr = "user=myuser password=qwerty dbname=taskdb sslmode=disable"
)

func GetDataFromDB(inmemory memory.Memory) error {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rowsOrders, err := db.Query("select * From orders")
	if err != nil {
		return fmt.Errorf("query from orders: %s", err)
	}
	defer rowsOrders.Close()

	for rowsOrders.Next() {
		var idDelivery int
		var idPayment string
		var data model.Order

		err := rowsOrders.Scan(&data.OrderUid, &data.TrackNumber, &data.Entry, &idDelivery, &idPayment, &data.Locale, &data.InternalSignature,
			&data.CustomerId, &data.DeliveryService, &data.Shardkey, &data.SmId, &data.DateCreated, &data.OofShard)
		if err != nil {
			return fmt.Errorf("scan rows orders: %s", err)
		}
		rowsDelivery, err := db.Query("select * From deliveries where id = $1", idDelivery)
		if err != nil {
			return fmt.Errorf("query from deliveries: %s", err)
		}

		rowsDelivery.Next()
		err = rowsDelivery.Scan(&idDelivery, &data.Delivery.Name, &data.Delivery.Phone, &data.Delivery.Zip, &data.Delivery.City,
			&data.Delivery.Address, &data.Delivery.Region, &data.Delivery.Email)
		if err != nil {
			return fmt.Errorf("scan rows deliveries: %s", err)
		}
		if rowsDelivery.Next() {
			return fmt.Errorf("more one row detected in deliveries")
		}
		rowsDelivery.Close()
		rowsPayment, err := db.Query("select * From payments where transaction = $1", idPayment)
		if err != nil {
			return fmt.Errorf("query from payments: %s", err)
		}
		rowsPayment.Next()
		err = rowsPayment.Scan(&data.Payment.Transaction, &data.Payment.RequestId, &data.Payment.Currency,
			&data.Payment.Provider, &data.Payment.Amount, &data.Payment.PaymentDt, &data.Payment.Bank,
			&data.Payment.DeliveryCost, &data.Payment.GoodsTotal, &data.Payment.CustomFee)
		if err != nil {
			return fmt.Errorf("scan rows payments: %s", err)
		}
		if rowsPayment.Next() {
			return fmt.Errorf("more one row detected in payments")
		}
		rowsPayment.Close()

		rowsItems, err := db.Query("select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"+
			" From itemstab where order_id = $1", data.OrderUid)
		if err != nil {
			return fmt.Errorf("query from items: %s", err)
		}
		for rowsItems.Next() {
			item := model.Items{}
			err := rowsItems.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			if err != nil {
				return fmt.Errorf("scan rows items: %s", err)
			}
			data.Items = append(data.Items, item)
		}
		rowsItems.Close()
		inmemory.Set(data.OrderUid, data)
	}
	return nil
}

func SendDB(data model.Order) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %s", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec("select addDatadelivery($1,$2,$3,$4,$5,$6,$7)", data.Delivery.Name,
		data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address,
		data.Delivery.Region, data.Delivery.Email)
	if err != nil {
		return fmt.Errorf("add data in deliveries: %s", err)
	}
	fmt.Println(result.RowsAffected())

	result, err = tx.Exec("select addDataPayments($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", data.Payment.Transaction,
		data.Payment.RequestId, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount,
		data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal,
		data.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("add data in payments: %s", err)
	}
	fmt.Println(result.RowsAffected())

	result, err = tx.Exec("select addDataOrders($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", data.OrderUid,
		data.TrackNumber, data.Entry, data.Payment.Transaction, data.Locale, data.InternalSignature, data.CustomerId,
		data.DeliveryService, data.Shardkey, data.SmId, data.DateCreated, data.OofShard)
	if err != nil {
		return fmt.Errorf("add data in orders: %s", err)
	}
	fmt.Println(result.RowsAffected())

	for ind, _ := range data.Items {
		result, err = tx.Exec("select addDataItems($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
			data.Items[ind].ChrtId, data.Items[ind].TrackNumber, data.Items[ind].Price, data.Items[ind].Rid,
			data.Items[ind].Name, data.Items[ind].Sale, data.Items[ind].Size, data.Items[ind].TotalPrice,
			data.Items[ind].NmId, data.Items[ind].Brand, data.Items[ind].Status, data.OrderUid)
		if err != nil {
			return fmt.Errorf("add data in items: %s", err)
		}
		fmt.Println(result.RowsAffected())
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
