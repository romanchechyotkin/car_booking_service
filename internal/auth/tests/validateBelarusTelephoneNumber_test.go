package tests

import (
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateBelarusTelephoneNumber(t *testing.T) {
	cases := []struct {
		name            string
		telephoneNumber string
		expErr          error
	}{
		{
			name:            "wrong len of telephone number",
			telephoneNumber: "+3755121234124124124",
			expErr:          auth.WrongTelephoneNumberError,
		},
		{
			name:            "wrong country code",
			telephoneNumber: "+374447534067",
			expErr:          auth.WrongTelephoneNumberError,
		},
		{
			name:            "wrong operator code",
			telephoneNumber: "+375457534067",
			expErr:          auth.WrongTelephoneNumberError,
		},
		{
			name:            "space between symbols",
			telephoneNumber: "+37544 7534067",
			expErr:          auth.WrongTelephoneNumberError,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := auth.ValidateBelarusTelephoneNumber(tCase.telephoneNumber)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
