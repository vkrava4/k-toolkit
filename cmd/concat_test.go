package cmd

import (
	"fmt"
	"github.com/vkrava4/k-toolkit/execution"
	"github.com/vkrava4/k-toolkit/util"
	"github.com/vkrava4/k-toolkit/validation"
	"testing"
)

func TestRunConcatCmd(t *testing.T) {
	type args struct {
		sources string
		output  string
	}
	tests := []struct {
		name   string
		result *execution.Result
		args   args
	}{
		// 1
		{
			name: "Should return validation.Result.IsValid=false with corresponding when sources blank",
			args: args{sources: " ", output: "test.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          false,
				ValidationErrors: []string{"Argument(s) can not be empty"},
			}},
		},

		// 2
		{
			name: "Should return validation.Result.IsValid=false with corresponding messages when neither source directory nor file can be found",
			args: args{sources: "/no/NOT_FOUND.XTX, /NO", output: "test.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid: false,
				ValidationErrors: []string{
					fmt.Sprintf("File or directory: '%s' can not be found", "/no/NOT_FOUND.XTX"),
					fmt.Sprintf("File or directory: '%s' can not be found", "/NO"),
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualResult = RunConcatCmd(tt.args.sources, tt.args.output)
			if actualResult.Success != tt.result.Success {
				t.Error(tt, fmt.Sprintf("Actual: '%t', Expected: '%t'", actualResult.Success, tt.result.Success))
			}

			if actualResult.ValidationResult.IsValid != tt.result.ValidationResult.IsValid {
				t.Error(tt, fmt.Sprintf("Actual: '%t', Expected: '%t'", actualResult.ValidationResult.IsValid, tt.result.ValidationResult.IsValid))
			}

			if !util.EqualSliceS(actualResult.ValidationResult.ValidationErrors, tt.result.ValidationResult.ValidationErrors) {
				t.Error(tt, fmt.Sprintf("Actual: '%q', Expected: %q", actualResult.ValidationResult.ValidationErrors, tt.result.ValidationResult.ValidationErrors))
			}
		})
	}
}
