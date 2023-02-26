package stdout

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCreateOutputText(t *testing.T) {
	SetOutputLevel(OutputLevelDebug)

	testData := []struct {
		outputLevel    OutputLevel
		sessionID      string
		format         string
		variadicArgs   []interface{}
		expectedOutput string
	}{
		{
			outputLevel:    OutputLevelDebug,
			sessionID:      "",
			format:         "content",
			variadicArgs:   nil,
			expectedOutput: "<b><cyan>[DEBUG]</cyan></b>: content",
		},
		{
			outputLevel:    OutputLevelInfo,
			sessionID:      "",
			format:         "content",
			variadicArgs:   nil,
			expectedOutput: "<b><green>[INFO]</green></b>: content",
		},
		{
			outputLevel:    OutputLevelWarn,
			sessionID:      "",
			format:         "content",
			variadicArgs:   nil,
			expectedOutput: "<b><yellow>[WARN]</yellow></b>: content",
		},
		{
			outputLevel:    OutputLevelError,
			sessionID:      "12345",
			format:         "content",
			variadicArgs:   nil,
			expectedOutput: "<b><red>[ERROR]</red>[12345]</b>: content",
		},
	}

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(innerT *testing.T) {
			var result string
			// Handling variadic parameters.
			// Note that passing "args = nil" means passing a slice with a single nil element.
			if td.variadicArgs != nil {
				result = createOutputText(td.outputLevel, td.sessionID, td.format, td.variadicArgs)
			} else {
				result = createOutputText(td.outputLevel, td.sessionID, td.format)
			}

			assert.Equal(innerT, td.expectedOutput, result)
		})
	}
}
