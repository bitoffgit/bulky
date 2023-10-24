package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CreateAddress(body map[string]interface{}, authToken string) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.bitoff.io/a/user/address", strings.NewReader(string(jsonBody)))
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
		return nil, fmt.Errorf("failed to create address with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func GetAddresses(authToken string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.bitoff.io/a/user/address", nil)
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

	var bodyResp []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&bodyResp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error:", bodyResp)
		return nil, fmt.Errorf("failed to get addresses with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func GetDefaultAddress(authToken string) (map[string]interface{}, error) {
	addresses, err := GetAddresses(authToken)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	for _, address := range addresses {
		if address["default"].(bool) {
			return address, nil
		}
	}

	return nil, fmt.Errorf("no default address found")
}

func SetAddressToDefault(body map[string]interface{}, authToken string) (map[string]interface{}, error) {
	body["default"] = true

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.bitoff.io/a/user/%v/address/default", body["id"]), strings.NewReader(string(jsonBody)))
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
		return nil, fmt.Errorf("failed to set address to default with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}

func UpdateAddress(body map[string]interface{}, authToken string) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.bitoff.io/a/user/%v/address", body["id"]), strings.NewReader(string(jsonBody)))
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
		return nil, fmt.Errorf("failed to update address with status code %d", resp.StatusCode)
	}

	return bodyResp, nil
}
