package tests

import (
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidationForEmptyPasswordAndFullName(t *testing.T) {
	cases := []struct {
		name     string
		password string
		fullName string
		expErr   error
	}{
		{
			name:     "wrong empty password",
			password: "          ",
			fullName: "Qwerty",
			expErr:   auth.WrongEnteredPasswordError,
		},
		{
			name:     "without password contains space symbol 1",
			password: " qwerty",
			fullName: "Qwerty",
			expErr:   auth.WrongEnteredPasswordError,
		},
		{
			name:     "without password contains space symbol 2",
			password: "qwerty ",
			fullName: "Qwerty",
			expErr:   auth.WrongEnteredPasswordError,
		},
		{
			name:     "without password contains space symbol 3",
			password: "qwe rty",
			fullName: "Qwerty",
			expErr:   auth.WrongEnteredPasswordError,
		},
		{
			name:     "without full name",
			password: "qwerty",
			fullName: "           ",
			expErr:   auth.EmptyFullNameError,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := auth.ValidateForEmptyPasswordAndFullName(tCase.password, tCase.fullName)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
