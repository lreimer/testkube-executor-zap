package runner

import (
	"errors"
	"os"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

type Params struct {
	Datadir string // RUNNER_DATADIR
}

func NewRunner() *ZapRunner {
	return &ZapRunner{
		params: Params{
			Datadir: os.Getenv("RUNNER_DATADIR"),
		},
	}
}

type ZapRunner struct {
	params Params
}

func (r *ZapRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// check that the datadir exists
	_, err = os.Stat(r.params.Datadir)
	if errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	return testkube.ExecutionResult{
		Status: testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output: "success",
	}, nil
}
