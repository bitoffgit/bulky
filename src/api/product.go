package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetProduct(asin string) (map[string]interface{}, int, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.bitoff.io/p/products/%s", asin), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	var bodyResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}

	if resp.StatusCode == 302 {
		return nil, resp.StatusCode, fmt.Errorf(
			"Amazon has a problem with product %s, It redirects to %s \n https://bitoff.io/product/%s",
			asin,
			bodyResp["productId"],
			bodyResp["productId"],
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("failed to get product with status code %d", resp.StatusCode)
	}

	return bodyResp, resp.StatusCode, nil
}

func UpdateProduct(asin string) (map[string]interface{}, int, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.bitoff.io/p/products/%s/price", asin), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	var bodyResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error:", bodyResp)
		return nil, resp.StatusCode, fmt.Errorf("failed to update product with status code %d", resp.StatusCode)
	}

	return bodyResp, resp.StatusCode, nil
}
