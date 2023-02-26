package stdout

import (
	"fmt"
)

type OutputLevel int

const (
	OutputLevelDebug OutputLevel = 0
	OutputLevelInfo  OutputLevel = 1
	OutputLevelWarn  OutputLevel = 2
	OutputLevelError OutputLevel = 3
)

var currentOutputLevel OutputLevel = OutputLevelInfo

func GetOutputLevel() OutputLevel {
	return currentOutputLevel
}

func SetOutputLevel(newLevel OutputLevel) {
	currentOutputLevel = newLevel
}

// Output prints the input content to stdout, highlighting it according to input level.
//
// Example output: [INFO]: Your input.
func Output(level OutputLevel, format string, args ...interface{}) {
	if level < currentOutputLevel {
		return
	}
	styled := createOutputText(level, "", format, args...) + "\n"
	PrintfStyled(styled)
}

// OutputWID outputs the input content with the given session ID.
//
// Example output: [INFO][64981af8-1b83-477b-a474-7dea8ca76b8f]: Your input.
func OutputWID(level OutputLevel, sessionID, format string, args ...interface{}) {
	if level < currentOutputLevel {
		return
	}
	styled := createOutputText(level, sessionID, format, args...) + "\n"
	PrintfStyled(styled)
}

func createOutputText(level OutputLevel, sessionID, format string, args ...interface{}) string {
	var mainHeader string
	var specifier Specifier

	switch level {
	case OutputLevelDebug:
		mainHeader = "[DEBUG]"
		specifier = SpecifierCyan
	case OutputLevelInfo:
		mainHeader = "[INFO]"
		specifier = SpecifierGreen
	case OutputLevelWarn:
		mainHeader = "[WARN]"
		specifier = SpecifierYellow
	case OutputLevelError:
		mainHeader = "[ERROR]"
		specifier = SpecifierRed
	default:
		mainHeader = "[UNKNOWN]"
		specifier = SpecifierCyan
	}

	boldOpener := getOpener(SpecifierBold)
	boldCloser := getCloser(SpecifierBold)
	specifierOpener := getOpener(specifier)
	specifierCloser := getCloser(specifier)

	var header string
	if sessionID == "" {
		header = boldOpener + specifierOpener + mainHeader + specifierCloser + boldCloser + ": "
	} else {
		header = boldOpener + specifierOpener + mainHeader + specifierCloser + "[" + sessionID + "]" + boldCloser + ": "
	}

	styled := fmt.Sprintf(header+format, args...)
	return styled
}
