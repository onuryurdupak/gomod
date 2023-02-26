package core

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestConditionalInt(t *testing.T) {
	testData := []struct {
		condition bool
		onTrue    int
		onFalse   int
		expected  int
	}{
		{
			condition: false,
			onTrue:    3,
			onFalse:   4,
			expected:  4,
		},
		{
			condition: true,
			onTrue:    -12,
			onFalse:   24,
			expected:  -12,
		},
	}

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(innerT *testing.T) {
			result := ConditionalInt(td.condition, td.onTrue, td.onFalse)
			assert.Equal(innerT, td.expected, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	type testStruct struct{}

	dataValueStruct := testStruct{}
	dataValidPointer := &testStruct{}
	var dataNilPointer *testStruct
	var dataNilInterface interface{}
	var dataValidInterfaceWithNilPointer interface{} = nil
	var dataValidInterfaceWithValidPointer interface{} = &testStruct{}

	testData := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    dataValueStruct,
			expected: false,
		},
		{
			input:    dataValidPointer,
			expected: false,
		},
		{
			input:    dataNilPointer,
			expected: true,
		},
		{
			input:    dataNilInterface,
			expected: true,
		},
		{
			input:    dataValidInterfaceWithNilPointer,
			expected: true,
		},
		{
			input:    dataValidInterfaceWithValidPointer,
			expected: false,
		},
	}

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(innerT *testing.T) {
			result := IsNil(td.input)
			assert.Equal(innerT, td.expected, result)
		})
	}
}
