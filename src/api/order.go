package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetOrder(orderId, authToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.bitoff.io/o/%s", orderId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bodyResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if message, ok := bodyResp["message"].(string); ok {
			fmt.Println("Get order has failed:", message)
		} else {
			fmt.Println("Get order has failed:", bodyResp)
		}
		return nil, fmt.Errorf("get order has failed with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func CreateOrder(body map[string]interface{}, authToken string) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.bitoff.io/o/submit", strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var bodyResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if resp.StatusCode != 201 {
		if message, ok := bodyResp["message"].(string); ok {
			fmt.Println("Order creation failed:", message)
		} else {
			fmt.Println("Order creation failed:", bodyResp)
		}
		return nil, fmt.Errorf("crafting order creation failed with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func UpdateOrder(orderId string, discount int, authToken string) error {
	jsonBody, err := json.Marshal(map[string]int{"off": discount})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	url := fmt.Sprintf("https://api.bitoff.io/o/%s/edit/apply", orderId)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("order modification failed with status code %d", resp.StatusCode)
	}

	return nil
}

func CancelOrder(orderId string, note string, authToken string) error {
	jsonBody, err := json.Marshal(map[string]string{"note": note})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	url := fmt.Sprintf("https://api.bitoff.io/o/%s/cancel", orderId)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("order cancellation failed with status code %d", resp.StatusCode)
	}

	return nil
}

func GetOrdersByFilter(status string, fromDate string, toDate string, authToken string) ([]map[string]interface{}, error) {
	status = url.QueryEscape(status)
	url := fmt.Sprintf(
		"https://api.bitoff.io/a/users/order-list?status=%s&pg=0&from_date=%s&to_date=%s",
		status,
		fromDate,
		toDate,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var bodyResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error:", bodyResp)
		return nil, fmt.Errorf("failed to get orders with status code %d", resp.StatusCode)
	}

	var result []map[string]interface{}

	for _, order := range bodyResp["items"].([]interface{}) {
		result = append(result, order.(map[string]interface{}))
	}

	return result, nil
}
