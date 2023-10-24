package main

import (
	"fmt"
	"src/src/actions"
	"src/src/utils"
	"strings"
)

func main() {
	var filePath string

	var authToken string
	fmt.Print("Please enter your auth token: ")
	fmt.Scanln(&authToken)

	var action string
	fmt.Println("1 - Create orders")
	fmt.Println("2 - Update orders")
	fmt.Println("3 - Cancel orders")
	fmt.Println("4 - Get orders detail")
	fmt.Println("5 - Filter orders by status")
	fmt.Print("Choose the action number: ")
	fmt.Scanln(&action)
	fmt.Println("")

	if utils.StringInSlice(action, []string{"1", "2", "3", "4"}) {
		fmt.Print("Please enter .csv file path: ")
		fmt.Scanln(&filePath)
		fmt.Println("")
	}

	if action == "1" {
		actions.CreateOrders(filePath, authToken)
	} else if action == "2" {
		actions.UpdateOrders(filePath, authToken)
	} else if action == "3" {
		actions.CancelOrders(filePath, authToken)
	} else if action == "4" {
		actions.GetOrders(filePath, authToken)
	} else if action == "5" {
		var fromDate string
		var toDate string
		var statusNumber string

		fmt.Print("From date YYYY-MM-DD: ")
		fmt.Scanln(&fromDate)
		fmt.Print("To date YYYY-MM-DD: ")
		fmt.Scanln(&toDate)
		fmt.Println("")

		fmt.Println("1 - Place Order Failed")
		fmt.Println("2 - Waiting for Payment")
		fmt.Println("3 - Waiting for Earner")
		fmt.Println("4 - Accepted")
		fmt.Println("5 - Processing")
		fmt.Println("6 - Canceled")
		fmt.Println("7 - Partial Shipping")
		fmt.Println("8 - Shipping")
		fmt.Println("9 - Partially Delivered")
		fmt.Println("10 - Delivered")
		fmt.Println("11 - Completed")
		fmt.Println("12 - Issue Founded")
		fmt.Println("13 - Confirm Pending")
		fmt.Print("Choose the status number (default is 1): ")
		fmt.Scanln(&statusNumber)

		reorder := false

		if statusNumber == "1" {
			r := "no"
			fmt.Print("Do you want the CSV file to be generated to create failed orders again [y/N]: ")
			fmt.Scanln(&r)
			if strings.ToLower(r) == "y" || strings.ToLower(r) == "yes" {
				reorder = true
			}
		}

		actions.GetOrdersDetail(authToken, utils.GetOrderStatus(statusNumber), fromDate, toDate, reorder)
	} else {
		fmt.Println("Undefined action number!")
	}
}
