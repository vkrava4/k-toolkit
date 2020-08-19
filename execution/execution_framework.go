package execution

import "github.com/vkrava4/k-toolkit/validation"

type Result struct {
	Success          bool
	ValidationResult *validation.Result
}
