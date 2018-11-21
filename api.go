package bitflyer

import (
	"encoding/json"
)

type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func (c *Client) GetBalance() ([]Balance, error) {
	r := &Request{
		Method:   "GET",
		Endpoint: "/v1/me/getbalance",
	}
	var res []Balance
	data, err := c.call(r)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
