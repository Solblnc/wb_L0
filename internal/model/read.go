package model

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewOrder(filename string) (*Order, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot reas order from file %s, error: %s", filename, err)
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal order from %s, error: %s", filename, err)
	}

	return &order, nil
}
