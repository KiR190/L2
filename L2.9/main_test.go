package unpack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"unpack"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		out   string
		isErr bool
	}{
		{"normal case", "a4bc2d5e", "aaaabccddddde", false},
		{"no numbers", "abcd", "abcd", false},
		{"empty string", "", "", false},
		{"starts with number", "45", "", true},
		{"escaped digits", "qwe\\4\\5", "qwe45", false},
		{"escaped digit with repeat", "qwe\\45", "qwe44444", false},
		{"zero repeat", "a0b", "b", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := unpack.Unpack(tc.in)

			if tc.isErr {
				require.Error(t, err, "expected an error for input %q", tc.in)
			} else {
				require.NoError(t, err, "unexpected error for input %q", tc.in)
			}

			assert.Equal(t, tc.out, got, "unexpected output for input %q", tc.in)
		})
	}
}
