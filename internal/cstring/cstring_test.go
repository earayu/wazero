package cstring

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNullTerminatedStrings(t *testing.T) {
	emptyWASIStringArray := &NullTerminatedStrings{NullTerminatedValues: [][]byte{}}
	tests := []struct {
		name     string
		input    []string
		expected *NullTerminatedStrings
	}{
		{
			name:     "nil",
			expected: emptyWASIStringArray,
		},
		{
			name:     "none",
			input:    []string{},
			expected: emptyWASIStringArray,
		},
		{
			name:  "two",
			input: []string{"a", "bc"},
			expected: &NullTerminatedStrings{
				NullTerminatedValues: [][]byte{
					{'a', 0},
					{'b', 'c', 0},
				},
				TotalBufSize: 5,
			},
		},
		{
			name:  "two and empty string",
			input: []string{"a", "", "bc"},
			expected: &NullTerminatedStrings{
				NullTerminatedValues: [][]byte{
					{'a', 0},
					{0},
					{'b', 'c', 0},
				},
				TotalBufSize: 6,
			},
		},
		{
			name: "utf-8",
			// "😨", "🤣", and "️🏃‍♀️" have 4, 4, and 13 bytes respectively
			input: []string{"😨🤣🏃\u200d♀️", "foo", "bar"},
			expected: &NullTerminatedStrings{
				NullTerminatedValues: [][]byte{
					[]byte("😨🤣🏃\u200d♀️\x00"),
					{'f', 'o', 'o', 0},
					{'b', 'a', 'r', 0},
				},
				TotalBufSize: 30,
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			s, err := NewNullTerminatedStrings(100, "", tc.input...)
			require.NoError(t, err)
			require.Equal(t, tc.expected, s)
		})
	}
}

func TestNewNullTerminatedStrings_Errors(t *testing.T) {
	t.Run("invalid utf-8", func(t *testing.T) {
		_, err := NewNullTerminatedStrings(100, "arg", "\xff\xfe\xfd", "foo", "bar")
		require.EqualError(t, err, "arg[0] is not a valid UTF-8 string")
	})
	t.Run("arg[0] too large", func(t *testing.T) {
		_, err := NewNullTerminatedStrings(1, "arg", "a", "bc")
		require.EqualError(t, err, "arg[0] will exceed max buffer size 1")
	})
	t.Run("empty arg too large due to null terminator", func(t *testing.T) {
		_, err := NewNullTerminatedStrings(2, "arg", "a", "", "bc")
		require.EqualError(t, err, "arg[1] will exceed max buffer size 2")
	})
}