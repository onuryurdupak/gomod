package os

import (
	"testing"
)

func Test_Run_Command(t *testing.T) {
	type testData struct {
		input        string
		outputsError bool
	}

	tests := []testData{
		{"echo 1", false},
		{"echo 2", false},
		{"echo x", false},
		{"pwd", false},
		{"whoami", false},
		{"afasdsa", true},
		{"zczxcxz", true},
	}

	for _, test := range tests {
		_, err := Exec(test.input)
		hasError := err != nil
		if hasError != test.outputsError {
			t.FailNow()
		}
	}
}
