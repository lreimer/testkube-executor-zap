package runner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

const ApiScan = `api:
   target: https://www.qaware.de`

func TestRun(t *testing.T) {
	// setup
	tempDir := os.TempDir()
	os.Setenv("RUNNER_DATADIR", tempDir)
	os.Setenv("ZAP_HOME", "../../zap/")

	t.Run("Run successful API scan", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.TestName = "simple-api-scan"
		execution.TestType = "zap/api"
		execution.Content = testkube.NewStringTestContent("")
		writeTestContent(t, tempDir, "../../examples/test-api-pass.yaml")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
		assert.Len(t, result.Steps, 2)
		assert.Equal(t, result.Steps[0].Name, "PASS: Vulnerable JS Library [10003]")
		assert.Equal(t, result.Steps[0].Status, "success")
	})

	t.Run("Run API scan with PASS and WARN", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.TestName = "warn-api-scan"
		execution.TestType = "zap/api"
		execution.Content = testkube.NewStringTestContent("")
		writeTestContent(t, tempDir, "../../examples/test-api-warn.yaml")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
		assert.Len(t, result.Steps, 2)
		assert.Equal(t, result.Steps[1].Name, "WARN-NEW: Re-examine Cache-control Directives [10015] x 12 ")
		assert.Equal(t, result.Steps[1].Status, "success")
	})

	t.Run("Run API scan with WARN and FailOnWarn", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.TestName = "fail-on-warn-api-scan"
		execution.TestType = "zap/api"
		execution.Content = testkube.NewStringTestContent("")
		writeTestContent(t, tempDir, "../../examples/test-api-fail-on-warn.yaml")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.Error(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusError)
		assert.Len(t, result.Steps, 2)
		assert.Equal(t, result.Steps[1].Name, "WARN-NEW: Re-examine Cache-control Directives [10015] x 12 ")
		assert.Equal(t, result.Steps[1].Status, "error")
	})

	t.Run("Run API scan with FAIL", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.TestName = "fail-api-scan"
		execution.TestType = "zap/api"
		execution.Content = testkube.NewStringTestContent("")
		writeTestContent(t, tempDir, "../../examples/test-api-fail.yaml")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.Error(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusError)
		assert.Len(t, result.Steps, 1)
		assert.Equal(t, result.Steps[0].Name, "FAIL: Unknown issue")
		assert.Equal(t, result.Steps[0].Status, "error")
	})
}

func writeTestContent(t *testing.T, dir string, configFile string) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		assert.FailNow(t, "Unable to read ZAP config file")
	}

	err = ioutil.WriteFile(filepath.Join(dir, "test-content"), data, 0644)
	if err != nil {
		assert.FailNow(t, "Unable to write ZAP test-content file")
	}
}
