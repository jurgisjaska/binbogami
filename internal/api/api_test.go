package api

import (
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Response
	}{
		{
			name:     "empty input",
			input:    "",
			expected: &Response{Status: statusError, Message: ""},
		},
		{
			name:     "normal case",
			input:    "ERROR MESSAGE",
			expected: &Response{Status: statusError, Message: "error message"},
		},
		{
			name:     "with numbers",
			input:    "ERROR123 MESSAGE",
			expected: &Response{Status: statusError, Message: "error123 message"},
		},
		{
			name:     "with special characters",
			input:    "ERROR!@# MESSAGE",
			expected: &Response{Status: statusError, Message: "error!@# message"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Error(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}
