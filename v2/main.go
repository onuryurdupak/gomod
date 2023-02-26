package main

import (
	"github.com/google/uuid"
	"github.com/onuryurdupak/gomod/v2/stdout"
)

func main() {
	stdout.SetOutputLevel(stdout.OutputLevelDebug)
	stdout.Output(stdout.OutputLevelDebug, "Starting application.")
	stdout.PrintfStyled("<red>Hello</red> <yellow>World</yellow><green>!</green>\n")
	stdout.OutputWID(stdout.OutputLevelInfo, uuid.NewString(), "Shutting down application.")
}
