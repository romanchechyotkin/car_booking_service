package tests

import (
	"github.com/romanchechyotkin/car_booking_service/internal/car"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateForEmptyStrings(t *testing.T) {
	cases := []struct {
		name     string
		brand    string
		model    string
		expBrand string
		expModel string
		expErr   error
	}{
		{
			name:     "both empty 1",
			brand:    "",
			model:    "",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "both empty 2",
			brand:    "    ",
			model:    "    ",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "one empty",
			brand:    "qwewreq",
			model:    "",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "one empty",
			brand:    "qwewreq",
			model:    "    ",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "one empty",
			brand:    "",
			model:    "  1",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "one empty",
			brand:    "       ",
			model:    "qwe",
			expBrand: "",
			expModel: "",
			expErr:   car.EmptyString,
		},
		{
			name:     "one empty",
			brand:    " 1 ",
			model:    "qwe",
			expBrand: "1",
			expModel: "qwe",
			expErr:   nil,
		},
		{
			name:     "one empty",
			brand:    "       1 ",
			model:    " land cruser ",
			expBrand: "1",
			expModel: "land cruser",
			expErr:   nil,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			brand, model, _ := car.ValidateForEmptyStrings(tCase.brand, tCase.model)
			require.Equal(t, tCase.expModel, model)
			require.Equal(t, tCase.expBrand, brand)
		})
	}
}
