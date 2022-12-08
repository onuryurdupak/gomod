package stdout

import "fmt"

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

func OutputDebug(format string, args ...interface{}) {
	if OutputLevelDebug < currentOutputLevel {
		return
	}
	styled := fmt.Sprintf("<b><cyan>Info:</cyan></b> "+format, args...)
	styled += "\n"
	PrintfStyled(styled)
}

func OutputInfo(format string, args ...interface{}) {
	if OutputLevelInfo < currentOutputLevel {
		return
	}
	styled := fmt.Sprintf("<b><green>Info:</green></b> "+format, args...)
	styled += "\n"
	PrintfStyled(styled)
}

func OutputWarn(format string, args ...interface{}) {
	if OutputLevelWarn < currentOutputLevel {
		return
	}
	styled := fmt.Sprintf("<b><yellow>Warn:</yellow></b> "+format, args...)
	styled += "\n"
	PrintfStyled(styled)
}

func OutputError(format string, args ...interface{}) {
	if OutputLevelError < currentOutputLevel {
		return
	}
	styled := fmt.Sprintf("<b><red>Error:</red></b> "+format, args...)
	styled += "\n"
	PrintfStyled(styled)
}
