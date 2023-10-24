package actions

import (
	"fmt"
	"os"
	"src/src/api"
	"src/src/utils"
	"time"
)

func makeOrder(order []map[string]interface{}, authToken string) map[string]interface{} {
	defaultAddress, err := api.GetDefaultAddress(authToken)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	defaultAddress = utils.MergeMaps(
		defaultAddress,
		utils.ExtractAddress(order[0]),
	)

	if _, err := api.UpdateAddress(defaultAddress, authToken); err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	address, err := api.GetDefaultAddress(authToken)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	cartItems, err := api.GetCartItems(authToken)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	if _, ok := cartItems["data"].([]interface{}); !ok {
		for _, item := range cartItems["data"].(map[string]interface{})["items"].([]interface{}) {
			itemID := int(item.(map[string]interface{})["id"].(float64))
			if err := api.DeleteCartItem(itemID, authToken); err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			fmt.Println("Delete cart item:", item.(map[string]interface{})["product"].(map[string]interface{})["id"])
		}
	}

	for _, orderItem := range order {
		product, status, err := api.GetProduct(orderItem["asin"].(string))
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		if status == 200 {
			for {
				product, status, err = api.UpdateProduct(orderItem["asin"].(string))
				if err != nil {
					fmt.Println("Error:", err)
					return nil
				}
				if status == 200 {
					break
				}
				fmt.Println(fmt.Sprintf("Waiting for updating product %s price and...", orderItem["asin"]))
				time.Sleep(2 * time.Second)
			}
		} else {
			for status == 202 {
				fmt.Println(fmt.Sprintf("Waiting for fetching product %s...", orderItem["asin"]))
				time.Sleep(2 * time.Second)
				product, status, err = api.GetProduct(orderItem["asin"].(string))
				if err != nil {
					fmt.Println("Error:", err)
					return nil
				}
			}
		}

		if _, err := api.AddToCart(map[string]interface{}{
			"offer":   product["merchant"],
			"product": orderItem["asin"].(string),
		}, authToken); err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		fmt.Println("Added to cart:", orderItem["asin"].(string))

		if orderItem["quantity"].(int) > 1 {
			fmt.Println("Update quantity")
			cartItems, err := api.GetCartItems(authToken)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			cartItemID, _ := utils.FindCartItemID(cartItems, orderItem["asin"].(string))
			quantity := orderItem["quantity"].(int)
			if quantity > 5 {
				quantity = 5
			}

			for i := 0; i < quantity-1; i++ {
				if _, err := api.IncreaseCartItemQuantity(cartItemID, authToken); err != nil {
					fmt.Println("Error:", err)
					return nil
				}
			}
		}
	}

	result, _ := api.CreateOrder(map[string]interface{}{
		"off":             order[0]["discount"].(string),
		"fast":            false,
		"currency":        order[0]["currency"].(string),
		"address":         address["id"].(float64),
		"username":        0,
		"can_see_address": false,
	}, authToken)

	return result
}

func CreateOrders(filePath, authToken string) {
	succeed := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	orders, err := utils.ParseCreateOrdersCsv(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	addresses, err := api.GetAddresses(authToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(addresses) == 0 {
		_, err := api.CreateAddress(utils.ExtractAddress(orders[0][0]), authToken)
		if err != nil {
			fmt.Println(err)
			return
		}
		addresses, err = api.GetAddresses(authToken)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if !addresses[0]["default"].(bool) {
		_, err := api.SetAddressToDefault(addresses[0], authToken)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, order := range orders {
		fmt.Println(fmt.Sprintf("-----Process started for order group: %s------", order[0]["group_order"]))
		//fmt.Println()
		result := makeOrder(order, authToken)
		if result != nil {
			for _, item := range order {
				item["bitoff_order_id"] = result["orderId"]
			}
			succeed = append(succeed, order...)
			fmt.Printf("Order Created: %s\n", result["orderId"])
		} else {
			failed = append(failed, order...)
		}
		//fmt.Println()
		fmt.Println(fmt.Sprintf("-----Process finished for order group: %s-----", order[0]["group_order"]))
		fmt.Println("---------------------------------------------")
	}

	if len(failed) > 0 {
		utils.WriteFailedCreateOrdersInCsv(failed)
	}

	if len(succeed) > 0 {
		utils.WriteSucceedCreateOrdersInCsv(succeed)
	}
}

func UpdateOrders(filePath, authToken string) {
	succeed := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	orders, err := utils.ParseUpdateOrdersCsv(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, order := range orders {
		err = api.UpdateOrder(
			utils.ConvertToString(order["orderId"]),
			utils.ParseInt(utils.ConvertToString(order["discount"])),
			authToken,
		)
		if err == nil {
			succeed = append(succeed, order)
			fmt.Printf("Order Updated: %s\n", order["orderId"])
		} else {
			fmt.Println(err)
			failed = append(failed, order)
		}
	}

	if len(failed) > 0 {
		utils.WriteFailedUpdateOrdersInCsv(failed)
	}

	if len(succeed) > 0 {
		utils.WriteSucceedUpdateOrdersInCsv(succeed)
	}
}

func CancelOrders(filePath, authToken string) {
	succeed := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	orders, err := utils.ParseCancelOrdersCsv(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, order := range orders {
		err = api.CancelOrder(
			utils.ConvertToString(order["orderId"]),
			utils.ConvertToString(order["note"]),
			authToken,
		)
		if err == nil {
			succeed = append(succeed, order)
			fmt.Printf("Order Canceled: %s\n", order["orderId"])
		} else {
			fmt.Println(err)
			failed = append(failed, order)
		}
	}

	if len(failed) > 0 {
		utils.WriteFailedCancelOrdersInCsv(failed)
	}

	if len(succeed) > 0 {
		utils.WriteSucceedCancelOrdersInCsv(succeed)
	}
}

func GetOrders(filePath, authToken string) {
	results := make([]map[string]interface{}, 0)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	orders, err := utils.ParseDetailOrdersCsv(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, order := range orders {
		result, err := api.GetOrder(order["orderId"].(string), authToken)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, result["data"].(map[string]interface{}))
	}

	if len(results) > 0 {
		utils.WriteOrdersInCsv(results)
	}
}

func GetOrdersDetail(authToken, status, fromDate, toDate string, reorder bool) {
	orders, err := api.GetOrdersByFilter(status, fromDate, toDate, authToken)

	if err != nil {
		return
	}

	if len(orders) == 0 {
		fmt.Println(fmt.Sprintf("There are no orders in this date range with the status [%s]", status))
		return
	}

	fmt.Println(fmt.Sprintf("%d orders found with status [%s], wait for creating CSV file", len(orders), status))
	utils.WriteOrderListDetailInCsv(orders)

	if reorder {
		var ordersCollection []map[string]interface{}
		for _, order := range orders {
			orderDetail, err := api.GetOrder(order["order_id"].(string), authToken)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ordersCollection = append(ordersCollection, orderDetail)
		}
		utils.WriteIssueOrdersInCsv(ordersCollection)
	}
}
