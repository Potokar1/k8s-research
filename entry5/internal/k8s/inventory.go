package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// inventoryRequest makes an http request to a pod to retrieve its inventory
func inventoryRequest(ctx context.Context, url string) (map[string]int, error) {
	url = url + "/inventory"

	// create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// make an HTTP GET request to the inventory endpoint
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve inventory, status code: %d", resp.StatusCode)
	}

	// parse the response body to get the inventory
	var inventory map[string]int
	err = json.NewDecoder(resp.Body).Decode(&inventory)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// GetInventory retrieves the inventory given the query parameters
func GetInventory(ctx context.Context, kingdom, town, shop string) error {

	return nil
}
