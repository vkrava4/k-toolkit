package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		var srcAbsPath, errAbsPath = filepath.Abs(src)
		if errAbsPath != nil {
			invalidate(v, errAbsPath.Error())
		}

		if _, isNotExistErr := os.Stat(srcAbsPath); os.IsNotExist(isNotExistErr) {
			invalidate(v, fmt.Sprintf("File or directory: '%s' can not be found", srcAbsPath))
		}
	}

	return v
}

func (v *Result) ShouldNotExistInPath(string) *Result {
	if !v.IsValid {
		return v
	}

	//panic("implement me")

	return v
}

func (v *Result) PrintPretty() *Result {
	panic("implement me")
}

func invalidate(v *Result, message string) {
	v.IsValid = false
	v.ValidationErrors = append(v.ValidationErrors, message)
}
