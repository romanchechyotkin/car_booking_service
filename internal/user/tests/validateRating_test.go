package tests

import (
	"github.com/romanchechyotkin/car_booking_service/internal/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateCarNumbers(t *testing.T) {
	cases := []struct {
		name   string
		rating float32
		expErr error
	}{
		{
			name:   "less",
			rating: 0,
			expErr: user.WrongRating,
		},
		{
			name:   "more",
			rating: 20,
			expErr: user.WrongRating,
		},
		{
			name:   "more",
			rating: 20.0,
			expErr: user.WrongRating,
		},
		{
			name:   "less",
			rating: 0,
			expErr: user.WrongRating,
		},
		{
			name:   "less",
			rating: 0.0,
			expErr: user.WrongRating,
		},
		{
			name:   "5.1",
			rating: 0.9,
			expErr: user.WrongRating,
		},
		{
			name:   "less",
			rating: 44.1,
			expErr: user.WrongRating,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := user.ValidateRating(tCase.rating)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
