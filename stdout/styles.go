package stdout

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	seekModeOpener byte = 0
	seekModeCloser byte = 1

	SpecifierBold      Specifier = "b"
	SpecifierUnderline Specifier = "u"
	SpecifierRed       Specifier = "red"
	SpecifierGreen     Specifier = "green"
	SpecifierBlue      Specifier = "blue"
	SpecifierYellow    Specifier = "yellow"
	SpecifierCyan      Specifier = "cyan"
	SpecifierMagenta   Specifier = "magenta"
)

type Specifier string

type styleData struct {
	specifier Specifier
	style     *color.Color
}

var styles = []styleData{
	{SpecifierBold, color.New(color.Bold)},
	{SpecifierUnderline, color.New(color.Underline)},
	{SpecifierRed, color.New(color.FgRed)},
	{SpecifierGreen, color.New(color.FgGreen)},
	{SpecifierBlue, color.New(color.FgBlue)},
	{SpecifierYellow, color.New(color.FgYellow)},
	{SpecifierCyan, color.New(color.FgCyan)},
	{SpecifierMagenta, color.New(color.FgMagenta)},
}

// ProcessStyle formats input string according to html-like tag placements.
// Output string can be sent to stdout which will be displayed as a style applied text.
//
// Examples:
//
// <b>This will be printed in bold.</b>
//
// <u>This will be printed underlined./<u>
//
// <b><u>This will be printed both bold and underlined.</u></b>
//
// <yellow>This will be printed in yellow.</yellow>
func ProcessStyle(in string) (string, error) {
	var err error
	for _, sd := range styles {
		in, err = processStyle(in, sd)
		if err != nil {
			return "", err
		}
	}
	return in, nil
}

func PrintfStyled(format string, args ...interface{}) {
	rawString := fmt.Sprintf(format, args...)
	processedString, _ := ProcessStyle(rawString)
	fmt.Print(processedString)
}

// RemoveStyle removes style tags from input string and returns it.
func RemoveStyle(in string) string {
	for _, sd := range styles {
		in = removeStyle(in, sd)
	}
	return in
}

func getOpener(specifier Specifier) string {
	return "<" + string(specifier) + ">"
}

func getCloser(specifier Specifier) string {
	return "</" + string(specifier) + ">"
}

func processStyle(in string, styleData styleData) (string, error) {
	sb := strings.Builder{}
	var builtString string

	opener := getOpener(styleData.specifier)
	closer := getCloser(styleData.specifier)

	openerSize := len(opener)
	closerSize := len(closer)

	cursor := 0
	lastSplit := 0
	seekMode := seekModeOpener
	for {
		var selection string
		var endRange int
		if seekMode == seekModeOpener {
			endRange = cursor + openerSize
		} else if seekMode == seekModeCloser {
			endRange = cursor + closerSize
		} else {
			return "", fmt.Errorf("unexpected seek mode: %d", seekMode)
		}

		if endRange > len(in) {
			if seekMode == seekModeOpener {
				sb.WriteString(in[lastSplit:])
			} else if seekMode == seekModeCloser {
				sb.WriteString(styleData.style.Sprint(in[lastSplit:]))
			} else {
				return "", fmt.Errorf("unexpected seek mode: %d", seekMode)
			}

			builtString = sb.String()
			builtString = strings.Replace(builtString, opener, "", -1)
			builtString = strings.Replace(builtString, closer, "", -1)

			break
		}

		selection = in[cursor:endRange]

		if seekMode == seekModeOpener && selection == opener {
			appendText := in[lastSplit:endRange]
			sb.WriteString(appendText)
			lastSplit = endRange
			seekMode = seekModeCloser
			cursor += openerSize
		} else if seekMode == seekModeCloser && selection == closer {
			appendText := styleData.style.Sprintf("%s", in[lastSplit:endRange])
			sb.WriteString(appendText)
			lastSplit = endRange
			seekMode = seekModeOpener
			cursor += closerSize
		} else {
			cursor++
		}
	}
	return builtString, nil
}

func removeStyle(in string, styleData styleData) string {
	opener := getOpener(styleData.specifier)
	closer := getCloser(styleData.specifier)
	return strings.Replace(strings.Replace(in, opener, "", -1), closer, "", -1)
}
