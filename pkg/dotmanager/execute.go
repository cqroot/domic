package dotmanager

import "github.com/cqroot/dotm/pkg/dotfile"

type ExecuteResult struct {
	Dot   Dot
	HasOp bool
	Err   error
}

func (d DotManager) Apply() []ExecuteResult {
	result := make([]ExecuteResult, 0)

	for _, dot := range d.Dots() {
		hasOp, err := dotfile.Link(&dot)
		result = append(result, ExecuteResult{
			Dot:   dot,
			HasOp: hasOp,
			Err:   err,
		})
	}

	return result
}

func (d DotManager) Revoke() []ExecuteResult {
	result := make([]ExecuteResult, 0)

	for _, dot := range d.Dots() {
		hasOp, err := dotfile.Unlink(&dot)
		result = append(result, ExecuteResult{
			Dot:   dot,
			HasOp: hasOp,
			Err:   err,
		})
	}

	return result
}
