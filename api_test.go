package bitflyer

import (
	"fmt"
	"testing"
)

func TestGetBalance(t *testing.T) {
	c := NewClient("apikey", "seceretkey")
	b, err := c.GetBalance()
	if err != nil {
		fmt.Errorf("err: %s", err.Error())
	}
	fmt.Printf("%+v", b)
}
