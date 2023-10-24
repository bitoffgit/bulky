package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseInt(str string) int {
	if strings.Contains(str, ".") {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return int(value)
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func ConvertToString(data interface{}) (res string) {
	switch v := data.(type) {
	case float64:
		res = strconv.FormatFloat(data.(float64), 'f', 2, 64)
	case float32:
		res = strconv.FormatFloat(float64(data.(float32)), 'f', 2, 32)
	case int:
		res = strconv.FormatInt(int64(data.(int)), 10)
	case int64:
		res = strconv.FormatInt(data.(int64), 10)
	case uint:
		res = strconv.FormatUint(uint64(data.(uint)), 10)
	case uint64:
		res = strconv.FormatUint(data.(uint64), 10)
	case uint32:
		res = strconv.FormatUint(uint64(data.(uint32)), 10)
	case json.Number:
		res = data.(json.Number).String()
	case string:
		res = data.(string)
	case []byte:
		res = string(v)
	case bool:
		res = "false"
		if v {
			res = "true"
		}
	default:
		res = ""
	}
	return res
}

func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ExtractAddress(order map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"first_name": order["first_name"],
		"last_name":  order["last_name"],
		"street":     order["street1"],
		"building":   order["street2"],
		"country":    order["country"],
		"state":      order["state"],
		"city":       order["city"],
		"phone":      order["phone"],
		"zip_code":   order["zip"],
	}
}

func FindCartItemID(cartItems map[string]interface{}, asin string) (int, error) {
	items := cartItems["data"].(map[string]interface{})["items"].([]interface{})
	for _, item := range items {
		product := item.(map[string]interface{})["product"].(map[string]interface{})
		if product["id"].(string) == asin {
			return ParseInt(
				ConvertToString(item.(map[string]interface{})["id"]),
			), nil
		}
	}

	return 0, fmt.Errorf("no cart item found with ASIN %s", asin)
}

func GetOrderStatus(statusNumber string) string {
	statuses := map[string]string{
		"1":  "Place Order Failed",
		"2":  "Waiting for Payment",
		"3":  "Waiting for Earner",
		"4":  "Accepted",
		"5":  "Processing",
		"6":  "Canceled",
		"7":  "Partial Shipping",
		"8":  "Shipping",
		"9":  "Partially Delivered",
		"10": "Delivered",
		"11": "Completed",
		"12": "Issue Founded",
		"13": "Confirm Pending",
	}
	return statuses[statusNumber]
}

func WriteFailedCreateOrdersInCsv(failedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("failed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"asin",
		"quantity",
		"discount",
		"currency",
		"first_name",
		"last_name",
		"street1",
		"street2",
		"city",
		"state",
		"zip",
		"country",
		"phone",
		"group_order",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range failedOrders {
		data = append(data, []string{
			ConvertToString(order["asin"]),
			ConvertToString(order["quantity"]),
			ConvertToString(order["discount"]),
			ConvertToString(order["currency"]),
			ConvertToString(order["first_name"]),
			ConvertToString(order["last_name"]),
			ConvertToString(order["street1"]),
			ConvertToString(order["street2"]),
			ConvertToString(order["city"]),
			ConvertToString(order["state"]),
			ConvertToString(order["zip"]),
			ConvertToString(order["country"]),
			ConvertToString(order["phone"]),
			ConvertToString(order["group_order"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteSucceedCreateOrdersInCsv(succeedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("succeed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"asin",
		"quantity",
		"discount",
		"currency",
		"first_name",
		"last_name",
		"street1",
		"street2",
		"city",
		"state",
		"zip",
		"country",
		"phone",
		"group_order",
		"bitoff_order_id",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range succeedOrders {
		data = append(data, []string{
			ConvertToString(order["asin"]),
			ConvertToString(order["quantity"]),
			ConvertToString(order["discount"]),
			ConvertToString(order["currency"]),
			ConvertToString(order["first_name"]),
			ConvertToString(order["last_name"]),
			ConvertToString(order["street1"]),
			ConvertToString(order["street2"]),
			ConvertToString(order["city"]),
			ConvertToString(order["state"]),
			ConvertToString(order["zip"]),
			ConvertToString(order["country"]),
			ConvertToString(order["phone"]),
			ConvertToString(order["group_order"]),
			ConvertToString(order["bitoff_order_id"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteFailedUpdateOrdersInCsv(failedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("failed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order_id",
		"discount",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range failedOrders {
		data = append(data, []string{
			ConvertToString(order["orderId"]),
			ConvertToString(order["discount"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteSucceedUpdateOrdersInCsv(succeedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("succeed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order_id",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range succeedOrders {
		data = append(data, []string{
			ConvertToString(order["orderId"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteFailedCancelOrdersInCsv(failedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("failed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order_id",
		"note",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range failedOrders {
		data = append(data, []string{
			ConvertToString(order["orderId"]),
			ConvertToString(order["note"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteSucceedCancelOrdersInCsv(succeedOrders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("succeed-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order_id",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range succeedOrders {
		data = append(data, []string{
			ConvertToString(order["orderId"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteOrderListDetailInCsv(orders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("orders-detail-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order_id",
		"currency",
		"status",
		"discount",
		"date",
		"fast_release",
		"price",
		"usdt_price",
		"usdt_origin_price",
		"origin_price",
		"origin_profit",
		"profit",
		"score",
		"source",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range orders {
		t := time.Unix(int64(order["date"].(float64)), 0)
		data = append(data, []string{
			ConvertToString(order["order_id"]),
			ConvertToString(order["currency"]),
			ConvertToString(order["status"]),
			strconv.FormatInt(int64(order["off"].(float64)), 10),
			t.Format(time.UnixDate),
			ConvertToString(order["fast_release"]),
			ConvertToString(order["price"]),
			ConvertToString(order["usdt_amount"]),
			ConvertToString(order["usdt_origin_amount"]),
			ConvertToString(order["origin_price"]),
			ConvertToString(order["origin_profit"]),
			ConvertToString(order["profit"]),
			ConvertToString(order["score"]),
			ConvertToString(order["source"]),
		})
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteOrdersInCsv(orders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("orders-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"order id",
		"currency",
		"status",
		"discount",
		"date",
		"items",
		"total price as USD",
		"total price as BTC",
		"you saved",
	}

	data := [][]string{
		headerRow,
	}

	for _, order := range orders {
		var orderItems string
		items := make(map[string]int)
		t := time.Unix(int64(order["at"].(float64)), 0)
		invoice := order["invoice"].(map[string]interface{})

		for _, item := range order["items"].([]interface{}) {
			productId := item.(map[string]interface{})["product"].(map[string]interface{})["id"].(string)

			if _, ok := items[productId]; ok {
				items[productId] += 1
				continue
			}

			items[productId] = 1
		}

		for asin, count := range items {
			orderItems += fmt.Sprintf("[Product: %s - Count: %d],", asin, count)
		}

		orderItems = strings.TrimRight(orderItems, ",")

		var totalPriceAsBTC string

		if order["currency"] == "btc" {
			totalPriceAsBTC = strconv.FormatFloat(invoice["btc"].(map[string]interface{})["amount"].(float64), 'f', 8, 64)
		}

		detail := []string{
			order["id"].(string),
			order["currency"].(string),
			order["status"].(string),
			strconv.FormatInt(int64(order["off"].(float64)), 10),
			t.Format(time.UnixDate),
			orderItems,
			ConvertToString(invoice["usd"].(map[string]interface{})["total"]),
			totalPriceAsBTC,
			ConvertToString(invoice["usd"].(map[string]interface{})["profit"]),
		}

		data = append(data, detail)
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WriteIssueOrdersInCsv(orders []map[string]interface{}) {
	currentDateTime := time.Now()

	f, err := os.Create(fmt.Sprintf("orders-%s.csv", currentDateTime.Format("2006-01-02-15-04-05")))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	headerRow := []string{
		"asin",
		"quantity",
		"discount",
		"currency",
		"first_name",
		"last_name",
		"street1",
		"street2",
		"city",
		"state",
		"zip",
		"country",
		"phone",
		"group_order",
	}

	data := [][]string{
		headerRow,
	}

	groupId := 1

	for _, order := range orders {
		orderData := order["data"].(map[string]interface{})
		orderAddress := orderData["address"].(map[string]interface{})
		orderItems := orderData["items"].([]interface{})

		products := make(map[string]int)

		for _, item := range orderItems {
			product := item.(map[string]interface{})["product"]
			key := product.(map[string]interface{})["id"].(string)
			if _, ok := products[product.(map[string]interface{})["id"].(string)]; ok {
				products[key] += 1
			} else {
				products[key] = 1
			}
		}

		for asin, quantity := range products {
			data = append(data, []string{
				asin,
				ConvertToString(quantity),
				strconv.FormatInt(int64(orderData["off"].(float64)), 10),
				ConvertToString(orderData["currency"]),
				ConvertToString(orderAddress["first_name"]),
				ConvertToString(orderAddress["last_name"]),
				ConvertToString(orderAddress["street"]),
				ConvertToString(orderAddress["building"]),
				ConvertToString(orderAddress["city"]),
				ConvertToString(orderAddress["state"]),
				ConvertToString(orderAddress["zip_code"]),
				ConvertToString(orderAddress["country"]),
				ConvertToString(orderAddress["phone"]),
				ConvertToString(groupId),
			})
		}

		groupId += 1
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}
