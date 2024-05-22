package types

import (
	"database/sql/driver"
)

type Payment string

func (p Payment) String() string {
	return string(p)
}

const (
	Cash       Payment = "cash"
	CreditCard Payment = "credit-card"
	Visa       Payment = "visa"
	Mastercard Payment = "mastercard"
	Mir        Payment = "mir"
)

type PaymentArray []Payment

func (p PaymentArray) Value() (driver.Value, error) {
	var res string

	res += "{"
	for _, payment := range p {
		res = res + payment.String() + ","
	}
	res = res[:len(res)-1]
	res += "}"

	return res, nil
}

func (p PaymentArray) String() []string {
	res := make([]string, 0, len(p))

	for _, v := range p {
		res = append(res, v.String())
	}

	return res
}