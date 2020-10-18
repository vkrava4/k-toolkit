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
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vkrava4/k-toolkit/execution"
	"github.com/vkrava4/k-toolkit/util"
	"github.com/vkrava4/k-toolkit/validation"
	"os"
	"path/filepath"
	"strings"
)

// concatCmd represents the concat command
var concatCmd = &cobra.Command{
	Use:   "concat",
	Short: "Concatenates provided set of files or files in given directory(es)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var executionResult *execution.Result
		if len(args) > 1 {
			executionResult = RunConcatCmd(args[0], args[1])
		} else {
			executionResult = RunConcatCmd("", args[0])
		}

		if executionResult.Success && executionResult.ValidationResult.IsValid {
			fmt.Println()
			fmt.Println("\n" + executionResult.Message)
			os.Exit(0)
		} else {
			executionResult.PrintPretty()
			os.Exit(-1)
		}
	},
}

var (
	fileSuffix  string
	isCascading bool

	greenColor = color.New(color.FgGreen)
	redColor   = color.New(color.FgRed)
)

func RunConcatCmd(sources string, output string) *execution.Result {
	var executionResult = &execution.Result{}
	var sourcesPath []string

	for _, src := range strings.Split(sources, ",") {
		var srcAbsPath, _ = filepath.Abs(strings.TrimSpace(src))
		sourcesPath = append(sourcesPath, srcAbsPath)
	}

	var validationResult = Validate(sourcesPath, output)
	if !validationResult.IsValid {
		executionResult.Success = false
	} else {
		var files, errDirectoriesWalk = util.DirectoriesWalk(sourcesPath, isCascading, fileSuffix)
		if errDirectoriesWalk != nil {
			validationResult.IsValid = false
			validationResult.ValidationErrors = append(validationResult.ValidationErrors, errDirectoriesWalk.Error())
		} else {
			var absOutput, errAbsOutput = filepath.Abs(output)
			if errAbsOutput != nil {
				validationResult.IsValid = false
				validationResult.ValidationErrors = append(validationResult.ValidationErrors, errAbsOutput.Error())
			} else {
				for _, file := range files {
					fmt.Println(file)
				}

				if util.AskConfirmation(fmt.Sprintf("%d files found eligble for concatination. Would you like to procceed?", len(files))) {
					var outputFile, errOutputFileCreate = os.Create(absOutput)

					if errOutputFileCreate != nil {
						validationResult.IsValid = false
						validationResult.ValidationErrors = append(validationResult.ValidationErrors, errOutputFileCreate.Error())
					} else {
						var totalWritten int64
						for _, file := range files {
							var written, errWrite = util.ConcatenateFiles(file, outputFile)
							if errWrite != nil {
								fmt.Println(fmt.Sprintf("[ %s ] %s", file, redColor.Sprint("ERROR")))

								executionResult.Success = false
								executionResult.Message = errWrite.Error()
								break
							} else {
								totalWritten += written
								fmt.Println(fmt.Sprintf("[ %s ] %s", greenColor.Sprint("OK"), fmt.Sprintf("Merged %d bytes from: %s", written, file)))
							}
						}
						_ = outputFile.Close()

						executionResult.Success = true
						executionResult.Message = fmt.Sprintf("[ %s ] %s", greenColor.Sprint("OK"), fmt.Sprintf("Total merged %d bytes to output file: %s", totalWritten, absOutput))
					}
				}
			}
		}
	}

	executionResult.ValidationResult = validationResult
	return executionResult
}

func Validate(sourcesPath []string, output string) *validation.Result {
	return validation.Init().
		ShouldNotBeNilOrBlank(sourcesPath).
		ShouldExistAsFileOrDirectory(sourcesPath).
		ShouldNotBeBlankS(output).
		ShouldNotExistInPath(output).
		ShouldContainAnyFilesWithPattern(sourcesPath, isCascading, fileSuffix)
}

func init() {
	rootCmd.AddCommand(concatCmd)

	concatCmd.Flags().StringVarP(&fileSuffix, "file-suffix", "s", "", "A suffix of files which should be included for concatenation")
	concatCmd.Flags().BoolVarP(&isCascading, "cascade", "c", false, "A suffix of files which should be included for concatenation")
}
