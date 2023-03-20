package tests

import (
	"github.com/romanchechyotkin/car_booking_service/internal/car"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateCarNumbers(t *testing.T) {
	cases := []struct {
		name    string
		numbers string
		expErr  error
	}{
		{
			name:    "word",
			numbers: "Qwerty",
			expErr:  car.WrongCarNumbersLen,
		},
		{
			name:    "empty",
			numbers: "        ",
			expErr:  car.WrongSymbolCarNumbers,
		},
		{
			name:    "wrong region1",
			numbers: "1234AA-0",
			expErr:  car.WrongRegionEnteredCarNumbers,
		},
		{
			name:    "wrong region2",
			numbers: "1234AA-8",
			expErr:  car.WrongRegionEnteredCarNumbers,
		},
		{
			name:    "wrong num part",
			numbers: "a111AA-7",
			expErr:  car.WrongNumbersPartCarNumbers,
		},
		{
			name:    "wrong letters part",
			numbers: "1234Aa-7",
			expErr:  car.WrongLettersPartEnteredCarNumbers,
		},
		{
			name:    "wrong letters part",
			numbers: "1234aA-7",
			expErr:  car.WrongLettersPartEnteredCarNumbers,
		},
		{
			name:    "wrong letters part",
			numbers: "1234aa-7",
			expErr:  car.WrongLettersPartEnteredCarNumbers,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := car.ValidateCarNumbers(tCase.numbers)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
