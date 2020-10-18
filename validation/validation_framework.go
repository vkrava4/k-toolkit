package validation

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vkrava4/k-toolkit/util"
	"os"
	"strings"
)

var (
	redColor = color.New(color.FgRed)
)

type Result struct {
	IsValid          bool
	ValidationErrors []string
}

type Framework interface {
	ShouldNotBeNilOrBlank()
	ShouldNotBeNilOrBlankS()
	ShouldExistAsFileOrDirectory()
	ShouldNotExistAsFile()
	ShouldContainAnyFilesWithPattern()

	PrintPretty()
}

func Init() *Result {
	return &Result{IsValid: true}
}

func (v *Result) ShouldNotBeBlankS(source string) *Result {
	if !v.IsValid {
		return v
	}

	if len(strings.TrimSpace(source)) == 0 {
		invalidate(v, "Argument(s) can not be empty")
	}

	return v
}

func (v *Result) ShouldNotBeNilOrBlank(sources []string) *Result {
	if !v.IsValid {
		return v
	}

	if sources != nil {
		for _, src := range sources {
			if len(strings.TrimSpace(src)) == 0 {
				invalidate(v, "Argument(s) can not be empty")
			}
		}
	} else {
		invalidate(v, "Argument(s) can not be set to 'nil' value or be empty")
	}

	return v
}

func (v *Result) ShouldExistAsFileOrDirectory(sources []string) *Result {
	if !v.IsValid {
		return v
	}

	for _, src := range sources {
		if _, isNotExistErr := os.Stat(src); os.IsNotExist(isNotExistErr) {
			invalidate(v, fmt.Sprintf("File or directory: '%s' can not be found", src))
		}
	}

	return v
}

func (v *Result) ShouldNotExistInPath(path string) *Result {
	if !v.IsValid {
		return v
	}

	if util.FileExists(path) {
		invalidate(v, fmt.Sprintf("Output file %s already exists", path))
	}

	return v
}

func (v *Result) ShouldContainAnyFilesWithPattern(sources []string, cascading bool, pattern string) *Result {
	if !v.IsValid {
		return v
	}

	var files []string
	var err error

	for _, src := range sources {
		if util.FileExists(src) {
			files = append(files, src)
		} else {
			var directoriesWalkFiles, errDirectoriesWalk = util.DirectoriesWalk(sources, cascading, pattern)
			err = errDirectoriesWalk
			if len(directoriesWalkFiles) > 0 {
				files = append(files, directoriesWalkFiles...)
			}
		}
	}

	if err != nil {
		invalidate(v, err.Error())
	} else if len(files) == 1 {
		invalidate(v, "Provided directory(es) contain only one file matching concatenation criteria")
	} else if len(files) == 0 {
		invalidate(v, "Provided directory(es) does not contain any files matched for concatenation")
	}

	return v
}

func (v *Result) PrintPretty() {
	for _, message := range v.ValidationErrors {
		fmt.Println(redColor.Sprintf(" - %s", message))
	}
}

func invalidate(v *Result, message string) {
	v.IsValid = false
	v.ValidationErrors = append(v.ValidationErrors, message)
}
