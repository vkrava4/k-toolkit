package cmd

import (
	"fmt"
	"github.com/vkrava4/k-toolkit/execution"
	"github.com/vkrava4/k-toolkit/util"
	"github.com/vkrava4/k-toolkit/validation"
	"os"
	"testing"
)

var (
	existingTestDirPath                = "test_root"
	existingTestDirPathWithoutFiles    = "test_root/test_empty_dir"
	existingTestDirPathWithOneFile     = "test_root/test_dir_with_one_file"
	existingTestDirPathWithoutFilesSub = "test_root/test_empty_dir/sub_dir"
	existingOutputFilePath             = "test_root/existing.out1"
	existingOutput2FilePath            = "test_root/existing1.out"
	existingOneFileInDirectory         = existingTestDirPathWithOneFile + "/existing_one_file.out"
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
			name: "Should return validation.Result.IsValid=true sources blank (the current directory will be picked as source)",
			args: args{sources: " ", output: "test.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          true,
				ValidationErrors: []string{},
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

		// 3
		{
			name: "Should return validation.Result.IsValid=true when source directory exists with several files and output file does not",
			args: args{sources: existingTestDirPath, output: "test.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          true,
				ValidationErrors: []string{},
			}},
		},

		// 4
		{
			name: "Should return validation.Result.IsValid=false with corresponding messages when output file path is blank",
			args: args{sources: existingTestDirPath, output: " "},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid: false,
				ValidationErrors: []string{
					"Argument(s) can not be empty",
				},
			}},
		},

		// 5
		{
			name: "Should return validation.Result.IsValid=false with corresponding messages when output file is already exist",
			args: args{sources: existingTestDirPath, output: existingOutputFilePath},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid: false,
				ValidationErrors: []string{
					fmt.Sprintf("Output file %s already exists", existingOutputFilePath),
				},
			}},
		},

		// 6
		{
			name: "Should return validation.Result.IsValid=false with corresponding messages when source directory and its sub-directories have no files",
			args: args{sources: existingTestDirPathWithoutFiles, output: "new.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          false,
				ValidationErrors: []string{"Provided directory(es) does not contain any files matched for concatenation"},
			}},
		},

		// 7
		{
			name: "Should return validation.Result.IsValid=false with corresponding messages when source directory and its sub-directories have oly one file",
			args: args{sources: existingTestDirPathWithOneFile, output: "new.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          false,
				ValidationErrors: []string{"Provided directory(es) contain only one file matching concatenation criteria"},
			}},
		},

		// 3
		{
			name: "Should return validation.Result.IsValid=true when source files exist and output file does not",
			args: args{sources: existingOutputFilePath + "," + existingOneFileInDirectory, output: "test.out"},
			result: &execution.Result{Success: false, ValidationResult: &validation.Result{
				IsValid:          true,
				ValidationErrors: []string{},
			}},
		},
	}

	setup()
	defer destroy()

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

func setup() {
	destroy()

	_ = os.MkdirAll(existingTestDirPath, util.DefaultDirPermMode)
	_ = os.MkdirAll(existingTestDirPathWithoutFiles, util.DefaultDirPermMode)
	_ = os.MkdirAll(existingTestDirPathWithoutFilesSub, util.DefaultDirPermMode)
	_ = os.MkdirAll(existingTestDirPathWithOneFile, util.DefaultDirPermMode)
	_, _ = os.Create(existingOutputFilePath)
	_, _ = os.Create(existingOutput2FilePath)
	_, _ = os.Create(existingOneFileInDirectory)
}

func destroy() {
	_ = os.Remove(existingOutputFilePath)
	_ = os.Remove(existingOutput2FilePath)
	_ = os.Remove(existingOneFileInDirectory)
	_ = os.RemoveAll(existingTestDirPathWithOneFile)
	_ = os.RemoveAll(existingTestDirPathWithoutFiles)
	_ = os.RemoveAll(existingTestDirPath)
}
