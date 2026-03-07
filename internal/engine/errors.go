package engine

import "fmt"

type IndError struct {
	Code       string
	Command    string
	Message    string
	Suggestion string
}

func (e *IndError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *IndError) Render() string {
	if e == nil {
		return ""
	}

	msg := fmt.Sprintf("Error: %s\nCommand failed: %s", e.Code, e.Command)
	if e.Message != "" {
		msg += "\nReason: " + e.Message
	}
	if e.Suggestion != "" {
		msg += "\nSuggestion: " + e.Suggestion
	}
	return msg
}

func missingCommandError() *IndError {
	return &IndError{
		Code:       "IND_ERR_001",
		Command:    "ind",
		Message:    "missing command",
		Suggestion: "run \"ind docs\"",
	}
}

func unknownCommandError(command string) *IndError {
	return &IndError{
		Code:       "IND_ERR_002",
		Command:    command,
		Message:    "unknown INDUS command",
		Suggestion: "run \"ind docs\" and review docs/commands.html",
	}
}

func invalidArgumentError(command, reason string) *IndError {
	return &IndError{
		Code:       "IND_ERR_003",
		Command:    command,
		Message:    reason,
		Suggestion: "run \"ind doctor\"",
	}
}

func commandFailedError(command string, err error) *IndError {
	return &IndError{
		Code:       "IND_ERR_004",
		Command:    command,
		Message:    err.Error(),
		Suggestion: "run \"ind doctor\"",
	}
}

func registryError(err error) *IndError {
	return &IndError{
		Code:       "IND_ERR_005",
		Command:    "ind",
		Message:    err.Error(),
		Suggestion: "verify core/commands/registry.json and rerun \"ind doctor\"",
	}
}

func panicError(command string, recovered any) *IndError {
	return &IndError{
		Code:       "IND_ERR_006",
		Command:    command,
		Message:    fmt.Sprintf("panic recovered: %v", recovered),
		Suggestion: "run \"ind doctor\"",
	}
}
