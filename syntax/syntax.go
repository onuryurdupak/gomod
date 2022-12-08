package syntax

func ConditionalInt(condition bool, onTrue, onFalse int) int {
	if condition {
		return onTrue
	} else {
		return onFalse
	}
}
