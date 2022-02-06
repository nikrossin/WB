package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"task/pkg/model"
	"testing"
)

var bodyResp = []string{

	`{"order_uid":"10","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000",
	"zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment"
	:{"transaction":"10","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha",
	"delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":90,"track_number":"WBILMTESTTRACK",
	"price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,
	"brand":"Vivienne Sabo","status":202},{"chrt_id":7930,"track_number":"WBILMTESTTRACK","price":453,"rid":
	"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,
	"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":
	"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`,

	`{"order_uid":"101","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov",
	"phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15",
	"region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"101","request_id":"",
	"currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,
	"goods_total":317,"custom_fee":0},"items":[{"chrt_id":90,"track_number":"WBILMTESTTRACK","price":453,
	"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,
	"brand":"Vivienne Sabo","status":202},{"chrt_id":7930,"track_number":"WBILMTESTTRACK","price":453,
	"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,
	"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test",
	"delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`,

	"Error ID",
	"empty ID or not correct request",
}

func TestHandleJson(t *testing.T) {
	tableTest := []struct {
		id     string
		status int
		body   string
	}{
		{
			id:     "10",
			status: 200,
			body:   bodyResp[0],
		},
		{
			id:     "101",
			status: 200,
			body:   bodyResp[1],
		},
		{
			id:     "1111111111111111111111asdadsfd",
			status: 404,
			body:   bodyResp[2],
		},
		{
			id:     "",
			status: 400,
			body:   bodyResp[3],
		},
	}

	var dataCorrect model.Order
	var dataResp model.Order
	for _, testItem := range tableTest {
		resp, _ := http.Get("http://localhost:8080/json?id=" + testItem.id)
		if resp.StatusCode != testItem.status {
			t.Errorf("HandlerJson with id = %s; statusCode not correct", testItem.id)
		}
		err := json.Unmarshal([]byte(testItem.body), &dataCorrect)
		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			if (testItem.body + "\n") != string(body) {
				t.Errorf("HandlerJson with id = %s; Body not correct", testItem.id)
			}
		} else {
			json.Unmarshal(body, &dataResp)
			if !reflect.DeepEqual(dataResp, dataCorrect) {
				t.Errorf("HandlerJson with id = %s; Body not correct", testItem.id)
			}
		}

	}
}
func TestHandleBase(t *testing.T) {
	resp, _ := http.Get("http://localhost:8080/")
	if resp.StatusCode != 200 {
		t.Errorf("HandleBase; StatusCode not correct")
	}
}
