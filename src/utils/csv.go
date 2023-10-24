package utils

import (
	"encoding/csv"
	"os"
)

func ParseCreateOrdersCsv(file *os.File) ([][]map[string]interface{}, error) {
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]map[string]interface{})

	for _, row := range rows[1:] {
		groupOrder := row[13]
		if _, ok := result[groupOrder]; !ok {
			result[groupOrder] = make([]map[string]interface{}, 0)
		}

		order := map[string]interface{}{
			"asin":        row[0],
			"quantity":    ParseInt(row[1]),
			"discount":    row[2],
			"currency":    row[3],
			"first_name":  row[4],
			"last_name":   row[5],
			"street1":     row[6],
			"street2":     row[7],
			"city":        row[8],
			"state":       row[9],
			"zip":         row[10],
			"country":     row[11],
			"phone":       row[12],
			"group_order": groupOrder,
		}

		result[groupOrder] = append(result[groupOrder], order)
	}

	orders := make([][]map[string]interface{}, 0)
	for _, v := range result {
		orders = append(orders, v)
	}

	return orders, nil
}

func ParseUpdateOrdersCsv(file *os.File) ([]map[string]interface{}, error) {
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for _, row := range rows[1:] {
		order := map[string]interface{}{
			"orderId":  row[0],
			"discount": row[1],
		}

		result = append(result, order)
	}

	return result, nil
}

func ParseCancelOrdersCsv(file *os.File) ([]map[string]interface{}, error) {
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for _, row := range rows[1:] {
		order := map[string]interface{}{
			"orderId": row[0],
			"note":    row[1],
		}

		result = append(result, order)
	}

	return result, nil
}

func ParseDetailOrdersCsv(file *os.File) ([]map[string]interface{}, error) {
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for _, row := range rows[1:] {
		order := map[string]interface{}{
			"orderId": row[0],
		}

		result = append(result, order)
	}

	return result, nil
}
