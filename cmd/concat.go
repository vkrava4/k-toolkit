/*
Copyright © 2020 Vlad Krava <vkrava4@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vkrava4/k-toolkit/execution"
	"github.com/vkrava4/k-toolkit/validation"
	"path/filepath"
	"strings"
)

// concatCmd represents the concat command
var concatCmd = &cobra.Command{
	Use:   "concat",
	Short: "Concatenates provided set of files or files in given directory(es)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var executionResult *execution.Result
		if len(args) > 1 {
			executionResult = RunConcatCmd(args[0], args[1])
		} else {
			executionResult = RunConcatCmd("", args[0])
		}

		if executionResult.Success && executionResult.ValidationResult.IsValid {
			//	todo print success
		} else {

		}
	},
}

var (
	fileSuffix string
)

func RunConcatCmd(sources string, output string) *execution.Result {
	var executionResult = &execution.Result{}

	var sourcesSlice []string
	for _, src := range strings.Split(sources, ",") {
		var srcAbsPath, _ = filepath.Abs(strings.TrimSpace(src))
		sourcesSlice = append(sourcesSlice, srcAbsPath)
	}

	var validationResult = validation.Init().
		ShouldNotBeNilOrBlank(sourcesSlice).
		ShouldExistAsFileOrDirectory(sourcesSlice).
		ShouldNotBeBlankS(output).
		ShouldNotExistInPath(output).
		ShouldContainAnyFilesWithPattern(sourcesSlice, true, fileSuffix)

	if !validationResult.IsValid {
		executionResult.Success = false
	}

	executionResult.ValidationResult = validationResult
	return executionResult
}

func init() {
	rootCmd.AddCommand(concatCmd)

	concatCmd.Flags().StringVarP(&fileSuffix, "file-suffix", "s", fileSuffix, "A suffix of files which should be included for concatenation")
}
