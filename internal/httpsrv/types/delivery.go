package types

import (
	"database/sql/driver"
)

type Delivery string

func (d Delivery) String() string {
	return string(d)
}

const (
	PickUp Delivery = "pick-up"
	Waiter Delivery = "waiter"
)

type DeliveryArray []Delivery

func (d DeliveryArray) Value() (driver.Value, error) {
	var res string

	res += "{"
	for _, delivery := range d {
		res = res + delivery.String() + ","
	}
	res = res[:len(res)-1]
	res += "}"

	return res, nil
}

func (d DeliveryArray) String() []string {
	res := make([]string, 0, len(d))

	for _, v := range d {
		res = append(res, v.String())
	}

	return res
}