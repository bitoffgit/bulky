package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetCartItems(authToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.bitoff.io/o/cart", nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error:", bodyResp)
		return nil, fmt.Errorf("failed to get cart items with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func DeleteCartItem(itemId int, authToken string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.bitoff.io/o/cart/items/%d", itemId), nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to delete cart item with status code %d", resp.StatusCode)
	}

	return nil
}

func AddToCart(body map[string]interface{}, authToken string) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.bitoff.io/o/cart", strings.NewReader(string(jsonBody)))
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error:", bodyResp)
		return nil, fmt.Errorf("failed to add to cart with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func IncreaseCartItemQuantity(itemId int, authToken string) (map[string]interface{}, error) {
	jsonBody := map[string]interface{}{
		"value": 1,
	}
	jsonStr, err := json.Marshal(jsonBody)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	url := fmt.Sprintf("https://api.bitoff.io/o/cart/items/%d/qty", itemId)
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(jsonStr)))
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error:", bodyResp)
		return nil, fmt.Errorf("failed to increase cart item quantity with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}
