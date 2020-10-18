package execution

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vkrava4/k-toolkit/validation"
)

type Result struct {
	Success          bool
	Message          string
	ValidationResult *validation.Result
}

var (
	redColor = color.New(color.FgRed)
)

func (v *Result) PrintPretty() {
	if !v.ValidationResult.IsValid && len(v.ValidationResult.ValidationErrors) > 0 {

		fmt.Println(redColor.Sprint("Execution terminated due to validation errors: \n"))
		v.ValidationResult.PrintPretty()
		fmt.Println()
	} else if !v.Success {
		fmt.Println(redColor.Sprint("Execution failed: \n"))
		fmt.Println(redColor.Sprintf(" - %s\n", v.Message))
	}
}
